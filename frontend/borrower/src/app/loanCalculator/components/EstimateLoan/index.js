'use client';
import { Button, Form, Space, Steps } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';
import styles from './index.module.scss';
import BorrowerInfo from './BorrowerInfo';
import Property from './Property';
import Loan from './Loan';
import { useRouter } from 'next/navigation';
import { useLoanStore } from '@/store/loan';
import _ from 'lodash';
import { createRequest } from '@/lib/fetch';
import { changeUnitByPropertyType } from '@/lib/format';

// pre-approval,提交简易loan计算
const EstimateLoan = ({ setIsResult, setResult }) => {
	const { request } = createRequest();
	const router = useRouter();
	const [currentStep, setCurrentStep] = useState(0);
	const [oldCurrentStep, setOldCurrentStep] = useState(null);
	const [loading, setLoading] = useState(false);
	const [form] = Form.useForm();
	const {
		estimateLoan,
		changeEstimateLoan,
		savePreApprovedResult,
		preApprovedResult,
		clearExpiredPreApprovedInStorage,
	} = useLoanStore();
	const [unitInfo, setUnitInfo] = useState({});
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
				children: (
					<Property
						oldPropertyType={estimateLoan?.property?.property_type}
						unitInfo={unitInfo}
					/>
				),
			},
			{
				title: 'Loan information',
				key: 'loan',
				description: <div className="h-6"></div>,
				children: (
					<Loan
						loan_purpose={estimateLoan?.loan?.loan_purpose}
						loan_type={estimateLoan?.loan?.loan_type}
					/>
				),
			},
		],
		[
			estimateLoan?.loan?.loan_purpose,
			estimateLoan?.loan?.loan_type,
			estimateLoan?.property?.property_type,
			unitInfo,
		],
	);
	const currentStepData = useMemo(() => loanSteps[currentStep], [currentStep, loanSteps]);

	const nextBtnInfo = useMemo(() => {
		const nextData = loanSteps?.[currentStep + 1];
		const info = {
			type: nextData ? 'Next' : 'Submit',
			text: nextData ? `Next` : 'Submit loan request',
		};
		return info;
	}, [currentStep, loanSteps]);

	const handleNext = async () => {
		try {
			await form.validateFields();
			// 清除过期的preApproved,只保留24小时
			clearExpiredPreApprovedInStorage();
			if (nextBtnInfo.type === 'Next') {
				setOldCurrentStep(currentStep);
				return setCurrentStep(currentStep + 1);
			}
			setLoading(true);
			const res = await request('/loan/pre_approval', {
				method: 'POST',
				body: { ...estimateLoan, request_id: preApprovedResult?.primary?.request_id || '' },
			});
			res.email = estimateLoan?.borrower?.email;
			setResult(res);
			savePreApprovedResult(res?.primary);
			setIsResult(true);
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	const btn = (
		<div className="w-full flex justify-between pt-6 pb-10">
			<Space>
				{currentStep > 0 ? (
					<Button
						onClick={() => {
							setOldCurrentStep(currentStep);
							setCurrentStep(currentStep - 1);
						}}
					>
						Back
					</Button>
				) : (
					<Button onClick={() => router.back()}>Cancel</Button>
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
					type="link"
				>
					Reset
				</Button>
			</Space>
			<Button type="primary" onClick={handleNext} loading={loading}>
				{nextBtnInfo?.text}
			</Button>
		</div>
	);

	useEffect(() => {
		if (currentStep !== oldCurrentStep) {
			form?.setFieldsValue?.(estimateLoan);
			setUnitInfo(changeUnitByPropertyType(estimateLoan?.property?.property_type));
		}
	}, [currentStep, estimateLoan, form, oldCurrentStep]);
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
					scrollToFirstError
					onFieldsChange={(changeFields) => {
						// 只有输入验证通过了才可更新值
						const isErr = changeFields.every(
							(v) => v?.errors?.length || v?.warnings?.length || !v?.validated,
						);
						if (!isErr) {
							changeFields.map((v) => {
								if (!v?.errors?.length && !v?.warnings?.length) {
									_.set(estimateLoan, v?.name, v?.value);
									// property_type修改时更新unit
									if (v?.name?.includes('property_type')) {
										const info = changeUnitByPropertyType(v?.value);
										setUnitInfo(info);
										form.setFieldValue(['property', 'unit'], info.unit);
										_.set(estimateLoan, ['property', 'unit'], info.unit);
									}
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
