'use client';
import React, { useMemo, useState, useRef } from 'react';
import FilterTable from '@/components/FilterTable';
import { archivedColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
import { ReloadOutlined } from '@ant-design/icons';
import useOpt from './useOpt';
const ArchivedLoans = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	api,
	method = 'POST',
	dataSource,
	isNotAddress,
	loading,
	scrollX = 1700,
	successCallback,
	stage = 'Archived',
}) => {
	const { loanTypeOptions, lenderOptions } = useOptionsStore();
	const { request } = createRequest();
	const tableFilterRef = useRef();
	const { showDeleteModal, contextHolder, handleSelectedRow } = useOpt(tableFilterRef);
	const [type, setType] = useState('address');
	const [searchValue, setSearchValue] = useState('');
	const [loanType, setLoanType] = useState([]);
	const [lender, setLender] = useState([]);
	const filterList = useMemo(
		() => [
			{
				filterType: 'select-input',
				selectValue: type,
				value: searchValue,
				selectParamsKey: 'keyword_type',
				paramsKey: 'keyword',
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
				label: 'Loan type',
				popupMatchSelectWidth: 205,
				value: loanType,
				paramsKey: 'loan_type',
				placeholder: 'Loan type',
				onChange: (value) => setLoanType(value),
				treeData: transformArrayStructure({
					sourceArray: loanTypeOptions,
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
					setLoanType([]);
					setLender([]);
					setTimeout(() => {
						tableFilterRef?.current?.handleSearch?.();
					}, 300);
				},
			},
		],
		[searchValue, type, loanType, lender, loanTypeOptions, lenderOptions],
	);

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				columns={archivedColumns(
					isNotAddress,
					successCallback || tableFilterRef?.current?.handleSearch,
					showDeleteModal,
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
							total: dataSource.length,
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
export default ArchivedLoans;
