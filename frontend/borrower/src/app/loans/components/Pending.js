'use client';
import React, { useEffect, useMemo, useState, useRef } from 'react';
import FilterTable from '@/components/FilterTable';
import { pendingColumns } from './column';
import { assemblyAddress, transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { createRequest } from '@/lib/fetch';
import { message, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';

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
}) => {
	const [modal, contextHolder] = Modal.useModal();
	const { loanPendingStatus, loanPurposeOptions, loanTypeOptions, lenderOptions } =
		useOptionsStore();
	const { request } = createRequest();
	const tableFilterRef = useRef();
	const [type, setType] = useState('');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [lender, setLender] = useState([]);
	const [status, setStatus] = useState([]);
	const [confirmLoading, setConfirmLoading] = useState(false);
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
						label: 'Borrower',
						value: 'borrower',
					},
				],
			},
			{
				type: 'treeSelect',
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
				type: 'treeSelect',
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
				type: 'treeSelect',
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
				type: 'treeSelect',
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
		],
	);

	const handleDelete = async (id) => {
		try {
			setConfirmLoading(true);
			await request(`/loan/${id}`, { method: 'DELETE' });
			tableFilterRef?.current?.handleSearch();
			setConfirmLoading(false);
			message.open({
				type: 'success',
				content: 'Delete successful!',
			});
		} catch (error) {
			setConfirmLoading(false);
		}
	};

	const showDeleteModal = ({ property, id }) => {
		modal.confirm({
			title: 'Are you sure you want to delete this loan?',
			icon: <ExclamationCircleOutlined className="!text-base-warning" />,
			content: (
				<div>
					You are about to delete the loan associated with:
					<p className="font-bold pb-8">{assemblyAddress(property?.address, true)}</p>
					Once deleted , this loan will no longer appear in your portfolio and cannot be
					recovered.
				</div>
			),
			onOk: () => handleDelete(id),
			okText: 'Delete loan',
			okButtonProps: {
				danger: true,
			},
			confirmLoading,
			centered: true,
		});
	};

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				columns={pendingColumns(
					isNotAddress,
					tableFilterRef?.current?.handleSearch,
					showDeleteModal,
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
			{contextHolder}
		</>
	);
};
export default PendingPage;
