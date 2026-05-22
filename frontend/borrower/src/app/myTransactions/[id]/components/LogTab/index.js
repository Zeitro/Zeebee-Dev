'use client';

import React, { useCallback, useMemo, useState } from 'react';
import { InfiniteVirtuoso, CommonFormDate } from '@repo/shared-methods';
import ActivityLogItem from './ActivityLogItem';
import ActivitySummary from './ActivitySummary';
import { ACTIVITY_LOG_ITEM_HEIGHT, MOCK_ACTIVITY_SUMMARY, fetchActivityLogs } from './constants';

const LogTab = ({ transactionId }) => {
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

	const getData = useCallback(
		async ({ page, limit }) => fetchActivityLogs({ page, limit, ...filterParams }),
		[filterParams],
	);

	const renderItem = useCallback(
		(item, _index, { style } = {}) => <ActivityLogItem item={item} />,
		[],
	);

	return (
		<div className="flex h-full min-h-0 gap-4 px-4 py-4">
			<div className="flex min-h-[360px] min-w-0 flex-1 flex-col rounded-lg border border-base-border bg-white px-6 py-4">
				<InfiniteVirtuoso
					key={`${transactionId}-${activityType}-${dateRange?.join('-') || ''}`}
					boxClassName="min-h-0 flex-1"
					className="h-full"
					itemHeight={ACTIVITY_LOG_ITEM_HEIGHT}
					getData={getData}
					params={filterParams}
					renderItem={renderItem}
					defaultItemHeight={ACTIVITY_LOG_ITEM_HEIGHT}
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
					summary={MOCK_ACTIVITY_SUMMARY}
					activeKey={activityType}
					onSelect={setActivityType}
				/>
			</div>
		</div>
	);
};

export default LogTab;
