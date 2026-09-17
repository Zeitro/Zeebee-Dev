/**
 * 项目详情 URL：`?view=` 切换大视图，`?task=` 打开任务详情。不改其它 query。
 */
'use client';

import { useCallback } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { PROJECT_DETAIL_VIEW, PROJECT_DETAIL_VIEWS } from '../constants';

/**
 * @returns {{
 *   view: string,
 *   setView: (view: string) => void,
 *   taskId: number|undefined,
 *   setTaskQuery: (taskId?: number) => void,
 * }}
 */
export function useProjectDetailQuery() {
	const router = useRouter();
	const pathname = usePathname();
	const searchParams = useSearchParams();
	const rawView = searchParams.get('view');
	const view = PROJECT_DETAIL_VIEWS.includes(rawView) ? rawView : PROJECT_DETAIL_VIEW.KANBAN;
	const taskId = Number(searchParams.get('task')) || undefined;

	const replaceQuery = useCallback(
		(mutate) => {
			const params = new URLSearchParams(searchParams.toString());
			mutate(params);
			const qs = params.toString();
			router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false });
		},
		[pathname, router, searchParams],
	);

	const setView = useCallback(
		(nextView) => {
			replaceQuery((params) => {
				if (!nextView || nextView === PROJECT_DETAIL_VIEW.KANBAN) params.delete('view');
				else params.set('view', nextView);
			});
		},
		[replaceQuery],
	);

	const setTaskQuery = useCallback(
		(nextTaskId) => {
			replaceQuery((params) => {
				if (nextTaskId) params.set('task', String(nextTaskId));
				else params.delete('task');
			});
		},
		[replaceQuery],
	);

	return { view, setView, taskId, setTaskQuery };
}
