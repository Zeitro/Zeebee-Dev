/**
 * 任务详情顶栏：面包屑、参与者头像、任务日志、依赖菜单、标题、删除。
 */
'use client';

import {
	ArrowLeftOutlined,
	DeleteOutlined,
	FileTextOutlined,
	// LinkOutlined,
	// MoreOutlined,
} from '@ant-design/icons';
import { Avatar, Button, Dropdown, Input, Tooltip } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';
// import { TASK_DEPENDENCY_DIRECTION } from '@repo/project-management/enum/task';
// import { DEPENDENCY_DIRECTION_OPTIONS } from '../constants';
// import { getMemberInitial, getMemberLabel } from '../helpers';

/**
 * @param {{
 *   task: object,
 *   project?: object|null,
 *   backTask?: object|null,
 *   assigneeIds?: number[],
 *   members?: Array,
 *   onBack?: () => void,
 *   onTitleBlur: (title: string) => void,
 *   onDelete: () => void,
 *   onOpenLog?: () => void,
 *   onAddDependency?: (mode: typeof TASK_DEPENDENCY_DIRECTION.PRE_TASK | typeof TASK_DEPENDENCY_DIRECTION.POST_TASK) => void,
 * }} props
 */
export function TaskDrawerHeader({
	task,
	project,
	backTask,
	assigneeIds = [],
	members = [],
	onBack,
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
		<div className="shrink-0 px-6 pt-5 pb-4 border-b border-base-border bg-white">
			<div className="flex items-center justify-between gap-3 mb-4">
				<div className="flex items-center gap-2 min-w-0 flex-wrap">
					{backTask ? (
						<Button
							type="text"
							size="small"
							icon={<ArrowLeftOutlined />}
							className="!text-base-label hover:!text-base-primary"
							onClick={onBack}
						>
							{t('kanban.backToParent')}
						</Button>
					) : null}
					<span className="inline-flex items-center gap-1.5 text-xs text-base-label min-w-0">
						<span className="px-2 py-0.5 rounded-md bg-base-text2 text-base-gray truncate max-w-[180px]">
							{displayValue(project?.name)}
						</span>
						<span className="text-base-disabled">/</span>
						<span className="font-mono text-[11px] px-1.5 py-0.5 rounded bg-base-active text-base-primary">
							#{task.id}
						</span>
					</span>
				</div>

				<div className="flex items-center gap-0.5 shrink-0">
					{/* {participantItems.length > 0 ? (
						<Avatar.Group
							max={{
								count: 3,
								style: { backgroundColor: 'var(--blue-medium)', fontSize: 12 },
							}}
							size={28}
							className="mr-1"
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
					) : null} */}
					<Button
						type="text"
						size="small"
						icon={<FileTextOutlined />}
						className="!text-base-label"
						onClick={onOpenLog}
					>
						{t('kanban.taskLog')}
					</Button>
					<Button
						type="text"
						size="small"
						icon={<DeleteOutlined />}
						onClick={onDelete}
						danger
					>
						{t('common.delete')}
					</Button>
					{/* <Dropdown
						trigger={['click']}
						menu={{
							items: [
								{
									key: 'add-dependency',
									label: t('kanban.addDependencyTask'),
									icon: <LinkOutlined />,
									children: DEPENDENCY_DIRECTION_OPTIONS.map((item) => ({
										key: item.value,
										label: t(item.labelKey),
										icon: <LinkOutlined />,
									})),
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
								else if (
									key === TASK_DEPENDENCY_DIRECTION.PRE_TASK ||
									key === TASK_DEPENDENCY_DIRECTION.POST_TASK
								) {
									onAddDependency?.(key);
								}
							},
						}}
					>
						<Button type="text" size="small" icon={<MoreOutlined />} />
					</Dropdown> */}
				</div>
			</div>

			<Input
				key={task.id}
				variant="borderless"
				defaultValue={task.title}
				onBlur={(event) => onTitleBlur(event.target.value.trim())}
				onPressEnter={(event) => event.currentTarget.blur()}
				className="!px-0 !py-0 font-semibold text-[22px] leading-snug text-base-text"
			/>
		</div>
	);
}
