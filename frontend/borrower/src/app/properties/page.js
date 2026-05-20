'use client';

import React, { useCallback, useMemo, useRef, useState } from 'react';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import { Button, message, Space } from 'antd';
import { useRouter } from 'next/navigation';
import {
	BtnPopUpModal,
	DeleteBtn,
	FieldTips,
	FilterTable,
	PermissionSlot,
} from '@repo/shared-methods';
import AddModal from '@repo/shared-methods/componentsUi/SelectCard/Add/property';
import { useProperty } from '@repo/shared-methods/componentsUi/SelectCard/Item/useCommon';
import { des } from '@repo/shared-methods/enum';
import { useAuth, useGetOptions } from '@repo/shared-methods/hooks';
import {
	assemblyAddress,
	createRequest,
	displayValue,
	EMPTY_DISPLAY,
	formatAmount,
	formatDate,
	formatNumber,
	transformArrayStructure,
} from '@repo/shared-methods/utils';
import { PERMISSIONS } from '@config/permissions';
import NewTransactionModal from '@/app/myTransactions/components/NewTransactionModal';

const PROPERTY_API = '/property/page';

const DEFAULT_SORTER = { column: 'total_amount', sort: des };

function loanStatusCell(record) {
	const parts = [];
	if (record?.active_loan_count) {
		parts.push(`${record.active_loan_count} Active`);
	}
	if (record?.pending_loan_count) {
		parts.push(`${record.pending_loan_count} Pending`);
	}
	return parts.length ? parts.join(' · ') : EMPTY_DISPLAY;
}

