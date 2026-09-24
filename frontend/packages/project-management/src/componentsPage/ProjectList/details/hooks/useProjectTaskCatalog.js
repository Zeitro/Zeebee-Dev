/**
 * 项目任务目录：不带 stage 分页拉全量，供表格 / 时间线 / 工时 / 文档复用。
 * 不改看板列数据。超页只标记 truncated，不静默丢弃提示。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { createTaskApi } from '@repo/project-management/api/task';
import { mapKanbanTask } from '../../kanban/lib/mapKanbanBoard';
import {
	PROJECT_TASK_CATALOG_MAX_PAGES,
	PROJECT_TASK_CATALOG_PAGE_SIZE,
} from '../constants';

/**
 * @param {number|string|undefined} projectId
 * @returns {{ tasks: Array, loading: boolean, truncated: boolean }}
 */
export function useProjectTaskCatalog(projectId) {
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const [tasks, setTasks] = useState([]);
	const [loading, setLoading] = useState(true);
	const [truncated, setTruncated] = useState(false);

	useEffect(() => {
		if (!projectId) return undefined;
		let cancelled = false;

		const load = async () => {
			const collected = [];
			let page = 1;
			let totalPage = 1;
			let hitCap = false;

			while (page <= totalPage && page <= PROJECT_TASK_CATALOG_MAX_PAGES) {
				const res = await taskApi.pageByProject({
					page,
					limit: PROJECT_TASK_CATALOG_PAGE_SIZE,
					filter: { project_id: Number(projectId) },
				});
				const list = Array.isArray(res?.list) ? res.list : [];
				collected.push(...list.map((item) => mapKanbanTask(item)).filter(Boolean));
				totalPage = Number(res?.total_page) || 1;
				if (page >= PROJECT_TASK_CATALOG_MAX_PAGES && page < totalPage) hitCap = true;
				page += 1;
			}

			if (cancelled) return;
			setTasks(collected);
			setTruncated(hitCap);
			setLoading(false);
		};

		load().catch(() => {
			if (cancelled) return;
			setTasks([]);
			setTruncated(false);
			setLoading(false);
		});

		return () => {
			cancelled = true;
		};
	}, [projectId, taskApi]);

	return { tasks, loading, truncated };
}
