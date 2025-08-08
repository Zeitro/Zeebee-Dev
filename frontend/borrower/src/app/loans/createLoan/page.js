'use client';
import { Button, Col, Form, message, Space, Steps } from 'antd';
import React, { useEffect, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
import { useRouter } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import { createRequest } from '@/lib/fetch';
import { formatDate } from '@/lib/date';
const CreateLoanPage = () => {
	const router = useRouter();
	const { request } = createRequest();
	const [currentStep, setCurrentStep] = useState(0);
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

	useEffect(() => {
		// back时赋值
		const item = loanData?.[currentStep];
		if (item && JSON.stringify(item) !== '{}' && form) {
			form.setFieldsValue({ ...item });
		}
	}, [currentStep, loanData]);

	// next/Add loan
	const handleNext = async () => {
		try {
			const values = await form.validateFields();
			if (currentStep < 2) {
				setLoanData((v) => ({ ...v, [currentStep]: values }));
				setCurrentStep((v) => v + 1);
				return;
			}
			// add loan
			const body = {
				...(loanData[0] || {}),
				...(loanData[1] || {}),
				...(values || {}),
			};
			if (body.closing_date) {
				body.closing_date = formatDate(body.closing_date);
			}
			if (body.maturity_date) {
				body.maturity_date = formatDate(body.maturity_date);
			}
			if (body.paid_off_date) {
				body.paid_off_date = formatDate(body.paid_off_date);
			}
			body.applicants = body?.applicants.filter((i) => i.name);
			setLoading(true);
			await request('/loan', {
				method: 'POST',
				body,
			});
			message.open({
				type: 'success',
				content: 'Success in creation！',
			});
			setLoading(false);
			router.back();
		} catch (error) {
			setLoading(false);
		}
	};

	const stepBtn = (
		<Col md={24} className="py-6 !flex justify-between items-center">
			<Button type="text" onClick={() => router.back()}>
				Cancel
			</Button>
			<Space size={8}>
				{!!currentStep && (
					<Button onClick={() => setCurrentStep((v) => v - 1)}>Back</Button>
				)}
				<Button type="primary" onClick={handleNext} loading={loading}>
					{currentStep < 2 ? 'Next' : 'Add Loan'}
				</Button>
			</Space>
		</Col>
	);
	const stepContent = {
		0: <CreateLoan form={form} stepBtn={stepBtn} />,
		1: (
			<CreateProperty
				stage={loanData?.[0]?.stage}
				stepBtn={stepBtn}
				form={form}
				gutter={16}
			/>
		),
		2: <CreateTerm stage={loanData?.[0]?.stage} stepBtn={stepBtn} />,
	};

	return (
		<div className="flex flex-col h-full">
			<div className="px-4 py-6 shrink-0 flex items-center justify-center border-b border-[#f0f0f0]">
				<Steps
					size="default"
					current={currentStep}
					items={createLoanSteps}
					className="max-w-[880px]"
				/>
			</div>
			<div className="flex-1 overflow-hidden overflow-y-auto flex justify-center pt-[88px]">
				<Form
					form={form}
					name="basic"
					layout="vertical"
					style={{ maxWidth: 576 }}
					autoComplete="off"
					className=""
					size="large"
					disabled={loading}
				>
					{stepContent[currentStep]}
				</Form>
			</div>
		</div>
	);
};
export default CreateLoanPage;
