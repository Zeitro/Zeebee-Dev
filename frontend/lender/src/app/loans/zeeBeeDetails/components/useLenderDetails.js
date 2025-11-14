'use client';
import { Button, Form, message } from 'antd';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { createZeeBeeSteps, zeeBeeDetailsTabs, lenderDetailsSteps } from '@/enum/loan';
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
import {
	findItem,
	flattenArray,
	assemblyAddress,
	eventBus,
	createRequest,
} from '@repo/shared-methods/utils';
import { uploadAwsS3 } from '@/api/common';
import General from '../../components/CreateZeeBee/General';
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
	const [disabled, setDisabled] = useState(true); //控制application内的页面不可编辑，只有点击Edit才可修改
	const [unitInfo, setUnitInfo] = useState({
		// unit: {},
		// future_unit: {},
	});
	const { request } = createRequest();
	const id = searchParams.get('id');
	const newList = flattenArray(lenderDetailsSteps);

	// 获取详情
	const getLoanDetail = useCallback(async () => {
		const res = (await request(`/loan/${searchParams.get('id')}`)) || {};
		const str = assemblyAddress(res?.loan_property?.address);
		eventBus.emit('changeCrumbs', {
			2: `Loan Details(${str})`,
		});
		// res.loan_borrower = {
		// 	...(res?.loan_borrower || {}),
		// 	// tin: {
		// 	// 	tin_type: res?.loan_borrower?.tin_type,
		// 	// 	tin: res?.loan_borrower?.tin,
		// 	// },
		// };

		if (!res?.applicants?.length) {
			res.applicants = initialValues.applicants;
		}
		// 保存两份，一份用于展示和修改（可能出现修改了就离开了，并未保存），一份用于最新数据保存
		setRowData(res);
		setNewRowData(res);
		form.setFieldsValue(res);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [form, searchParams]);

	// 详情
	useEffect(() => {
		if (searchParams.has('id')) getLoanDetail();
		eventBus.on('getLoanDetail', getLoanDetail);
		return () => eventBus.off('getLoanDetail', getLoanDetail);
	}, [getLoanDetail, searchParams]);

	// 全量修改loan
	const handleSubmit = async (values) => {
		await request(`/loan/${searchParams.get('id')}`, {
			method: 'PUT',
			body: values,
		});
		message.open({
			type: 'success',
			content: 'ZeeBee loan updated successfully！',
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
		// params.loan_borrower.tin_type = rowData?.loan_borrower?.tin?.tin_type;
		// params.loan_borrower.tin = rowData?.loan_borrower?.tin?.tin;
		await handleSubmit({ ...params });
		setCurrentStep(nextInfo.key);
		scrollRef.current.scrollTo({
			top: 0,
		});
	};

	const handleSave = async () => {
		await form.validateFields();
		try {
			setConfirmLoading(true);
			const params = _.cloneDeep(newRowData);
			setRowData(newRowData);
			// params.loan_borrower.tin_type = newRowData?.loan_borrower?.tin?.tin_type;
			// params.loan_borrower.tin = newRowData?.loan_borrower?.tin?.tin;
			await handleSubmit({ ...params });
			eventBus.emit('getLoanDetail');
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

	const getContent = (pageKey, disabledOperate) => {
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
						<ApplicantInfo
							form={form}
							isHideTitle={true}
							disabled={disabled}
							isShowSelectCard={false}
						/>
						<BorrowerInfo
							form={form}
							isHideTitle={true}
							disabled={disabled}
							isShowSelectCard={false}
						/>
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
						<PropertyInfo
							form={form}
							unitInfo={unitInfo}
							isRequestProperty={disabled}
							isShowSelectCard={false}
						/>
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
						<Experience
							form={form}
							rowData={rowData}
							setIsChange={setIsChange}
							changeFieldsValue={(value, key) => {
								newRowData[key] = value;
								setNewRowData({ ...newRowData });
							}}
						/>
						<ThirdPartyService
							form={form}
							rowData={rowData}
							setIsChange={setIsChange}
							changeFieldsValue={(value, key) => {
								let newValue = { ...newRowData };
								newValue[key] = value;
								setNewRowData({ ...newValue });
							}}
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
						<Declaration form={form} disabled={disabled} />
						<Demographics form={form} />
					</CommonSave>
				);
			case 'reviewAndSign':
				return (
					<div className="px-6 py-6">
						<ReviewAndSign form={form} isNotRequest={true} rowData={rowData} />
					</div>
				);
			case 'Quoting':
				return (
					<Quoting
						id={id}
						setNextStep={() => setCurrentStep('Processing')}
						disabledOperate={disabledOperate}
					/>
				);
			case 'Processing':
				return (
					<Processing
						id={id}
						setNextStep={() => setCurrentStep('Underwriting')}
						disabledOperate={disabledOperate}
						rowData={rowData}
					/>
				);
			case 'Underwriting':
				return (
					<Processing
						isHide={true}
						isHaveRejectBtn={true}
						title="Underwriting"
						id={id}
						setRejectStep={() => setCurrentStep('Processing')}
						setNextStep={() => setCurrentStep('Closing')}
						disabledOperate={disabledOperate}
						rowData={rowData}
					/>
				);
			case 'Closing':
				return (
					<Closing
						// setNextStep={() => setCurrentStep('Active')}
						disabledOperate={disabledOperate}
						rowData={rowData}
						id={id}
					/>
				);
			case 'Active':
				return 'Active';
			case 'Funded':
				return 'Funded';
			default:
				return '';
		}
	};

	const getDetailsComponents = (disabledOperate) => {
		// disabledOperate 可查看，不可操作页面
		return {
			General: getContent('general'),
			borrowerInformation: getContent('borrowerInformation'),
			propertyAndLoan: getContent('propertyAndLoan'),
			experienceAndThirdService: getContent('experienceAndThirdService'),
			declarationAndDemographics: getContent('declarationAndDemographics'),
			reviewAndSign: getContent('reviewAndSign'),
			Quoting: getContent('Quoting', disabledOperate),
			Processing: getContent('Processing', disabledOperate),
			Underwriting: getContent('Underwriting', disabledOperate),
			Closing: getContent('Closing', disabledOperate),
			Funded: getContent('Funded'),
		};
	};

	// 初始化时，根据不同状态判断选中第几步
	const loanProgress = useMemo(() => {
		const pendingStatus = rowData?.pending_status;
		const stage = rowData?.stage;

		let index = 0;
		let selectTabKey = '';
		if (stage === 'Draft') {
			index = 0;
		} else if (stage === 'Pending') {
			index = newList.findIndex((i) => i.key === pendingStatus);
		} else if (stage === 'Active') {
			index = newList.findIndex((i) => i.key === 'Funded');
		} else if (stage === 'Rejected') {
			index = newList.findIndex((i) => i.key === 'Quoting');
		} else {
			index = newList.length - 1;
		}
		selectTabKey = newList[index]?.key;
		return {
			selectTabKey,
			index,
		};
	}, [newList, rowData?.pending_status, rowData?.stage]);

	const tabList = newList
		.map((i, index) => ({
			key: i.key,
			label: i.title,
			children: getDetailsComponents(index !== loanProgress.index)[i.key],
		}))
		.filter((i) => i.children);
	const selectIndex = useMemo(
		() => tabList.findIndex((i) => i.key === currentStep),
		[currentStep, tabList],
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
		getDetailsComponents,
		lenderStepBtn,
		setRowData,
		disabled,
		newRowData,
		setNewRowData,
		tabList,
		selectIndex,
		loanProgress,
	};
};

export default useLenderDetails;
