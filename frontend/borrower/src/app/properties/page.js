'use client';

import React, { useCallback, useMemo, useRef, useState } from 'react';
import {
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	ProfileOutlined,
	ReloadOutlined,
} from '@ant-design/icons';
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
import TruncatedTagGrid from './components/TruncatedTagGrid';

const PROPERTY_API = '/property/page';

// const DEFAULT_SORTER = { column: 'total_amount', sort: des };

const Properties = () => {
	const router = useRouter();
	const { hasPermission, getActiveByPermission } = useAuth();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const { borrowerOptions, borrowerPropertyStatusOptions } = useGetOptions([
		{ type: 'borrower' },
		{
			key: 'borrower_property_status',
			valueName: 'borrowerPropertyStatusOptions',
		},
	]);
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [borrower, setBorrower] = useState([]);
	const [propertyStatus, setPropertyStatus] = useState([]);
	const [transactionModalOpen, setTransactionModalOpen] = useState(false);
	const [transactionModalProperty, setTransactionModalProperty] = useState(null);

	const refreshList = useCallback(() => tableFilterRef.current?.handleSearch?.(), []);
	const { handleDelete, loading: deleteLoading } = useProperty({ handleSearch: refreshList });

	const getData = useCallback(
		(params) => {
			cancelMatch(PROPERTY_API);
			return request(PROPERTY_API, {
				method: 'POST',
				body: { ...params },
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
				filterType: 'treeSelect',
				label: 'Property status',
				value: propertyStatus,
				paramsKey: 'status',
				placeholder: 'Property status',
				onChange: (value) => setPropertyStatus(value),
				treeData: transformArrayStructure({
					sourceArray: borrowerPropertyStatusOptions,
					outputKey: 'title',
				})?.map((i) => ({ ...i, key: i.value })),
				style: { width: 200 },
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
						Add property
					</BtnPopUpModal>
				),
			},
		],
		[
			searchValue,
			borrower,
			borrowerOptions,
			cancelMatch,
			hasPermission,
			propertyStatus,
			borrowerPropertyStatusOptions,
		],
	);

	const columns = useMemo(
		() => [
			{
				title: 'Property address',
				dataIndex: ['address', 'address1'],
				key: 'address',
				sorter: true,
				fixed: 'left',
				width: 200,
				ellipsis: true,
				render: (_, record) => displayValue(record?.address?.address1),
			},
			{
				title: 'City',
				dataIndex: ['address', 'city'],
				key: 'city',
				sorter: true,
				fixed: 'left',
				width: 150,
				render: (_, record) => displayValue(record?.address?.city),
			},
			{
				title: 'State',
				dataIndex: ['address', 'state'],
				key: 'state',
				width: 100,
				sorter: true,
				fixed: 'left',
				render: (_, record) => displayValue(record?.address?.state),
			},
			{
				title: 'Property status',
				dataIndex: 'status',
				key: 'status',
				width: 300,
				ellipsis: true,
				render: (v) => <TruncatedTagGrid items={v} />,
			},
			{
				title: 'Transaction status',
				dataIndex: 'transaction_status',
				key: 'transaction_status',
				width: 300,
				ellipsis: true,
				render: (v) => <TruncatedTagGrid items={v} />,
			},
			{
				title: 'Loan status',
				key: 'loan_status',
				width: 180,
				align: 'center',
				render: (_, record) => {
					return (
						<div className="flex flex-col flex-wrap items-stretch gap-1 w-max m-auto">
							{!!record?.active_loan_count && (
								<div className="rounded-[2px] px-2 text-xs border border-[var(--green-light-100)] bg-green-light text-green-medium">
									{record?.active_loan_count || 0} Active Loans
								</div>
							)}
							{!!record?.pending_loan_count && (
								<div className="rounded-[2px] px-2 text-xs border border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]">
									{record?.pending_loan_count || 0} Pending Loans
								</div>
							)}
							{!record?.pending_loan_count && !record?.active_loan_count && (
								<span>--</span>
							)}
						</div>
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
				width: 90,
				sorter: true,
				render: (v) => displayValue(v),
			},
			// {
			// 	title: 'Beds',
			// 	dataIndex: 'bedrooms',
			// 	key: 'bedrooms',
			// 	width: 72,
			// 	sorter: true,
			// 	render: (v) => displayValue(v),
			// },
			// {
			// 	title: 'Baths',
			// 	dataIndex: 'bathrooms',
			// 	key: 'bathrooms',
			// 	width: 72,
			// 	sorter: true,
			// 	render: (v) => displayValue(v),
			// },
			{
				title: 'Total active loan',
				dataIndex: 'total_amount',
				key: 'total_amount',
				width: 180,
				sorter: true,
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
				width: 150,
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
				ellipsis: true,
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
				ellipsis: true,
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
					{ permission: PERMISSIONS?.PROPERTY?.DETAILS, width: 80 },
					{ permission: PERMISSIONS?.PROPERTY?.EDIT, width: 62 },
					{ permission: PERMISSIONS?.TRANSACTION?.ADD, width: 140 },
					{ permission: PERMISSIONS?.PROPERTY?.DELETE, width: 78 },
				],
				render: (_, record) => (
					<Space size={0}>
						<PermissionSlot code={PERMISSIONS?.PROPERTY?.DETAILS}>
							<Button
								icon={<ProfileOutlined />}
								type="link"
								size="small"
								onClick={(e) => {
									e.stopPropagation();
									router.push(`/properties/${record?.id}`);
								}}
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
								icon={<EditOutlined />}
								isStopPropagation={true}
							>
								Edit
							</BtnPopUpModal>
						</PermissionSlot>
						<PermissionSlot code={PERMISSIONS?.TRANSACTION?.ADD}>
							<Button
								icon={<PlusOutlined />}
								type="link"
								size="small"
								onClick={(e) => {
									e.stopPropagation();
									setTransactionModalProperty(record);
									setTransactionModalOpen(true);
								}}
							>
								Add Transaction
							</Button>
						</PermissionSlot>
						<PermissionSlot code={PERMISSIONS?.PROPERTY?.DELETE}>
							<DeleteBtn
								icon={<DeleteOutlined />}
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
								isStopPropagation={true}
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

	const handleSelectedRow = (row) => {
		if (hasPermission(PERMISSIONS?.PROPERTY?.DETAILS)) {
			router.push(`/properties/${row.id}`);
		}
	};

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				rowKey="id"
				columns={columns}
				filterList={filterList}
				filterName="Property filters"
				getData={getData}
				scroll={{ x: 1800 }}
				filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden flex flex-col min-h-0"
				onRow={(row) => ({
					onClick: () => {
						handleSelectedRow(row);
					},
				})}
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
					router.push('/myTransactions');
				}}
			/>
		</>
	);
};

export default Properties;
