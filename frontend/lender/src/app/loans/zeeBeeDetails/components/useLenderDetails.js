'use client';
import { Button, Col, Form, message, Space, Steps, Tabs } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { createZeeBeeSteps, zeeBeeDetailsTabs, lenderDetailsSteps } from '@/enum/loan';
import { createRequest } from '@/lib/fetch';
import ApplicantInfo from '../../components/CreateZeeBee/ApplicantInfo';
import BorrowerInfo from '../../components/CreateZeeBee/BorrowerInfo';
import PropertyInfo from '../../components/CreateZeeBee/PropertyInfo';
import Loan from '../../components/CreateZeeBee/Loan';
import Experience from '../../components/CreateZeeBee/Experience';
import ThirdPartyService from '../../components/CreateZeeBee/ThirdPartyService';
import Declaration from '../../components/CreateZeeBee/Declaration';
import Demographics from '../../components/CreateZeeBee/Demographics';
import ReviewAndSign from '../../components/CreateZeeBee/ReviewAndSign';
import Quoting from './Quoting';
import Processing from './Processing';
import Closing from './Closing';
import _ from 'lodash';
import { flattenObjectPaths, findItem, flattenArray } from '@/lib/common';
import { uploadAwsS3 } from '@/api/common';
import { LeaveAndSaveDraftBtn } from '@/components/NewModal/BtnPopUpModal';
import eventBus from '@/lib/bus';
import { assemblyAddress } from '@/lib/format';
// import CommonZeeBeeForm from '../../../zeeBeeDetails/components/CommonZeeBeeForm';
import General from '../../components/CreateZeeBee/General';
import { get } from 'http';
import HeaderBtn from './HeaderBtn';
import { EditOutlined } from '@ant-design/icons';
import SaveSvg from '@/assets/images/save.svg';
const initialValues = {
	applicants: [{ type: 'Individual', is_primary: true }],
};
const CommonSave = ({ children, title, setDisabled, handleUpdate, confirmLoading }) => {
	const [isSave, setIsSave] = useState(false);
	const handleClick = async () => {
		if (isSave) {
			await handleUpdate();
		}
		setDisabled(isSave);
		setIsSave(!isSave);
	};
	const changeHeaderBtnDom = () => {
		return (
			<Button
				type={isSave ? 'primary' : 'default'}
				onClick={handleClick}
				icon={isSave ? <SaveSvg /> : <EditOutlined />}
				loading={confirmLoading}
				disabled={false}
				size="middle"
			>
				{isSave ? 'Save' : 'Edit'}
			</Button>
		);
	};
	return (
		<>
			<HeaderBtn getContent={changeHeaderBtnDom}>{title}</HeaderBtn>
			<div className="px-6 py-6">{children}</div>
		</>
	);
};
/**
 * hooks组件，用于 新建/详情/草稿 三种ZeeBee Loan
 * @param {*} pageType 区分详情页 create/ details
 */
