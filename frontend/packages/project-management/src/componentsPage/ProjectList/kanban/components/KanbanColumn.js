/**
 * 看板列：按 workflow_state_id 收任务；整列可放置；仅初始列可加任务。
 * 列内用 Virtuoso 虚拟列表（只渲染可视区）+ endReached 加载下一页。
 */
'use client';

import { useCallback, useMemo } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Button, Spin, Tag } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { getWorkflowStateColor, isInitialState } from '../../../../enum/workflow';
import { KANBAN_COLUMN_MIN_WIDTH_CLASS, KANBAN_TASK_CARD_DEFAULT_HEIGHT } from './constants';
import { countSubtasks } from './helpers';
import { TaskCard } from './TaskCard';

/**
 * @param {{
 *   state: object,
 *   tasks: Array,
 *   members?: Array,
 *   taskTotal?: number,
 *   hasMore?: boolean,
 *   loadingMore?: boolean,
 *   onLoadMore?: (stageId: number|string) => void,
 *   onAddTask: () => void,
 *   onTaskClick: (task: object) => void,
 * }} props
 */
export function KanbanColumn({
	state,
	tasks,
	members = [],
	taskTotal,
	hasMore = false,
	loadingMore = false,
	onLoadMore,
	onAddTask,
	onTaskClick,
}) {
	const columnTasks = useMemo(
		() => tasks.filter((task) => task.workflow_state_id === state.id),
		[tasks, state.id],
	);
	const sortableIds = useMemo(() => columnTasks.map((task) => `task-${task.id}`), [columnTasks]);
	const displayCount = taskTotal ?? columnTasks.length;
	const stateColor = getWorkflowStateColor(state.color);
	const { setNodeRef, isOver } = useDroppable({
		id: `column-${state.id}`,
		data: { stateId: state.id },
	});

	const handleEndReached = useCallback(() => {
		if (!hasMore || loadingMore || !onLoadMore) return;
		onLoadMore(state.id);
	}, [hasMore, loadingMore, onLoadMore, state.id]);

	const renderTask = useCallback(
		(_index, task) => (
			<TaskCard
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
		),
		[members, onTaskClick, state, tasks],
	);

	const renderFooter = useCallback(() => {
		if (!loadingMore) return null;
		return (
			<div className="flex items-center justify-center py-3 min-h-8">
				<Spin size="small" />
			</div>
		);
	}, [loadingMore]);

	return (
		<div
			ref={setNodeRef}
			className={`flex-1 ${KANBAN_COLUMN_MIN_WIDTH_CLASS} rounded-[10px] px-3.5 py-3 flex flex-col h-full min-h-0 ${
				isOver ? 'bg-base-active' : 'bg-[var(--base-text2)]'
			}`}
		>
			<div className="flex items-center justify-between mb-3.5 shrink-0">
				<div className="flex items-center gap-2">
					<span
						className="w-2.5 h-2.5 rounded-full shrink-0"
						style={{ background: stateColor }}
					/>
					<span className="font-semibold text-sm">{state.name}</span>
					<Tag className="ml-1 rounded-[10px] text-xs leading-[18px] px-2">
						{displayCount}
					</Tag>
				</div>
				{isInitialState(state) ? (
					<Button type="text" size="small" icon={<PlusOutlined />} onClick={onAddTask} />
				) : null}
			</div>
			{/* Virtuoso 需父级限定高度：flex-1 + min-h-0 + overflow-hidden */}
			<div className="flex-1 min-h-0 overflow-hidden px-1.5 -mx-1.5">
				<SortableContext items={sortableIds} strategy={verticalListSortingStrategy}>
					<Virtuoso
						data={columnTasks}
						style={{ height: '100%' }}
						defaultItemHeight={KANBAN_TASK_CARD_DEFAULT_HEIGHT}
						overscan={4}
						endReached={handleEndReached}
						itemContent={renderTask}
						components={{ Footer: renderFooter }}
					/>
				</SortableContext>
			</div>
		</div>
	);
}
