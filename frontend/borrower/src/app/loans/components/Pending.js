'use client';
import React, { useEffect, useMemo, useState, useRef } from 'react';
import FilterTable from '@/components/FilterTable';
import { pendingColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
const PendingPage = ({
	hasFilterField,
	hasPagination,
	filterWrapClassName,
	api,
	method = 'POST',
	dataSource,
	isNotAddress,
	loading,
	scrollX=3500
}) => {
	const {
		loanPendingStatus,
		loanPurposeOptions,
		loanTypeOptions,
		lenderOptions,
	} = useOptionsStore();
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
						label: 'Zip Code',
						value: 'zipCode',
					},
				],
			},
			{
				type: 'treeSelect',
				label: 'Loan purpose',
				paramsKey: 'loan_purpose',
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
				paramsKey: 'loan_type',
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
				paramsKey: 'lender',
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
				type: 'treeSelect',
				label: 'Status',
				value: status,
				paramsKey: 'pending_status',
				placeholder: 'Status',
				popupMatchSelectWidth: 172,
				onChange: (value) => setStatus(value),
				treeData: transformArrayStructure({
					sourceArray: loanPendingStatus,
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
			status,
			loanPurposeOptions,
			loanTypeOptions,
			loanPendingStatus,
			lenderOptions,
			tableFilterRef?.current,
		],
	);

	return (
		<FilterTable
			ref={tableFilterRef}
			columns={pendingColumns(isNotAddress)}
			scroll={{ x: scrollX, y: 'calc(100% - 54px)' }}
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
						total: dataSource?.length,
					};
				}
				const filter = hasFilterField
					? {}
					: {
							...(params?.filter || {}),
							stage: 'Pending',
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
export default PendingPage;
