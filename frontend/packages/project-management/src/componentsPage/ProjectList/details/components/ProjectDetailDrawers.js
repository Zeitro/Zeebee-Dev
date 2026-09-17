/**
 * 项目详情共享操作层：新建任务、邀请、成员、项目动态、任务详情。
 * 由顶栏与各视图触发，不持有看板拖拽状态。
 */
'use client';

import CreateTaskModal from '../../kanban/components/CreateTaskModal';
import { InviteMemberModal } from '../../kanban/components/InviteMemberModal';
import { ProjectActivityDrawer } from '../../kanban/components/ProjectActivityDrawer';
import { ProjectMemberListDrawer } from '../../kanban/components/ProjectMemberListDrawer';
import { TaskDetailDrawer } from '../../kanban/components/TaskDetail';

/**
 * @param {object} props
 */
export function ProjectDetailDrawers({
	projectId,
	messageApi,
	createOpen,
	createParent,
	onCreateCancel,
	onCreateSuccess,
	inviteOpen,
	onInviteCancel,
	memberListOpen,
	projectMembers,
	canManage,
	onSearchMembers,
	onMemberClose,
	onInvite,
	onMembersChanged,
	activityOpen,
	tasks,
	members,
	onActivityClose,
	onActivityTaskClick,
	taskId,
	seedTask,
	onTaskClose,
	onTaskIdChange,
	onTaskMutated,
}) {
	return (
		<>
			<CreateTaskModal
				open={createOpen}
				projectId={projectId}
				parentTask={createParent}
				messageApi={messageApi}
				onCancel={onCreateCancel}
				onSuccess={onCreateSuccess}
			/>
			<InviteMemberModal
				open={inviteOpen}
				projectId={projectId}
				onCancel={onInviteCancel}
				onSuccess={onMembersChanged}
				messageApi={messageApi}
			/>
			<ProjectMemberListDrawer
				open={memberListOpen}
				projectId={projectId}
				projectMembers={projectMembers}
				canManage={canManage}
				onSearchMembers={onSearchMembers}
				onClose={onMemberClose}
				onInvite={onInvite}
				onSuccess={onMembersChanged}
				messageApi={messageApi}
			/>
			<ProjectActivityDrawer
				open={activityOpen}
				projectId={projectId}
				tasks={tasks}
				members={members}
				onClose={onActivityClose}
				onTaskClick={onActivityTaskClick}
			/>
			<TaskDetailDrawer
				taskId={taskId}
				projectId={projectId}
				seedTask={seedTask}
				onClose={onTaskClose}
				onTaskIdChange={onTaskIdChange}
				onMutated={onTaskMutated}
			/>
		</>
	);
}
