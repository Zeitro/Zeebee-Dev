/**
 * 项目详情入口：顶栏操作与看板/列表/时间线/工时切换。
 * 任务抽屉与成员类弹层由本页统一挂载；看板只渲染列。
 * 项目与任务数据暂走 useKanbanBoard，供顶栏与看板共用。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { Spin, message } from 'antd';
import { useParams } from 'next/navigation';
import { EmptyWorkflow } from '../kanban/components/EmptyWorkflow';
import { KanbanHeader } from '../kanban/components/KanbanHeader';
import { useKanbanBoard } from '../kanban/hooks/useKanbanBoard';
import { ProjectKanbanBoard } from '../kanban';
import { PROJECT_DETAIL_VIEW } from './constants';
import { ProjectDetailDrawers } from './components/ProjectDetailDrawers';
import { ProjectViewPlaceholder } from './components/ProjectViewPlaceholder';
import { ProjectViewSwitch } from './components/ProjectViewSwitch';
import { useProjectDetailQuery } from './hooks/useProjectDetailQuery';

/**
 * @param {{ PERMISSIONS?: object }} props
 */
export function ProjectDetails({ PERMISSIONS }) {
	const { id: projectId } = useParams();
	const [messageApi, messageContextHolder] = message.useMessage();
	const { loading, workflow, ...board } = useKanbanBoard(projectId, messageApi);
	const { view, setView, taskId, setTaskQuery } = useProjectDetailQuery();
	const seedTask = useMemo(
		() => (taskId ? board.tasks.find((item) => item.id === taskId) || null : null),
		[taskId, board.tasks],
	);

	const [createOpen, setCreateOpen] = useState(false);
	const [createParent, setCreateParent] = useState(null);
	const [memberListOpen, setMemberListOpen] = useState(false);
	const [inviteOpen, setInviteOpen] = useState(false);
	const [projectActivityOpen, setProjectActivityOpen] = useState(false);

	const currentProjectRole = board.projectMembers.find(
		(item) => item.organization_user_id === board.userId,
	)?.role;
	const canManage = ['ProjectAdmin', 'ProjectOwner'].includes(currentProjectRole);
	const hasWorkflow = Boolean(workflow?.states?.length);

	const openCreate = useCallback((parent) => {
		setCreateParent(parent || null);
		setCreateOpen(true);
	}, []);

	const handleTaskClick = useCallback(
		(task) => {
			if (!task?.id) return;
			setTaskQuery(task.id);
		},
		[setTaskQuery],
	);

	return (
		<>
			{messageContextHolder}
			{loading ? (
				<div className="flex flex-1 items-center justify-center py-24">
					<Spin size="large" />
				</div>
			) : (
				<div className="flex flex-col flex-1 px-8 pt-6 pb-6 overflow-hidden">
					<KanbanHeader
						project={board.project}
						taskCount={board.taskTotal || board.tasks.length}
						canManage={canManage}
						PERMISSIONS={PERMISSIONS}
						showCreate={hasWorkflow}
						extra={<ProjectViewSwitch value={view} onChange={setView} />}
						onCreateTask={() => openCreate()}
						onMembersClick={() => setMemberListOpen(true)}
						onInviteClick={() => setInviteOpen(true)}
						onActivityClick={() => setProjectActivityOpen(true)}
					/>
					{view === PROJECT_DETAIL_VIEW.KANBAN ? (
						hasWorkflow ? (
							<ProjectKanbanBoard
								workflow={workflow}
								tasks={board.tasks}
								stageTaskMeta={board.stageTaskMeta}
								members={board.members}
								loadMoreStageTasks={board.loadMoreStageTasks}
								transitionTask={board.transitionTask}
								onCreateTask={() => openCreate()}
								onTaskClick={handleTaskClick}
							/>
						) : (
							<EmptyWorkflow project={board.project} />
						)
					) : (
						<ProjectViewPlaceholder view={view} />
					)}
					<ProjectDetailDrawers
						projectId={projectId}
						messageApi={messageApi}
						createOpen={createOpen}
						createParent={createParent}
						onCreateCancel={() => {
							setCreateOpen(false);
							setCreateParent(null);
						}}
						onCreateSuccess={board.refreshBoard}
						inviteOpen={inviteOpen}
						onInviteCancel={() => setInviteOpen(false)}
						memberListOpen={memberListOpen}
						projectMembers={board.projectMembers}
						canManage={canManage}
						onSearchMembers={board.searchProjectMembers}
						onMemberClose={() => setMemberListOpen(false)}
						onInvite={() => setInviteOpen(true)}
						onMembersChanged={board.refreshProjectMembers}
						activityOpen={projectActivityOpen}
						tasks={board.tasks}
						members={board.members}
						onActivityClose={() => setProjectActivityOpen(false)}
						onActivityTaskClick={(task) => {
							setProjectActivityOpen(false);
							handleTaskClick(task);
						}}
						taskId={taskId}
						seedTask={seedTask}
						onTaskClose={() => setTaskQuery()}
						onTaskIdChange={setTaskQuery}
						onTaskMutated={board.refreshBoard}
					/>
				</div>
			)}
		</>
	);
}
