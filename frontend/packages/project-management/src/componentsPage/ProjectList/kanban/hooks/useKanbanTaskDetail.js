/**
 * 任务详情抽屉：当前任务、详情拉取/保存/删除、依赖与活动日志 API。
 * useKanbanBoard 挂载 Provider；TaskDetail 等子组件通过 useKanbanTaskDetail 消费。
 */
'use client';

import {
	createContext,
	useCallback,
	useContext,
	useEffect,
	useMemo,
	useRef,
	useState,
} from 'react';
import { Modal, message } from 'antd';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest } from '@repo/shared-methods/utils';
import { getProjectMembersApi } from '@repo/project-management/api/project';
import { createTaskApi, getTaskActivitiesApi, getTaskApi } from '@repo/project-management/api/task';
import { unwrapList } from '../components/helpers';
import { buildTaskUpdateBody, mapKanbanTaskDetail } from '../lib/mapKanbanBoard';
import { getMockTaskActivities } from '../mockData';
import { createTaskTypeApi } from '@repo/project-management/api/taskType';

/** @type {import('react').Context<ReturnType<typeof useKanbanTaskDetailState>|null>} */
const KanbanTaskDetailContext = createContext(null);

/**
 * 看板页挂载，向子树提供任务详情状态与 API。
 * @param {{ value: ReturnType<typeof useKanbanTaskDetailState>, children: import('react').ReactNode }} props
 */
export function KanbanTaskDetailProvider({ value, children }) {
	return (
		<KanbanTaskDetailContext.Provider value={value}>
			{children}
		</KanbanTaskDetailContext.Provider>
	);
}

/** TaskDetail 等子组件消费任务详情上下文 */
export function useKanbanTaskDetail() {
	const context = useContext(KanbanTaskDetailContext);
	if (!context) {
		throw new Error('useKanbanTaskDetail must be used within KanbanTaskDetailProvider');
	}
	return context;
}

/**
 * @param {{
 *   projectId?: number|string,
 *   projectMembers?: Array,
 *   onRefreshBoard?: () => Promise<void>,
 *   onRefreshProjectMembers?: () => Promise<void>,
 *   queryTaskId?: number,
 * }} options
 */