const useLenderDetails = () => {
	const searchParams = useSearchParams();
	const api = `/property/${searchParams.get('propertyId')}`;
	const [form] = Form.useForm();
	const scrollRef = useRef(null);
	const [currentStep, setCurrentStep] = useState(0); //顶部-总步骤
	const [stepNum, setStepNum] = useState(0); //侧边步骤
	const [loading, setLoading] = useState(false);
	const [errorList, setErrorList] = useState([]);
	const [isChange, setIsChange] = useState(false); // 记录是否操作过form
	const [confirmLoading, setConfirmLoading] = useState(false); // 保存
	const [rowData, setRowData] = useState({}); // loan详情
	const [newRowData, setNewRowData] = useState({}); // loan详情
	const [disabled, setDisabled] = useState(true);
	const { request } = createRequest();

	// 从property详情跳转的创建loan
	const getProperty = async () => {
		const data = await request(api);
		const newData = {
			loan_property: {
				// property_type: data.property_type,
				property_base: {
					unit: data.unit,
					address: data.address,
				},
			},
		};
		form.setFieldsValue(newData);
	};

	// 获取详情（草稿使用）
	const getLoanDetail = async () => {
		const res = (await request(`/loan/${searchParams.get('id')}`)) || {};
		const str = assemblyAddress(res?.loan_property?.property_base?.address);
		eventBus.emit('changeCrumbs', {
			2: `ZeeBee Loan Details(${str})`,
		});
		console.log(res);

		res.loan_borrower = {
			...(res?.loan_borrower || {}),
			tin: {
				tin_type: res?.loan_borrower?.tin_type,
				tin: res?.loan_borrower?.tin,
			},
		};

		if (!res?.applicants?.length) {
			res.applicants = initialValues.applicants;
		}
		setRowData(res);
		form.setFieldsValue(res);
	};

	// 获取property/草稿详情
	useEffect(() => {
		if (searchParams.has('propertyId')) getProperty();
		if (searchParams.has('id')) getLoanDetail();
	}, []);

	// 创建loan
	const handleSubmit = async (values) => {
		let url = '/loan/apply';
		let method = 'POST';
		if (searchParams.has('id')) {
			url = `/loan/${searchParams.get('id')}/update_apply`;
			method = 'PUT';
		}
		await request(url, {
			method,
			body: values,
		});
		message.open({
			type: 'success',
			content: 'ZeeBee loan created successfully！',
		});
	};

	const handlePdf = async (values) => {
		// 获取签名已生成的pdf数据流
		const res = await request('/pdf/apply_form', {
			method: 'POST',
			body: values,
		});
		const file = await res.blob();
		file.name = 'contract.pdf';
		const onProgress = (e) => {
			console.log('onProgress', e);
		};
		return await uploadAwsS3({ url: '/storage/presign_url', file, onProgress });
	};

	/**
	 * 以下为Lender系统内zeebee详情所需
	 */
	const handleContinue = async () => {
		const currentInfo = findItem({
			list: lenderDetailsSteps,
			key: currentStep,
			isPushParent: true,
			isPushIndex: true,
		});
		let nextInfo = null;
		if (currentInfo?.parentKey) {
			const parent = lenderDetailsSteps[currentInfo.parentIndex];
			if (parent?.children?.length === currentInfo?.index + 1) {
				nextInfo = lenderDetailsSteps[currentInfo.parentIndex + 1];
			} else {
				nextInfo = parent.children[currentInfo.index + 1];
			}
		} else {
			nextInfo = lenderDetailsSteps[currentInfo.index + 1];
		}
		await form.validateFields();
		const params = _.cloneDeep(rowData);
		params.loan_borrower.tin_type = rowData?.loan_borrower?.tin?.tin_type;
		params.loan_borrower.tin = rowData?.loan_borrower?.tin?.tin;
		await handleSubmit({ ...params });
		setCurrentStep(nextInfo.key);
		scrollRef.current.scrollTo({
			top: 0,
		});
		console.log(currentInfo, currentStep, rowData);
	};

	const handleSave = async () => {
		await form.validateFields();
		try {
			setConfirmLoading(true);
			const params = _.cloneDeep(newRowData);
			setRowData(newRowData);
			params.loan_borrower.tin_type = newRowData?.loan_borrower?.tin?.tin_type;
			params.loan_borrower.tin = newRowData?.loan_borrower?.tin?.tin;
			await handleSubmit({ ...params });
			setDisabled(true);
			setConfirmLoading(false);
		} catch (error) {
			console.log(error);

			setConfirmLoading(false);
		}
	};

	const lenderStepBtn = (
		<div className="w-[198px] flex flex-col gap-4">
			<Button type="primary" onClick={handleContinue}>
				Continue
			</Button>
		</div>
	);

	const getContent = (pageKey) => {
		const item = findItem({ list: lenderDetailsSteps, key: pageKey });
		switch (pageKey) {
			case 'general':
				return <General key="0-general" rowData={rowData} className="py-6 px-6" />;
			case 'borrowerInformation':
				return (
					<CommonSave
						setDisabled={setDisabled}
						title={item?.title}
						handleUpdate={handleSave}
						confirmLoading={confirmLoading}
					>
						<ApplicantInfo form={form} isHideTitle={true} />
						<BorrowerInfo form={form} isHideTitle={true} />
					</CommonSave>
				);
			case 'propertyAndLoan':
				return (
					<CommonSave
						setDisabled={setDisabled}
						title={item?.title}
						handleUpdate={handleSave}
						confirmLoading={confirmLoading}
					>
						<PropertyInfo form={form} />
						<Loan form={form} />
					</CommonSave>
				);
			case 'experienceAndThirdService':
				return (
					<CommonSave
						setDisabled={setDisabled}
						title={item?.title}
						handleUpdate={handleSave}
						confirmLoading={confirmLoading}
					>
						<Experience form={form} rowData={rowData} setIsChange={setIsChange} />
						<ThirdPartyService
							form={form}
							rowData={rowData}
							setIsChange={setIsChange}
						/>
					</CommonSave>
				);
			case 'declarationAndDemographics':
				return (
					<CommonSave
						setDisabled={setDisabled}
						title={item?.title}
						handleUpdate={handleSave}
						confirmLoading={confirmLoading}
					>
						<Declaration form={form} />
						<Demographics form={form} />
					</CommonSave>
				);
			case 'reviewAndSign':
				return (
					<div className="px-6 py-6">
						<ReviewAndSign form={form} />
					</div>
				);
			case 'Quoting':
				return <Quoting />;
			case 'Processing':
				return <Processing />;
			case 'Underwriting':
				return <Processing isHide={true} isHaveRejectBtn={true} title="Underwriting" />;
			case 'Closing':
				return <Closing />;
			default:
				return '';
		}
	};

	const lenderDetailsComponents = {
		General: getContent('general'),
		borrowerInformation: getContent('borrowerInformation'),
		propertyAndLoan: getContent('propertyAndLoan'),
		experienceAndThirdService: getContent('experienceAndThirdService'),
		declarationAndDemographics: getContent('declarationAndDemographics'),
		reviewAndSign: getContent('reviewAndSign'),
		Quoting: getContent('Quoting'),
		Processing: getContent('Processing'),
		Underwriting: getContent('Underwriting'),
		Closing: getContent('Closing'),
		Active: '',
		Funded: '',
	};

	const tabList = flattenArray(lenderDetailsSteps)
		.map((i) => ({
			key: i.key,
			label: i.title,
			children: lenderDetailsComponents[i.key],
		}))
		.filter((i) => i.children);
	const selectIndex = useMemo(
		() => tabList.findIndex((i) => i.key === currentStep),
		[currentStep],
	);

	useEffect(() => {
		setDisabled(selectIndex < 5);
	}, [selectIndex]);
	return {
		scrollRef,
		currentStep,
		setCurrentStep,
		stepNum,
		setStepNum,
		setIsChange,
		createZeeBeeSteps,
		form,
		initialValues,
		loading,
		rowData,
		zeeBeeDetailsTabs,
		lenderDetailsSteps,
		lenderDetailsComponents,
		lenderStepBtn,
		setRowData,
		disabled,
		newRowData,
		setNewRowData,
		tabList,
		selectIndex,
	};
};

export default useLenderDetails;
