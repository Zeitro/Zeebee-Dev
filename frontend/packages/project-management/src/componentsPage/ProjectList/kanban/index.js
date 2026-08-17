/**
 * 项目任务看板：读路由 id，拉工作流列与任务；拖拽改状态，点击打开详情。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { Spin } from 'antd';
import { useParams, usePathname, useRouter, useSearchParams } from 'next/navigation';
import {
	DndContext,
	DragOverlay,
	PointerSensor,
	pointerWithin,
	useSensor,
	useSensors,
} from '@dnd-kit/core';
import { displayValue } from '@repo/shared-methods/utils';
import CreateTaskModal from './components/CreateTaskModal';
import { EmptyWorkflow } from './components/EmptyWorkflow';
import { KanbanColumn } from './components/KanbanColumn';
import { KanbanHeader } from './components/KanbanHeader';
import { TaskDetailDrawer } from './components/TaskDetailDrawer';
import { useKanbanBoard } from './hooks/useKanbanBoard';

export function ProjectKanban() {
	const { id } = useParams();
	const router = useRouter();
	const pathname = usePathname();
	const searchParams = useSearchParams();
	const projectId = id;
	const queryTaskId = Number(searchParams.get('task')) || undefined;
	const {
		project,
		workflow,
		tasks,
		members,
		projectMembers,
		memberOptions,
		taskTypes,
		loading,
		activeTask,
		openTask,
		closeTask,
		transitionTask,
		createTask,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
	} = useKanbanBoard(projectId, queryTaskId);

	const [draggedTask, setDraggedTask] = useState(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [createParent, setCreateParent] = useState(null);
	const [targetStateId, setTargetStateId] = useState(null);

	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));
	const states = useMemo(
		() =>
			[...(workflow?.states || [])].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0)),
		[workflow],
	);

	const setTaskQuery = useCallback(
		(taskId) => {
			const params = new URLSearchParams(searchParams.toString());
			if (taskId) params.set('task', String(taskId));
			else params.delete('task');
			const qs = params.toString();
			router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false });
		},
		[pathname, router, searchParams],
	);

	const handleTaskClick = useCallback(
		(task) => {
			openTask(task);
			setTaskQuery(task.id);
		},
		[openTask, setTaskQuery],
	);

	const handleDrawerClose = useCallback(() => {
		closeTask();
		setTaskQuery();
	}, [closeTask, setTaskQuery]);

	const openCreate = useCallback((stateId, parent) => {
		setTargetStateId(stateId || null);
		setCreateParent(parent || null);
		setCreateOpen(true);
	}, []);

	const handleDragEnd = useCallback(
		async (event) => {
			setDraggedTask(null);
			const { active, over } = event;
			if (!over) return;
			const taskId = Number(String(active.id).replace('task-', ''));
			const task = tasks.find((item) => item.id === taskId);
			const newStateId = over.data.current?.stateId;
			if (task && newStateId) {
				await transitionTask(task, newStateId);
			}
		},
		[tasks, transitionTask],
	);

	if (loading) {
		return (
			<div className="flex items-center justify-center py-24">
				<Spin size="large" />
			</div>
		);
	}

	if (!workflow) {
		return <EmptyWorkflow project={project} />;
	}

	return (
		<div className="flex flex-col h-[calc(100vh-130px)] px-8 pt-6 pb-6">
			<KanbanHeader
				project={project}
				taskCount={tasks.length}
				members={members}
				projectMembers={projectMembers}
				onCreateTask={() => openCreate(null)}
			/>
			<div className="flex-1 overflow-x-auto p-4 bg-white">
				<DndContext
					sensors={sensors}
					collisionDetection={pointerWithin}
					onDragStart={(event) => {
						const taskId = Number(String(event.active.id).replace('task-', ''));
						setDraggedTask(tasks.find((item) => item.id === taskId) || null);
					}}
					onDragEnd={handleDragEnd}
				>
					<div className="flex gap-5 h-full">
						{states.map((state) => (
							<KanbanColumn
								key={state.id}
								state={state}
								tasks={tasks}
								members={members}
								onAddTask={() => openCreate(state.id)}
								onTaskClick={handleTaskClick}
							/>
						))}
					</div>
					<DragOverlay>
						{draggedTask ? (
							<div className="rounded-lg bg-white px-3.5 py-3 border border-base-border shadow-md opacity-80">
								{displayValue(draggedTask.title)}
							</div>
						) : null}
					</DragOverlay>
				</DndContext>
			</div>
			{createOpen ? (
				<CreateTaskModal
					parentTask={createParent}
					taskTypes={taskTypes}
					memberOptions={memberOptions}
					onCancel={() => {
						setCreateOpen(false);
						setCreateParent(null);
						setTargetStateId(null);
					}}
					onSubmit={async (values) => {
						await createTask(values, {
							parentId: createParent?.id,
							targetStateId,
						});
						setCreateOpen(false);
						setCreateParent(null);
						setTargetStateId(null);
					}}
				/>
			) : null}
			<TaskDetailDrawer
				open={Boolean(activeTask)}
				task={activeTask}
				project={project}
				workflow={workflow}
				tasks={tasks}
				taskTypes={taskTypes}
				memberOptions={memberOptions}
				members={members}
				onClose={handleDrawerClose}
				onTaskClick={handleTaskClick}
				onCreateSubtask={(parent) => openCreate(null, parent)}
				onTitleBlur={saveTitle}
				onPersistAssignees={persistAssignees}
				onSave={saveTask}
				onDelete={() => deleteTask(handleDrawerClose)}
			/>
		</div>
	);
}
