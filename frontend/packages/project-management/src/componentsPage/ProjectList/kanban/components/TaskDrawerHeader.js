/**
 * 任务详情顶栏：面包屑、标题失焦保存、删除菜单。
 */
'use client';

import {
	ArrowLeftOutlined,
	DeleteOutlined,
	MoreOutlined,
	ProjectOutlined,
} from '@ant-design/icons';
import { Button, Dropdown, Input } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';

/**
 * @param {{
 *   task: object,
 *   project?: object|null,
 *   parentTask?: object,
 *   onTaskClick: (task: object) => void,
 *   onTitleBlur: (title: string) => void,
 *   onDelete: () => void,
 * }} props
 */
export function TaskDrawerHeader({
	task,
	project,
	parentTask,
	onTaskClick,
	onTitleBlur,
	onDelete,
}) {
	const t = useTranslations();
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
				<Dropdown
					trigger={['click']}
					menu={{
						items: [
							{
								key: 'delete',
								label: t('common.delete'),
								icon: <DeleteOutlined />,
								danger: true,
							},
						],
						onClick: ({ key }) => {
							if (key === 'delete') onDelete();
						},
					}}
				>
					<Button type="text" size="small" icon={<MoreOutlined />} />
				</Dropdown>
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
