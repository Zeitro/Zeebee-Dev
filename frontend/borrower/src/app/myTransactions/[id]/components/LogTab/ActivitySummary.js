'use client';

import React from 'react';
import { displayValue } from '@repo/shared-methods/utils';
import { ACTIVITY_FILTER_OPTIONS } from './constants';

const ActivitySummary = ({ summary, activeKey, onSelect }) => {
	return (
		<div className="flex flex-col gap-4 rounded-md border border-base-border bg-white p-4">
			<p className="m-0 text-base font-medium leading-6 text-base-gray">Activity summary</p>
			<div className="flex flex-col gap-4 text-sm leading-[22px]">
				{ACTIVITY_FILTER_OPTIONS.map((option, index) => {
					const isActive = activeKey === option.key;
					const isLast = index === ACTIVITY_FILTER_OPTIONS.length - 1;
					return (
						<div
							key={option.key}
							type="button"
							onClick={() => onSelect?.(option.key)}
							className={[
								'flex w-full cursor-pointer items-center justify-between border-0 border-base-border bg-transparent p-0 py-2 text-left',
								isActive ? 'border-b-2' : !isLast ? 'border-b' : '',
							]
								.filter(Boolean)
								.join(' ')}
						>
							<span className="text-base-gray">{option.label}</span>
							<span className="font-medium text-base-primary">
								{displayValue(summary?.[option.key])}
							</span>
						</div>
					);
				})}
			</div>
		</div>
	);
};

export default ActivitySummary;
