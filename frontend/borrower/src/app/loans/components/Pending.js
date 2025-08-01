'use client';
import React, { useEffect, useMemo, useState } from 'react';
import FilterTable from '@/components/FilterTable';
import { pendingColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useLoanStore } from '@/store/loan';
const PendingPage = ({ hasFilterField, hasPagination, filterWrapClassName }) => {
	const { loanPendingStatus, loanPurposeOptions, loanTypeOptions } = useLoanStore();
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
						label: 'County',
						value: 'county',
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
				treeData: [
					{
						title: 'Bridge',
						value: 'bridge',
						key: 'bridge',
						isLeaf: true,
					},
					{
						title: 'Rehab',
						value: 'rehab',
						key: 'rehab',
						isLeaf: true,
					},
					{
						title: 'Ground up construction',
						value: 'groundUpConstruction',
						key: 'groundUpConstruction',
						isLeaf: true,
					},
				],
			},
			{
				type: 'treeSelect',
				label: 'Status',
				value: status,
				placeholder: 'Status',
				popupMatchSelectWidth: 172,
				onChange: (value) => setStatus(value),
				treeData: transformArrayStructure({
					sourceArray: loanPendingStatus,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
			},
		],
		[searchValue, type, loanPurpose, loanType, lender, status],
	);

	return (
		<FilterTable
			columns={pendingColumns}
			scroll={{ x: 3500, y: 'calc(100% - 54px)' }}
			dataSource={[{ status: 'quoting', id: 1 }]}
			filterList={filterList}
			hasFilterField={hasFilterField}
			pagination={hasPagination}
			filterWrapClassName={filterWrapClassName}
		/>
	);
};
export default PendingPage;
