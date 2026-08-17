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
import { InviteMemberModal } from './components/InviteMemberModal';
import { KanbanColumn } from './components/KanbanColumn';
import { KanbanHeader } from './components/KanbanHeader';
import { ProjectActivityDrawer } from './components/ProjectActivityDrawer';
import { ProjectMemberListDrawer } from './components/ProjectMemberListDrawer';
import { TaskActivityDrawer } from './components/TaskActivityDrawer';
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
		userId,
		openTask,
		closeTask,
		transitionTask,
		createTask,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
		inviteProjectMembers,
		removeProjectMember,
		updateProjectMemberRole,
		fetchProjectActivities,
		fetchTaskActivities,
		fetchTaskDependencies,
		fetchTaskDependents,
		addTaskDependency,
		addTaskDependent,
		updateTaskDependency,
		updateTaskDependent,
		deleteTaskDependency,
		deleteTaskDependent,
	} = useKanbanBoard(projectId, queryTaskId);

	const [draggedTask, setDraggedTask] = useState(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [createParent, setCreateParent] = useState(null);
	const [targetStateId, setTargetStateId] = useState(null);
	const [memberListOpen, setMemberListOpen] = useState(false);
	const [memberSearch, setMemberSearch] = useState('');
	const [inviteOpen, setInviteOpen] = useState(false);
	const [inviteLoading, setInviteLoading] = useState(false);
	const [projectActivityOpen, setProjectActivityOpen] = useState(false);
	const [projectActivities, setProjectActivities] = useState([]);
	const [taskLogOpen, setTaskLogOpen] = useState(false);
	const [taskActivities, setTaskActivities] = useState([]);

	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));
	const states = useMemo(
		() =>
			[...(workflow?.states || [])].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0)),
		[workflow],
	);

	const currentProjectRole = projectMembers.find((item) => item.user_id === userId)?.role;
	// const canManageMembers = userId === project?.owner_id || currentProjectRole === 'admin';
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
		setTaskLogOpen(false);
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

	const handleOpenProjectActivity = useCallback(async () => {
		setProjectActivityOpen(true);
		const activities = await fetchProjectActivities();
		setProjectActivities(activities);
	}, [fetchProjectActivities]);

	const handleOpenTaskLog = useCallback(async () => {
		if (!activeTask) return;
		setTaskLogOpen(true);
		const activities = await fetchTaskActivities(activeTask.id);
		setTaskActivities(activities);
	}, [activeTask, fetchTaskActivities]);

	const handleInviteMembers = useCallback(
		async (userIds) => {
			setInviteLoading(true);
			try {
				await inviteProjectMembers(userIds);
				setInviteOpen(false);
			} finally {
				setInviteLoading(false);
			}
		},
		[inviteProjectMembers],
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
				canManageMembers={canManageMembers}
				onCreateTask={() => openCreate(null)}
				onMembersClick={() => {
					setMemberSearch('');
					setMemberListOpen(true);
				}}
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
			<InviteMemberModal
				open={inviteOpen}
				members={members}
				projectMembers={projectMembers}
				loading={inviteLoading}
				onCancel={() => setInviteOpen(false)}
				onInvite={handleInviteMembers}
			/>
			<ProjectMemberListDrawer
				open={memberListOpen}
				project={project}
				members={members}
				projectMembers={projectMembers}
				canManage={canManageMembers}
				search={memberSearch}
				onSearchChange={setMemberSearch}
				onClose={() => setMemberListOpen(false)}
				onInvite={() => {
					setMemberListOpen(false);
					setInviteOpen(true);
				}}
				onRemove={removeProjectMember}
				onRoleChange={updateProjectMemberRole}
			/>
			<ProjectActivityDrawer
				open={projectActivityOpen}
				activities={projectActivities}
				tasks={tasks}
				members={members}
				onClose={() => setProjectActivityOpen(false)}
				onTaskClick={(task) => {
					setProjectActivityOpen(false);
					handleTaskClick(task);
				}}
			/>
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
				onOpenLog={handleOpenTaskLog}
				fetchTaskDependencies={fetchTaskDependencies}
				fetchTaskDependents={fetchTaskDependents}
				onAddDependency={addTaskDependency}
				onAddDependent={addTaskDependent}
				onUpdateDependency={updateTaskDependency}
				onUpdateDependent={updateTaskDependent}
				onDeleteDependency={deleteTaskDependency}
				onDeleteDependent={deleteTaskDependent}
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
