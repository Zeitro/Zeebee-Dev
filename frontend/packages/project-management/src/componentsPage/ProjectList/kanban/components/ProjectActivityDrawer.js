/**
 * 项目动态：POST `/operation/page` 拉取项目级操作记录，按时间倒序展示，滚动加载更多。
 */
'use client';

import { useCallback, useMemo } from 'react';
import { Button, Drawer, Empty, Spin, Timeline } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue, formatTime } from '@repo/shared-methods/utils';
import { formatOperationDescription } from '@repo/project-management/lib/operationLog';
import { useOperationActivityList } from '../hooks/useOperationActivityList';
import { getMemberLabel } from './helpers';

/**
 * @param {{
 *   open: boolean,
 *   projectId?: number|string,
 *   tasks?: Array,
 *   members?: Array,
 *   onClose: () => void,
 *   onTaskClick?: (task: object) => void,
 * }} props
 */
export function ProjectActivityDrawer({
	open,
	projectId,
	tasks = [],
	members = [],
	onClose,
	onTaskClick,
}) {
	const t = useTranslations();
	const { activities, loading, loadingMore, hasMore, sentinelRef } = useOperationActivityList({
		open,
		projectId,
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
				const taskId = act.metadata?.task_id;
				const task = tasks.find((item) => item.id === taskId);
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
								{task ? (
									<Button
										type="link"
										size="small"
										className="ml-2 p-0 h-auto text-blue-medium"
										onClick={() => onTaskClick?.(task)}
									>
										#{task.id} {task.title}
									</Button>
								) : null}
							</div>
							<div className="text-xs text-base-label mt-0.5">
								{displayValue(act.created_at, formatTime)}
							</div>
						</div>
					),
				};
			}),
		[activities, onTaskClick, resolveUserName, t, tasks],
	);

	return (
		<Drawer
			title={t('kanban.projectActivity')}
			placement="right"
			width={480}
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
