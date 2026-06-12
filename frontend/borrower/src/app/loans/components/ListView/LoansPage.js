'use client';
import React, { useEffect } from 'react';
import { FilterTable } from '@repo/shared-methods';
import useColumn from './useColumn';
import { createRequest, eventBus, getDateRange } from '@repo/shared-methods/utils';
import useOpt from './useOpt';

/** searchQuery 里除 sorter 外都进 filter */
function pickSearchQueryFilter(searchQuery) {
	if (!searchQuery) return {};
	const { sorter: _sorter, ...filter } = searchQuery;
	return filter;
}

/** 稳定序列化，避免父组件 inline 对象导致 useEffect 重复触发 */
function getSearchQueryKey(searchQuery) {
	if (!searchQuery) return '';
	const { sorter, ...filter } = searchQuery;
	const cleanFilter = Object.fromEntries(
		Object.entries(filter).filter(([, v]) => v != null && v !== ''),
	);
	if (!Object.keys(cleanFilter).length && !sorter) return '';
	return JSON.stringify({ filter: cleanFilter, sorter: sorter ?? null });
}

const LoansPage = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	method = 'POST',
	dataSource,
	isNotAddress,
	loading,
	scrollX = 'max-content',
	successCallback,
	columns,
	api,
	stage = 'Pending', // 默认是pending，可传草稿draft(因草稿展示和pending loans的一致，只有表格操作栏不同，所以需要区分)
	searchQuery, // 其他页面跳转来所需要的搜索参数
	/** 为 true 时等 searchQuery 有效后再请求（用于 searchQuery 异步就绪的场景，避免先打一遍无参请求） */
	waitForSearchQuery = false,
}) => {
	const searchQueryKey = getSearchQueryKey(searchQuery);
	const shouldManualFetch = waitForSearchQuery || Boolean(searchQueryKey);
	const { showDeleteModal, contextHolder, handleSelectedRow, filterList, tableFilterRef } =
		useOpt({ stage, searchQuery });
	const { getColumnsByStage } = useColumn({
		isNotAddress,
		successCallback: () => {
			return successCallback
				? successCallback?.()
				: tableFilterRef?.current?.handleSearch?.();
		},
		showDeleteModal,
		stage,
		searchQuery,
	});
	const { request } = createRequest();

	// 外部 / 异步 searchQuery：不走 FilterTable 自带初始化，仅在条件就绪后请求一次
	useEffect(() => {
		if (dataSource || !shouldManualFetch || !searchQueryKey) return;
		tableFilterRef?.current?.handleSearch?.();
	}, [dataSource, shouldManualFetch, searchQueryKey]);

	useEffect(() => {
		if (dataSource || stage !== 'Pending') return;
		const refresh = () => tableFilterRef?.current?.handleSearch?.();
		eventBus.on('searchPendingLoans', refresh);
		return () => eventBus.off('searchPendingLoans', refresh);
	}, [stage, dataSource]);

	return (
		<>
			<FilterTable
				key={stage}
				ref={tableFilterRef}
				columns={columns || getColumnsByStage()}
				filterName={`${stage} filters`}
				scroll={{ x: scrollX }}
				filterList={filterList}
				hasFilterField={hasFilterField}
				pagination={hasPagination}
				filterWrapClassName={filterWrapClassName}
				dataSource={dataSource}
				dataSourceLoading={loading}
				hasInitRequest={!dataSource && !shouldManualFetch}
				onRow={(row) => ({
					onClick: () => {
						handleSelectedRow(row, stage);
					},
				})}
				getData={(params) => {
					if (!api) {
						return {
							list: dataSource,
							total: dataSource?.length,
						};
					}

					const sqFilter = pickSearchQueryFilter(searchQuery);
					const paramsFilter = { ...(params?.filter || {}) };

					// 筛选栏有值用 params，否则用 searchQuery 初始值
					const dayNum = paramsFilter.dayCount ?? sqFilter.dayCount;
					const sourceType = paramsFilter.sourceType ?? sqFilter.sourceType;
					delete paramsFilter.dayCount;
					delete paramsFilter.sourceType;

					const filterData = {
						...sqFilter,
						...paramsFilter,
						stage,
						source: [],
						...(dayNum && dayNum !== 'All' ? getDateRange(dayNum) : {}),
					};

					if (sourceType === 'ZeeBee') {
						filterData.source = ['System', 'Preapproval'];
					} else if (sourceType === '3rdParty') {
						filterData.source = ['Borrower'];
					}

					// 表头排过序用 params.sorter，否则用 searchQuery.sorter
					const sorter = params?.sorter?.column
						? params.sorter
						: (searchQuery?.sorter ?? params?.sorter);

					return request(api, {
						method,
						body: {
							...params,
							sorter,
							filter: filterData,
						},
					});
				}}
			/>
			{contextHolder}
		</>
	);
};
export default LoansPage;
