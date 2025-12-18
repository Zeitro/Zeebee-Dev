'use client';
import React, { useMemo, useState, useRef, useEffect } from 'react';
import { FilterTable } from '@repo/shared-methods';
import { pendingColumns } from './column';
import { transformArrayStructure, createRequest, eventBus } from '@repo/shared-methods/utils';
import { useOptionsStore } from '@repo/shared-methods/store';
import { ReloadOutlined } from '@ant-design/icons';
import useOpt from './useOpt';
const PendingPage = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	api,
	method = 'POST',
	dataSource,
	isNotAddress,
	loading,
	scrollX = 2700,
	successCallback,
	stage = 'Pending', // 默认是pending，可传草稿draft(因草稿展示和pending loans的一致，只有表格操作栏不同，所以需要区分)
}) => {
	const tableFilterRef = useRef();
	const { loanPendingStatus, loanPurposeOptions, loanTypeOptions } = useOptionsStore();
	const { showDeleteModal, contextHolder, handleSelectedRow } = useOpt(tableFilterRef);
	const { request } = createRequest();
	const [type, setType] = useState('address');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [status, setStatus] = useState([]);
	const filterList = useMemo(
		() => [
			{
				filterType: 'select-input',
				selectValue: type,
				selectParamsKey: 'keyword_type',
				paramsKey: 'keyword',
				value: searchValue,
				selectProps: {
					onChange: (e) => setType(e),
				},
				onChange: ({ target: { value } }) => setSearchValue(value),
				options: [
					{
						label: 'Address',
						value: 'address',
					},
					{
						label: 'State',
						value: 'state',
					},
					{
						label: 'City',
						value: 'city',
					},
					{
						label: 'Borrower',
						value: 'borrower',
					},
				],
			},
			{
				filterType: 'treeSelect',
				label: 'Loan purpose',
				paramsKey: 'loan_purpose',
				placeholder: 'Loan purpose',
				value: loanPurpose,
				onChange: (value) => setLoanPurpose(value),
				treeData: transformArrayStructure({
					sourceArray: loanPurposeOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
			},
			{
				filterType: 'treeSelect',
				label: 'Loan type',
				paramsKey: 'loan_type',
				style: { width: 160 },
				popupMatchSelectWidth: 205,
				value: loanType,
				placeholder: 'Loan type',
				onChange: (value) => setLoanType(value),
				treeData: transformArrayStructure({
					sourceArray: loanTypeOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
			},
			{
				filterType: 'treeSelect',
				label: 'Status',
				value: status,
				paramsKey: 'pending_status',
				placeholder: 'Status',
				style: { width: 160 },
				popupMatchSelectWidth: 172,
				onChange: (value) => setStatus(value),
				treeData: transformArrayStructure({
					sourceArray: loanPendingStatus,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
			},
			{
				filterType: 'button',
				onClick: () => {
					tableFilterRef?.current?.handleSearch?.();
				},
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					setType('address');
					setSearchValue('');
					setLoanPurpose([]);
					setLoanType([]);
					setStatus([]);
					setTimeout(() => {
						tableFilterRef?.current?.handleSearch?.();
					}, 300);
				},
			},
		],
		[
			searchValue,
			type,
			loanPurpose,
			loanType,
			status,
			loanPurposeOptions,
			loanTypeOptions,
			loanPendingStatus,
		],
	);

	// 导入3rd party后调用，更新列表
	useEffect(() => {
		const searchPendingLoans = () => tableFilterRef?.current?.handleSearch?.();
		eventBus.on('searchPendingLoans', searchPendingLoans);
		return () => eventBus.off('searchPendingLoans', searchPendingLoans);
	}, []);

	return (
		<>
			<FilterTable
				key={stage}
				ref={tableFilterRef}
				columns={pendingColumns(
					isNotAddress,
					// eslint-disable-next-line react-hooks/refs
					() =>
						successCallback
							? successCallback?.()
							: tableFilterRef?.current?.handleSearch?.(),
					showDeleteModal,
					stage,
				)}
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
export default PendingPage;