export function useKanbanTaskDetailState({
	projectId,
	projectMembers = [],
	onRefreshBoard,
	onRefreshProjectMembers,
	queryTaskId,
}) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const userId = useUserStore(
		(state) => state.userInfo?.organization_user_id ?? state.userInfo?.id,
	);
	const queryTaskIdRef = useRef(queryTaskId);
	const [activeTask, setActiveTask] = useState(null);
	const [taskTypes, setTaskTypes] = useState([]);

	useEffect(() => {
		queryTaskIdRef.current = queryTaskId;
	}, [queryTaskId]);

	const fetchTaskTypes = useCallback(async () => {
		try {
			const typeRes = await taskTypeApi.list();
			const typeList = Array.isArray(typeRes?.list) ? typeRes.list : unwrapList(typeRes);
			setTaskTypes(typeList);
		} catch {
			setTaskTypes([]);
		}
	}, [taskTypeApi]);

	/** 看板刷新后同步当前详情任务（保留 assignments / assignee_ids） */
	const syncActiveTask = useCallback((list) => {
		setActiveTask((prev) => {
			const targetId = prev?.id || queryTaskIdRef.current;
			if (!targetId) return prev;
			const fresh = list.find((item) => item.id === targetId);
			if (!fresh) return prev;
			return {
				...fresh,
				assignments: fresh.assignments ?? prev?.assignments,
				assignee_ids: fresh.assignee_ids ?? prev?.assignee_ids,
			};
		});
	}, []);

	const openTask = useCallback((task) => setActiveTask(task), []);
	const closeTask = useCallback(() => setActiveTask(null), []);

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
			if (memberAdded) await onRefreshProjectMembers?.();
		},
		[onRefreshProjectMembers, projectId, projectMembers, request, userId],
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

	const fetchTaskActivities = useCallback(
		async (taskId) => {
			if (!taskId) return [];
			try {
				const data = await request(getTaskActivitiesApi(taskId), {
					method: 'GET',
					isHideTips: true,
				}).catch(() => getMockTaskActivities(taskId));
				return unwrapList(data);
			} catch {
				return getMockTaskActivities(taskId);
			}
		},
		[request],
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
	 * @param {object} sourceTask 当前详情
	 * @param {object} [overrides] 表单或 dependencies 等局部变更
	 * @returns {Promise<object|null>}
	 */
	const updateTaskDetail = useCallback(
		async (sourceTask, overrides = {}) => {
			if (!sourceTask?.id) return null;
			const body = buildTaskUpdateBody(sourceTask, overrides, { taskTypes });
			const data = await taskApi.update(sourceTask.id, body);
			const mapped = mapKanbanTaskDetail(data);
			if (mapped) {
				setActiveTask((prev) => (prev?.id === mapped.id ? { ...prev, ...mapped } : prev));
				await onRefreshBoard?.();
			}
			return mapped;
		},
		[onRefreshBoard, taskApi, taskTypes],
	);

	const persistAssignees = useCallback(
		async (sourceTask, ids) => {
			if (!sourceTask?.id) return null;
			try {
				const updated = await updateTaskDetail(sourceTask, { assignee_ids: ids });
				await syncAssigneeParticipants(ids);
				return updated;
			} catch {
				message.open({
					type: 'error',
					content: t('kanban.assigneeSaveFailed'),
				});
				return null;
			}
		},
		[syncAssigneeParticipants, t, updateTaskDetail],
	);

	const saveTask = useCallback(
		async (sourceTask, values) => {
			if (!sourceTask?.id || !values) return null;
			try {
				const updated = await updateTaskDetail(sourceTask, values);
				message.open({
					type: 'success',
					content: t('kanban.saveSuccess'),
				});
				return updated;
			} catch {
				message.open({
					type: 'error',
					content: t('kanban.saveFailed'),
				});
				return null;
			}
		},
		[t, updateTaskDetail],
	);

	const saveTitle = useCallback(
		async (title) => {
			if (!activeTask || !title || title === activeTask.title) return;
			try {
				await request(getTaskApi(activeTask.id), {
					method: 'PUT',
					isHideTips: true,
					body: { title, user_id: userId || 0 },
				});
				message.open({
					type: 'success',
					content: t('kanban.saveSuccess'),
				});
				await onRefreshBoard?.();
			} catch {
				message.open({
					type: 'error',
					content: t('kanban.saveFailed'),
				});
			}
		},
		[activeTask, onRefreshBoard, request, t, userId],
	);

	const deleteTask = useCallback(
		(onDeleted) => {
			if (!activeTask) return;
			const resolvedProjectId = projectId ?? activeTask.project_id;
			if (!resolvedProjectId) return;
			Modal.confirm({
				title: t('kanban.deleteConfirmTitle'),
				content: t('kanban.deleteConfirmContent', { title: activeTask.title }),
				okText: t('kanban.deleteConfirmOk'),
				okType: 'danger',
				cancelText: t('common.cancel'),
				onOk: async () => {
					await taskApi.delete(activeTask.id, {
						project_id: Number(resolvedProjectId),
					});
					message.open({
						type: 'success',
						content: t('kanban.deleteSuccess'),
					});
					setActiveTask(null);
					onDeleted?.();
					await onRefreshBoard?.();
				},
			});
		},
		[activeTask, onRefreshBoard, projectId, taskApi, t],
	);

	return {
		activeTask,
		taskTypes,
		openTask,
		closeTask,
		syncActiveTask,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
		fetchTaskDetail,
		fetchTaskTypes,
		fetchTaskActivities,
		fetchTaskCandidates,
		updateTaskDetail,
	};
}
