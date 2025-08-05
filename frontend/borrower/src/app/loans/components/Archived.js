'use client';
import React, { useEffect, useMemo, useState, useRef } from 'react';
import FilterTable from '@/components/FilterTable';
import { archivedColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
const ArchivedLoans = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	api = '/loan/page',
}) => {
	const { loanTypeOptions, lenderOptions, getLender, getOptions } = useOptionsStore();
	const { request } = createRequest();
	const tableFilterRef = useRef();
	const [type, setType] = useState('');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [lender, setLender] = useState([]);
	const [status, setStatus] = useState([]);
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
			loanTypeOptions,
			lenderOptions,
			tableFilterRef?.current,
		],
	);

	useEffect(() => {
		!loanTypeOptions?.length && getOptions('loan_type', 'loanTypeOptions');
	}, [loanTypeOptions]);

	useEffect(() => {
		!lenderOptions?.length && getLender();
	}, [lenderOptions]);

	return (
		<FilterTable
			ref={tableFilterRef}
			columns={archivedColumns}
			scroll={{ x: 2000, y: 'calc(100% - 54px)' }}
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
							stage: 'Archived',
						},
					},
				});
			}}
		/>
	);
};
export default ArchivedLoans;
