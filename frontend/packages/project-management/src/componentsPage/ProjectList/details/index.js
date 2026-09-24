/**
 * 项目详情入口：顶栏是项目操作；任务页签右侧才放视图切换和新建任务。
 * 表格 / 时间线 / 工时 / 文档在各自组件内拉数；本页只编排顶栏、看板与共享弹层。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { Spin, message } from 'antd';
import { useParams } from 'next/navigation';
import { KanbanHeader } from '../kanban/components/KanbanHeader';
import { useKanbanBoard } from '../kanban/hooks/useKanbanBoard';
import { PROJECT_DETAIL_TAB } from './constants';
import { ProjectDetailDrawers } from './components/ProjectDetailDrawers';
import { ProjectDocumentsView } from './components/ProjectDocumentsView';
import { ProjectPageTabs } from './components/ProjectPageTabs';
import { ProjectTaskToolbar } from './components/ProjectTaskToolbar';
import { ProjectTaskViews } from './components/ProjectTaskViews';
import { useProjectDetailQuery } from './hooks/useProjectDetailQuery';

/**
 * @param {{ PERMISSIONS?: object }} props
 */
export function ProjectDetails({ PERMISSIONS }) {
	const { id: projectId } = useParams();
	const [messageApi, messageContextHolder] = message.useMessage();
	const { loading, workflow, ...board } = useKanbanBoard(projectId, messageApi);
	const { tab, setTab, view, setView, taskId, setTaskQuery } = useProjectDetailQuery();
	const seedTask = useMemo(
		() => (taskId ? board.tasks.find((item) => item.id === taskId) || null : null),
		[taskId, board.tasks],
	);

	const [createOpen, setCreateOpen] = useState(false);
	const [createParent, setCreateParent] = useState(null);
	const [memberListOpen, setMemberListOpen] = useState(false);
	const [inviteOpen, setInviteOpen] = useState(false);
	const [projectActivityOpen, setProjectActivityOpen] = useState(false);
	const [catalogRevision, setCatalogRevision] = useState(0);

	const currentProjectRole = board.projectMembers.find(
		(item) => item.organization_user_id === board.userId,
	)?.role;
	const canManage = ['ProjectAdmin', 'ProjectOwner'].includes(currentProjectRole);
	const hasWorkflow = Boolean(workflow?.states?.length);
	const isTaskTab = tab === PROJECT_DETAIL_TAB.TASKS;

	const openCreate = useCallback((parent) => {
		setCreateParent(parent || null);
		setCreateOpen(true);
	}, []);

	const bumpCatalog = useCallback(() => {
		setCatalogRevision((current) => current + 1);
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
						showCreate={false}
						onMembersClick={() => setMemberListOpen(true)}
						onInviteClick={() => setInviteOpen(true)}
						onActivityClick={() => setProjectActivityOpen(true)}
					/>
					<ProjectPageTabs
						value={tab}
						onChange={setTab}
						extra={
							isTaskTab ? (
								<ProjectTaskToolbar
									view={view}
									onViewChange={setView}
									showCreate={hasWorkflow}
									onCreateTask={() => openCreate()}
								/>
							) : null
						}
					/>
					{isTaskTab ? (
						<ProjectTaskViews
							view={view}
							projectId={projectId}
							project={board.project}
							workflow={workflow}
							tasks={board.tasks}
							stageTaskMeta={board.stageTaskMeta}
							members={board.members}
							loadMoreStageTasks={board.loadMoreStageTasks}
							transitionTask={board.transitionTask}
							catalogRevision={catalogRevision}
							onCreateTask={() => openCreate()}
							onTaskClick={handleTaskClick}
						/>
					) : (
						<ProjectDocumentsView
							key={`${projectId}-${catalogRevision}`}
							projectId={projectId}
						/>
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
						onCreateSuccess={() => {
							board.refreshBoard();
							bumpCatalog();
						}}
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
						onTaskMutated={() => {
							board.refreshBoard();
							bumpCatalog();
						}}
					/>
				</div>
			)}
		</>
	);
}
