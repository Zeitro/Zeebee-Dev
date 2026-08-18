/**
 * 任务详情顶栏：面包屑、参与者头像、任务日志、依赖菜单、标题、删除。
 */
'use client';

import {
	ArrowLeftOutlined,
	DeleteOutlined,
	FileTextOutlined,
	LinkOutlined,
	MoreOutlined,
	ProjectOutlined,
} from '@ant-design/icons';
import { Avatar, Button, Dropdown, Input, Tooltip } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';
import { getMemberInitial, getMemberLabel } from '../helpers';

/**
 * @param {{
 *   task: object,
 *   project?: object|null,
 *   parentTask?: object,
 *   assigneeIds?: number[],
 *   members?: Array,
 *   onTaskClick: (task: object) => void,
 *   onTitleBlur: (title: string) => void,
 *   onDelete: () => void,
 *   onOpenLog?: () => void,
 *   onAddDependency?: (mode: 'predecessor' | 'successor') => void,
 * }} props
 */
export function TaskDrawerHeader({
	task,
	project,
	parentTask,
	assigneeIds = [],
	members = [],
	onTaskClick,
	onTitleBlur,
	onDelete,
	onOpenLog,
	onAddDependency,
}) {
	const t = useTranslations();

	const participantItems = [];
	const seen = new Set();
	if (task.creator_id) {
		seen.add(task.creator_id);
		participantItems.push({ userId: task.creator_id, color: 'var(--purple-dark)' });
	}
	assigneeIds.forEach((id) => {
		if (!seen.has(id)) {
			seen.add(id);
			participantItems.push({ userId: id, color: 'var(--blue-medium)' });
		}
	});

	return (
		<div className="px-6 pt-[18px] pb-2.5">
			<div className="flex items-center justify-between mb-1.5">
				<div className="flex items-center gap-1.5 text-xs text-base-label">
					{parentTask ? (
						<Button
							type="text"
							size="small"
							icon={<ArrowLeftOutlined />}
							onClick={() => onTaskClick(parentTask)}
						>
							{t('kanban.backToParent')}
						</Button>
					) : null}
					<span>{displayValue(project?.name)}</span>
					<span>/</span>
					<span>#{task.id}</span>
				</div>
				<div className="flex items-center gap-1">
					{participantItems.length > 0 ? (
						<Avatar.Group
							max={{
								count: 3,
								style: { backgroundColor: 'var(--blue-medium)', fontSize: 12 },
							}}
							size={28}
						>
							{participantItems.map((item) => {
								const member = members.find((row) => row.id === item.userId);
								const name = getMemberLabel(member, item.userId);
								return (
									<Tooltip key={item.userId} title={name}>
										<Avatar size={28} style={{ backgroundColor: item.color }}>
											{getMemberInitial(member, item.userId)}
										</Avatar>
									</Tooltip>
								);
							})}
						</Avatar.Group>
					) : null}
					<Button
						type="text"
						size="small"
						icon={<FileTextOutlined />}
						onClick={onOpenLog}
					>
						{t('kanban.taskLog')}
					</Button>
					<Dropdown
						trigger={['click']}
						menu={{
							items: [
								{
									key: 'add-dependency',
									label: t('kanban.addDependencyTask'),
									icon: <LinkOutlined />,
									children: [
										{
											key: 'add-predecessor',
											label: t('kanban.predecessor'),
											icon: <LinkOutlined />,
										},
										{
											key: 'add-successor',
											label: t('kanban.successor'),
											icon: <LinkOutlined />,
										},
									],
								},
								{ type: 'divider' },
								{
									key: 'delete',
									label: t('common.delete'),
									icon: <DeleteOutlined />,
									danger: true,
								},
							],
							onClick: ({ key }) => {
								if (key === 'delete') onDelete();
								else if (key === 'add-predecessor') onAddDependency?.('predecessor');
								else if (key === 'add-successor') onAddDependency?.('successor');
							},
						}}
					>
						<Button type="text" size="small" icon={<MoreOutlined />} />
					</Dropdown>
				</div>
			</div>
			<Input
				key={task.id}
				variant="borderless"
				defaultValue={task.title}
				onBlur={(event) => onTitleBlur(event.target.value.trim())}
				onPressEnter={(event) => event.currentTarget.blur()}
				className="!px-0 font-semibold text-xl"
			/>
			<div className="flex items-center gap-2 mt-1.5 text-[13px] text-base-label">
				<ProjectOutlined />
				<span>{displayValue(project?.name)}</span>
			</div>
		</div>
	);
}
