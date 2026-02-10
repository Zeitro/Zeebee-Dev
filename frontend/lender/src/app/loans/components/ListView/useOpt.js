'use client';
import React, { useMemo, useRef, useState } from 'react';
import { message, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import {
	assemblyAddress,
	createRequest,
	transformArrayStructure,
} from '@repo/shared-methods/utils';
import { useRouter } from 'next/navigation';
import { useGlobalStore } from '@/store/global';
import { useOptionsStore } from '@repo/shared-methods/store';
import { ReloadOutlined } from '@ant-design/icons';

export default function useOpt({ stage }) {
	const [modal, contextHolder] = Modal.useModal();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const { request } = createRequest();
	const router = useRouter();
	const { changeSidebarCollapsed } = useGlobalStore();

	const tableFilterRef = useRef();
	const { loanPendingStatus, loanPurposeOptions, loanTypeOptions } = useOptionsStore();

	const [type, setType] = useState('address');
	const [searchValue, setSearchValue] = useState('');
	const [loanPurpose, setLoanPurpose] = useState([]);
	const [loanType, setLoanType] = useState([]);
	const [status, setStatus] = useState([]);

	const filterList = useMemo(() => {
		const purposeSetting = {
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
		};
		const statusSetting = {
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
		};
		let setting = [
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
			...(stage !== 'Archived' ? [purposeSetting] : []),
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
			...(['Archived', 'Active'].includes(stage) ? [] : [statusSetting]),
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
					setStatus([]);
					setTimeout(() => {
						tableFilterRef?.current?.handleSearch?.();
					}, 300);
				},
			},
		];
		return setting;
	}, [
		searchValue,
		type,
		loanPurpose,
		loanType,
		status,
		loanPurposeOptions,
		loanTypeOptions,
		loanPendingStatus,
		stage,
	]);

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
		} catch {
			setConfirmLoading(false);
		}
	};

	const showDeleteModal = ({ loan_property, id }) => {
		modal.confirm({
			title: 'Are you sure you want to delete this loan?',
			icon: <ExclamationCircleOutlined className="!text-base-warning" />,
			content: (
				<div>
					You are about to delete the loan associated with:
					<p className="font-bold pb-8">
						{assemblyAddress(loan_property?.address, true)}
					</p>
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

	const handleSelectedRow = (record) => {
		let url;
		if (record?.source === 'Borrower') {
			url = `/loans/${record?.id}`;
		} else {
			url = `/loans/zeeBeeDetails?id=${record?.id}`;
		}
		router.push(url);
		setTimeout(() => changeSidebarCollapsed(true), 300);
	};
	return { showDeleteModal, contextHolder, handleSelectedRow, tableFilterRef, filterList };
}
