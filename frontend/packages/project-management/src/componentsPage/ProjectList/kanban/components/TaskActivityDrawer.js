/**
 * 单任务操作日志抽屉。
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
 *   members?: Array,
 *   onClose: () => void,
 * }} props
 */
export function TaskActivityDrawer({ open, activities = [], members = [], onClose }) {
	const t = useTranslations();

	const resolveUserName = (id) => {
		const member = members.find((item) => item.id === id);
		return getMemberLabel(member, id);
	};

	return (
		<Drawer
			title={t('kanban.taskLog')}
			placement="right"
			width={420}
			open={open}
			onClose={onClose}
		>
			{activities.length === 0 ? (
				<Empty description={t('kanban.activityEmptyList')} image={Empty.PRESENTED_IMAGE_SIMPLE} />
			) : (
				<Timeline
					className="mt-3"
					items={activities.map((act) => {
						const userName = act.user_id ? resolveUserName(act.user_id) : t('kanban.system');
						return {
							children: (
								<div>
									<div className="text-xs">
										<strong>{userName}</strong>{' '}
										{formatActivityDescription(act, t, resolveUserName)}
									</div>
									<div className="text-xs text-base-label mt-0.5">
										{dayjs(act.created_at).format('YYYY-MM-DD HH:mm')}
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
