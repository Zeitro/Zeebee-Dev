'use client';
import { Button, Col, Form, message, Space } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { createZeeBeeSteps, zeeBeeDetailsTabs } from '@/enum/loan';
import {
	createRequest,
	flattenObjectPaths,
	eventBus,
	assemblyAddress,
} from '@repo/shared-methods/utils';
import ApplicantInfo from '../ApplicantInfo';
import BorrowerInfo from '../BorrowerInfo';
import PropertyInfo from '../PropertyInfo';
import LoanDetail from '../Loan';
import Experience from '../Experience';
import ThirdPartyService from '../ThirdPartyService';
import Declaration from '../Declaration';
import Demographics from '../Demographics';
import General from '../../../zeeBeeDetails/components/General';
import Documents from '../../../zeeBeeDetails/components/Documents';
import ReviewAndSign, { contractSignatureZones } from '../ReviewAndSign';
import _ from 'lodash';
import { uploadAwsS3 } from '@/api/common';
import { LeaveAndSaveDraftBtn } from '@repo/shared-methods';
const initialValues = {
	applicants: [{ type: 'Individual', is_primary: true }],
};
/**
 * hooks组件，用于 新建/详情/草稿 三种ZeeBee Loan
 * @param {*} pageType 区分详情页 create/ details
 */
