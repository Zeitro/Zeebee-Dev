/**
 * 看板数据与写操作：拉项目/工作流/任务/成员；拖拽走 transition；不写入全局 store。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Modal, message } from 'antd';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest } from '@repo/shared-methods/utils';
// MOCK: 接口就绪后恢复下方 import
// import { MEMBER_API } from '@repo/project-management/api/member';
import {
	getProjectMemberApi,
	getProjectMemberRoleApi,
	getProjectMembersApi,
} from '@repo/project-management/api/project';
import {
	getTaskActivitiesApi,
	getTaskApi,
	getTaskDependenciesApi,
	getTaskDependencyApi,
	getTaskDependentApi,
	getTaskDependentsApi,
	getTaskTransitionApi,
	TASK_API,
	// TASK_TYPES_API,
} from '@repo/project-management/api/task';
// import { getProjectWorkflowApi } from '@repo/project-management/api/workflow';
import { isInitialState } from '@repo/project-management/enum/workflow';
// import { TASK_LIST_PAGE_SIZE } from '../components/constants';
import { toIsoDateTime, buildQuery, unwrapList } from '../components/helpers';
import {
	getMockKanbanBoard,
	getMockKanbanMembers,
	getMockProjectActivities,
	getMockTaskActivities,
	getMockTaskDependencies,
	getMockTaskDependents,
} from '../mockData';

/**
 * @param {number|string|undefined} projectId
 * @param {number|undefined} queryTaskId URL `?task=`，首次拉数时打开对应详情
 */
