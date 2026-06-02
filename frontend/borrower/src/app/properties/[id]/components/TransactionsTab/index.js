'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { Button, Spin, Timeline } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import TransactionCard from './TransactionCard';
import { useTransactionList } from './useTransactionList';
import styles from './index.module.scss';
import { PermissionSlot } from '@repo/shared-methods';
import { PERMISSIONS } from '@config/permissions';
import { useGetOptions } from '@repo/shared-methods/hooks';

const getTimelineStatus = (status) => {
	if (status === 'Pending') return 'pending';
	if (status === 'Cancelled') return 'cancelled';
	return 'default';
};

const TransactionsTab = ({ propertyId, onNewTransaction, onRegisterRefresh }) => {
	const { list, loading, refreshTransactionList } = useTransactionList(propertyId);
	const [expandedTransactionId, setExpandedTransactionId] = useState(null);
	const { transactionStatusOptions } = useGetOptions([
		{ key: 'borrower_property_transaction_status', valueName: 'transactionStatusOptions' },
	]);
	useEffect(() => {
		onRegisterRefresh?.(refreshTransactionList);
	}, [onRegisterRefresh, refreshTransactionList]);

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
							onRefreshList={refreshTransactionList}
							transactionStatusOptions={transactionStatusOptions}
						/>
					),
					className: statusType === 'pending' ? styles.itemPending : '',
				};
			}),
		[expandedTransactionId, list, refreshTransactionList, transactionStatusOptions],
	);

	return (
		<div className={[styles.transactionsTab].join(' ')}>
			<div className="mb-[22px] flex items-center justify-between gap-4 border-b border-base-border ml-[62px] pb-2 mr-4">
				<h2 className="m-0 text-base font-medium leading-6 text-foreground">
					{list.length} Transaction history
				</h2>
				<PermissionSlot code={PERMISSIONS?.TRANSACTION?.ADD}>
					<Button type="primary" icon={<PlusOutlined />} onClick={onNewTransaction}>
						Add transaction
					</Button>
				</PermissionSlot>
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
