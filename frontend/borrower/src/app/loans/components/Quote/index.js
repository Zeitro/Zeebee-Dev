'use client';
import { useCallback, useState } from 'react';
import { NewModal, FieldTips } from '@repo/shared-methods';
import { get } from 'lodash';
import { Button, Checkbox, Col, Row, Space, Spin } from 'antd';
import { HomeTwoTone } from '@ant-design/icons';
import { assemblyAddress, formatAmount } from '@repo/shared-methods/utils';
import { useUserStore } from '@/store/user';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

export const PreviewQuote = ({ setOpen, rowData, getRowData, onSave }) => {
	const { userInfo } = useUserStore();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [options, setOptions] = useState([
		// { label: 'Apple', value: 'Apple', className: 'label-1' },
		// { label: 'Pear', value: 'Pear', className: 'label-2' },
		// { label: 'Orange', value: 'Orange', className: 'label-3' },
	]);
	const [value, setValue] = useState([]);
	const [step, setStep] = useState(1);
	const [data, setData] = useState({});
	const handleSaveAndSendQuote = () => {
		if (step === 1) {
			setStep(2);
			return;
		}
		// 保存3rd party loan并发送邮件给lender询价
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

	const updateRowData = useCallback(async () => {
		const row = (await getRowData?.()) || rowData || {};
		setData({ ...row });
	}, [getRowData, rowData]);

	useDeferredSyncEffect({
		predicate: () => step === 1,
		deferred: () => {
			updateRowData();
		},
		deps: [step, updateRowData],
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
												checked={item.value}
												className="border border-base-border !px-[10px] !py-[8px] rounded-md w-full"
											>
												{item.label}
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
						<div className="pb-[14px]">Hi lender,</div>
						<div className=" pb-6">
							Can you help provide us with a quote based on the following information?
						</div>
						<div className="flex items-start pb-[6px]">
							<HomeTwoTone style={{ fontSize: 18 }} />
							<span className="font-bold pl-1">Property address:</span>
						</div>
						<div className="font-medium">
							{assemblyAddress(data?.loan_property?.address, true)}
						</div>
						{list.map((item, index) => (
							<div key={index} className="flex items-center gap-3 pt-[14px]">
								<div className="w-[40%] font-bold shrink-0">{item.title}:</div>
								<div className="flex-1">
									{item?.render?.(get(data, item.value)) ||
										get(data, item.value) ||
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

// 公共询价按钮
const QuoteBtn = ({ disabled, rowData, getRowData, onSave }) => {
	const [open, setOpen] = useState(false);
	const handleClick = () => {
		setOpen(true);
	};
	return (
		<FieldTips title={disabled ? 'Only for pending stage & quoting status loan' : ''}>
			<Button type="primary" onClick={handleClick} ghost disabled={disabled}>
				Save & quote
			</Button>
			{open && (
				<PreviewQuote
					setOpen={setOpen}
					rowData={rowData}
					onSave={onSave}
					getRowData={getRowData}
				/>
			)}
		</FieldTips>
	);
};

export default QuoteBtn;
