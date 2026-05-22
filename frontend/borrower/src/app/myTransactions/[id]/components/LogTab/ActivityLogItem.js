'use client';

import React from 'react';
import { Avatar } from 'antd';
import { displayValue, formatAmount, formatDate, formatTime } from '@repo/shared-methods/utils';

const formatFieldValue = (field) => {
	const raw = field?.value;
	if (field?.format === 'amount') {
		return displayValue(raw, formatAmount);
	}
	return displayValue(raw);
};

const formatActivityDateTime = (value) => {
	if (!value) return displayValue(null);
	const datePart = formatDate(value);
	const timePart = formatTime(value);
	if (!datePart || datePart === displayValue(null)) return displayValue(null);
	if (!timePart || timePart === displayValue(null)) return datePart;
	return `${datePart} ${timePart}`;
};

const ActivityLogItem = ({ item }) => {
	return (
		<div className="flex flex-col gap-4 border-b border-base-border py-4 mb-2">
			<div className="flex items-center justify-between gap-4">
				<p className="m-0 min-w-0 text-base font-medium leading-6 text-foreground">
					{displayValue(item?.title)}
				</p>
				<div className="flex shrink-0 items-center gap-4">
					<div className="flex items-center gap-2">
						<Avatar
							size={24}
							className="!text-xs !font-normal shrink-0"
							style={{
								backgroundColor: 'var(--orange-light)',
								color: 'var(--orange-dark)',
							}}
						>
							{displayValue(item?.operator_initials)}
						</Avatar>
						<span className="whitespace-nowrap text-sm leading-[22px] text-foreground">
							{displayValue(item?.operator_name)}
						</span>
					</div>
					<span className="whitespace-nowrap text-sm leading-[22px] text-base-gray">
						{formatActivityDateTime(item?.created_at)}
					</span>
				</div>
			</div>
			{(item?.fields || []).map((field, index) => (
				<div
					key={`${field?.label}-${index}`}
					className="flex items-center gap-6 text-sm leading-[22px]"
				>
					<span className="w-[120px] shrink-0 text-base-gray">
						{displayValue(field?.label)}
					</span>
					<span className="min-w-0 font-medium text-foreground">
						{formatFieldValue(field)}
					</span>
				</div>
			))}
		</div>
	);
};

export default ActivityLogItem;