export function useKanbanBoard(projectId, queryTaskId) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const userId = useUserStore((state) => state.userInfo?.id);
	const queryTaskIdRef = useRef(queryTaskId);
	const [project, setProject] = useState(null);
	const [workflow, setWorkflow] = useState(null);
	const [tasks, setTasks] = useState([]);
	const [members, setMembers] = useState([]);
	const [projectMembers, setProjectMembers] = useState([]);
	const [taskTypes, setTaskTypes] = useState([]);
	const [loading, setLoading] = useState(true);
	const [activeTask, setActiveTask] = useState(null);

	const fetchBoard = useCallback(async () => {
		if (!projectId) return;
		try {
			// MOCK: start — 接口就绪后删除本段并恢复下方 API
			const {
				project: mockProject,
				workflow: mockWorkflow,
				tasks: mockTasks,
			} = getMockKanbanBoard(projectId);
			setWorkflow(mockWorkflow);
			setTasks(mockTasks);
			setProject(mockProject);
			setActiveTask((prev) => {
				const targetId = prev?.id || queryTaskIdRef.current;
				if (!targetId) return prev;
				const fresh = mockTasks.find((item) => item.id === targetId);
				if (!fresh) return prev;
				return {
					...fresh,
					assignments: fresh.assignments ?? prev?.assignments,
					assignee_ids: fresh.assignee_ids ?? prev?.assignee_ids,
				};
			});
			/* API:
			const [workflowRes, taskRes, projectRes] = await Promise.all([
				request(getProjectWorkflowApi(projectId), {
					method: 'GET',
					isHideTips: true,
				}).catch(() => null),
				request(
					`${TASK_API}${buildQuery({
						project_id: projectId,
						page: 1,
						page_size: TASK_LIST_PAGE_SIZE,
						user_id: userId,
					})}`,
					{ method: 'GET', isHideTips: true },
				).catch(() => []),
				request(`${getProjectApi(projectId)}${buildQuery({ user_id: userId })}`, {
					method: 'GET',
					isHideTips: true,
				}).catch(() => null),
			]);
			const list = unwrapList(taskRes);
			setWorkflow(workflowRes || null);
			setTasks(list);
			setProject(projectRes || null);
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
			*/
		} catch {
			message.error(t('kanban.loadFailed'));
			setTasks([]);
		} finally {
			setLoading(false);
		}
	}, [projectId, t]);

	const fetchMembers = useCallback(async () => {
		if (!projectId) return;
		try {
			// MOCK: start — 接口就绪后删除本段并恢复下方 API
			const {
				members: mockMembers,
				projectMembers: mockProjectMembers,
				taskTypes: mockTaskTypes,
			} = getMockKanbanMembers(projectId);
			setMembers(mockMembers);
			setProjectMembers(mockProjectMembers);
			setTaskTypes(mockTaskTypes);
			/* API:
			const [memberRes, projectMemberRes, typeRes] = await Promise.all([
				request(MEMBER_API, { method: 'GET', isHideTips: true }).catch(() => []),
				request(getProjectMembersApi(projectId), { method: 'GET', isHideTips: true }).catch(
					() => [],
				),
				request(TASK_TYPES_API, { method: 'GET', isHideTips: true }).catch(() => []),
			]);
			const allMembers = unwrapList(memberRes);
			setMembers(allMembers.filter((item) => !item.status || item.status === 'active'));
			setProjectMembers(unwrapList(projectMemberRes));
			setTaskTypes(unwrapList(typeRes));
			*/
		} catch {
			setMembers([]);
			setProjectMembers([]);
			setTaskTypes([]);
		}
	}, [projectId]);

	useEffect(() => {
		queryTaskIdRef.current = queryTaskId;
	}, [queryTaskId]);

	useEffect(() => {
		// eslint-disable-next-line react-hooks/set-state-in-effect -- 挂载/切换项目时拉看板
		fetchBoard();
		fetchMembers();
	}, [fetchBoard, fetchMembers]);

	const memberOptions = useMemo(() => {
		const viewerIds = new Set(
			projectMembers.filter((item) => item.role === 'viewer').map((item) => item.user_id),
		);
		const projectMemberIds = new Set(projectMembers.map((item) => item.user_id));
		return members.filter((item) => projectMemberIds.has(item.id) && !viewerIds.has(item.id));
	}, [members, projectMembers]);

	const openTask = useCallback((task) => setActiveTask(task), []);
	const closeTask = useCallback(() => setActiveTask(null), []);

	const transitionTask = useCallback(
		async (task, targetStateId) => {
			if (!task || !targetStateId || targetStateId === task.workflow_state_id) return;
			try {
				await request(getTaskTransitionApi(task.id), {
					method: 'POST',
					isHideTips: true,
					body: { target_state_id: targetStateId, user_id: userId || 0 },
				});
				message.success(t('kanban.transitionSuccess'));
				await fetchBoard();
			} catch (error) {
				message.error(error?.message || t('kanban.transitionFailed'));
			}
		},
		[fetchBoard, request, t, userId],
	);

	const createTask = useCallback(
		async (values, { parentId, targetStateId } = {}) => {
			const initialState = workflow?.states?.find((state) => isInitialState(state));
			try {
				await request(TASK_API, {
					method: 'POST',
					isHideTips: true,
					body: {
						project_id: Number(projectId),
						parent_id: parentId || null,
						title: values.title,
						description: values.description,
						priority: Number(values.priority) || 1,
						type: values.type || taskTypes[0]?.name,
						creator_id: userId || 0,
						assignee_id: values.assignee_id ? Number(values.assignee_id) : undefined,
						start_date: toIsoDateTime(values.start_date),
						due_date: toIsoDateTime(values.due_date),
						estimated_hours: values.estimated_hours
							? Number(values.estimated_hours)
							: undefined,
						workflow_state_id: targetStateId || initialState?.id || null,
					},
				});
				message.success(t('kanban.createSuccess'));
				await fetchBoard();
			} catch (error) {
				message.error(error?.message || t('kanban.createFailed'));
				throw error;
			}
		},
		[fetchBoard, projectId, request, t, taskTypes, userId, workflow],
	);

	const refreshProjectMembers = useCallback(async () => {
		await fetchMembers();
	}, [fetchMembers]);

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
								role: 'executor',
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

	const persistAssignees = useCallback(
		async (ids) => {
			if (!activeTask) return;
			try {
				await request(getTaskApi(activeTask.id), {
					method: 'PUT',
					isHideTips: true,
					body: { assignee_ids: ids, user_id: userId || 0 },
				});
				setActiveTask((prev) =>
					prev && prev.id === activeTask.id
						? { ...prev, assignee_ids: ids, assignee_id: ids[0] || null }
						: prev,
				);
				await syncAssigneeParticipants(ids);
			} catch {
				message.error(t('kanban.assigneeSaveFailed'));
			}
		},
		[activeTask, request, syncAssigneeParticipants, t, userId],
	);

	const inviteProjectMembers = useCallback(
		async (userIds) => {
			if (!projectId || !userIds?.length) return { successCount: 0, failCount: 0 };
			let successCount = 0;
			let failCount = 0;
			for (const uid of userIds) {
				try {
					await request(getProjectMembersApi(projectId), {
						method: 'POST',
						isHideTips: true,
						body: {
							operator_id: userId || 0,
							user_id: uid,
							role: 'member',
						},
					});
					successCount += 1;
				} catch {
					failCount += 1;
				}
			}
			if (successCount > 0) {
				message.success(t('kanban.inviteSuccess', { count: successCount }));
				await refreshProjectMembers();
			}
			if (failCount > 0) message.error(t('kanban.inviteFailed', { count: failCount }));
			return { successCount, failCount };
		},
		[projectId, refreshProjectMembers, request, t, userId],
	);

	const removeProjectMember = useCallback(
		async (memberUserId) => {
			if (!projectId) return;
			try {
				await request(
					`${getProjectMemberApi(projectId, memberUserId)}${buildQuery({
						operator_id: userId || 0,
					})}`,
					{ method: 'DELETE', isHideTips: true },
				);
				message.success(t('kanban.memberRemoved'));
				await refreshProjectMembers();
			} catch (error) {
				message.error(error?.message || t('kanban.memberRemoveFailed'));
			}
		},
		[projectId, refreshProjectMembers, request, t, userId],
	);

	const updateProjectMemberRole = useCallback(
		async (memberUserId, role) => {
			if (!projectId) return;
			try {
				await request(getProjectMemberRoleApi(projectId, memberUserId), {
					method: 'PUT',
					isHideTips: true,
					body: { user_id: userId || 0, role },
				});
				message.success(t('kanban.memberRoleUpdated'));
				await refreshProjectMembers();
			} catch (error) {
				message.error(error?.message || t('kanban.memberRoleUpdateFailed'));
			}
		},
		[projectId, refreshProjectMembers, request, t, userId],
	);

	const fetchProjectActivities = useCallback(async () => {
		// MOCK: 接口就绪后改为逐任务 request(getTaskActivitiesApi)
		return getMockProjectActivities(tasks);
		/*
		const all = [];
		for (const task of tasks) {
			try {
				const data = await request(getTaskActivitiesApi(task.id), {
					method: 'GET',
					isHideTips: true,
				});
				const items = unwrapList(data).map((act) => ({ ...act, task_id: task.id }));
				all.push(...items);
			} catch {
				// ignore per-task errors
			}
		}
		all.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
		return all;
		*/
	}, [tasks]);

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

	const fetchTaskDependencies = useCallback(
		async (taskId) => {
			if (!taskId) return [];
			try {
				const data = await request(getTaskDependenciesApi(taskId), {
					method: 'GET',
					isHideTips: true,
				}).catch(() => getMockTaskDependencies(taskId));
				return unwrapList(data);
			} catch {
				return getMockTaskDependencies(taskId);
			}
		},
		[request],
	);

	const fetchTaskDependents = useCallback(
		async (taskId) => {
			if (!taskId) return [];
			try {
				const data = await request(getTaskDependentsApi(taskId), {
					method: 'GET',
					isHideTips: true,
				}).catch(() => getMockTaskDependents(taskId));
				return unwrapList(data);
			} catch {
				return getMockTaskDependents(taskId);
			}
		},
		[request],
	);

	const addTaskDependency = useCallback(
		async (taskId, dependsOnTaskId) => {
			await request(getTaskDependenciesApi(taskId), {
				method: 'POST',
				isHideTips: true,
				body: {
					depends_on_task_id: dependsOnTaskId,
					dependency_type: 'finish_to_start',
				},
			});
			message.success(t('kanban.dependencyAdded'));
		},
		[request, t],
	);

	const addTaskDependent = useCallback(
		async (taskId, dependentTaskId) => {
			await request(getTaskDependentsApi(taskId), {
				method: 'POST',
				isHideTips: true,
				body: {
					dependent_task_id: dependentTaskId,
					dependency_type: 'finish_to_start',
				},
			});
			message.success(t('kanban.dependencyAdded'));
		},
		[request, t],
	);

	const updateTaskDependency = useCallback(
		async (taskId, dependencyId, dependencyType) => {
			await request(getTaskDependencyApi(taskId, dependencyId), {
				method: 'PUT',
				isHideTips: true,
				body: { dependency_type: dependencyType },
			});
		},
		[request],
	);

	const updateTaskDependent = useCallback(
		async (taskId, dependencyId, dependencyType) => {
			await request(getTaskDependentApi(taskId, dependencyId), {
				method: 'PUT',
				isHideTips: true,
				body: { dependency_type: dependencyType },
			});
		},
		[request],
	);

	const deleteTaskDependency = useCallback(
		async (taskId, dependencyId) => {
			await request(getTaskDependencyApi(taskId, dependencyId), {
				method: 'DELETE',
				isHideTips: true,
			});
			message.success(t('kanban.dependencyRemoved'));
		},
		[request, t],
	);

	const deleteTaskDependent = useCallback(
		async (taskId, dependencyId) => {
			await request(getTaskDependentApi(taskId, dependencyId), {
				method: 'DELETE',
				isHideTips: true,
			});
			message.success(t('kanban.dependencyRemoved'));
		},
		[request, t],
	);

	const saveTask = useCallback(
		async (values) => {
			if (!activeTask) return;
			try {
				if (
					values.workflow_state_id &&
					values.workflow_state_id !== activeTask.workflow_state_id
				) {
					await request(getTaskTransitionApi(activeTask.id), {
						method: 'POST',
						isHideTips: true,
						body: {
							target_state_id: Number(values.workflow_state_id),
							user_id: userId || 0,
						},
					});
				}
				await request(getTaskApi(activeTask.id), {
					method: 'PUT',
					isHideTips: true,
					body: {
						description: values.description,
						priority: values.priority,
						type: values.type,
						assignee_ids: values.assignee_ids || [],
						estimated_hours: values.estimated_hours || null,
						start_date: toIsoDateTime(values.start_date),
						due_date: toIsoDateTime(values.due_date),
						user_id: userId || 0,
					},
				});
				message.success(t('kanban.saveSuccess'));
				await fetchBoard();
			} catch {
				message.error(t('kanban.saveFailed'));
			}
		},
		[activeTask, fetchBoard, request, t, userId],
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
				message.success(t('kanban.saveSuccess'));
				await fetchBoard();
			} catch {
				message.error(t('kanban.saveFailed'));
			}
		},
		[activeTask, fetchBoard, request, t, userId],
	);

	const deleteTask = useCallback(
		(onDeleted) => {
			if (!activeTask) return;
			Modal.confirm({
				title: t('kanban.deleteConfirmTitle'),
				content: t('kanban.deleteConfirmContent', { title: activeTask.title }),
				okText: t('kanban.deleteConfirmOk'),
				okType: 'danger',
				cancelText: t('common.cancel'),
				onOk: async () => {
					await request(getTaskApi(activeTask.id), {
						method: 'DELETE',
						isHideTips: true,
					});
					message.success(t('kanban.deleteSuccess'));
					setActiveTask(null);
					onDeleted?.();
					await fetchBoard();
				},
			});
		},
		[activeTask, fetchBoard, request, t],
	);

	return {
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
		refreshProjectMembers,
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
	};
}
