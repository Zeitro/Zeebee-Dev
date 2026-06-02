'use client';

import React from 'react';
import { Skeleton } from 'antd';
import styles from '../../index.module.scss';

const skeletonActive = { active: true };

/**
 * 交易详情页首屏骨架：对齐 page 头部 + Tabs 导航 + 默认 Document Tab 左右栏布局。
 * 纯展示，无数据请求。
 */
const TransactionDetailSkeleton = () => (
	<div className="flex h-full flex-1 flex-col overflow-hidden bg-white">
		<div className="flex shrink-0 flex-col gap-4 px-6 pb-4 pt-6">
			<div className="flex flex-wrap items-start justify-between gap-4">
				<div className="flex min-w-0 flex-1 flex-col gap-3">
					<div className="flex flex-wrap items-center gap-8">
						<Skeleton.Input
							{...skeletonActive}
							size="large"
							style={{ width: 220, minWidth: 220 }}
						/>
						<Skeleton.Input
							{...skeletonActive}
							size="default"
							style={{ width: 140, minWidth: 140 }}
						/>
					</div>
					<Skeleton
						{...skeletonActive}
						title={false}
						paragraph={{ rows: 1, width: ['28%', '22%', '26%'] }}
					/>
				</div>
				<Skeleton.Button {...skeletonActive} size="middle" style={{ width: 88 }} />
			</div>
		</div>

		<div
			className={[styles.transactionDetailTabs, 'flex flex-1 flex-col overflow-hidden'].join(
				' ',
			)}
		>
			<div className="shrink-0 bg-white px-6 pb-3 pt-1">
				<Skeleton
					{...skeletonActive}
					title={false}
					paragraph={{ rows: 1, width: ['18%', '14%', '12%', '10%'] }}
				/>
			</div>
			<div className="flex flex-1 gap-4 overflow-hidden bg-base-text2 p-4">
				<div className="flex w-[311px] shrink-0 flex-col gap-3">
					<div className="flex gap-2">
						<Skeleton.Button {...skeletonActive} block size="middle" />
						<Skeleton.Button {...skeletonActive} size="middle" style={{ width: 48 }} />
					</div>
					<div className="rounded-lg border border-base-border bg-white p-6">
						<Skeleton {...skeletonActive} paragraph={{ rows: 2 }} />
					</div>
					<div className="rounded-lg border border-base-border bg-white p-4">
						<Skeleton {...skeletonActive} title={false} paragraph={{ rows: 6 }} />
					</div>
				</div>
				<div className="min-h-0 flex-1 rounded-lg border border-base-border bg-white p-6">
					<Skeleton {...skeletonActive} paragraph={{ rows: 10 }} />
				</div>
			</div>
		</div>
	</div>
);

export default TransactionDetailSkeleton;
