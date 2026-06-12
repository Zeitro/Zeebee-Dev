'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { Spin } from 'antd';
import { createRequest, displayValue } from '@repo/shared-methods/utils';

async function fetchActivitySummary({ transactionId, dateRange, request }) {
	const params = new URLSearchParams();
	if (dateRange?.[0] && dateRange?.[1]) {
		params.set('start_date', dateRange[0]);
		params.set('end_date', dateRange[1]);
	}
	const query = params.toString() ? `?${params.toString()}` : '';
	return await request(
		`/borrower_property_transaction/${transactionId}/operations/stats${query}`,
		{ method: 'GET' },
	);
}

const ActivitySummary = ({ transactionId, dateRange, activeKey, onSelect, allActivityType }) => {
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const [summary, setSummary] = useState([]);
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState(null);

	useEffect(() => {
		if (!transactionId) {
			return undefined;
		}

		let cancelled = false;
		const statsUrl = `/borrower_property_transaction/${transactionId}/operations/stats`;

		const loadSummary = async () => {
			cancelMatch(statsUrl);
			setLoading(true);
			try {
				const { items, total } = await fetchActivitySummary({
					transactionId,
					dateRange,
					request,
				});
				if (!cancelled) {
					let allActivity = { type: allActivityType, count: total };
					setSummary([allActivity, ...items]);
					setRowData(allActivity);
				}
			} catch {
				if (!cancelled) {
					setSummary([]);
				}
			} finally {
				if (!cancelled) {
					setLoading(false);
				}
			}
		};

		loadSummary();

		return () => {
			cancelled = true;
			cancelMatch(statsUrl);
		};
	}, [cancelMatch, dateRange, request, transactionId, allActivityType]);

	return (
		<div className="relative flex flex-col gap-4 rounded-md border border-base-border bg-white p-4">
			{loading ? (
				<div className="absolute inset-0 z-10 flex items-center justify-center rounded-md bg-white/60">
					<Spin size="small" />
				</div>
			) : null}
			<p className="m-0 text-base font-medium leading-6 text-base-gray">Activity summary</p>
			<div className="flex flex-col gap-4 text-sm leading-[22px]">
				{summary?.map?.((option, index) => {
					const isActive = activeKey === option.type;
					const isLast = index === summary.length - 1;
					return (
						<div
							key={option.type}
							tabIndex={0}
							onClick={() => {
								onSelect?.(option.type);
								setRowData(option);
							}}
							onKeyDown={(event) => {
								if (event.key === 'Enter' || event.key === ' ') {
									event.preventDefault();
									onSelect?.(option.type);
									setRowData(option);
								}
							}}
							className={[
								'flex w-full cursor-pointer items-center justify-between border-0 border-base-border bg-transparent p-0 py-2 text-left hover:[&>span:first-child]:text-foreground',
								isActive ? 'border-b-2' : !isLast ? 'border-b' : '',
								rowData?.type === option.type
									? '[&>span:first-child]:text-foreground'
									: '',
							]
								.filter(Boolean)
								.join(' ')}
						>
							<span className="text-base-gray">{option.type}</span>
							<span className="font-medium text-base-primary">
								{displayValue(option.count)}
							</span>
						</div>
					);
				})}
			</div>
		</div>
	);
};

export default ActivitySummary;
