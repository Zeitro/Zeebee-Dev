/**
 * 看板视图：列拖拽改状态。数据与顶栏由项目详情入口传入，本组件不拉项目、不挂操作弹层。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import {
	DndContext,
	DragOverlay,
	PointerSensor,
	pointerWithin,
	useSensor,
	useSensors,
} from '@dnd-kit/core';
import { KanbanColumn } from './components/KanbanColumn';
import { TaskCardDragPreview } from './components/TaskCard';
import { countSubtasks } from './components/helpers';

/**
 * @param {{
 *   workflow?: object|null,
 *   tasks?: Array,
 *   stageTaskMeta?: Record<string, object>,
 *   members?: Array,
 *   loadMoreStageTasks?: (stageId: number|string) => void,
 *   transitionTask?: (task: object, stateId: number) => Promise<void>|void,
 *   onCreateTask?: () => void,
 *   onTaskClick?: (task: object) => void,
 * }} props
 */
export function ProjectKanbanBoard({
	workflow,
	tasks = [],
	stageTaskMeta,
	members = [],
	loadMoreStageTasks,
	transitionTask,
	onCreateTask,
	onTaskClick,
}) {
	const [draggedTask, setDraggedTask] = useState(null);
	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));
	const states = useMemo(
		() =>
			[...(workflow?.states || [])].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0)),
		[workflow],
	);

	const handleDragEnd = useCallback(
		async (event) => {
			setDraggedTask(null);
			const { active, over } = event;
			if (!over) return;
			const dragTaskId = Number(String(active.id).replace('task-', ''));
			const task = tasks.find((item) => item.id === dragTaskId);
			const newStateId = over.data.current?.stateId;
			if (task && newStateId) {
				await transitionTask?.(task, newStateId);
			}
		},
		[tasks, transitionTask],
	);

	return (
		<div className="flex-1 overflow-x-auto bg-white min-h-0">
			<DndContext
				sensors={sensors}
				collisionDetection={pointerWithin}
				onDragStart={(event) => {
					const dragTaskId = Number(String(event.active.id).replace('task-', ''));
					setDraggedTask(tasks.find((item) => item.id === dragTaskId) || null);
				}}
				onDragEnd={handleDragEnd}
			>
				<div className="flex gap-5 h-full min-h-0">
					{states.map((state) => {
						const meta = stageTaskMeta?.[state.id];
						return (
							<KanbanColumn
								key={state.id}
								state={state}
								tasks={tasks}
								members={members}
								taskTotal={meta?.total}
								hasMore={Boolean(meta?.hasMore)}
								loadingMore={Boolean(meta?.loadingMore)}
								onLoadMore={loadMoreStageTasks}
								onAddTask={onCreateTask}
								onTaskClick={onTaskClick}
							/>
						);
					})}
				</div>
				<DragOverlay dropAnimation={null}>
					{draggedTask ? (
						<TaskCardDragPreview
							task={draggedTask}
							members={members}
							workflowState={states.find(
								(state) => state.id === draggedTask.workflow_state_id,
							)}
							parentTitle={
								draggedTask.parent_id
									? tasks.find((item) => item.id === draggedTask.parent_id)?.title
									: undefined
							}
							subtaskCount={countSubtasks(tasks, draggedTask.id)}
						/>
					) : null}
				</DragOverlay>
			</DndContext>
		</div>
	);
}
