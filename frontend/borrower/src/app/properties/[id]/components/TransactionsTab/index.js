'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { Button, Spin, Timeline } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
import TransactionCard from './TransactionCard';
import styles from './index.module.scss';

const TRANSACTION_LIST_API = '/transaction/page';

const getTimelineStatus = (status) => {
	if (status === 'Pending') return 'pending';
	if (status === 'Cancelled') return 'cancelled';
	return 'default';
};

const TransactionsTab = ({ propertyId, onNewTransaction, refreshKey = 0 }) => {
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [expandedTransactionId, setExpandedTransactionId] = useState(null);
	const [list, setList] = useState([
		{ id: 1, type: 'Purchase', status: 'Pending' },
		{ id: 2, type: 'Sale', status: 'Completed' },
	]);

	useEffect(() => {
		if (!propertyId) return;
		let cancelled = false;
		void (async () => {
			await Promise.resolve();
			if (cancelled) return;
			try {
				setLoading(true);
				cancelMatch(TRANSACTION_LIST_API);
				await request(TRANSACTION_LIST_API, {
					method: 'POST',
					body: {
						page: 1,
						limit: 100,
						filter: {
							property_id: Number(propertyId),
						},
					},
				});
				// if (!cancelled) setList(res?.list ?? []);
			} catch {
				// if (!cancelled) setList([]);
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [cancelMatch, propertyId, refreshKey, request]);

	const timelineItems = useMemo(
		() =>
			list.map((item, index) => {
				const statusType = getTimelineStatus(item?.status);
				return {
					key: item.id ?? index,
					dot: (
						<span
							className={[
								styles.timelineDot,
								statusType === 'pending'
									? styles.timelineDotPending
									: statusType === 'cancelled'
										? styles.timelineDotCancelled
										: styles.timelineDotDefault,
							].join(' ')}
						/>
					),
					children: (
						<TransactionCard
							item={item}
							isExpanded={expandedTransactionId === item?.id}
							onToggle={() =>
								setExpandedTransactionId((prev) =>
									prev === item?.id ? null : item?.id,
								)
							}
						/>
					),
					className: statusType === 'pending' ? styles.itemPending : '',
				};
			}),
		[expandedTransactionId, list],
	);

	return (
		<div className={[styles.transactionsTab].join(' ')}>
			<div className="mb-[22px] flex items-center justify-between gap-4 border-b border-base-border ml-[62px] pb-2 mr-4">
				<h2 className="m-0 text-base font-medium leading-6 text-foreground">
					{list.length} Transaction history
				</h2>
				<Button type="primary" icon={<PlusOutlined />} onClick={onNewTransaction}>
					New transaction
				</Button>
			</div>
			<Spin spinning={loading} wrapperClassName="mr-4 ml-5 pt-2">
				{list.length ? (
					<Timeline mode="left" items={timelineItems} />
				) : (
					<div className="py-12 text-center text-base-label">No transactions yet.</div>
				)}
			</Spin>
		</div>
	);
};

export default TransactionsTab;
