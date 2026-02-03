'use client';
import React, { useState } from 'react';
import PaidSvg from '@/assets/images/paid.svg';
import AccountCircle from '@/assets/images/account_circle.svg';
import SellSvg from '@/assets/images/sell.svg';
import Calender from '@/assets/images/calendar_today.svg';
import { MultiVirtuosoDndKit } from '@repo/shared-methods';
import { pendingStatusColor } from '@/enum/loan';
import {
	assemblyAddress,
	createRequest,
	formatAmount,
	formatDate,
} from '@repo/shared-methods/utils';
import { message } from 'antd';
import { useRouter } from 'next/navigation';
import { throttle } from 'lodash';
const Pipeline = () => {
	const router = useRouter();
	const { request } = createRequest();
	const [items, setItems] = useState([
		{ value: 'Quoting', label: 'Quoting', total: 0, totalAmount: 0, list: [] },
		{ value: 'Processing', label: 'Processing', total: 0, totalAmount: 0, list: [] },
		{ value: 'Underwriting', label: 'Underwriting', total: 0, totalAmount: 0, list: [] },
		{ value: 'Closing', label: 'Closing', total: 0, totalAmount: 0, list: [] },
	]);

	const handleDragEnd = async ({ movingItem, key: pending_status }) => {
		await request(`/loan/${movingItem.id}/pending_status`, {
			method: 'PATCH',
			body: { pending_status },
		});
		message.open({
			type: 'success',
			content: 'Update pending status successful！',
		});
	};

	const handleClick = throttle(
		(e, item) => {
			e.preventDefault();
			e.stopPropagation();
			router.push(`/loans/zeeBeeDetails?id=${item?.id}`);
		},
		1000,
		{ trailing: false },
	);

	return (
		<MultiVirtuosoDndKit
			itemHeight={196}
			items={items}
			setItems={setItems}
			enableDrag={true}
			parentClassName="!h-[calc(100%-84px)]  py-4 px-6 gap-4 overflow-x-auto"
			columnClassName="w-[338px] min-w-[338px] max-w-[338px]"
			boxClassName="flex-1 bg-[#F6F6F6] py-4 px-4"
			renderItem={(item, idx, options) => {
				const address = assemblyAddress(item?.loan_property?.address);
				const total_amount = formatAmount(item?.total_amount);
				const borrower = item?.applicants?.map((item) => item.full_name).join(', ');
				const created_time = formatDate(item?.created_at);
				return (
					<div
						className="bg-white py-2 px-4 rounded-xl shadow-[0_2px_8px_0_rgba(0,0,0,0.15)] mb-4 flex flex-col gap-4"
						{...options}
						onClick={(e) => handleClick(e, item)}
					>
						<div className="text-base font-medium pb-2 common-ellipsis" title={address}>
							{address}
						</div>
						{/* <div className="pl-[2px] text-sm before:content-[''] before:w-[14px] before:h-[14px] before:bg-[#1677FF] before:inline-block before:rounded-full flex items-center gap-2">
						New application
					</div> */}
						<div
							className="flex items-center text-base-gray  gap-2"
							title={`Total amount: ${total_amount}`}
						>
							<PaidSvg className="flex-shrink-0" />
							{total_amount}
						</div>
						<div
							className="flex items-center text-base-gray  gap-2"
							title={`Borrower: ${borrower}`}
						>
							<AccountCircle className="flex-shrink-0" />
							<span className="common-ellipsis">{borrower || '--'}</span>
						</div>
						<div
							className="flex items-center  text-base-gray  gap-2"
							title={`Loan type: ${item?.loan_type} (Loan purpose: ${item?.loan_purpose})`}
						>
							<SellSvg title="Loan type" className="flex-shrink-0" />
							{item?.loan_type} ({item?.loan_purpose})
						</div>
						<div
							className="flex items-center text-base-gray  gap-2"
							title={`Created time: ${created_time}`}
						>
							<Calender className="flex-shrink-0" />
							{created_time || '--'}
						</div>
					</div>
				);
			}}
			getData={(params) => {
				return request('/loan/page', {
					method: 'POST',
					body: {
						filter: {
							stage: 'Pending',
							pending_status: [params.itemKey],
						},
						page: params.page,
						limit: params.limit,
					},
				});
			}}
			onDragEnd={handleDragEnd}
			renderHeader={(pending_status, total) => (
				<div className="py-4 px-4 flex items-center gap-2 flex-shrink-0">
					<div
						style={{
							background: pendingStatusColor?.[pending_status]?.bg,
							color: pendingStatusColor?.[pending_status]?.color,
							borderColor: pendingStatusColor?.[pending_status]?.border,
						}}
						className="px-2 rounded-sm text-sm/[22px] font-medium border"
					>
						{pending_status} {total}
					</div>
					{/* <span className="text-sm text-base-gray">$ 12,345,4678</span> */}
				</div>
			)}
		/>
	);
};
export default Pipeline;
