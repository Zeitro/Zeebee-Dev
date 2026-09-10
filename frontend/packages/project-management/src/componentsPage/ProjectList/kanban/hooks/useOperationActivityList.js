/**
 * 操作日志分页：POST `/operation/page`，支持项目级与任务级（额外 task_id）筛选。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { unwrapList } from '../components/helpers';

export const OPERATION_PAGE_API = '/operation/page';
export const PROJECT_ENTITY_TYPE = 'ProjectManagementProject';
export const ACTIVITY_PAGE_SIZE = 20;

/**
 * @param {*} res
 * @returns {{ list: Array, total: number, totalPage: number }}
 */
export function normalizeActivityPage(res) {
	const list = unwrapList(res);
	const total = Number(res?.total) || 0;
	const totalPage =
		Number(res?.total_page) || (total > 0 ? Math.ceil(total / ACTIVITY_PAGE_SIZE) : 0);
	return { list, total, totalPage };
}

/**
 * @param {{
 *   open: boolean,
 *   projectId?: number|string,
 *   taskId?: number|string,
 *   requireTaskId?: boolean,
 * }} options
 */
export function useOperationActivityList({ open, projectId, taskId, requireTaskId = false }) {
	const { request } = useMemo(() => createRequest(), []);
	const [activities, setActivities] = useState([]);
	const [loading, setLoading] = useState(false);
	const [loadingMore, setLoadingMore] = useState(false);
	const [hasMore, setHasMore] = useState(false);
	const pageRef = useRef(1);
	const hasMoreRef = useRef(false);
	const loadingRef = useRef(false);
	const sentinelRef = useRef(null);

	const fetchActivitiesPage = useCallback(
		async (page) => {
			if (!projectId) return { list: [], total: 0, totalPage: 0 };
			const filter = {
				entity_id: Number(projectId),
				entity_type: PROJECT_ENTITY_TYPE,
			};
			if (taskId != null && taskId !== '') {
				filter.task_id = Number(taskId);
			}
			const data = await request(OPERATION_PAGE_API, {
				method: 'POST',
				isHideTips: true,
				body: {
					page,
					limit: ACTIVITY_PAGE_SIZE,
					filter,
				},
			});
			return normalizeActivityPage(data);
		},
		[projectId, request, taskId],
	);

	const applyPageResult = useCallback((page, { list, totalPage }, append) => {
		setActivities((prev) => (append ? [...prev, ...list] : list));
		pageRef.current = page;
		const more = page < totalPage;
		hasMoreRef.current = more;
		setHasMore(more);
	}, []);

	useEffect(() => {
		if (!open || !projectId) return undefined;
		if (requireTaskId && (taskId == null || taskId === '')) return undefined;

		let cancelled = false;
		(async () => {
			if (cancelled || loadingRef.current) return;
			loadingRef.current = true;
			setLoading(true);
			pageRef.current = 1;
			hasMoreRef.current = false;
			setHasMore(false);
			try {
				const result = await fetchActivitiesPage(1);
				if (cancelled) return;
				applyPageResult(1, result, false);
			} catch {
				if (!cancelled) {
					setActivities([]);
					hasMoreRef.current = false;
					setHasMore(false);
				}
			} finally {
				if (!cancelled) {
					loadingRef.current = false;
					setLoading(false);
				}
			}
		})();

		return () => {
			cancelled = true;
		};
	}, [open, projectId, taskId, requireTaskId, fetchActivitiesPage, applyPageResult]);

	const loadMore = useCallback(async () => {
		if (!projectId || loadingRef.current || !hasMoreRef.current) return;
		loadingRef.current = true;
		setLoadingMore(true);
		const nextPage = pageRef.current + 1;
		try {
			const result = await fetchActivitiesPage(nextPage);
			applyPageResult(nextPage, result, true);
		} catch {
			hasMoreRef.current = false;
			setHasMore(false);
		} finally {
			loadingRef.current = false;
			setLoadingMore(false);
		}
	}, [applyPageResult, fetchActivitiesPage, projectId]);

	useEffect(() => {
		if (!open || !hasMore || loading || loadingMore) return undefined;

		const sentinel = sentinelRef.current;
		if (!sentinel) return undefined;

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0]?.isIntersecting) loadMore();
			},
			{ rootMargin: '80px' },
		);
		observer.observe(sentinel);
		return () => observer.disconnect();
	}, [open, hasMore, loading, loadingMore, loadMore, activities.length]);

	return {
		activities,
		loading,
		loadingMore,
		hasMore,
		sentinelRef,
	};
}
