/**
 * 看板数据与写操作：拉项目/工作流/任务/成员；拖拽走 transition；不写入全局 store。
 * 任务详情逻辑见 useKanbanTaskDetail。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { message } from 'antd';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest } from '@repo/shared-methods/utils';
import { getProjectDetailApi } from '@repo/project-management/api/project';
import {
	createProjectMemberApi,
	PROJECT_MEMBER_LIST_PAGE_SIZE,
} from '@repo/project-management/api/projectMember';
import { createTaskTypeApi } from '@repo/project-management/api/taskType';
import { createTaskApi } from '@repo/project-management/api/task';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { TASK_LIST_PAGE_SIZE } from '../components/constants';
import { unwrapList } from '../components/helpers';
import { mapKanbanProjectDetail, mapKanbanTask } from '../lib/mapKanbanBoard';
import { useKanbanTaskDetailState } from './useKanbanTaskDetail';

/**
 * @param {number|string|undefined} projectId
 * @param {number|undefined} queryTaskId URL `?task=`，首次拉数时打开对应详情
 */
export function useKanbanBoard(projectId, queryTaskId, messageApi) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const userId = useUserStore((state) => state.userInfo?.organization_user_id);
	const [project, setProject] = useState(null);
	const [workflow, setWorkflow] = useState(null);
	const [tasks, setTasks] = useState([]);
	const [taskTotal, setTaskTotal] = useState(0);
	const [projectMembers, setProjectMembers] = useState([]);
	const [taskTypes, setTaskTypes] = useState([]);
	const [loading, setLoading] = useState(true);

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

	const refreshProjectMembers = useCallback(async () => {
		await fetchProjectMembers();
	}, [fetchProjectMembers]);

	const fetchBoard = useCallback(async () => {
		if (!projectId) return [];
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
			return list;
		} catch {
			message.error(t('kanban.loadFailed'));
			setProject(null);
			setWorkflow(null);
			setTasks([]);
			setTaskTotal(0);
			return [];
		} finally {
			setLoading(false);
		}
	}, [projectId, request, t, taskApi]);

	// const fetchTaskTypes = useCallback(async () => {
	// 	try {
	// 		const typeRes = await taskTypeApi.list();
	// 		const typeList = Array.isArray(typeRes?.list) ? typeRes.list : unwrapList(typeRes);
	// 		setTaskTypes(typeList);
	// 	} catch {
	// 		setTaskTypes([]);
	// 	}
	// }, [taskTypeApi]);

	const refreshBoardRef = useRef(async () => {});

	const taskDetailContextValue = useKanbanTaskDetailState({
		projectId,
		projectMembers,
		queryTaskId,
		messageApi,
		onRefreshBoard: () => refreshBoardRef.current(),
		onRefreshProjectMembers: refreshProjectMembers,
	});
	const { syncActiveTask } = taskDetailContextValue;

	const fetchBoardAndSyncTask = useCallback(async () => {
		const list = await fetchBoard();
		syncActiveTask(list || []);
		return list;
	}, [fetchBoard, syncActiveTask]);

	useEffect(() => {
		refreshBoardRef.current = fetchBoardAndSyncTask;
	}, [fetchBoardAndSyncTask]);

	useEffect(() => {
		// eslint-disable-next-line react-hooks/set-state-in-effect -- 挂载/切换项目时拉看板
		fetchBoardAndSyncTask();
		fetchProjectMembers();
		// fetchTaskTypes();
	}, [fetchBoardAndSyncTask, fetchProjectMembers]);

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

	const transitionTask = useCallback(
		async (task, targetStateId) => {
			if (!task || !targetStateId || targetStateId === task.workflow_state_id) return;
			try {
				await taskApi.transition(task.id, {
					project_id: Number(projectId),
					target_stage_id: Number(targetStateId),
				});
				messageApi?.success(t('kanban.transitionSuccess'));
				await fetchBoardAndSyncTask();
			} catch (error) {
				messageApi?.error(error?.message || t('kanban.transitionFailed'));
			}
		},
		[fetchBoardAndSyncTask, projectId, t, taskApi, messageApi],
	);

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

	return {
		project,
		workflow,
		tasks,
		taskTotal,
		members,
		projectMembers,
		memberOptions,
		// taskTypes,
		loading,
		userId,
		transitionTask,
		refreshBoard: fetchBoardAndSyncTask,
		refreshProjectMembers,
		searchProjectMembers,
		taskDetailContextValue,
	};
}
