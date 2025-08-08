'use client';
import React, { useMemo, useState, useRef } from 'react';
import FilterTable from '@/components/FilterTable';
import { activeColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
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
}) => {
	const { loanPurposeOptions, loanTypeOptions, lenderOptions } = useOptionsStore();
	const { request } = createRequest();
	const tableFilterRef = useRef();
	const [type, setType] = useState('');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [lender, setLender] = useState([]);
	const filterList = useMemo(
		() => [
			{
				type: 'select-input',
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
				type: 'treeSelect',
				label: 'Loan purpose',
				popupMatchSelectWidth: 172,
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
				type: 'treeSelect',
				label: 'Loan type',
				popupMatchSelectWidth: 195,
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
				type: 'treeSelect',
				label: 'Lender',
				placeholder: 'Lender',
				paramsKey: 'lender',
				popupMatchSelectWidth: 172,
				value: lender,
				onChange: (value) => setLender(value),
				treeData: transformArrayStructure({
					sourceArray: lenderOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
			},
			{
				type: 'button',
				onClick: () => {
					tableFilterRef?.current?.handleSearch?.();
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
			tableFilterRef?.current,
		],
	);

	return (
		<FilterTable
			ref={tableFilterRef}
			columns={activeColumns(isNotAddress)}
			scroll={{ x: scrollX }}
			filterList={filterList}
			hasFilterField={hasFilterField}
			pagination={hasPagination}
			filterWrapClassName={filterWrapClassName}
			dataSource={dataSource}
			dataSourceLoading={loading}
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
							stage: 'Active',
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
	);
};
export default ActiveLoans;
