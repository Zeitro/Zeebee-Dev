'use client';
import { Button, Form, Steps } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';
import styles from './index.module.scss';
import BorrowerInfo from './components/BorrowerInfo';
import Property from './components/Property';
import Loan from './components/Loan';
import { useRouter } from 'next/navigation';
import { useLoanStore } from '@/store/loan';
import _ from 'lodash';
import { createRequest } from '@/lib/fetch';
// 创建或编辑loan时的步骤

const EstimateLoan = () => {
	const request = createRequest();
	const router = useRouter();
	const [currentStep, setCurrentStep] = useState(0);
	const [oldCurrentStep, setOldCurrentStep] = useState(null);
	const [form] = Form.useForm();
	const { estimateLoan, clearEstimateLoan, changeEstimateLoan } = useLoanStore();
	const loan_purpose = Form.useWatch('loan_purpose', form);
	const loanSteps = useMemo(
		() => [
			{
				title: 'Borrower information',
				key: 'Borrower',
				description: <div className="h-6"></div>,
				children: <BorrowerInfo />,
			},
			{
				title: 'Property information',
				key: 'Property',
				description: <div className="h-6"></div>,
				children: <Property loan_type={estimateLoan.loan_type} />,
			},
			{
				title: 'Loan information',
				key: 'loan',
				description: <div className="h-6"></div>,
				children: (
					<Loan
						loan_purpose={estimateLoan.loan_purpose}
						loan_type={estimateLoan.loan_type}
					/>
				),
			},
		],
		[loan_purpose, estimateLoan.loan_type],
	);
	const currentStepData = useMemo(() => loanSteps[currentStep], [currentStep]);

	const nextBtnInfo = useMemo(() => {
		const nextData = loanSteps?.[currentStep + 1];
		const info = {
			type: nextData ? 'Next' : 'Submit',
			text: nextData ? `Next` : 'Submit loan request',
		};
		return info;
	}, [currentStep]);

	const handleNext = async () => {
		await form.validateFields();
		if (nextBtnInfo.type === 'Next') {
			setOldCurrentStep(currentStep);
			return setCurrentStep(currentStep + 1);
		}
		// const res = await request('', {
		// 	method: 'POST',
		// 	body: estimateLoan,
		// });
		// console.log(res);
		router.replace('/estimateLoan/preApproved');
	};

	const btn = (
		<div className="w-full flex justify-between pt-6 pb-10">
			<div>
				{currentStep > 0 ? (
					<Button
						onClick={() => {
							setOldCurrentStep(currentStep);
							setCurrentStep(currentStep - 1);
						}}
						type="link"
					>
						Back
					</Button>
				) : (
					<Button onClick={() => router.back()} type="link">
						Cancel
					</Button>
				)}
				<Button
					onClick={() => {
						form.resetFields();
						const mergedObj = _.mergeWith(
							{},
							estimateLoan,
							form.getFieldsValue(),
							(objVal, srcVal) => {
								return srcVal === undefined ? null : undefined; // 将undefined转为null
							},
						);
						changeEstimateLoan(mergedObj);
					}}
				>
					Reset
				</Button>
			</div>
			<Button type="primary" onClick={handleNext}>
				{nextBtnInfo?.text}
			</Button>
		</div>
	);

	useEffect(() => {
		if (currentStep !== oldCurrentStep) {
			form?.setFieldsValue?.(estimateLoan);
		}
	}, [currentStep, estimateLoan, oldCurrentStep]);
	return (
		<div className={['flex h-full items-stretch', styles.estimateLoan].join(' ')}>
			<div className="flex-1 pl-[5%] overflow-hidden">
				<div className="text-xl font-medium pt-10 pb-6">Check your rate</div>
				<Steps
					direction="vertical"
					size="default"
					current={currentStep}
					items={loanSteps.map((step) => {
						// delete step.children;
						return step;
					})}
					// className="max-w-[880px]"
				/>
			</div>
			<div className="flex-2 w-1/2 border-l border-r border-[#91CAFF] bg-white py-6 flex justify-center overflow-hidden overflow-y-auto">
				<Form
					form={form}
					layout="vertical"
					autoComplete="off"
					className="w-[80%] m-auto"
					size="large"
					// disabled={!['documents', 'general'].includes(currentStep)}
					// initialValues={estimateLoan}
					scrollToFirstError
					// onValuesChange={(v) => console.log(v)}
					onFieldsChange={(changeFields) => {
						console.log(changeFields);
						// 只有输入验证通过了才可更新值
						const isErr = changeFields.every(
							(v) => v?.errors?.length || v?.warnings?.length || !v?.validated,
						);
						if (!isErr) {
							changeFields.map((v) => {
								if (!v?.errors?.length && !v?.warnings?.length) {
									_.set(estimateLoan, v?.name, v?.value);
								}
							});
							changeEstimateLoan(estimateLoan);
						}
					}}
				>
					{currentStepData.children}
					{btn}
				</Form>
			</div>
			<div className="flex-1 hidden sm:hidden md:hidden lg:block overflow-hidden"></div>
		</div>
	);
};
export default EstimateLoan;
