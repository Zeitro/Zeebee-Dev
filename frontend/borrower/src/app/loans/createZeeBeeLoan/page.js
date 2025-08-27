'use client';
import { App, Button, Col, Form, message, Space, Steps } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
import { useRouter, useSearchParams } from 'next/navigation';
import { createZeeBeeSteps } from '@/enum/loan';
import { createRequest } from '@/lib/fetch';
import { formatDate } from '@/lib/date';
import { loan_pending } from '@/enum/common';
import ApplicantInfo from '../components/CreateZeeBee/ApplicantInfo';
import BorrowerInfo from '../components/CreateZeeBee/BorrowerInfo';
import PropertyInfo from '../components/CreateZeeBee/PropertyInfo';
import LoanDetail from '../components/CreateZeeBee/Loan';
import Experience from '../components/CreateZeeBee/Experience';
import ThirdPartyService from '../components/CreateZeeBee/ThirdPartyService';
import Declaration from '../components/CreateZeeBee/Declaration';
import Demographics from '../components/CreateZeeBee/Demographics';
import ReviewAndSign from '../components/CreateZeeBee/ReviewAndSign';
const CreateZeeBeeLoan = () => {
	const router = useRouter();
	const searchParams = useSearchParams();
	const api = `/property/${searchParams.get('propertyId')}`;
	const { request } = createRequest();
	const scrollRef = useRef(null);
	const [currentStep, setCurrentStep] = useState(0); //顶部-总步骤
	const [stepNum, setStepNum] = useState(0); //侧边步骤
	const [loading, setLoading] = useState(false);
	const [form] = Form.useForm();
	// 存储每个步骤的数据
	const [loanData, setLoanData] = useState({
		0: {},
		1: {},
		2: {
			applicants: [
				{ is_primary: true, type: 'Individual' },
				{ is_primary: false, type: 'Individual' },
			],
		},
	});

	// useEffect(() => {
	// 	// back时赋值
	// 	const item = loanData?.[currentStep];
	// 	if (item && JSON.stringify(item) !== '{}' && form) {
	// 		form.setFieldsValue({ ...item });
	// 	}
	// }, [currentStep, loanData]);

	// // 获取property
	// useEffect(() => {
	// 	if (searchParams.has('propertyId')) getProperty();
	// }, []);

	const stepContent = [
		[
			<ApplicantInfo key="0-ApplicantInfo" form={form} />,
			<BorrowerInfo key="1-BorrowerInfo" form={form} />,
		],
		[
			<PropertyInfo key="0-PropertyInfo" form={form} />,
			<LoanDetail key="1-LoanDetail" form={form} />,
		],
		[
			<Experience key="0-Experience" form={form} />,
			<ThirdPartyService key="1-ThirdPartyService" form={form} />,
		],
		[
			<Declaration key="0-Declaration" form={form} />,
			<Demographics key="1-Demographics" form={form} />,
		],
		<ReviewAndSign key="0-ReviewAndSign" form={form} />,
	];

	// 获取下一步步骤类型（Submit创建按钮/NextMaxStep下一大步骤/NextMinStep下一小步骤）
	const diffSteps = () => {
		const lastStep = stepContent[stepContent.length - 1];
		// 某大步骤下有多个小步骤（数组），如为数组需判断是否最后一小步骤
		const isLastStepNumByArr = lastStep?.length ? stepNum + 1 === lastStep?.length : true;
		if (currentStep + 1 === stepContent.length && isLastStepNumByArr) {
			return 'Submit';
		}
		const parentStepData = createZeeBeeSteps[currentStep];
		const childArr = parentStepData?.children;
		if (childArr ? stepNum + 1 === childArr?.length : !childArr) {
			return 'NextMaxStep';
		}

		return 'NextMinStep';
	};

	// // next/Add loan
	const handleNext = async () => {
		try {
			// const values = await form.validateFields();
			// console.log(values);
			const isType = diffSteps();
			if (isType === 'Submit') {
				// 		// add loan
				// 		const body = {
				// 			...(loanData[0] || {}),
				// 			...(loanData[1] || {}),
				// 			...(values || {}),
				// 		};
				// 		if (body.closing_date) {
				// 			body.closing_date = formatDate(body.closing_date);
				// 		}
				// 		if (body.maturity_date) {
				// 			body.maturity_date = formatDate(body.maturity_date);
				// 		}
				// 		if (body.paid_off_date) {
				// 			body.paid_off_date = formatDate(body.paid_off_date);
				// 		}
				// 		body.applicants = body?.applicants.filter((i) => i.name);
				// 		setLoading(true);
				// 		await request('/loan', {
				// 			method: 'POST',
				// 			body,
				// 		});
				// 		message.open({
				// 			type: 'success',
				// 			content: 'Success in creation！',
				// 		});
				// 		setLoading(false);
				// 		router.back();
				return;
			}
			// setLoanData((v) => ({ ...v, [currentStep]: values }));
			if (isType === 'NextMaxStep') {
				setCurrentStep((v) => v + 1);
				setStepNum(0);
			} else {
				setStepNum((v) => v + 1);
			}

			scrollRef.current.scrollTo({
				top: 0,
			});
		} catch (error) {
			setLoading(false);
		}
	};

	// 下一步按钮文案
	const getNextBtnText = () => {
		const isType = diffSteps();
		const parentStepData = createZeeBeeSteps[currentStep];
		switch (isType) {
			case 'Submit':
				return 'Submit application';
			case 'NextMaxStep':
				return `Next ${createZeeBeeSteps[currentStep + 1]?.title}`;
			case 'NextMinStep':
				return `Next ${parentStepData?.children[stepNum + 1]?.title}`;
			default:
				break;
		}
	};

	const handleBack = () => {
		if (stepNum) {
			setStepNum((v) => v - 1);
		} else {
			const prevStepData = stepContent[currentStep - 1];
			setCurrentStep((v) => v - 1);
			setStepNum(prevStepData?.length - 1);
		}
	};

	// // 从property详情跳转的创建loan
	// const getProperty = async () => {
	// 	const data = await request(api);
	// 	const newData = {
	// 		loan_property: {
	// 			property_type: data.property_type,
	// 			property_base: {
	// 				unit: data.unit,
	// 				address: data.address,
	// 			},
	// 		},
	// 	};
	// 	form.setFieldsValue(newData);
	// };

	const stepBtn = useMemo(() => {
		const cancelBtn = <Button onClick={() => router.back()}>Cancel</Button>;
		const backBtn = !(!currentStep && !stepNum) && (
			<Button className="w-[90px]" onClick={handleBack}>
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
				{cancelBtn}
				<Space size={8}>
					{backBtn}
					{nextBtn}
				</Space>
			</Col>
		);
	}, [currentStep, stepNum, loading]);

	return (
		<div className="flex flex-col h-full">
			<div className="px-12 py-4 shrink-0 flex items-center justify-center border-b border-base-border">
				<Steps
					size="default"
					current={currentStep}
					items={createZeeBeeSteps}
					// className="max-w-[880px]"
				/>
			</div>
			<div className="flex-1 flex items-stretch overflow-hidden">
				{createZeeBeeSteps[currentStep]?.children && (
					<div className="w-[250px] py-6 px-6 shrink-0 border-r border-base-border">
						<Steps
							size="small"
							progressDot
							direction="vertical"
							current={stepNum}
							items={createZeeBeeSteps[currentStep]?.children}
						/>
					</div>
				)}
				<div
					className={`flex-1 overflow-hidden overflow-y-auto py-${stepContent.length === currentStep ? 6 : 8} px-8`}
					ref={scrollRef}
				>
					<Form
						form={form}
						name="basic"
						layout="vertical"
						autoComplete="off"
						className=""
						size="large"
						// disabled={loading}
						initialValues={{
							borrowers: [{}],
						}}
					>
						{stepContent[currentStep]?.length
							? stepContent[currentStep]?.[stepNum]
							: stepContent[currentStep]}
						{stepBtn}
					</Form>
				</div>
			</div>
		</div>
	);
};
export default CreateZeeBeeLoan;
