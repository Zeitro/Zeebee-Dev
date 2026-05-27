'use client';

import React, { useCallback, useMemo, useState } from 'react';
import { InfiniteVirtuoso, CommonFormDate } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import ActivityLogItem from './ActivityLogItem';
import ActivitySummary from './ActivitySummary';
const LogTab = ({ transactionId }) => {
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const [activityType, setActivityType] = useState('all');
	const [dateRange, setDateRange] = useState(null);

	const filterParams = useMemo(
		() => ({
			transactionId,
			activityType,
			dateRange: dateRange || undefined,
		}),
		[activityType, dateRange, transactionId],
	);

	async function fetchActivityLogs({ transactionId, activityType, dateRange }) {
		if (!transactionId || !request) {
			return { list: [], total: 0 };
		}

		const params = new URLSearchParams({});
		if (activityType && activityType !== 'all') {
			params.set('operation_type', activityType);
		}
		if (dateRange?.[0] && dateRange?.[1]) {
			params.set('start_date', dateRange[0]);
			params.set('end_date', dateRange[1]);
		}

		const data = await request(
			`/borrower_property_transaction/${transactionId}/operations?${params.toString()}`,
			{ method: 'GET' },
		);
		const list = data?.list ?? [];
		return { list, total: list.length };
	}

	const getData = useCallback(
		async ({ transactionId: txId, activityType: type, dateRange: range }) => {
			const operationsUrl = `/borrower_property_transaction/${txId}/operations`;
			cancelMatch(operationsUrl);
			return fetchActivityLogs({
				transactionId: txId,
				activityType: type,
				dateRange: range,
			});
		},
		[cancelMatch, request],
	);

	const renderItem = useCallback(
		(item, _index, { style } = {}) => <ActivityLogItem item={item} style={style} />,
		[],
	);

	return (
		<div className="flex h-full min-h-0 gap-4 px-4 py-4">
			<div className="flex min-h-[360px] min-w-0 flex-1 flex-col rounded-lg border border-base-border bg-white px-6 py-4">
				<InfiniteVirtuoso
					key={`${transactionId}-${activityType}-${dateRange?.join('-') || ''}`}
					boxClassName="min-h-0 flex-1"
					minItemHeight={100}
					enableLoadMore={false}
					getData={getData}
					params={filterParams}
					renderItem={renderItem}
				/>
			</div>

			<div className="flex w-[260px] shrink-0 flex-col gap-2">
				<CommonFormDate
					type="RangePicker"
					size="middle"
					value={dateRange}
					onChange={setDateRange}
					allowClear
				/>
				<ActivitySummary
					transactionId={transactionId}
					dateRange={dateRange}
					activeKey={activityType}
					onSelect={setActivityType}
				/>
			</div>
		</div>
	);
};

export default LogTab;
