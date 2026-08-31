/**
 * 看板数据与写操作：拉项目/工作流/任务/成员；拖拽走 transition；不写入全局 store。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Modal, message } from 'antd';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest } from '@repo/shared-methods/utils';
import { getProjectMembersApi } from '@repo/project-management/api/project';
import { getProjectDetailApi } from '@repo/project-management/api/project';
import {
	createProjectMemberApi,
	PROJECT_MEMBER_LIST_PAGE_SIZE,
} from '@repo/project-management/api/projectMember';
import { createTaskTypeApi } from '@repo/project-management/api/taskType';
import {
	createTaskApi,
	getTaskActivitiesApi,
	getTaskApi,
	getTaskDependenciesApi,
	getTaskDependencyApi,
	getTaskDependentApi,
	getTaskDependentsApi,
	getTaskTransitionApi,
} from '@repo/project-management/api/task';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { TASK_LIST_PAGE_SIZE } from '../components/constants';
import { toIsoDateTime, unwrapList } from '../components/helpers';
import { mapKanbanProjectDetail, mapKanbanTask } from '../lib/mapKanbanBoard';
import { getMockTaskActivities, getMockTaskDependencies, getMockTaskDependents } from '../mockData';

/**
 * @param {number|string|undefined} projectId
 * @param {number|undefined} queryTaskId URL `?task=`，首次拉数时打开对应详情
 */
export function useKanbanBoard(projectId, queryTaskId) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const userId = useUserStore(
		(state) => state.userInfo?.organization_user_id ?? state.userInfo?.id,
	);
	const queryTaskIdRef = useRef(queryTaskId);
	const [project, setProject] = useState(null);
	const [workflow, setWorkflow] = useState(null);
	const [tasks, setTasks] = useState([]);
	const [taskTotal, setTaskTotal] = useState(0);
	const [projectMembers, setProjectMembers] = useState([]);
	const [taskTypes, setTaskTypes] = useState([]);
	const [loading, setLoading] = useState(true);
	const [activeTask, setActiveTask] = useState(null);

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

	const fetchBoard = useCallback(async () => {
		if (!projectId) return;
		try {
			const projectRes = await request(getProjectDetailApi(projectId), {
				method: 'GET',
				isHideTips: true,
			}).catch(() => null);

			const {
				project: mappedProject,
				workflow: mappedWorkflow,
				taskTotal: projectTaskTotal,
			} = mapKanbanProjectDetail(projectRes);

			const stages = mappedWorkflow?.states || [];
			const taskResults = stages.length
				? await Promise.all(
						stages.map((stage) =>
							taskApi
								.pageByProject({
									page: 1,
									limit: TASK_LIST_PAGE_SIZE,
									filter: {
										project_id: Number(projectId),
										stage_id: stage.id,
									},
								})
								.catch(() => ({ list: [], total: 0 })),
						),
					)
				: [];

			const list = taskResults
				.flatMap((res) => (Array.isArray(res?.list) ? res.list : []))
				.map((item) => mapKanbanTask(item))
				.filter(Boolean);

			const totalFromStages = taskResults.reduce(
				(sum, res) => sum + (Number(res?.total) || 0),
				0,
			);

			setProject(mappedProject);
			setWorkflow(mappedWorkflow);
			setTasks(list);
			setTaskTotal(projectTaskTotal || totalFromStages || list.length);
			syncActiveTask(list);
		} catch {
			message.error(t('kanban.loadFailed'));
			setProject(null);
			setWorkflow(null);
			setTasks([]);
			setTaskTotal(0);
		} finally {
			setLoading(false);
		}
	}, [projectId, request, syncActiveTask, t, taskApi]);

	const fetchProjectMembers = useCallback(async () => {
		if (!projectId) return;
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
			const projectMemberList = unwrapList(projectMemberRes)
				.map(mapProjectMember)
				.filter(Boolean);
			setProjectMembers(projectMemberList);
		} catch {
			setProjectMembers([]);
		}
	}, [projectId, projectMemberApi]);

	/** 抽屉内搜索专用，不写入 projectMembers，避免影响权限等全量数据 */
	const searchProjectMembers = useCallback(
		async (keyword) => {
			if (!projectId) return [];
			const projectMemberRes = await projectMemberApi.page({
				page: 1,
				limit: PROJECT_MEMBER_LIST_PAGE_SIZE,
				filter: {
					project_id: Number(projectId),
					keyword: keyword.trim(),
					role: '',
				},
			});
			return unwrapList(projectMemberRes).map(mapProjectMember).filter(Boolean);
		},
		[projectId, projectMemberApi],
	);

	const fetchTaskTypes = useCallback(async () => {
		try {
			const typeRes = await taskTypeApi.list();
			const typeList = Array.isArray(typeRes?.list) ? typeRes.list : unwrapList(typeRes);
			setTaskTypes(typeList);
		} catch {
			setTaskTypes([]);
		}
	}, [taskTypeApi]);

	useEffect(() => {
		queryTaskIdRef.current = queryTaskId;
	}, [queryTaskId]);

	useEffect(() => {
		// eslint-disable-next-line react-hooks/set-state-in-effect -- 挂载/切换项目时拉看板
		fetchBoard();
		fetchProjectMembers();
		fetchTaskTypes();
	}, [fetchBoard, fetchProjectMembers, fetchTaskTypes]);

	/** 看板展示用：由项目成员派生，供任务卡片/详情等按 id 查找 */
	const members = useMemo(
		() =>
			projectMembers.map((item) => ({
				id: item.user_id,
				name: item.full_name || item.name,
				email: item.email,
			})),
		[projectMembers],
	);

	const memberOptions = useMemo(
		() =>
			projectMembers.map((item) => ({
				id: item.user_id,
				name: item.full_name || item.name,
				email: item.email,
			})),
		[projectMembers],
	);

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

	const refreshProjectMembers = useCallback(async () => {
		await fetchProjectMembers();
	}, [fetchProjectMembers]);

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
		taskTotal,
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
		refreshBoard: fetchBoard,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
		refreshProjectMembers,
		searchProjectMembers,
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
