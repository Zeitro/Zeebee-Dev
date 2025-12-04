'use client';
import React, { useState } from 'react';
import PaidSvg from '@/assets/images/Paid.svg';
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
const Pipeline = () => {
	const { request } = createRequest();
	const [items, setItems] = useState({
		Quoting: [],
		Processing: [],
		Underwriting: [],
		Closing: [],
	});

	const handleDragEnd = async (activeRow, pending_status) => {
		await request(`/loan/${activeRow.id}/pending_status`, {
			method: 'PATCH',
			body: { pending_status },
		});
		message.open({
			type: 'success',
			content: 'Update pending status successful！',
		});
	};

	return (
		<MultiVirtuosoDndKit
			width={338}
			items={items}
			setItems={setItems}
			enableDrag={true}
			parentClassName="!h-[calc(100%-84px)]  py-4 px-6 gap-4 overflow-x-auto"
			columnClassName="w-[338px] min-w-[338px] max-w-[338px]"
			boxClassName="flex-1 bg-[#F6F6F6] py-4 px-4"
			renderItem={(item, idx, options) => (
				<div
					className="bg-white py-2 px-4 rounded-xl shadow-[0_2px_8px_0_rgba(0,0,0,0.15)] mb-4 flex flex-col gap-4"
					{...options}
				>
					<div
						className="text-base font-medium pb-2 common-ellipsis"
						title={assemblyAddress(item?.loan_property?.address)}
					>
						{assemblyAddress(item?.loan_property?.address)}
					</div>
					<div className="pl-[2px] text-sm before:content-[''] before:w-[14px] before:h-[14px] before:bg-[#1677FF] before:inline-block before:rounded-full flex items-center gap-2">
						New application
					</div>
					<div className="flex items-center text-base-gray  gap-2">
						<PaidSvg />
						{formatAmount(item?.total_amount)}
					</div>
					<div className="flex items-center text-base-gray  gap-2">
						<AccountCircle />
						Jacob Jones
					</div>
					<div className="flex items-center  text-base-gray  gap-2">
						<SellSvg />
						{item?.loan_type}({item?.loan_purpose})
					</div>
					<div className="flex items-center text-base-gray  gap-2">
						<Calender />
						{formatDate(item?.closing_date) || '--'}
					</div>
				</div>
			)}
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