const useCreate = (pageType = 'create') => {
	const isDetails = pageType === 'details';
	const router = useRouter();
	const searchParams = useSearchParams();
	const api = `/property/${searchParams.get('propertyId')}`;
	const [form] = Form.useForm();
	const scrollRef = useRef(null);
	const [currentStep, setCurrentStep] = useState(0); //顶部-总步骤
	const [stepNum, setStepNum] = useState(0); //侧边步骤
	const [loading, setLoading] = useState(false);
	const [isChange, setIsChange] = useState(false); // 记录是否操作过form
	const [confirmLoading, setConfirmLoading] = useState(false); // 保存草稿
	const [rowData, setRowData] = useState({}); // loan详情
	const [unitInfo, setUnitInfo] = useState({
		// unit: {},
		// future_unit: {},
	});
	const [loanUuid, setLoanUuid] = useState(null); // loan uuid,提前创建id用于pdf中展示
	const isClickBackBtn = useRef(false);
	const { request } = createRequest();

	// 从property详情跳转的创建loan
	const getProperty = async () => {
		const data = await request(api);
		const newData = {
			loan_property: {
				unit: data.unit,
				address: data.address,
			},
		};
		form.setFieldsValue(newData);
	};

	// 获取详情（草稿使用）
	const getLoanDetail = async () => {
		const res = await request(`/loan/${searchParams.get('id')}`);
		console.log(res);
		const str = assemblyAddress(res?.loan_property?.address);
		eventBus.emit('changeCrumbs', {
			2: isDetails ? `ZeeBee Loan Details(${str})` : `Edit Draft(${str})`,
		});
		res.loan_borrower.tin = {
			tin_type: res?.loan_borrower?.tin_type,
			tin: res?.loan_borrower?.tin,
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
		eventBus.on('getLoanDetail', getLoanDetail);
		return () => eventBus.off('getLoanDetail', getLoanDetail);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const stepContent = useMemo(
		() => [
			[
				<ApplicantInfo key="0-ApplicantInfo" form={form} />,
				<BorrowerInfo key="1-BorrowerInfo" form={form} />,
			],
			[
				<PropertyInfo
					key="0-PropertyInfo"
					form={form}
					unitInfo={unitInfo}
					isRequestProperty={searchParams.has('propertyId') || isDetails}
				/>,
				<LoanDetail key="1-LoanDetail" form={form} />,
			],
			[
				<Experience
					key="0-Experience"
					form={form}
					rowData={rowData}
					setIsChange={setIsChange}
				/>,
				<ThirdPartyService
					key="1-ThirdPartyService"
					form={form}
					rowData={rowData}
					setIsChange={setIsChange}
				/>,
			],
			[
				<Declaration key="0-Declaration" form={form} />,
				<Demographics key="1-Demographics" form={form} />,
			],
			<ReviewAndSign
				key="0-ReviewAndSign"
				form={form}
				rowData={isDetails ? rowData : null}
				loanUuid={loanUuid || rowData?.code}
				isNotRequest={isDetails}
			/>,
		],
		[form, isDetails, rowData, searchParams, unitInfo, loanUuid],
	);

	// zeeBeeDetails内独有的内容
	const zeeBeeDetailsContent = [
		<General key="0-general" rowData={rowData} />,
		<Documents key="0-documents" rowData={rowData} />,
	];
	const getZeeBeeDetailsTabsItem = () => {
		return zeeBeeDetailsTabs.map((item, index) => ({
			key: `${item.value}-0`,
			label: item.label,
			children: zeeBeeDetailsContent[index],
		}));
	};

	// 将stepContent的多维数组，转换成一维，key对应的就是当前步骤（currentStep + stepNum）
	const getTabsItem = (content = stepContent, isShowLabel = false, getContent) => {
		let newArr = [];
		content.map((item, index) => {
			const isLast = content.length - 1 === index;
			if (Array.isArray(item)) {
				item.map((item2, index2) => {
					const isLastStep = isLast && index2 === item.length - 1;
					newArr.push({
						key: `${index}-${index2}`,
						label: isShowLabel ? item2.title : '',
						children: getContent ? getContent(item2) : item2,
						destroyOnHidden: isLastStep,
					});
				});
			} else {
				newArr.push({
					key: `${index}-0`,
					label: isShowLabel ? item.title : '',
					children: getContent ? getContent(item) : item,
					destroyOnHidden: isLast,
				});
			}
		});
		return newArr;
	};

	// 获取下一步步骤类型（Submit创建按钮/NextMaxStep下一大步骤/NextMinStep下一小步骤）
	const diffSteps = () => {
		// 最后一步
		const lastStep = stepContent[stepContent.length - 1];
		// 倒数第二步
		const secondToLastStep = stepContent[stepContent.length - 2];

		// 某大步骤下有多个小步骤（数组），如为数组需判断是否最后一小步骤
		const isLastStepNumByArr = lastStep?.length ? stepNum + 1 === lastStep?.length : true;
		if (currentStep + 1 === stepContent.length && isLastStepNumByArr) {
			return 'Submit';
		}
		if (stepContent.length - 2 === currentStep && stepNum + 1 === secondToLastStep?.length) {
			return 'NextMaxStep-CreateUUID';
		}
		const parentStepData = createZeeBeeSteps[currentStep];
		const childArr = parentStepData?.children;
		if (childArr ? stepNum + 1 === childArr?.length : !childArr) {
			return 'NextMaxStep';
		}

		return 'NextMinStep';
	};

	// 创建loan
	const handleSubmit = async (values) => {
		try {
			let url = '/loan/apply';
			let method = 'POST';
			if (searchParams.has('id')) {
				url = `/loan/${searchParams.get('id')}/update_apply`;
				method = 'PUT';
			} else {
				values.code = loanUuid;
			}
			await request(url, {
				method,
				body: values,
			});
			message.open({
				type: 'success',
				content: 'ZeeBee loan created successfully！',
			});
		} catch (error) {
			setLoading(false);
			return Promise.reject(error);
		}
	};

	const handlePdf = async (values) => {
		try {
			// 获取签名已生成的pdf数据流
			if (!searchParams.has('id')) {
				values.code = loanUuid;
			}
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
		} catch (error) {
			setLoading(false);
			return Promise.reject(error);
		}
	};

	// 添加id- (草稿时必须传详情接口返回的id)
	const addId = (values) => {
		const ArrayAddId = (arr, arr2) => {
			return arr.map((item, index) => {
				const i = arr2?.[index];
				if (i?.id) {
					item.id = i?.id;
				}
				return item;
			});
		};
		if (rowData?.loan_property?.id && values.loan_property) {
			values.loan_property.id = rowData?.loan_property?.id;
		}
		if (rowData?.applicants?.length && values.applicants?.length) {
			values.applicants = ArrayAddId(values.applicants, rowData?.applicants);
		}
		if (rowData?.loan_borrower?.id && values.loan_borrower) {
			values.loan_borrower.id = rowData?.loan_borrower?.id;
		}
		if (rowData?.loan_demographic?.id && values.loan_demographic) {
			values.loan_demographic.id = rowData?.loan_demographic?.id;
		}
		return values;
	};

	const handleCreateUUID = async () => {
		const { uuid } = await request('/uuid', {
			method: 'GET',
		});
		setLoanUuid(uuid);
	};

	const getFieldsKeyForStep = (currentIndex = currentStep, stepIndex = stepNum) => {
		const componentFieldMap = {
			'0-0': ['loan_borrower'],
			'0-1': ['applicants'],
			'1-0': [
				['loan_property', 'address', 'address1'],
				['loan_property', 'address', 'address2'],
				['loan_property', 'address', 'city'],
				['loan_property', 'address', 'state'],
				['loan_property', 'address', 'zip_code'],
				'loan_purpose',
				'loan_type',
				['loan_property', 'property_type'],
				['loan_property', 'living_area_sqft'],
				['loan_property', 'unit'],
				['loan_property', 'bedrooms'],
				['loan_property', 'bathrooms'],
				['loan_property', 'future_property_type'],
				['loan_property', 'future_living_area_sqft'],
				['loan_property', 'future_unit'],
				['loan_property', 'future_bedrooms'],
				['loan_property', 'future_bathrooms'],
			],
			'1-1': [],
			'2-0': ['transaction_experiences', 'real_estate_experiences'],
			'2-1': ['loan_third_services'],
			'3-0': [['loan_demographic', 'property']],
			'3-1': ['loan_demographic'],
			'4-0': ['signatures'],
		};
		return componentFieldMap[`${currentIndex}-${stepIndex}`];
	};

	// 获取除去其他步骤字段外的剩余字段
	const getRemainingFields = (allFieldNames) => {
		// 收集所有其他步骤的字段前缀
		const allPrefixes = new Set();

		stepContent.forEach((step, stepIndex) => {
			if (Array.isArray(step)) {
				step.forEach((_, subStepIndex) => {
					if (stepIndex === currentStep && subStepIndex === stepNum) {
						// 跳过当前步骤
						return;
					}
					const prefixes = getFieldsKeyForStep(stepIndex, subStepIndex);
					if (prefixes) {
						prefixes.forEach((prefix) => {
							allPrefixes.add(prefix);
						});
					}
				});
			} else {
				if (stepIndex === currentStep && stepNum === 0) {
					// 跳过当前步骤
					return;
				}
				const prefixes = getFieldsKeyForStep(stepIndex, 0);
				if (prefixes) {
					prefixes.forEach((prefix) => {
						allPrefixes.add(prefix);
					});
				}
			}
		});
		// 返回不属于其他步骤前缀的字段
		return allFieldNames.filter((fieldName) => {
			return ![...allPrefixes].some((prefix) => {
				if (Array.isArray(prefix)) {
					return prefix.every((p) => fieldName.includes(p));
				}
				// 如果prefix是字符串
				return fieldName.includes(prefix);
			});
		});
	};

	// 获取当前步骤需要验证的字段
	const getCurrentStepFields = async () => {
		// 获取当前步骤定义的字段前缀
		const fieldPrefixes = getFieldsKeyForStep();

		// 获取表单中所有字段
		let allValues,
			applicantsErr = [];
		try {
			allValues = await form.validateFields({ validateOnly: true });
		} catch (error) {
			allValues = error?.values;
			// applicants的字段单独处理，因为是list验证不到
			applicantsErr =
				error.errorFields?.map((i) => i.name).filter((i) => i.includes('applicants')) || [];
		}

		const allFieldNames = [...flattenObjectPaths(allValues), ...applicantsErr];

		// 如果字段前缀为空数组，表示验证除去其他步骤字段外的所有剩余字段
		if (fieldPrefixes && fieldPrefixes.length === 0) {
			return getRemainingFields(allFieldNames);
		}

		// 根据字段前缀筛选需要验证的字段
		if (fieldPrefixes && fieldPrefixes.length > 0) {
			return allFieldNames.filter((fieldName) => {
				return fieldPrefixes.some((prefix) => {
					// 如果prefix是数组，需要匹配所有前缀
					if (Array.isArray(prefix)) {
						return prefix.every((p) => fieldName.includes(p));
					}
					// 如果prefix是字符串
					return fieldName.includes(prefix);
				});
			});
		}

		// 默认验证所有字段
		return allFieldNames;
	};
	// next/Add loan
	const handleNext = async () => {
		try {
			let values;
			const isType = diffSteps();
			// 解决：没填数据就触发了错误验证，又直接切换到上一步骤，又点击下一步时，validateFields会报错，提示刚才错误验证的字段没有填，导致下一页切换不了
			if (isType === 'Submit') {
				values = await form.validateFields();
				if (values?.signatures?.remainingSignatureCount || !values?.signatures) {
					message.open({
						type: 'error',
						content: `There are still ${values?.signatures?.remainingSignatureCount || contractSignatureZones.length} locations without a signature!`,
					});
					return;
				}
				setLoading(true);
				const body = _.cloneDeep(addId(values));

				body.loan_borrower.tin_type = values?.loan_borrower?.tin?.tin_type;
				body.loan_borrower.tin = values?.loan_borrower?.tin?.tin;
				_.unset(body, ['signatures', 'remainingSignatureCount']);
				const fileUrl = await handlePdf(body);
				await handleSubmit({ ...body, apply_form_file: fileUrl });
				setLoading(false);
				router.back();
				return;
			}
			const fieldsKeys = await getCurrentStepFields();
			values = await form.validateFields(fieldsKeys, { recursive: true });
			if (isType === 'NextMaxStep-CreateUUID' && !rowData?.code && !loanUuid) {
				await handleCreateUUID();
			}

			if (['NextMaxStep', 'NextMaxStep-CreateUUID'].includes(isType)) {
				setCurrentStep((v) => v + 1);
				setStepNum(0);
			} else {
				setStepNum((v) => v + 1);
			}

			scrollRef.current.scrollTo({
				top: 0,
			});
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	// 保存/更新草稿
	const handleSaveDraft = async () => {
		const allValues = form.getFieldsValue();
		const values = { ...rowData, ...allValues };
		setConfirmLoading(true);
		values.loan_borrower.tin_type = values?.loan_borrower?.tin?.tin_type;
		values.loan_borrower.tin = values?.loan_borrower?.tin?.tin;
		// 草稿时删除签名信息，因为合同的签名及时间默认以最新时间为准，所以不需要保存
		_.unset(values, ['signatures']);
		const isDraft = searchParams.has('id');
		let url = isDraft ? `/loan/${searchParams.get('id')}/draft` : '/loan/draft';
		let method = isDraft ? 'PUT' : 'POST';
		let content = isDraft ? 'Draft updated successfully!' : 'Draft created successfully!';
		const body = isDraft ? addId(values) : values;
		await request(url, {
			method,
			body,
		});
		message.open({
			type: 'success',
			content,
		});
		setConfirmLoading(false);
		router.back();
	};

	// 下一步按钮文案
	const getNextBtnText = () => {
		const isType = diffSteps();
		const parentStepData = createZeeBeeSteps[currentStep];
		switch (isType) {
			case 'Submit':
				return 'Submit application';
			case 'NextMaxStep':
			case 'NextMaxStep-CreateUUID':
				return `Next ${createZeeBeeSteps[currentStep + 1]?.title}`;
			case 'NextMinStep':
				return `Next ${parentStepData?.children[stepNum + 1]?.title}`;
			default:
				break;
		}
	};

	const handleBack = () => {
		isClickBackBtn.current = true;
		const errors = form.getFieldsError()?.filter((i) => i?.errors?.length);
		if (errors?.length) {
			// 光设置重新验证不行，需要结合validateFields处的操作才可避免
			let errorList = errors.map((i) => i?.name);
			form.resetFields(errorList);
		}
		if (stepNum) {
			setStepNum((v) => (v - 1 < 0 ? 0 : v - 1));
		} else {
			const prevStepData = stepContent[currentStep - 1];
			setCurrentStep((v) => (v - 1 < 0 ? 0 : v - 1));
			setStepNum(prevStepData?.length - 1 < 0 ? 0 : prevStepData?.length - 1);
		}
	};

	const stepBtn = useMemo(() => {
		const cancelBtn = (
			<LeaveAndSaveDraftBtn
				handleSaveDraft={handleSaveDraft}
				isChange={isChange}
				router={router}
				disabled={loading}
				confirmLoading={confirmLoading}
			/>
		);
		const draftBtn = (
			<Button
				type="primary"
				ghost
				className="w-[90px]"
				onClick={handleSaveDraft}
				disabled={loading || !isChange}
			>
				Draft
			</Button>
		);
		const backBtn = !(!currentStep && !stepNum) && (
			<Button className="w-[90px]" onClick={handleBack} disabled={loading}>
				Back
			</Button>
		);
		const nextBtn = (
			<Button type="primary" onClick={handleNext} loading={loading}>
				{getNextBtnText()}
			</Button>
		);
		if (currentStep === stepContent.length - 1) {
			return (
				<div className="w-[198px] m-auto flex flex-col gap-4">
					{nextBtn}
					<Space size={24}>
						{cancelBtn}
						{backBtn}
					</Space>
				</div>
			);
		}
		return (
			<Col xs={24} sm={24} className="py-6 !flex justify-between items-center">
				<Space size={8}>
					{cancelBtn}
					{draftBtn}
				</Space>
				<Space size={8}>
					{backBtn}
					{nextBtn}
				</Space>
			</Col>
		);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [currentStep, stepNum, loading, isChange, confirmLoading, stepContent.length]);

	return {
		getTabsItem,
		scrollRef,
		stepContent,
		currentStep,
		setCurrentStep,
		stepNum,
		setStepNum,
		stepBtn,
		setIsChange,
		createZeeBeeSteps,
		zeeBeeDetailsTabs,
		getNextBtnText,
		form,
		initialValues,
		loading,
		getZeeBeeDetailsTabsItem,
		rowData,
		unitInfo,
		setUnitInfo,
	};
};

export default useCreate;
