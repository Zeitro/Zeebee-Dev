/**
 * 详情抽屉内子任务列表：展示直属子任务。
 */
'use client';

import { CheckCircleFilled, RightOutlined } from '@ant-design/icons';
import { Button, Tag } from 'antd';
import { displayValue } from '@repo/shared-methods/utils';
import { TASK_PHASE } from '@repo/project-management/enum/task';
import { getSubtasks } from '../helpers';

/**
 * @param {{
 *   parentId: number,
 *   tasks?: Array,
 *   workflow?: object|null,
 *   onTaskClick: (task: object) => void,
 * }} props
 */
export function SubtaskList({ parentId, tasks = [], workflow, onTaskClick }) {
	const subtasks = getSubtasks(tasks, parentId);
	if (!subtasks.length) return null;

	return (
		<div className="py-3 px-4">
			{subtasks.map((item) => {
				const state = workflow?.states?.find((row) => row.id === item.workflow_state_id);
				const isDone = item?.phase === TASK_PHASE.FINISHED;
				return (
					<Button
						key={item.id}
						type="text"
						block
						className="group flex w-full items-center !justify-start gap-2.5 !px-0 !py-2.5 h-auto text-left text-[13px] hover:!bg-base-active rounded-none border-b border-base-border/40 last:border-b-0"
						onClick={() => onTaskClick(item)}
					>
						<CheckCircleFilled
							style={{
								color: isDone ? 'var(--base-success)' : 'var(--base-disabled)',
							}}
						/>
						<span className="min-w-0 flex-1 truncate text-left text-base-gray group-hover:text-base-primary transition-colors">
							{displayValue(item.title)}
							<span className="text-[11px] text-base-disabled font-mono pl-1">
								#{item.id}
							</span>
						</span>
						<span className="ml-auto flex shrink-0 items-center">
							{state?.name ? (
								<Tag className="m-0 text-[11px] leading-[18px]">{state.name}</Tag>
							) : null}
							<RightOutlined className="text-[10px] text-base-disabled opacity-0 group-hover:opacity-100 transition-opacity" />
						</span>
					</Button>
				);
			})}
		</div>
	);
}
