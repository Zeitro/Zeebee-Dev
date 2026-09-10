/**
 * 我的任务列表页：FilterTable 负责筛选/分页，点击标题跳转看板任务详情。
 * 进入页时拉优先级 / 阶段下拉（`project_management_priority`、`project_management_stage_phase`）。
 */
'use client';

import { useEffect } from 'react';
import { FilterTable } from '@repo/shared-methods';
import { useOptionsStore } from '@repo/project-management/store/options';
import { useMyTaskList } from './hooks/useMyTaskList';

export function ProjectMyTask({ PERMISSIONS }) {
	const getOptions = useOptionsStore((state) => state.getOptions);
	useEffect(() => {
		const { priorityOptions, stageOptions } = useOptionsStore.getState();
		if (!priorityOptions.length) {
			getOptions('project_management_priority', 'priorityOptions').catch(() => {});
		}
		if (!stageOptions.length) {
			getOptions('project_management_stage_phase', 'stageOptions').catch(() => {});
		}
	}, [getOptions]);
	const { tableFilterRef, filterList, columns, getData, messageContextHolder } =
		useMyTaskList(PERMISSIONS);

	return (
		<>
			{messageContextHolder}
			<FilterTable
				filterWrapClassName="flex-1 overflow-hidden rounded-lg p-6"
				ref={tableFilterRef}
				columns={columns}
				filterList={filterList}
				hasFilterField
				getData={getData}
				// settingPagination={{ pageSize: 20 }}
				// locale={{ emptyText }}
			/>
		</>
	);
}

export default ProjectMyTask;
