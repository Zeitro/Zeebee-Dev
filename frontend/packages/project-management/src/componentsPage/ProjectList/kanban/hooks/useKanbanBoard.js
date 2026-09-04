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
 * 解析单列任务分页结果，供首屏与滚动加载共用。
 * @param {*} res
 * @param {number} page
 * @returns {{ list: Array, page: number, total: number, hasMore: boolean, loadingMore: boolean }}
 */
function normalizeStageTaskPage(res, page) {
	const list = Array.isArray(res?.list) ? res.list : [];
	const total = Number(res?.total) || 0;
	const totalPage =
		Number(res?.total_page) || (total > 0 ? Math.ceil(total / TASK_LIST_PAGE_SIZE) : 0);
	const hasMore = totalPage > 0 ? page < totalPage : list.length >= TASK_LIST_PAGE_SIZE;
	return {
		list,
		page,
		total: total || list.length,
		hasMore,
		loadingMore: false,
	};
}

/** @param {{ page: number, total: number, hasMore: boolean, loadingMore: boolean }} pageInfo */
function toStageMeta(pageInfo) {
	return {
		page: pageInfo.page,
		total: pageInfo.total,
		hasMore: pageInfo.hasMore,
		loadingMore: pageInfo.loadingMore,
	};
}

/**
 * 拉取单列某一页任务并映射为看板结构（首屏与滚动加载共用）。
 * @param {{ pageByProject: (body: object) => Promise<*> }} taskApi
 * @param {number|string} projectId
 * @param {number|string} stageId
 * @param {number} page
 * @returns {Promise<{ tasks: Array, meta: { page: number, total: number, hasMore: boolean, loadingMore: boolean } }>}
 */
async function fetchStageTasksPage(taskApi, projectId, stageId, page) {
	const res = await taskApi.pageByProject({
		page,
		limit: TASK_LIST_PAGE_SIZE,
		filter: {
			project_id: Number(projectId),
			stage_id: stageId,
		},
	});
	const pageInfo = normalizeStageTaskPage(res, page);
	return {
		tasks: pageInfo.list.map((item) => mapKanbanTask(item)).filter(Boolean),
		meta: toStageMeta(pageInfo),
	};
}

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
	/** @type {[Record<string|number, { page: number, total: number, hasMore: boolean, loadingMore: boolean }>, Function]} */
	const [stageTaskMeta, setStageTaskMeta] = useState({});
	const stageTaskMetaRef = useRef({});
	const loadingStagesRef = useRef(new Set());
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
			const stagePages = stages.length
				? await Promise.all(
						stages.map((stage) =>
							fetchStageTasksPage(taskApi, projectId, stage.id, 1).catch(() => ({
								tasks: [],
								meta: toStageMeta(
									normalizeStageTaskPage({ list: [], total: 0 }, 1),
								),
							})),
						),
					)
				: [];

			const nextMeta = {};
			const list = [];
			stages.forEach((stage, index) => {
				const { tasks: stageTasks, meta } = stagePages[index];
				nextMeta[stage.id] = meta;
				list.push(...stageTasks);
			});

			const totalFromStages = Object.values(nextMeta).reduce(
				(sum, meta) => sum + (Number(meta?.total) || 0),
				0,
			);

			stageTaskMetaRef.current = nextMeta;
			loadingStagesRef.current = new Set();
			setStageTaskMeta(nextMeta);
			setProject(mappedProject);
			setWorkflow(mappedWorkflow);
			setTasks(list);
			setTaskTotal(projectTaskTotal || totalFromStages || list.length);
			return list;
		} catch {
			message.error(t('kanban.loadFailed'));
			stageTaskMetaRef.current = {};
			loadingStagesRef.current = new Set();
			setStageTaskMeta({});
			setProject(null);
			setWorkflow(null);
			setTasks([]);
			setTaskTotal(0);
			return [];
		} finally {
			setLoading(false);
		}
	}, [projectId, request, t, taskApi]);

	/** 列内滚动加载下一页任务 */
	const loadMoreStageTasks = useCallback(
		async (stageId) => {
			if (!projectId || stageId == null || stageId === '') return;
			const meta = stageTaskMetaRef.current[stageId];
			if (!meta?.hasMore || loadingStagesRef.current.has(stageId)) return;

			loadingStagesRef.current.add(stageId);
			setStageTaskMeta((prev) => ({
				...prev,
				[stageId]: { ...prev[stageId], ...meta, loadingMore: true },
			}));

			const nextPage = (meta.page || 1) + 1;
			try {
				const { tasks: mapped, meta: nextMeta } = await fetchStageTasksPage(
					taskApi,
					projectId,
					stageId,
					nextPage,
				);

				setTasks((prev) => {
					const existingIds = new Set(prev.map((item) => item.id));
					const appended = mapped.filter((item) => !existingIds.has(item.id));
					return appended.length ? [...prev, ...appended] : prev;
				});

				stageTaskMetaRef.current[stageId] = nextMeta;
				setStageTaskMeta((prev) => ({
					...prev,
					[stageId]: nextMeta,
				}));
			} catch {
				const failedMeta = { ...meta, hasMore: false, loadingMore: false };
				stageTaskMetaRef.current[stageId] = failedMeta;
				setStageTaskMeta((prev) => ({
					...prev,
					[stageId]: failedMeta,
				}));
			} finally {
				loadingStagesRef.current.delete(stageId);
			}
		},
		[projectId, taskApi],
	);

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

			const prevStateId = task.workflow_state_id;
			const nextStateId = Number(targetStateId);
			const targetState = workflow?.states?.find((state) => state.id === nextStateId);

			// 乐观更新：拖拽松手时立刻换列，避免 transform 清除后弹回原列
			setTasks((prev) =>
				prev.map((item) =>
					item.id === task.id
						? {
								...item,
								workflow_state_id: nextStateId,
								...(targetState ? { workflow_state: targetState } : {}),
							}
						: item,
				),
			);

			try {
				await taskApi.transition(task.id, {
					project_id: Number(projectId),
					target_stage_id: nextStateId,
				});
				messageApi?.success(t('kanban.transitionSuccess'));
				fetchBoardAndSyncTask();
			} catch (error) {
				setTasks((prev) =>
					prev.map((item) =>
						item.id === task.id
							? {
									...item,
									workflow_state_id: prevStateId,
									workflow_state: task.workflow_state,
								}
							: item,
					),
				);
				messageApi?.error(error?.message || t('kanban.transitionFailed'));
			}
		},
		[fetchBoardAndSyncTask, projectId, t, taskApi, messageApi, workflow],
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
		stageTaskMeta,
		members,
		projectMembers,
		memberOptions,
		// taskTypes,
		loading,
		userId,
		transitionTask,
		loadMoreStageTasks,
		refreshBoard: fetchBoardAndSyncTask,
		refreshProjectMembers,
		searchProjectMembers,
		taskDetailContextValue,
	};
}
