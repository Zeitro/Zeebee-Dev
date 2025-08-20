'use client';
import { App, Button, Col, Form, message, Space, Steps } from 'antd';
import React, { useEffect, useState } from 'react';
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

const CreateZeeBeeLoan = () => {
	const router = useRouter();
	const searchParams = useSearchParams();
	const api = `/property/${searchParams.get('propertyId')}`;
	const { request } = createRequest();
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
		[<ApplicantInfo key="0" form={form} />, <BorrowerInfo key="1" form={form} />],
		// 1: (
		// 	<CreateProperty
		// 		stage={loanData?.[0]?.stage}
		// 		stepBtn={stepBtn}
		// 		form={form}
		// 		gutter={16}
		// 		isRequestProperty={searchParams.has('propertyId')}
		// 	/>
		// ),
		// 2: <CreateTerm stage={loanData?.[0]?.stage} stepBtn={stepBtn} form={form} />,
	];

	// // next/Add loan
	const handleNext = async () => {
		try {
			// const values = await form.validateFields();
			// console.log(values);
			if (
				currentStep + 1 === createZeeBeeSteps.length &&
				stepNum + 1 === stepContent[stepContent.length - 1]?.length
			) {
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

			const parentStepData = createZeeBeeSteps[currentStep];
			const childArr = parentStepData?.children;
			if (stepNum + 1 === childArr?.length) {
				setCurrentStep((v) => v + 1);
				return;
			}
			setStepNum((v) => v + 1);
		} catch (error) {
			setLoading(false);
		}
	};

	// 下一步按钮文案
	const getNextBtnText = () => {
		if (
			currentStep + 1 === stepContent.length &&
			stepNum + 1 === stepContent[stepContent.length - 1]?.length
		) {
			return 'Submit';
		}
		const parentStepData = createZeeBeeSteps[currentStep];
		const childArr = parentStepData?.children;
		if (stepNum + 1 === childArr?.length) {
			return `Next ${createZeeBeeSteps[currentStep + 1]?.title}`;
		}
		return `Next ${childArr[stepNum + 1]?.title}`;
	};

	const handleBack = () => {
		const stepItemArr = stepContent[currentStep];
		if (stepItemArr?.length >= stepNum + 1 && stepItemArr?.length > 1) {
			setStepNum((v) => v - 1);
		} else {
			setCurrentStep((v) => v - 1);
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

	const stepBtn = (
		<Col xs={24} sm={24} className="py-6 !flex justify-between items-center">
			<Button type="text" onClick={() => router.back()}>
				Cancel
			</Button>
			<Space size={8}>
				{!!currentStep && <Button onClick={handleBack}>Back</Button>}
				<Button type="primary" onClick={handleNext} loading={loading}>
					{getNextBtnText()}
				</Button>
			</Space>
		</Col>
	);

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
				<div className="w-[250px] py-6 px-6 shrink-0 border-r border-base-border">
					<Steps
						size="small"
						progressDot
						direction="vertical"
						current={stepNum}
						items={createZeeBeeSteps[currentStep]?.children}
					/>
				</div>
				<div className="flex-1 overflow-hidden overflow-y-auto py-8 px-8">
					<Form
						form={form}
						name="basic"
						layout="vertical"
						style={{ maxWidth: 520 }}
						autoComplete="off"
						className=""
						size="large"
						// disabled={loading}
						// initialValues={{
						// 	stage: loan_pending,
						// }}
					>
						{stepContent[currentStep]?.[stepNum]}
						{stepBtn}
					</Form>
				</div>
			</div>
		</div>
	);
};
export default CreateZeeBeeLoan;
