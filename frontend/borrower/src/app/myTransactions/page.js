'use client';

import React, { useCallback, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { DeleteOutlined, EditOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import { Button, message, Space, Typography } from 'antd';
import { DeleteBtn, FieldTips, FilterTable, PermissionSlot } from '@repo/shared-methods';
import { useGetOptions, useAuth } from '@repo/shared-methods/hooks';
import { createRequest, displayValue, formatDate } from '@repo/shared-methods/utils';
import NewTransactionModal from './components/NewTransactionModal';
import { getStatusTagClass } from './components/utils';
import { PERMISSIONS } from '@config/permissions';
import DocumentSvg from '@/assets/images/document.svg';

/** 与后端约定一致的分页接口；若路径不同，只改此处 */
const TRANSACTION_LIST_API = '/borrower_property_transaction/page';

function statusTag(status) {
	return (
		<span
			className={[
				'inline-flex whitespace-nowrap rounded-[2px] border px-2 text-xs leading-5',
				getStatusTagClass(status),
			].join(' ')}
		>
			{status}
		</span>
	);
}

const MyTransactionsPage = () => {
	const router = useRouter();
	const { hasPermission, getActiveByPermission } = useAuth();
	const { transactionTypeOptions, transactionStatusOptions } = useGetOptions([
		{ key: 'borrower_property_transaction_type', valueName: 'transactionTypeOptions' },
		{ key: 'borrower_property_transaction_status', valueName: 'transactionStatusOptions' },
	]);
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const tableFilterRef = useRef();
	const [keyword, setKeyword] = useState('');
	const [status, setStatus] = useState('');
	const [transactionType, setTransactionType] = useState('');
	const [newTransactionOpen, setNewTransactionOpen] = useState(false);
	const [deleteLoading, setDeleteLoading] = useState(false);
	const [messageApi, contextHolder] = message.useMessage();
	const getData = useCallback(
		(params) => {
			cancelMatch(TRANSACTION_LIST_API);
			const reqBody = {
				...params,
				filter: {
					...(params?.filter || {}),
					keyword: params?.filter?.keyword?.trim?.() || '',
					type: params?.filter?.type || '',
					status: params?.filter?.status || '',
				},
			};
			return request(TRANSACTION_LIST_API, {
				method: 'POST',
				body: reqBody,
			});
		},
		[request, cancelMatch],
	);

	const handleDelete = async (id) => {
		try {
			setDeleteLoading(true);
			await request(`/borrower_property_transaction/${id}`, { method: 'DELETE' });
			messageApi.success('Delete transaction successfully');
			tableFilterRef.current?.handleSearch?.();
		} finally {
			setDeleteLoading(false);
		}
	};

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by ID, address, city, state',
				paramsKey: 'keyword',
				value: keyword,
				onChange: ({ target: { value } }) => setKeyword(value),
				colProps: { xs: 24, sm: 24, md: 12, lg: 5, xl: 5 },
				style: { width: '100%', maxWidth: '100%' },
			},
			{
				filterType: 'select',
				placeholder: 'Transaction type',
				paramsKey: 'type',
				options: transactionTypeOptions,
				value: transactionType || undefined,
				onChange: (v) => setTransactionType(v ?? ''),
				allowClear: true,
				colProps: { xs: 24, sm: 12, md: 6, lg: 4, xl: 4 },
				style: { width: '100%', minWidth: 140 },
			},
			{
				filterType: 'select',
				placeholder: 'All status',
				paramsKey: 'status',
				options: transactionStatusOptions,
				value: status || undefined,
				onChange: (v) => setStatus(v ?? ''),
				allowClear: true,
				colProps: { xs: 24, sm: 12, md: 6, lg: 4, xl: 4 },
				style: { width: '100%', minWidth: 140 },
			},
			{
				filterType: 'button',
				text: 'Search',
				onClick: () => {
					cancelMatch(TRANSACTION_LIST_API);
					tableFilterRef.current?.handleSearch?.();
				},
				// colProps: { xs: 12, sm: 6, md: 'auto' },
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					setKeyword('');
					setStatus('');
					setTransactionType('');
					cancelMatch(TRANSACTION_LIST_API);
					tableFilterRef.current?.handleSearch?.({
						filter: { keyword: '', type: '', status: '' },
					});
				},
				colProps: {
					className: 'flex-1',
				},
			},
			{
				filterType: 'button',
				text: 'Add transaction',
				icon: <PlusOutlined />,
				onClick: () => setNewTransactionOpen(true),
				hidden: !hasPermission(PERMISSIONS?.TRANSACTION?.ADD),
			},
		],
		[
			keyword,
			status,
			transactionType,
			transactionTypeOptions,
			transactionStatusOptions,
			cancelMatch,
			hasPermission,
		],
	);

	const columns = useMemo(() => {
		const documentColumn = {
			title: 'Documents',
			dataIndex: 'id',
			key: 'id',
			width: 120,
			fixed: 'right',
			align: 'center',
			render: (v) => (
				<Button
					size="small"
					type="link"
					icon={<DocumentSvg style={{ width: 20, height: 20 }} />}
					onClick={() => router.push(`/myTransactions/${v}`)}
				/>
			),
		};
		const arr = [
			// {
			// 	title: 'ID',
			// 	dataIndex: 'code',
			// 	key: 'code',
			// 	sorter: true,
			// 	width: 140,
			// 	render: (v) => displayValue(v),
			// },
			{
				title: 'Property address',
				dataIndex: ['address', 'address1'],
				key: 'address',
				sorter: true,
				render: (_, record) => {
					const text = displayValue(record?.address?.address1);
					const isHavePermission = hasPermission(PERMISSIONS?.PROPERTY?.DETAILS);
					return isHavePermission ? (
						<Typography.Link
							onClick={(e) => {
								e.stopPropagation();
								e.preventDefault();
								window.open(`/properties/${record.borrower_property_id}`, '_blank');
							}}
						>
							<FieldTips title="Click to go to the property details page">
								{text}
							</FieldTips>
						</Typography.Link>
					) : (
						text
					);
				},
			},
			{
				title: 'City',
				dataIndex: ['address', 'city'],
				key: 'city',
				sorter: true,
				// width: 120,
				render: (_, record) => displayValue(record?.address?.city),
			},
			{
				title: 'State',
				dataIndex: ['address', 'state'],
				key: 'state',
				// width: 80,
				sorter: true,
				render: (_, record) => displayValue(record?.address?.state),
			},
			{
				title: 'Type',
				dataIndex: 'type',
				key: 'type',
				// width: 120,
				sorter: true,
				render: (_, record) => displayValue(record?.type),
			},
			{
				title: 'Status',
				dataIndex: 'status',
				key: 'status',
				// width: 120,
				sorter: true,
				render: (v) => (v ? statusTag(v) : displayValue(v)),
			},
			{
				title: 'COE date',
				dataIndex: 'coe_date',
				key: 'coe_date',
				// width: 120,
				sorter: true,
				render: (_, record) => {
					const formatted = record.coe_date ? formatDate(record.coe_date) : '';
					return displayValue(formatted);
				},
			},
			...(hasPermission(PERMISSIONS?.TRANSACTION?.DETAILS) ? [documentColumn] : []),
			getActiveByPermission({
				buttonOptions: [
					{ permission: PERMISSIONS?.TRANSACTION?.EDIT, width: 65 },
					{ permission: PERMISSIONS?.TRANSACTION?.DELETE, width: 80 },
				],
				render: (_, record) => (
					<Space size={0} onClick={(e) => e.stopPropagation()}>
						<PermissionSlot code={PERMISSIONS?.TRANSACTION?.EDIT}>
							<Button
								icon={<EditOutlined />}
								size="small"
								type="link"
								onClick={() =>
									router.push(
										`/myTransactions/${record.id}?tab=details&isEditing=true`,
									)
								}
							>
								Edit
							</Button>
						</PermissionSlot>
						<PermissionSlot code={PERMISSIONS?.TRANSACTION?.DELETE}>
							<DeleteBtn
								icon={<DeleteOutlined />}
								content={`ID: ${record?.code}`}
								onOk={() => handleDelete(record.id)}
								loading={deleteLoading}
							/>
						</PermissionSlot>
					</Space>
				),
			}),
		];
		return arr;
	}, [router, hasPermission, getActiveByPermission, deleteLoading]);

	const handleSelectedRow = (row) => {
		if (hasPermission(PERMISSIONS?.TRANSACTION?.DETAILS)) {
			router.push(`/myTransactions/${row.id}`);
		}
	};

	return (
		<>
			{contextHolder}
			<FilterTable
				ref={tableFilterRef}
				columns={columns}
				filterList={filterList}
				filterName="Transaction filters"
				getData={getData}
				scroll={{ x: 'max-content' }}
				filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden flex flex-col min-h-0"
				onRow={(row) => ({
					onClick: () => {
						handleSelectedRow(row);
					},
				})}
			/>
			<NewTransactionModal
				open={newTransactionOpen}
				onClose={() => setNewTransactionOpen(false)}
				onSaved={() => {
					tableFilterRef.current?.handleSearch?.();
				}}
			/>
		</>
	);
};

export default MyTransactionsPage;
