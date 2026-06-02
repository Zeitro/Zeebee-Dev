'use client';

import React from 'react';
import { Progress } from 'antd';
import { CaretRightFilled, FileTextOutlined } from '@ant-design/icons';
import { displayValue, formatDate, formatTime } from '@repo/shared-methods/utils';
import { getChecklistProgress } from '@/app/myTransactions/[id]/components/DocumentsTab/checklist';
import TransactionCardExpandedContent from './TransactionCardExpandedContent';
import TransactionStatusSelect from '@/app/myTransactions/[id]/components/TransactionStatusSelect';

function getProgressStatus(status, completed, total) {
	if (total > 0 && completed >= total) return 'success';
	if (status === 'Pending') return 'active';
	return 'normal';
}

const TransactionCard = ({
	item,
	isExpanded,
	onToggle,
	onRefreshList,
	transactionStatusOptions,
}) => {
	const status = item?.status;
	const { completed, total } = getChecklistProgress(item?.checklist);
	const progressPercent = total > 0 ? Math.round((completed / total) * 100) : 0;

	const handleHeaderClick = () => {
		onToggle?.();
	};

	const handleHeaderKeyDown = (e) => {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			onToggle?.();
		}
	};

	return (
		<div className="flex min-w-0 flex-1 flex-col overflow-hidden rounded-lg border border-base-border bg-background">
			<div
				role="button"
				tabIndex={0}
				className="flex cursor-pointer flex-col items-stretch gap-6 px-4 py-6 lg:flex-row lg:items-center lg:justify-between"
				onClick={handleHeaderClick}
				onKeyDown={handleHeaderKeyDown}
				aria-expanded={isExpanded}
			>
				<div className="flex min-w-0 flex-1 items-center gap-4">
					<CaretRightFilled
						className={[
							'shrink-0 text-base text-base-label transition-transform duration-200',
							isExpanded ? 'rotate-90' : '',
						].join(' ')}
					/>
					<div className="flex min-w-0 flex-col gap-2">
						<div className="flex flex-wrap items-center gap-4">
							<p className="m-0 text-xl font-medium leading-7 text-foreground">
								{item?.type ? `${item.type} transaction` : displayValue(null)}
							</p>
							<div onClick={(e) => e.stopPropagation()}>
								<TransactionStatusSelect
									value={status}
									options={transactionStatusOptions}
									id={item?.id}
									onSuccess={onRefreshList}
									size="small"
								/>
							</div>
						</div>
						<div className="flex flex-wrap gap-6 text-sm font-medium leading-[22px] text-base-label">
							<span>
								<span className="mr-2">ID</span>
								{displayValue(item?.code)}
							</span>
							<span>
								<span className="mr-2">COE</span>
								{displayValue(item?.coe_date, formatDate)}
							</span>
							<span>
								<span className="mr-2">Updated</span>
								{displayValue(item?.updated_at, formatTime)}
							</span>
						</div>
					</div>
				</div>
				<div className="flex w-full shrink-0 flex-col gap-1 lg:w-[200px]">
					<div className="flex items-center justify-between gap-2">
						<span className="inline-flex items-center gap-2 text-sm font-medium leading-[22px] text-base-label">
							<FileTextOutlined />
							Documents
						</span>
						<span className="text-sm font-medium leading-[22px] text-foreground">
							{total > 0 ? `${completed}/${total}` : displayValue(null)}
						</span>
					</div>
					<Progress
						percent={progressPercent}
						showInfo={false}
						strokeColor={
							total > 0 && completed >= total
								? 'var(--green-medium)'
								: status === 'Pending'
									? 'var(--orange-dark)'
									: 'var(--base-border)'
						}
						trailColor="var(--base-text2)"
						size="small"
						status={getProgressStatus(status, completed, total)}
					/>
				</div>
			</div>
			{isExpanded && item?.id != null ? (
				<TransactionCardExpandedContent
					transactionId={item.id}
					initialRowData={item}
					onRefreshList={onRefreshList}
				/>
			) : null}
		</div>
	);
};

export default TransactionCard;
