'use client';
import { Button, Col, Form, message, Row, Space, Steps, Tabs } from 'antd';
import React, { useEffect, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
// import styles from './loans.module.scss';
import { useRouter } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import { createRequest } from '@/lib/fetch';
import dayjs from 'dayjs';
const CreateLoanPage = () => {
	const router = useRouter();
	const { request } = createRequest();
	const [currentStep, setCurrentStep] = useState(0);
	const [form] = Form.useForm();
	// 存储每个步骤的数据
	const [loanData, setLoanData] = useState({
		0: {},
		1: {},
		2: {},
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
		const values = await form.validateFields();
		if (currentStep < 2) {
			setLoanData((v) => ({ ...v, [currentStep]: values }));
			setCurrentStep((v) => v + 1);
			return;
		}
		// add loan
		console.log('innadd');
		const body = {
			...(loanData[0] || {}),
			...(loanData[1] || {}),
			...(values || {}),
		};
		if (body.closing_date) {
			body.closing_date = dayjs(body.closing_date).format('MM/DD/YYYY');
		}
		if (body.maturity_date) {
			body.maturity_date = dayjs(body.maturity_date).format('MM/DD/YYYY');
		}
		if (body.paid_off_date) {
			body.paid_off_date = dayjs(body.paid_off_date).format('MM/DD/YYYY');
		}
		const res = await request('/loan', {
			method: 'POST',
			body,
		});
		message.success('Success in creation！');
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
				<Button type="primary" onClick={handleNext}>
					{currentStep < 2 ? 'Next' : 'Add Loan'}
				</Button>
			</Space>
		</Col>
	);
	const stepContent = {
		0: <CreateLoan form={form} stepBtn={stepBtn} />,
		1: <CreateProperty stage={loanData?.[0]?.stage} stepBtn={stepBtn} />,
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
				>
					{stepContent[currentStep]}
				</Form>
			</div>
		</div>
	);
};
export default CreateLoanPage;