const Properties = () => {
	const router = useRouter();
	const { hasPermission, getActiveByPermission } = useAuth();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const { borrowerOptions } = useGetOptions([{ type: 'borrower' }]);
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [borrower, setBorrower] = useState([]);
	const [transactionModalOpen, setTransactionModalOpen] = useState(false);
	const [transactionModalProperty, setTransactionModalProperty] = useState(null);

	const refreshList = useCallback(() => tableFilterRef.current?.handleSearch?.(), []);
	const { handleDelete, loading: deleteLoading } = useProperty({ handleSearch: refreshList });

	const getData = useCallback(
		(params) => {
			cancelMatch(PROPERTY_API);
			const sorter = params?.sorter?.column ? params.sorter : DEFAULT_SORTER;
			return request(PROPERTY_API, {
				method: 'POST',
				body: { ...params, sorter },
			});
		},
		[request, cancelMatch],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by any field',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
			},
			{
				filterType: 'treeSelect',
				label: 'Borrower',
				value: borrower,
				paramsKey: 'applicant',
				placeholder: 'Borrower',
				onChange: (value) => setBorrower(value),
				treeData: transformArrayStructure({
					sourceArray: borrowerOptions,
					outputKey: 'title',
				})?.map((i) => ({ ...i, key: i.value })),
			},
			{
				filterType: 'button',
				onClick: () => {
					cancelMatch(PROPERTY_API);
					tableFilterRef.current?.handleSearch?.();
				},
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					setSearchValue('');
					setBorrower([]);
					cancelMatch(PROPERTY_API);
					setTimeout(() => tableFilterRef.current?.handleSearch?.(), 300);
				},
				colProps: {
					className: 'flex-1',
				},
			},
			{
				filterType: 'button',
				hidden: !hasPermission(PERMISSIONS?.PROPERTY?.ADD),
				children: (
					<BtnPopUpModal
						getContent={(setOpen) => (
							<AddModal
								setOpen={setOpen}
								handleSearch={() => tableFilterRef.current?.handleSearch?.()}
							/>
						)}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add a property
					</BtnPopUpModal>
				),
			},
		],
		[searchValue, borrower, borrowerOptions, cancelMatch, hasPermission],
	);

	const columns = useMemo(
		() => [
			{
				title: 'Address',
				key: 'address',
				dataIndex: 'address',
				fixed: 'left',
				width: 300,
				sorter: true,
				render: (_, record) => {
					const text = displayValue(assemblyAddress(record?.address));
					return (
						<FieldTips title={text === EMPTY_DISPLAY ? undefined : text}>
							<span
								className="text-primary cursor-pointer"
								onClick={() => router.push(`/properties/${record?.id}`)}
							>
								{text}
							</span>
						</FieldTips>
					);
				},
			},
			{
				title: 'Property type',
				dataIndex: 'property_type',
				key: 'property_type',
				width: 200,
				sorter: true,
				render: (v) => displayValue(v),
			},
			{
				title: 'Living area',
				dataIndex: 'living_area_sqft',
				key: 'living_area_sqft',
				width: 160,
				sorter: true,
				render: (v) =>
					v != null && v !== '' ? displayValue(formatNumber(v)) : EMPTY_DISPLAY,
			},
			{
				title: 'Units',
				dataIndex: 'unit',
				key: 'unit',
				width: 72,
				sorter: true,
				render: (v) => displayValue(v),
			},
			{
				title: 'Beds',
				dataIndex: 'bedrooms',
				key: 'bedrooms',
				width: 72,
				sorter: true,
				render: (v) => displayValue(v),
			},
			{
				title: 'Baths',
				dataIndex: 'bathrooms',
				key: 'bathrooms',
				width: 72,
				sorter: true,
				render: (v) => displayValue(v),
			},
			{
				title: 'Total active loan',
				dataIndex: 'total_amount',
				key: 'total_amount',
				width: 180,
				sorter: true,
				defaultSortOrder: 'descend',
				render: (v) => displayValue(v, formatAmount),
			},
			{
				title: 'Outstanding loan',
				dataIndex: 'total_principal_amount',
				key: 'total_principal_amount',
				width: 180,
				sorter: true,
				render: (v) => displayValue(v, formatAmount),
			},
			{
				title: 'Avg active rate',
				dataIndex: 'avg_rate',
				key: 'avg_rate',
				width: 180,
				sorter: true,
				render: (v) => displayValue(v, (n) => `${n}%`),
			},
			{
				title: 'Monthly payment',
				dataIndex: 'total_monthly_payment_amount',
				key: 'total_monthly_payment_amount',
				width: 180,
				sorter: true,
				render: (v) => displayValue(v, formatAmount),
			},
			{
				title: 'Loan status',
				key: 'loan_status',
				width: 200,
				render: (_, record) => loanStatusCell(record),
			},
			{
				title: 'Nearest maturity',
				dataIndex: 'nearest_maturity_date',
				key: 'nearest_maturity_date',
				width: 180,
				sorter: true,
				render: (v) => displayValue(v, formatDate),
			},
			{
				title: 'Active borrower',
				dataIndex: 'applicants',
				key: 'applicants',
				width: 200,
				render: (v) => {
					const text = displayValue(v?.length ? v.join(', ') : null);
					return (
						<FieldTips title={text === EMPTY_DISPLAY ? undefined : text}>
							{text}
						</FieldTips>
					);
				},
			},
			{
				title: 'Active lender',
				dataIndex: 'lenders',
				key: 'lenders',
				width: 200,
				render: (v) => {
					const text = displayValue(v?.length ? v.join(', ') : null);
					return (
						<FieldTips title={text === EMPTY_DISPLAY ? undefined : text}>
							{text}
						</FieldTips>
					);
				},
			},
			getActiveByPermission({
				key: 'action',
				fixed: 'right',
				buttonOptions: [
					{ permission: PERMISSIONS?.PROPERTY?.DETAILS, width: 62 },
					{ permission: PERMISSIONS?.PROPERTY?.EDIT, width: 50 },
					{ width: 110 },
					{ permission: PERMISSIONS?.PROPERTY?.DELETE, width: 62 },
				],
				render: (_, record) => (
					<Space size="small" onClick={(e) => e.stopPropagation()}>
						<PermissionSlot code={PERMISSIONS?.PROPERTY?.DETAILS}>
							<Button
								type="link"
								size="small"
								onClick={() => router.push(`/properties/${record?.id}`)}
							>
								Details
							</Button>
						</PermissionSlot>
						<PermissionSlot code={PERMISSIONS?.PROPERTY?.EDIT}>
							<BtnPopUpModal
								getContent={(setOpen) => (
									<AddModal
										setOpen={setOpen}
										rowData={record}
										handleSearch={() => tableFilterRef.current?.getList?.()}
									/>
								)}
								type="link"
								size="small"
							>
								Edit
							</BtnPopUpModal>
						</PermissionSlot>
						<Button
							type="link"
							size="small"
							onClick={() => {
								setTransactionModalProperty(record);
								setTransactionModalOpen(true);
							}}
						>
							New Transaction
						</Button>
						<PermissionSlot code={PERMISSIONS?.PROPERTY?.DELETE}>
							<DeleteBtn
								type="link"
								size="small"
								title="Are you sure delete this Property?"
								content={assemblyAddress(record?.address)}
								onOk={() => handleDelete(record.id)}
								confirmLoading={deleteLoading}
								danger={!record?.loan_count}
								className={record?.loan_count ? '!text-base-disabled' : undefined}
								predicate={() => {
									if (record?.loan_count) {
										message.destroy();
										message.open({
											type: 'warning',
											content:
												'This property has an associated loan, so it cannot be deleted.',
										});
										return false;
									}
									return true;
								}}
							>
								Delete
							</DeleteBtn>
						</PermissionSlot>
					</Space>
				),
			}),
		],
		[router, handleDelete, deleteLoading, getActiveByPermission],
	);

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				rowKey="id"
				columns={columns}
				filterList={filterList}
				filterName="Property filters"
				getData={getData}
				scroll={{ x: 'max-content' }}
				filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden flex flex-col min-h-0"
			/>
			<NewTransactionModal
				open={transactionModalOpen}
				initialProperty={transactionModalProperty}
				onClose={() => {
					setTransactionModalOpen(false);
					setTransactionModalProperty(null);
				}}
				onSaved={() => {
					setTransactionModalOpen(false);
					setTransactionModalProperty(null);
				}}
			/>
		</>
	);
};

export default Properties;
