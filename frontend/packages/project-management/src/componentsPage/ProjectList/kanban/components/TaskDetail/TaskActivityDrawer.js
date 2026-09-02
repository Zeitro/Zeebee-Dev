/**
 * 单任务操作日志：与项目动态共用 `/operation/page`，filter 额外传 task_id。
 */
'use client';

import { useCallback, useMemo } from 'react';
import { Drawer, Empty, Spin, Timeline } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue, formatTime } from '@repo/shared-methods/utils';
import { formatOperationDescription } from '@repo/project-management/lib/operationLog';
import { useOperationActivityList } from '../../hooks/useOperationActivityList';
import { getMemberLabel } from '../helpers';

/**
 * @param {{
 *   open: boolean,
 *   projectId?: number|string,
 *   taskId?: number|string,
 *   members?: Array,
 *   onClose: () => void,
 * }} props
 */
export function TaskActivityDrawer({ open, projectId, taskId, members = [], onClose }) {
	const t = useTranslations();
	const { activities, loading, loadingMore, hasMore, sentinelRef } = useOperationActivityList({
		open,
		projectId,
		taskId,
		requireTaskId: true,
	});

	const resolveUserName = useCallback(
		(id) => {
			const member = members.find((item) => item.id === id);
			return getMemberLabel(member, id);
		},
		[members],
	);

	const timelineItems = useMemo(
		() =>
			activities.map((act) => {
				const userName =
					act.operator_name ||
					(act.operator_id ? resolveUserName(act.operator_id) : t('kanban.system'));
				const description = formatOperationDescription(act, t, {
					resolveUserName,
				});
				return {
					key: act.id,
					children: (
						<div>
							<div className="text-xs">
								<strong>{userName}</strong> {description}
							</div>
							<div className="text-xs text-base-label mt-0.5">
								{displayValue(act.created_at, formatTime)}
							</div>
						</div>
					),
				};
			}),
		[activities, resolveUserName, t],
	);

	return (
		<Drawer
			title={t('kanban.taskLog')}
			placement="right"
			width={420}
			open={open}
			onClose={onClose}
		>
			<Spin spinning={loading}>
				{!loading && activities.length === 0 ? (
					<Empty
						description={t('kanban.activityEmptyList')}
						image={Empty.PRESENTED_IMAGE_SIMPLE}
					/>
				) : (
					<>
						<Timeline className="mt-3" items={timelineItems} />
						{hasMore || loadingMore ? (
							<div
								ref={sentinelRef}
								className="flex items-center justify-center py-3 min-h-8"
							>
								{loadingMore ? <Spin size="small" /> : null}
							</div>
						) : null}
					</>
				)}
			</Spin>
		</Drawer>
	);
}
