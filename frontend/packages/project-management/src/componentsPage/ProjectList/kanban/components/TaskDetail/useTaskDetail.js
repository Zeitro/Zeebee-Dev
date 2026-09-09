/**
 * 任务详情数据与读写：按 projectId 自拉项目/工作流/成员；按 taskId 拉详情并保存/删除。
 * 与看板壳无关；mutation 后通过 onMutated 通知宿主。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { message } from 'antd';
import { useTranslations } from 'next-intl';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest } from '@repo/shared-methods/utils';
import { getProjectDetailApi, getProjectMembersApi } from '@repo/project-management/api/project';
import {
	createProjectMemberApi,
	PROJECT_MEMBER_LIST_PAGE_SIZE,
} from '@repo/project-management/api/projectMember';
import { createTaskApi, getTaskApi } from '@repo/project-management/api/task';
import { createTaskTypeApi } from '@repo/project-management/api/taskType';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { unwrapList } from '../helpers';
import {
	buildTaskUpdateBody,
	mapKanbanProjectDetail,
	mapKanbanTaskDetail,
} from '../../lib/mapKanbanBoard';

/**
 * @param {{
 *   projectId?: number|string,
 *   enabled?: boolean,
 *   onMutated?: () => Promise<void>|void,
 * }} options
 */
