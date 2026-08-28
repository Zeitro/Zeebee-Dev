/**
 * 项目动态：POST `/operation/page` 拉取项目级操作记录，按时间倒序展示。
 */
'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { Button, Drawer, Empty, Spin, Timeline } from 'antd';
import { useTranslations } from 'next-intl';
import { dayjs, createRequest } from '@repo/shared-methods/utils';
import { getMemberLabel, unwrapList } from './helpers';

const OPERATION_PAGE_API = '/operation/page';
const PROJECT_ENTITY_TYPE = 'ProjectManagementProject';
const ACTIVITY_PAGE_SIZE = 20;

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
	const { request } = useMemo(() => createRequest(), []);
	const [activities, setActivities] = useState([]);
	const [loading, setLoading] = useState(false);

	const resolveUserName = useCallback(
		(id) => {
			const member = members.find((item) => item.id === id);
			return getMemberLabel(member, id);
		},
		[members],
	);

	const fetchActivities = useCallback(async () => {
		if (!projectId) return [];
		const data = await request(OPERATION_PAGE_API, {
			method: 'POST',
			isHideTips: true,
			body: {
				page: 1,
				limit: ACTIVITY_PAGE_SIZE,
				filter: {
					entity_id: Number(projectId),
					entity_type: PROJECT_ENTITY_TYPE,
				},
			},
		});
		return unwrapList(data);
	}, [projectId, request]);

	useEffect(() => {
		if (!open || !projectId) return undefined;

		let cancelled = false;
		(async () => {
			setLoading(true);
			try {
				const list = await fetchActivities();
				if (!cancelled) setActivities(list);
			} catch {
				if (!cancelled) setActivities([]);
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();

		return () => {
			cancelled = true;
		};
	}, [open, projectId, fetchActivities]);

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
					<Timeline
						className="mt-3"
						items={activities.map((act) => {
							const taskId = act.metadata?.task_id;
							const task = tasks.find((item) => item.id === taskId);
							const userName =
								act.operator_name ||
								(act.operator_id
									? resolveUserName(act.operator_id)
									: t('kanban.system'));
							const description = act.operation || '';
							return {
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
											{dayjs(act.created_at).format('MM-DD HH:mm')}
										</div>
									</div>
								),
							};
						})}
					/>
				)}
			</Spin>
		</Drawer>
	);
}
