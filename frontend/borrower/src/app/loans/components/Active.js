'use client';
import React, { useMemo, useState, useRef } from 'react';
import { ReloadOutlined } from '@ant-design/icons';
import FilterTable from '@/components/FilterTable';
import { activeColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
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
	scrollX = 2600,
	stage = 'Active',
}) => {
	const { loanPurposeOptions, loanTypeOptions, lenderOptions } = useOptionsStore();
	const { request } = createRequest();
	const tableFilterRef = useRef();
	const { showDeleteModal, contextHolder, handleSelectedRow } = useOpt(tableFilterRef);
	const [type, setType] = useState('address');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [lender, setLender] = useState([]);
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
				filterType: 'treeSelect',
				label: 'Lender',
				placeholder: 'Lender',
				paramsKey: 'lender',
				value: lender,
				onChange: (value) => setLender(value),
				treeData: transformArrayStructure({
					sourceArray: lenderOptions,
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
					setLender([]);
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
			lender,
			loanPurposeOptions,
			loanTypeOptions,
			lenderOptions,
		],
	);

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				columns={activeColumns(isNotAddress, showDeleteModal)}
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
