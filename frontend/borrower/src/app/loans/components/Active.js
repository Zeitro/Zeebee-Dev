'use client';
import React, { useEffect, useMemo, useState } from 'react';
import FilterTable from '@/components/FilterTable';
import { activeColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useLoanStore } from '@/store/loan';
import { createRequest } from '@/lib/fetch';
const ActiveLoans = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	api = '/loan/page',
}) => {
	const { loanPurposeOptions, loanTypeOptions, lenderOptions, getOptions, getLender } =
		useLoanStore();
	const { request } = createRequest();
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
						label: 'Zip Code',
						value: 'zipCode',
					},
				],
			},
			{
				type: 'treeSelect',
				label: 'Loan purpose',
				popupMatchSelectWidth: 172,
				placeholder: 'Loan purpose',
				value: loanPurpose,
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
				popupMatchSelectWidth: 172,
				value: lender,
				onChange: (value) => setLender(value),
				treeData: transformArrayStructure({
					sourceArray: lenderOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
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

	useEffect(() => {
		!loanPurposeOptions?.length && getOptions('loan_purpose', 'loanPurposeOptions');
	}, [loanPurposeOptions]);

	useEffect(() => {
		console.log('iii');

		!loanTypeOptions?.length && getOptions('loan_type', 'loanTypeOptions');
	}, [loanTypeOptions]);

	useEffect(() => {
		!lenderOptions?.length && getLender();
	}, [lenderOptions]);

	return (
		<FilterTable
			columns={activeColumns}
			scroll={{ x: 3500, y: 'calc(100% - 54px)' }}
			filterList={filterList}
			hasFilterField={hasFilterField}
			pagination={hasPagination}
			filterWrapClassName={filterWrapClassName}
			getData={(params) => {
				return request(api, {
					method: 'POST',
					body: {
						...params,
						filter: {
							...(params?.filter || {}),
							stage: 'Active',
						},
					},
				});
			}}
		/>
	);
};
export default ActiveLoans;
