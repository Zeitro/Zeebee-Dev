/**
 * 项目任务看板：读路由 id，拉工作流列与任务；拖拽改状态，点击打开详情。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { Spin, message } from 'antd';
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
import { InviteMemberModal } from './components/InviteMemberModal';
import { KanbanColumn } from './components/KanbanColumn';
import { KanbanHeader } from './components/KanbanHeader';
import { ProjectActivityDrawer } from './components/ProjectActivityDrawer';
import { ProjectMemberListDrawer } from './components/ProjectMemberListDrawer';
import { TaskActivityDrawer, TaskDetailDrawer } from './components/TaskDetail';
import { useKanbanBoard } from './hooks/useKanbanBoard';
import { KanbanTaskDetailProvider, useKanbanTaskDetail } from './hooks/useKanbanTaskDetail';
import { useTaskNavStack } from './hooks/useTaskNavStack';

export function ProjectKanban() {
	const { id } = useParams();
	const searchParams = useSearchParams();
	const projectId = id;
	const queryTaskId = Number(searchParams.get('task')) || undefined;
	const { loading, workflow, taskDetailContextValue, ...boardProps } = useKanbanBoard(
		projectId,
		queryTaskId,
	);

	if (loading) {
		return (
			<div className="flex items-center justify-center py-24">
				<Spin size="large" />
			</div>
		);
	}

	if (!workflow?.states?.length) {
		return <EmptyWorkflow project={boardProps.project} />;
	}

	return (
		<KanbanTaskDetailProvider value={taskDetailContextValue}>
			<ProjectKanbanView projectId={projectId} workflow={workflow} {...boardProps} />
		</KanbanTaskDetailProvider>
	);
}

/**
 * @param {object} props
 */
function ProjectKanbanView({
	projectId,
	project,
	workflow,
	tasks,
	taskTotal,
	members,
	projectMembers,
	userId,
	transitionTask,
	refreshBoard,
	refreshProjectMembers,
	searchProjectMembers,
}) {
	const router = useRouter();
	const pathname = usePathname();
	const searchParams = useSearchParams();
	const { activeTask, openTask, closeTask, fetchTaskActivities } = useKanbanTaskDetail();
	const [draggedTask, setDraggedTask] = useState(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [createParent, setCreateParent] = useState(null);
	const [memberListOpen, setMemberListOpen] = useState(false);
	const [inviteOpen, setInviteOpen] = useState(false);
	const [projectActivityOpen, setProjectActivityOpen] = useState(false);
	const [taskLogOpen, setTaskLogOpen] = useState(false);
	const [taskActivities, setTaskActivities] = useState([]);
	const [messageApi, messageContextHolder] = message.useMessage();

	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));
	const states = useMemo(
		() =>
			[...(workflow?.states || [])].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0)),
		[workflow],
	);

	const currentProjectRole = projectMembers.find(
		(item) => item.user_id === userId || item.organization_user_id === userId,
	)?.role;
	// const canManageMembers = userId === project?.owner_id || currentProjectRole === 'ProjectAdmin';
	const canManageMembers = true;

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

	const { previousTask, openFromBoard, openFromDrawer, backFromDrawer, reset } = useTaskNavStack({
		tasks,
		activeTask,
		openTask,
		setTaskQuery,
	});

	const handleTaskClick = openFromBoard;

	const handleDrawerClose = useCallback(() => {
		closeTask();
		setTaskQuery();
		setTaskLogOpen(false);
		reset();
	}, [closeTask, reset, setTaskQuery]);

	const openCreate = useCallback((parent) => {
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

	const handleOpenProjectActivity = useCallback(() => {
		setProjectActivityOpen(true);
	}, []);

	const handleOpenTaskLog = useCallback(async () => {
		if (!activeTask) return;
		setTaskLogOpen(true);
		const activities = await fetchTaskActivities(activeTask.id);
		setTaskActivities(activities);
	}, [activeTask, fetchTaskActivities]);

	return (
		<div className="flex flex-col flex-1 px-8 pt-6 pb-6">
			{messageContextHolder}
			<KanbanHeader
				project={project}
				taskCount={taskTotal || tasks.length}
				canManageMembers={canManageMembers}
				onCreateTask={() => openCreate()}
				onMembersClick={() => setMemberListOpen(true)}
				onInviteClick={() => setInviteOpen(true)}
				onActivityClick={handleOpenProjectActivity}
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
								onAddTask={() => openCreate()}
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
			<CreateTaskModal
				open={createOpen}
				projectId={projectId}
				parentTask={createParent}
				messageApi={messageApi}
				onCancel={() => {
					setCreateOpen(false);
					setCreateParent(null);
				}}
				onSuccess={refreshBoard}
			/>
			<InviteMemberModal
				open={inviteOpen}
				projectId={projectId}
				onCancel={() => setInviteOpen(false)}
				onSuccess={refreshProjectMembers}
				messageApi={messageApi}
			/>
			<ProjectMemberListDrawer
				open={memberListOpen}
				projectId={projectId}
				projectMembers={projectMembers}
				canManage={canManageMembers}
				onSearchMembers={searchProjectMembers}
				onClose={() => setMemberListOpen(false)}
				onInvite={() => {
					setInviteOpen(true);
				}}
				onSuccess={refreshProjectMembers}
				messageApi={messageApi}
			/>
			<ProjectActivityDrawer
				open={projectActivityOpen}
				projectId={projectId}
				tasks={tasks}
				members={members}
				onClose={() => setProjectActivityOpen(false)}
				onTaskClick={(task) => {
					setProjectActivityOpen(false);
					handleTaskClick(task);
				}}
			/>
			<TaskDetailDrawer
				project={project}
				workflow={workflow}
				tasks={tasks}
				onClose={handleDrawerClose}
				previousTask={previousTask}
				onTaskClick={openFromDrawer}
				onBack={backFromDrawer}
				onCreateSubtask={(parent) => openCreate(parent)}
				onOpenLog={handleOpenTaskLog}
				messageApi={messageApi}
			/>
			<TaskActivityDrawer
				open={taskLogOpen}
				activities={taskActivities}
				members={members}
				onClose={() => setTaskLogOpen(false)}
			/>
		</div>
	);
}
