/**
 * 我的任务列表页：FilterTable 负责筛选/分页，点击标题跳转看板任务详情。
 */
'use client';

import { FilterTable } from '@repo/shared-methods';
import { useMyTaskList } from './hooks/useMyTaskList';

export function ProjectMyTask() {
	const { tableFilterRef, filterList, columns, getData, messageContextHolder, emptyText } =
		useMyTaskList();

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
				settingPagination={{ pageSize: 15 }}
				locale={{ emptyText }}
			/>
		</>
	);
}

export default ProjectMyTask;
