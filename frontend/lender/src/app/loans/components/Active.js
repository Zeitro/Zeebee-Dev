'use client';
import React, { useMemo, useState, useRef } from 'react';
import { ReloadOutlined } from '@ant-design/icons';
import { FilterTable } from '@repo/shared-methods';
import { activeColumns } from './column';
import { useOptionsStore } from '@/store/options';
import { transformArrayStructure, createRequest } from '@repo/shared-methods/utils';
import useOpt from './useOpt';
const ActiveLoans = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	api,
	method = 'POST',
	dataSource,
	isNotAddress,
	loading,
	scrollX = 2650,
	successCallback,
	stage = 'Active',
}) => {
	const { loanPurposeOptions, loanTypeOptions } = useOptionsStore();
	const { request } = createRequest();
	const tableFilterRef = useRef();
	const { showDeleteModal, contextHolder, handleSelectedRow } = useOpt(tableFilterRef);
	const [type, setType] = useState('address');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const filterList = useMemo(
		() => [
			{
				filterType: 'select-input',
				selectValue: type,
				value: searchValue,
				selectProps: {
					onChange: (e) => setType(e),
				},
				selectParamsKey: 'keyword_type',
				paramsKey: 'keyword',
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
				placeholder: 'Loan purpose',
				value: loanPurpose,
				paramsKey: 'loan_purpose',
				onChange: (value) => setLoanPurpose(value),
				treeData: transformArrayStructure({
					sourceArray: loanPurposeOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
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
					setLoanPurpose([]);
					setLoanType([]);
					setTimeout(() => {
						tableFilterRef?.current?.handleSearch?.();
					}, 300);
				},
			},
		],
		[searchValue, type, loanPurpose, loanType, loanPurposeOptions, loanTypeOptions],
	);

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				columns={activeColumns(
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
export default ActiveLoans;
