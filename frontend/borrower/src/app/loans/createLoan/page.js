'use client';
import { Button, Col, Form, message, Space, Steps } from 'antd';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
import { useRouter, useSearchParams } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import { createRequest, changeUnitByPropertyType } from '@repo/shared-methods/utils';
import { loan_pending } from '@/enum/common';
import { BtnPopUpModal } from '@repo/shared-methods';
const CreateLoanPage = () => {
	const router = useRouter();
	const searchParams = useSearchParams();
	const api = `/property/${searchParams.get('propertyId')}`;
	const { request } = createRequest();
	const [currentStep, setCurrentStep] = useState(0);
	const [loading, setLoading] = useState(false);
	const [form] = Form.useForm();
	// 存储每个步骤的数据
	const [loanData, setLoanData] = useState({
		0: {},
		1: {},
		2: {},
	});
	const [unitInfo, setUnitInfo] = useState({});

	useEffect(() => {
		// back时赋值
		const item = loanData?.[currentStep];
		if (item && JSON.stringify(item) !== '{}' && form) {
			form.setFieldsValue({ ...item });
		}
	}, [currentStep, form, loanData]);

	// 获取property
	useEffect(() => {
		if (searchParams.has('propertyId')) getProperty();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	// next/Add loan
	const handleNext = useCallback(async () => {
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

			body.applicants = body?.applicants.filter(
				(i) => i.first_name || i.mid_name || i.last_name,
			);
			setLoading(true);
			await request('/loan', {
				method: 'POST',
				body,
			});
			message.open({
				type: 'success',
				content: '3rd party loan created successfully！',
			});
			setLoading(false);
			router.back();
		} catch {
			setLoading(false);
		}
	}, [currentStep, form, loanData, request, router]);

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

	const stepBtn = useMemo(
		() => (
			<Col xs={24} sm={24} className="py-6 !flex justify-between items-center">
				<BtnPopUpModal
					config={{
						onOk: () => router.back(),
						title: 'Tips',
						content:
							'The data is not saved yet, are you sure you want to leave the current page?',
						okText: 'Leave',
					}}
					modalType="confirm"
					type="text"
					size="middle"
				>
					Cancel
				</BtnPopUpModal>
				<Space size={8}>
					{!!currentStep && (
						<Button onClick={() => setCurrentStep((v) => v - 1)}>Back</Button>
					)}
					<Button type="primary" onClick={handleNext} loading={loading}>
						{currentStep < 2 ? 'Next' : 'Add Loan'}
					</Button>
				</Space>
			</Col>
		),
		[currentStep, handleNext, loading, router],
	);
	const stepContent = useMemo(
		() => ({
			0: <CreateLoan form={form} stepBtn={stepBtn} isRadioStyle={true} />,
			1: (
				<CreateProperty
					stage={loanData?.[0]?.stage}
					stepBtn={stepBtn}
					form={form}
					gutter={16}
					isRequestProperty={searchParams.has('propertyId')}
					unitInfo={unitInfo}
				/>
			),
			2: (
				<CreateTerm
					stage={loanData?.[0]?.stage}
					stepBtn={stepBtn}
					form={form}
					gutter={16}
				/>
			),
		}),
		[form, loanData, searchParams, stepBtn, unitInfo],
	);

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
					initialValues={{
						stage: loan_pending,
						applicants: [
							{ is_primary: true, type: 'Individual' },
							{ is_primary: false, type: 'Individual' },
						],
					}}
					onValuesChange={(v) => {
						// 根据property_type设置unit
						const property_type = v?.loan_property?.property_type;
						if (property_type) {
							const info = changeUnitByPropertyType(property_type);

							setUnitInfo({ ...info });
							form.setFieldValue(
								['loan_property', 'property_base', 'unit'],
								info.unit,
							);
						}
					}}
				>
					{stepContent[currentStep]}
				</Form>
			</div>
		</div>
	);
};
export default CreateLoanPage;
