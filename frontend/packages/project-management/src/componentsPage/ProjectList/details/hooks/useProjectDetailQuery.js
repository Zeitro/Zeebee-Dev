/**
 * 项目详情 URL：`?tab=` 任务/文档，`?view=` 任务视图，`?task=` 打开任务详情。不改其它 query。
 */
'use client';

import { useCallback } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import {
	PROJECT_DETAIL_TAB,
	PROJECT_DETAIL_TABS,
	PROJECT_DETAIL_VIEW,
	PROJECT_DETAIL_VIEWS,
} from '../constants';

/**
 * @returns {{
 *   tab: string,
 *   setTab: (tab: string) => void,
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
	const rawTab = searchParams.get('tab');
	const tab = PROJECT_DETAIL_TABS.includes(rawTab) ? rawTab : PROJECT_DETAIL_TAB.TASKS;
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

	const setTab = useCallback(
		(nextTab) => {
			replaceQuery((params) => {
				if (!nextTab || nextTab === PROJECT_DETAIL_TAB.TASKS) params.delete('tab');
				else params.set('tab', nextTab);
			});
		},
		[replaceQuery],
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

	return { tab, setTab, view, setView, taskId, setTaskQuery };
}