export function useTaskDetail({ projectId, enabled = true, onMutated }) {
	const t = useTranslations();
	const [messageApi, messageContextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);
	const userId = useUserStore(
		(state) => state.userInfo?.organization_user_id ?? state.userInfo?.id,
	);

	const [project, setProject] = useState(null);
	const [workflow, setWorkflow] = useState(null);
	const [projectMembers, setProjectMembers] = useState([]);
	const [taskTypes, setTaskTypes] = useState([]);
	const [deleting, setDeleting] = useState(false);

	const members = useMemo(
		() =>
			projectMembers.map((item) => ({
				id: item.user_id,
				name: item.full_name || item.name,
				email: item.email,
			})),
		[projectMembers],
	);

	const canManage = useMemo(() => {
		const role = projectMembers.find((item) => item.organization_user_id === userId)?.role;
		return ['ProjectAdmin', 'ProjectOwner'].includes(role);
	}, [projectMembers, userId]);

	const refreshProjectMembers = useCallback(async () => {
		if (!projectId) {
			setProjectMembers([]);
			return;
		}
		try {
			const projectMemberRes = await projectMemberApi.page({
				page: 1,
				limit: PROJECT_MEMBER_LIST_PAGE_SIZE,
				filter: {
					project_id: Number(projectId),
					keyword: '',
					role: '',
				},
			});
			setProjectMembers(unwrapList(projectMemberRes).map(mapProjectMember).filter(Boolean));
		} catch {
			setProjectMembers([]);
		}
	}, [projectId, projectMemberApi]);

	const refreshProjectContext = useCallback(async () => {
		if (!projectId) {
			setProject(null);
			setWorkflow(null);
			return;
		}
		try {
			const projectRes = await request(getProjectDetailApi(projectId), {
				method: 'GET',
				isHideTips: true,
			}).catch(() => null);
			const mapped = mapKanbanProjectDetail(projectRes);
			setProject(mapped.project);
			setWorkflow(mapped.workflow);
		} catch {
			setProject(null);
			setWorkflow(null);
		}
	}, [projectId, request]);

	useDeferredSyncEffect({
		predicate: () => Boolean(enabled && projectId),
		deferred: () => {
			refreshProjectContext();
			refreshProjectMembers();
		},
		deps: [enabled, projectId, refreshProjectContext, refreshProjectMembers],
		isUseMicrotask: true,
	});

	const fetchTaskTypes = useCallback(async () => {
		try {
			const typeRes = await taskTypeApi.list();
			const typeList = Array.isArray(typeRes?.list) ? typeRes.list : unwrapList(typeRes);
			setTaskTypes(typeList);
		} catch {
			setTaskTypes([]);
		}
	}, [taskTypeApi]);

	const syncAssigneeParticipants = useCallback(
		async (assigneeIds) => {
			if (!projectId || !Array.isArray(assigneeIds)) return;
			const existingMemberIds = new Set(projectMembers.map((item) => item.user_id));
			let memberAdded = false;
			for (const uid of assigneeIds) {
				if (!existingMemberIds.has(uid)) {
					try {
						await request(getProjectMembersApi(projectId), {
							method: 'POST',
							isHideTips: true,
							body: {
								operator_id: userId || 0,
								user_id: uid,
								role: 'TaskExecutor',
							},
						});
						memberAdded = true;
					} catch {
						/* 无权限或已是成员 */
					}
				}
			}
			if (memberAdded) await refreshProjectMembers();
		},
		[projectId, projectMembers, refreshProjectMembers, request, userId],
	);

	const fetchTaskDetail = useCallback(
		async (taskId, projId) => {
			if (!taskId || !projId) return null;
			try {
				const data = await taskApi.detail(taskId, projId);
				return mapKanbanTaskDetail(data);
			} catch {
				return null;
			}
		},
		[taskApi],
	);

	const fetchTaskCandidates = useCallback(
		async ({ project_id, task_id, keyword = '' }) => {
			if (!project_id || !task_id) return [];
			try {
				const data = await taskApi.candidate({
					project_id: Number(project_id),
					task_id: Number(task_id),
					keyword: keyword.trim(),
				});
				return unwrapList(data);
			} catch {
				return [];
			}
		},
		[taskApi],
	);

	/**
	 * PUT 全量更新任务详情；响应用于直接刷新详情，无需再 GET。
	 * @param {object} sourceTask
	 * @param {object} [overrides]
	 * @returns {Promise<object|null>}
	 */
	const updateTaskDetail = useCallback(
		async (sourceTask, overrides = {}) => {
			if (!sourceTask?.id) return null;
			const body = buildTaskUpdateBody(sourceTask, overrides, { taskTypes });
			const data = await taskApi.update(sourceTask.id, body);
			const mapped = mapKanbanTaskDetail(data);
			if (mapped) await onMutated?.();
			return mapped;
		},
		[onMutated, taskApi, taskTypes],
	);

	const persistAssignees = useCallback(
		async (sourceTask, ids) => {
			if (!sourceTask?.id) return null;
			try {
				const updated = await updateTaskDetail(sourceTask, { assignee_ids: ids });
				await syncAssigneeParticipants(ids);
				return updated;
			} catch {
				messageApi.error(t('kanban.assigneeSaveFailed'));
				return null;
			}
		},
		[messageApi, syncAssigneeParticipants, t, updateTaskDetail],
	);

	const saveTask = useCallback(
		async (sourceTask, values) => {
			if (!sourceTask?.id || !values) return null;
			try {
				const updated = await updateTaskDetail(sourceTask, values);
				messageApi.success(t('kanban.saveSuccess'));
				return updated;
			} catch {
				messageApi.error(t('kanban.saveFailed'));
				return null;
			}
		},
		[messageApi, t, updateTaskDetail],
	);

	/**
	 * @param {object} sourceTask
	 * @param {string} title
	 */
	const saveTitle = useCallback(
		async (sourceTask, title) => {
			if (!sourceTask?.id || !title || title === sourceTask.title) return;
			try {
				await request(getTaskApi(sourceTask.id), {
					method: 'PUT',
					isHideTips: true,
					body: { title, user_id: userId || 0 },
				});
				messageApi.success(t('kanban.saveSuccess'));
				await onMutated?.();
			} catch {
				messageApi.error(t('kanban.saveFailed'));
			}
		},
		[messageApi, onMutated, request, t, userId],
	);

	/**
	 * @param {object} sourceTask
	 * @param {() => void} [onDeleted]
	 */
	const deleteTask = useCallback(
		async (sourceTask, onDeleted) => {
			if (!sourceTask?.id || deleting) return;
			const resolvedProjectId = projectId ?? sourceTask.project_id;
			if (!resolvedProjectId) return;
			setDeleting(true);
			try {
				await taskApi.delete(sourceTask.id, Number(resolvedProjectId));
				messageApi.success(t('kanban.deleteSuccess'));
				onDeleted?.();
				await onMutated?.();
			} catch (error) {
				messageApi.error(error?.message || t('kanban.deleteFailed'));
				throw error;
			} finally {
				setDeleting(false);
			}
		},
		[deleting, messageApi, onMutated, projectId, taskApi, t],
	);

	return {
		messageApi,
		messageContextHolder,
		project,
		workflow,
		members,
		projectMembers,
		canManage,
		taskTypes,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
		deleting,
		fetchTaskDetail,
		fetchTaskTypes,
		fetchTaskCandidates,
		updateTaskDetail,
	};
}
