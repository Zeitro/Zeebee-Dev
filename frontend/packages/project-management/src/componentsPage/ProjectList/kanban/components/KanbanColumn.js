/**
 * 看板列：按 workflow_state_id 收任务；整列可放置；仅初始列可加任务。
 */
'use client';

import { PlusOutlined } from '@ant-design/icons';
import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Button, Tag } from 'antd';
import { getWorkflowStateColor, isInitialState } from '../../../../enum/workflow';
import { KANBAN_COLUMN_MIN_WIDTH_CLASS } from './constants';
import { countSubtasks } from './helpers';
import { TaskCard } from './TaskCard';

/**
 * @param {{
 *   state: object,
 *   tasks: Array,
 *   members?: Array,
 *   onAddTask: () => void,
 *   onTaskClick: (task: object) => void,
 * }} props
 */
export function KanbanColumn({ state, tasks, members = [], onAddTask, onTaskClick }) {
	const columnTasks = tasks.filter((task) => task.workflow_state_id === state.id);
	const stateColor = getWorkflowStateColor(state.color);
	const { setNodeRef, isOver } = useDroppable({
		id: `column-${state.id}`,
		data: { stateId: state.id },
	});

	return (
		<div
			ref={setNodeRef}
			className={`flex-1 ${KANBAN_COLUMN_MIN_WIDTH_CLASS} rounded-[10px] px-3.5 py-3 flex flex-col h-full ${
				isOver ? 'bg-base-active' : 'bg-[var(--base-text2)]'
			}`}
		>
			<div className="flex items-center justify-between mb-3.5">
				<div className="flex items-center gap-2">
					<span
						className="w-2.5 h-2.5 rounded-full shrink-0"
						style={{ background: stateColor }}
					/>
					<span className="font-semibold text-sm">{state.name}</span>
					<Tag className="ml-1 rounded-[10px] text-xs leading-[18px] px-2">
						{columnTasks.length}
					</Tag>
				</div>
				{isInitialState(state) ? (
					<Button type="text" size="small" icon={<PlusOutlined />} onClick={onAddTask} />
				) : null}
			</div>
			<div className="flex-1 overflow-y-auto min-h-20">
				<SortableContext
					items={columnTasks.map((task) => `task-${task.id}`)}
					strategy={verticalListSortingStrategy}
				>
					{columnTasks.map((task) => (
						<TaskCard
							key={task.id}
							task={task}
							members={members}
							workflowState={state}
							parentTitle={
								task.parent_id
									? tasks.find((item) => item.id === task.parent_id)?.title
									: undefined
							}
							subtaskCount={countSubtasks(tasks, task.id)}
							onClick={() => onTaskClick(task)}
						/>
					))}
				</SortableContext>
			</div>
		</div>
	);
}
