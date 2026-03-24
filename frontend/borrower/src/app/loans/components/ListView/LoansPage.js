'use client';
import React, { useEffect } from 'react';
import { FilterTable } from '@repo/shared-methods';
import useColumn from './useColumn';
import { createRequest, eventBus } from '@repo/shared-methods/utils';
import useOpt from './useOpt';
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
}) => {
	const { showDeleteModal, contextHolder, handleSelectedRow, filterList, tableFilterRef } =
		useOpt({ stage });
	const { getColumnsByStage } = useColumn({
		isNotAddress,
		successCallback: () => {
			return successCallback
				? successCallback?.()
				: tableFilterRef?.current?.handleSearch?.();
		},
		showDeleteModal,
		stage,
	});
	const { request } = createRequest();

	// 导入3rd party后调用，更新列表
	useEffect(() => {
		if (stage === 'Pending') {
			const searchPendingLoans = () => tableFilterRef?.current?.handleSearch?.();
			eventBus.on('searchPendingLoans', searchPendingLoans);
			return () => eventBus.off('searchPendingLoans', searchPendingLoans);
		}
	}, [stage]);

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
					const filter = hasFilterField
						? {}
						: {
								...(params?.filter || {}),
								stage,
							};
					const payload = {
						method,
					};
					if (method === 'POST') {
						payload.body = {
							...params,
							filter,
						};
					}
					return request(api, payload);
				}}
			/>
			{contextHolder}
		</>
	);
};
export default LoansPage;
