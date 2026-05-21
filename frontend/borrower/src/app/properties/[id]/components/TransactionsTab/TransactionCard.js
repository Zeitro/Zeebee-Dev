'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { Progress } from 'antd';
import { CaretRightFilled, FileTextOutlined } from '@ant-design/icons';
import { displayValue, formatDate, formatTime } from '@repo/shared-methods/utils';

function formatTransactionId(item) {
	if (item?.transaction_no) return item.transaction_no;
	if (item?.id != null && item?.id !== '') return `TC${item.id}`;
	return displayValue(null);
}

function getDocumentProgress(item) {
	const completed =
		item?.document_completed_count ?? item?.documents_completed ?? item?.document_count ?? 0;
	const total = item?.document_total_count ?? item?.documents_total ?? item?.document_total ?? 0;
	return { completed, total };
}

function getProgressStatus(status, completed, total) {
	if (total > 0 && completed >= total) return 'success';
	if (status === 'Pending') return 'active';
	return 'normal';
}

const TransactionCard = ({ item }) => {
	const router = useRouter();
	const status = item?.status;
	const tagClass =
		status === 'Pending'
			? 'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]'
			: status === 'Cancelled'
				? 'border-base-error bg-gray-light text-base-error'
				: 'border-base-border bg-gray-light text-foreground';
	const { completed, total } = getDocumentProgress(item);
	const progressPercent = total > 0 ? Math.round((completed / total) * 100) : 0;
	const updatedRaw = item?.updated_at ?? item?.updatedAt;
	const updatedDate = updatedRaw ? formatDate(updatedRaw) : displayValue(null);
	const updatedTime = updatedRaw ? formatTime(updatedRaw) : null;
	const coeRaw = item?.coe_date ?? item?.coeDate;

	return (
		<div
			role="button"
			tabIndex={0}
			className="flex min-w-0 flex-1 cursor-pointer flex-col items-stretch gap-6 rounded-lg border border-base-border bg-background px-4 py-6 lg:flex-row lg:items-center lg:justify-between"
			onClick={() => item?.id != null && router.push(`/myTransactions/${item.id}`)}
			onKeyDown={(e) => {
				if (e.key === 'Enter' && item?.id != null) {
					router.push(`/myTransactions/${item.id}`);
				}
			}}
		>
			<div className="flex min-w-0 flex-1 items-center gap-4">
				<CaretRightFilled className="shrink-0 text-base text-base-label" />
				<div className="flex min-w-0 flex-col gap-2">
					<div className="flex flex-wrap items-center gap-4">
						<p className="m-0 text-xl font-medium leading-7 text-foreground">
							{item?.type ? `${item.type} transaction` : displayValue(null)}
						</p>
						{status ? (
							<span
								className={[
									'inline-flex whitespace-nowrap rounded-[2px] border px-2 text-xs leading-5',
									tagClass,
								].join(' ')}
							>
								{status}
							</span>
						) : null}
					</div>
					<div className="flex flex-wrap gap-6 text-sm font-medium leading-[22px] text-base-label">
						<span>
							<span className="mr-2">ID</span>
							{formatTransactionId(item)}
						</span>
						<span>
							<span className="mr-2">COE</span>
							{displayValue(coeRaw, formatDate)}
						</span>
						<span>
							<span className="mr-2">Updated</span>
							{updatedDate}
							{updatedTime ? ` ${updatedTime}` : ''}
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
	);
};

export default TransactionCard;
