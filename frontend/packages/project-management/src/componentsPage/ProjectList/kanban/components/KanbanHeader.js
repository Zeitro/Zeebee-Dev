/**
 * 看板顶栏：返回项目列表、项目名与任务数、成员头像、新建任务。
 */
'use client';

import { ArrowLeftOutlined, PlusOutlined, SettingOutlined, TeamOutlined } from '@ant-design/icons';
import { Avatar, Button, Tooltip } from 'antd';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';
import { getMemberInitial, getMemberLabel } from './helpers';

/**
 * @param {{
 *   project?: object|null,
 *   taskCount?: number,
 *   members?: Array,
 *   projectMembers?: Array,
 *   onCreateTask?: () => void,
 *   showCreate?: boolean,
 * }} props
 */
export function KanbanHeader({
	project,
	taskCount = 0,
	members = [],
	projectMembers = [],
	onCreateTask,
	showCreate = true,
}) {
	const t = useTranslations();
	const router = useRouter();
	const projectId = project?.id;

	return (
		<div className="flex items-center justify-between mb-4">
			<div className="flex items-center gap-3 min-w-0">
				<Button
					icon={<ArrowLeftOutlined />}
					onClick={() => router.push('/project/projectList')}
				>
					{t('common.back')}
				</Button>
				<h2 className="m-0 text-lg font-semibold truncate">
					{displayValue(project?.name)}
					<span className="ml-1 text-sm font-normal text-base-label">
						{t('kanban.taskCount', { count: taskCount })}
					</span>
				</h2>
			</div>
			<div className="flex items-center gap-3 shrink-0">
				<div className="flex items-center gap-1.5" title={t('kanban.members')}>
					<TeamOutlined className="text-base-label" />
					<Avatar.Group
						max={{
							count: 8,
							style: { backgroundColor: 'var(--blue-medium)', fontSize: 12 },
						}}
						size={28}
					>
						{projectMembers.map((item) => {
							const member = members.find((row) => row.id === item.user_id);
							const name = getMemberLabel(member, item.user_id);
							const isOwner = item.user_id === project?.owner_id;
							return (
								<Tooltip
									key={item.id || item.user_id}
									title={`${name}（${isOwner ? t('kanban.owner') : item.role}）`}
								>
									<Avatar
										size={28}
										style={{
											backgroundColor: isOwner
												? 'var(--purple-dark)'
												: 'var(--blue-medium)',
										}}
									>
										{getMemberInitial(member, item.user_id)}
									</Avatar>
								</Tooltip>
							);
						})}
					</Avatar.Group>
				</div>
				{projectId != null ? (
					<Button
						icon={<SettingOutlined />}
						onClick={() => router.push(`/project/${projectId}/settings`)}
					>
						{t('projectList.openSettings')}
					</Button>
				) : null}
				{showCreate ? (
					<Button type="primary" icon={<PlusOutlined />} onClick={onCreateTask}>
						{t('kanban.newTask')}
					</Button>
				) : null}
			</div>
		</div>
	);
}
