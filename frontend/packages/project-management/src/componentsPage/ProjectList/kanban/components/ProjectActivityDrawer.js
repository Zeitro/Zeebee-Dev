/**
 * 项目动态：聚合各任务 activity，按时间倒序展示。
 */
'use client';

import { Drawer, Empty, Timeline } from 'antd';
import { useTranslations } from 'next-intl';
import { dayjs } from '@repo/shared-methods/utils';
import { formatActivityDescription, getMemberLabel } from './helpers';

/**
 * @param {{
 *   open: boolean,
 *   activities?: Array,
 *   tasks?: Array,
 *   members?: Array,
 *   onClose: () => void,
 *   onTaskClick?: (task: object) => void,
 * }} props
 */
export function ProjectActivityDrawer({
	open,
	activities = [],
	tasks = [],
	members = [],
	onClose,
	onTaskClick,
}) {
	const t = useTranslations();

	const resolveUserName = (id) => {
		const member = members.find((item) => item.id === id);
		return getMemberLabel(member, id);
	};

	return (
		<Drawer
			title={t('kanban.projectActivity')}
			placement="right"
			width={480}
			open={open}
			onClose={onClose}
		>
			{activities.length === 0 ? (
				<Empty description={t('kanban.activityEmptyList')} image={Empty.PRESENTED_IMAGE_SIMPLE} />
			) : (
				<Timeline
					className="mt-3"
					items={activities.map((act) => {
						const task = tasks.find((item) => item.id === act.task_id);
						const userName = act.user_id ? resolveUserName(act.user_id) : t('kanban.system');
						const description = formatActivityDescription(act, t, resolveUserName);
						return {
							children: (
								<div>
									<div className="text-xs">
										<strong>{userName}</strong> {description}
										{task ? (
											<button
												type="button"
												className="ml-2 text-blue-medium hover:underline"
												onClick={() => onTaskClick?.(task)}
											>
												#{task.id} {task.title}
											</button>
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
		</Drawer>
	);
}
