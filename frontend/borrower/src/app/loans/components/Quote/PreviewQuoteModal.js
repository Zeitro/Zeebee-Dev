'use client';
import { useState } from 'react';
import { NewModal } from '@repo/shared-methods';
import { get } from 'lodash';
import { Button, Checkbox, Col, message, Row, Space, Spin } from 'antd';
import { HomeTwoTone } from '@ant-design/icons';
import { assemblyAddress, formatAmount, createRequest } from '@repo/shared-methods/utils';
import { useUserStore } from '@/store/user';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

// 询价（选择lender发送邮件询价，邮箱未授权只有admin有权限设置）
const PreviewQuote = ({ setOpen, rowData }) => {
	const { userInfo } = useUserStore();
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [options, setOptions] = useState([]);
	const [value, setValue] = useState([]);
	const [step, setStep] = useState(1);
	const handleSaveAndSendQuote = async () => {
		if (step === 1) {
			setStep(2);
			return;
		}
		try {
			// 保存3rd party loan并发送邮件给lender询价
			if (!value?.length) {
				return message.open({
					content: 'Please select a lender',
					type: 'warning',
				});
			}
			setConfirmLoading(true);
			await request('/loan/send_quote_email', {
				method: 'POST',
				body: {
					...rowData,
					lender_ids: value,
				},
			});
			setConfirmLoading(false);
			message.open({
				content: `Quote email sent successfully to ${value.length} lender${
					value.length > 1 ? 's' : ''
				}`,
				type: 'success',
			});
			setOpen(false);
		} catch (error) {
			// 邮箱未授权
			if (error.code === 401002) {
				if (userInfo?.is_admin || userInfo?.is_owner) {
					message.open({
						content: 'Email is not authorized for this organization. Page jumping...',
						type: 'warning',
					});
					setTimeout(() => {
						setConfirmLoading(false);
						window.open('/setting/system', '_blank');
					}, 500);
					return;
				}
				message.open({
					content:
						'Email is not authorized for this organization.Please contact your admin.',
					type: 'warning',
				});
			}
			setConfirmLoading(false);
			setOpen(false);
		}
	};
	const list = [
		{ title: 'Loan purpose', value: 'loan_purpose' },
		{ title: 'Loan type', value: 'loan_type' },
		{ title: 'Total loan amount', value: 'total_amount', render: (v) => formatAmount(v) },
		{ title: 'Business strategy', value: 'business_strategy' },
		{ title: 'Property type', value: ['loan_property', 'property_type'] },
		{ title: 'Number of units', value: ['loan_property', 'unit'] },
		{
			title: 'Purchase price',
			value: ['loan_property', 'purchase_amount'],
			render: (v) => formatAmount(v),
		},
		{
			title: 'Estimated after repair value',
			value: ['loan_property', 'arv'],
			render: (v) => formatAmount(v),
		},
		{
			title: 'Borrower',
			value: 'applicants',
			render: (v) =>
				v
					?.map(
						(item) =>
							`${item?.first_name} ${item.mid_name ? item.mid_name + ' ' : ''} ${item?.last_name}`,
					)
					.join(', '),
		},
		{ title: 'Target closing date', value: 'closing_date' },
		{ title: 'Holdback', value: 'holdback_amount', render: (v) => formatAmount(v) },
	];

	const getLender = async () => {
		setLoading(true);
		const { list } = await request('/organization/lenders');
		setOptions(list);
		setLoading(false);
	};

	useDeferredSyncEffect({
		// predicate: () => step === 1,
		deferred: () => {
			getLender();
		},
		deps: [],
		isUseMicrotask: true,
	});
	return (
		<NewModal
			width={572}
			title="Preview: Loan Quote Email"
			confirmLoading={confirmLoading}
			onCancel={() => setOpen(false)}
			footer={[
				<Button key="cancel" onClick={() => setOpen(false)} size="middle">
					Cancel
				</Button>,
				<Space key="save">
					{step === 2 && (
						<Button onClick={() => setStep(1)} size="middle">
							Back
						</Button>
					)}
					<Button
						type="primary"
						onClick={handleSaveAndSendQuote}
						size="middle"
						loading={confirmLoading}
					>
						{step === 1 ? 'Next' : 'Send quote'}
					</Button>
				</Space>,
			]}
		>
			<div className="h-[64vh] max-h-[700px] overflow-hidden overflow-y-auto">
				{step === 2 ? (
					<>
						<div className="font-medium text-xl text-[#1E2939] pb-2">
							Select a lender
						</div>
						{loading ? (
							<div className="h-[calc(100%-36px)] flex items-center justify-center">
								<Spin />
							</div>
						) : (
							<Checkbox.Group
								style={{ width: '100%' }}
								value={value}
								onChange={(e) => setValue(e)}
							>
								<Row gutter={[16, 16]}>
									{options.map((item, index) => (
										<Col span={24} key={index}>
											<Checkbox
												key={index}
												value={item.id}
												className="border border-base-border !px-[10px] !py-[8px] rounded-md w-full"
											>
												{item.name}
											</Checkbox>
										</Col>
									))}
								</Row>
							</Checkbox.Group>
						)}
					</>
				) : (
					// 预览邮件
					<div className="text-[#1E2939] text-sm">
						<div className="pb-[14px]">Hi xxx team,</div>
						<div className=" pb-6">
							Can you help provide us with a quote based on the following information?
						</div>
						<div className="flex items-start pb-[6px]">
							<HomeTwoTone style={{ fontSize: 18 }} />
							<span className="font-bold pl-1">Property address:</span>
						</div>
						<div className="font-medium">
							{assemblyAddress(rowData?.loan_property?.address, true)}
						</div>
						{list.map((item, index) => (
							<div key={index} className="flex items-center gap-3 pt-[14px]">
								<div className="w-[40%] font-bold shrink-0">{item.title}:</div>
								<div className="flex-1">
									{item?.render?.(get(rowData, item.value)) ||
										get(rowData, item.value) ||
										'--'}
								</div>
							</div>
						))}
						<div className="pb-[2px] pt-[38px]">Best regards</div>
						<div>{userInfo?.organization_name}</div>
					</div>
				)}
			</div>
		</NewModal>
	);
};
export default PreviewQuote;
