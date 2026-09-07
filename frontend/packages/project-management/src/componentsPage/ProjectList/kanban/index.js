/**
 * 项目任务看板：读路由 id，拉工作流列与任务；拖拽改状态；URL `?task=` 打开详情壳。
 * 详情自持数据与交互；本页只同步 taskId 与 mutation 后刷新看板。
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
import CreateTaskModal from './components/CreateTaskModal';
import { EmptyWorkflow } from './components/EmptyWorkflow';
import { InviteMemberModal } from './components/InviteMemberModal';
import { KanbanColumn } from './components/KanbanColumn';
import { KanbanHeader } from './components/KanbanHeader';
import { ProjectActivityDrawer } from './components/ProjectActivityDrawer';
import { ProjectMemberListDrawer } from './components/ProjectMemberListDrawer';
import { TaskCardDragPreview } from './components/TaskCard';
import { countSubtasks } from './components/helpers';
import { TaskDetailDrawer } from './components/TaskDetail';
import { useKanbanBoard } from './hooks/useKanbanBoard';

export function ProjectKanban() {
	const { id } = useParams();
	const projectId = id;
	const [messageApi, messageContextHolder] = message.useMessage();
	const { loading, workflow, ...boardProps } = useKanbanBoard(projectId, messageApi);

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
		<>
			{messageContextHolder}
			<ProjectKanbanView
				projectId={projectId}
				workflow={workflow}
				messageApi={messageApi}
				{...boardProps}
			/>
		</>
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
	stageTaskMeta,
	members,
	projectMembers,
	userId,
	transitionTask,
	loadMoreStageTasks,
	refreshBoard,
	refreshProjectMembers,
	searchProjectMembers,
	messageApi,
}) {
	const router = useRouter();
	const pathname = usePathname();
	const searchParams = useSearchParams();
	const taskId = Number(searchParams.get('task')) || undefined;
	const seedTask = useMemo(
		() => (taskId ? tasks.find((item) => item.id === taskId) || null : null),
		[taskId, tasks],
	);

	const [draggedTask, setDraggedTask] = useState(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [createParent, setCreateParent] = useState(null);
	const [memberListOpen, setMemberListOpen] = useState(false);
	const [inviteOpen, setInviteOpen] = useState(false);
	const [projectActivityOpen, setProjectActivityOpen] = useState(false);

	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));
	const states = useMemo(
		() =>
			[...(workflow?.states || [])].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0)),
		[workflow],
	);

	const currentProjectRole = projectMembers.find(
		(item) => item.organization_user_id === userId,
	)?.role;

	// 项目管理员或项目创建者可操作邀请成员/成员查看/查看项目日志/修改工作流
	const canManage = ['ProjectAdmin', 'ProjectOwner'].includes(currentProjectRole);

	const setTaskQuery = useCallback(
		(nextTaskId) => {
			const params = new URLSearchParams(searchParams.toString());
			if (nextTaskId) params.set('task', String(nextTaskId));
			else params.delete('task');
			const qs = params.toString();
			router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false });
		},
		[pathname, router, searchParams],
	);

	const handleTaskClick = useCallback(
		(task) => {
			if (!task?.id) return;
			setTaskQuery(task.id);
		},
		[setTaskQuery],
	);

	const handleDrawerClose = useCallback(() => {
		setTaskQuery();
	}, [setTaskQuery]);

	const openCreate = useCallback((parent) => {
		setCreateParent(parent || null);
		setCreateOpen(true);
	}, []);

	const handleDragEnd = useCallback(
		async (event) => {
			setDraggedTask(null);
			const { active, over } = event;
			if (!over) return;
			const dragTaskId = Number(String(active.id).replace('task-', ''));
			const task = tasks.find((item) => item.id === dragTaskId);
			const newStateId = over.data.current?.stateId;
			if (task && newStateId) {
				await transitionTask(task, newStateId);
			}
		},
		[tasks, transitionTask],
	);

	return (
		<div className="flex flex-col flex-1 px-8 pt-6 pb-6 overflow-hidden">
			<KanbanHeader
				project={project}
				taskCount={taskTotal || tasks.length}
				canManage={canManage}
				onCreateTask={() => openCreate()}
				onMembersClick={() => setMemberListOpen(true)}
				onInviteClick={() => setInviteOpen(true)}
				onActivityClick={() => setProjectActivityOpen(true)}
			/>
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
									onAddTask={() => openCreate()}
									onTaskClick={handleTaskClick}
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
										? tasks.find((item) => item.id === draggedTask.parent_id)
												?.title
										: undefined
								}
								subtaskCount={countSubtasks(tasks, draggedTask.id)}
							/>
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
				canManage={canManage}
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
				taskId={taskId}
				projectId={projectId}
				seedTask={seedTask}
				onClose={handleDrawerClose}
				onTaskIdChange={setTaskQuery}
				onMutated={refreshBoard}
			/>
		</div>
	);
}
