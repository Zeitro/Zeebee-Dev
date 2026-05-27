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
	const data = await request(
		`/borrower_property_transaction/${transactionId}/operations/stats${query}`,
		{ method: 'GET' },
	);
	return data?.items || [];
}

const ActivitySummary = ({ transactionId, dateRange, activeKey, onSelect }) => {
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const [summary, setSummary] = useState([]);
	const [loading, setLoading] = useState(false);

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
				const nextSummary = await fetchActivitySummary({
					transactionId,
					dateRange,
					request,
				});
				if (!cancelled) {
					setSummary(nextSummary);
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
	}, [cancelMatch, dateRange, request, transactionId]);

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
					const isActive = activeKey === option.key;
					const isLast = index === summary.length - 1;
					return (
						<div
							key={option.type}
							tabIndex={0}
							onClick={() => onSelect?.(option.type)}
							onKeyDown={(event) => {
								if (event.key === 'Enter' || event.key === ' ') {
									event.preventDefault();
									onSelect?.(option.type);
								}
							}}
							className={[
								'flex w-full cursor-pointer items-center justify-between border-0 border-base-border bg-transparent p-0 py-2 text-left',
								isActive ? 'border-b-2' : !isLast ? 'border-b' : '',
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
