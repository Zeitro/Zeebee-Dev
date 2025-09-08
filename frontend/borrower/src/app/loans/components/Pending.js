'use client';
import React, { useMemo, useState, useRef } from 'react';
import FilterTable from '@/components/FilterTable';
import { pendingColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
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
	scrollX = 2400,
	successCallback,
	stage = 'Pending', // 默认是pending，可传草稿draft(因草稿展示和pending loans的一致，只有表格操作栏不同，所以需要区分)
}) => {
	const tableFilterRef = useRef();
	// const [modal, contextHolder] = Modal.useModal();
	const { loanPendingStatus, loanPurposeOptions, loanTypeOptions, lenderOptions } =
		useOptionsStore();
	const { showDeleteModal, contextHolder } = useOpt(tableFilterRef);
	const { request } = createRequest();
	const [type, setType] = useState('address');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [lender, setLender] = useState([]);
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
				label: 'Lender',
				paramsKey: 'lender',
				placeholder: 'Lender',
				style: { width: 160 },
				popupMatchSelectWidth: 172,
				value: lender,
				onChange: (value) => setLender(value),
				treeData: transformArrayStructure({
					sourceArray: lenderOptions,
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
					setLender([]);
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
			lender,
			status,
			loanPurposeOptions,
			loanTypeOptions,
			loanPendingStatus,
			lenderOptions,
		],
	);

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				columns={pendingColumns(
					isNotAddress,
					successCallback || tableFilterRef?.current?.handleSearch,
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
