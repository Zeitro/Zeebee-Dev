'use client';

import React, { useCallback, useMemo, useRef, useState } from 'react';
import { FolderOpenOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import { Button, message, Space, Tag, Typography } from 'antd';
import { FilterTable } from '@repo/shared-methods';
import { createRequest, formatDate } from '@repo/shared-methods/utils';
import NewTransactionModal from './components/NewTransactionModal';

/** 与后端约定一致的分页接口；若路径不同，只改此处 */
const TRANSACTION_LIST_API = '/transaction/page';

const STATUS_OPTIONS = [
	{ label: 'Pending', value: 'Pending' },
	{ label: 'Cancelled', value: 'Cancelled' },
	{ label: 'Closed', value: 'Closed' },
];

function statusTag(status) {
	const color = status === 'Pending' ? 'gold' : status === 'Cancelled' ? 'red' : 'success';
	return (
		<Tag bordered={false} color={color}>
			{status}
		</Tag>
	);
}

const MyTransactionsPage = () => {
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const tableFilterRef = useRef();
	const [keyword, setKeyword] = useState('');
	const [status, setStatus] = useState('');
	const [newTransactionOpen, setNewTransactionOpen] = useState(false);

	const getData = useCallback(
		(params) => {
			cancelMatch(TRANSACTION_LIST_API);
			return request(TRANSACTION_LIST_API, {
				method: 'POST',
				body: params,
			});
		},
		[request, cancelMatch],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by property, city, state, type or status...',
				paramsKey: 'keyword',
				value: keyword,
				onChange: ({ target: { value } }) => setKeyword(value),
				colProps: { xs: 24, sm: 24, md: 12, lg: 10, xl: 10 },
				style: { width: '100%', maxWidth: '100%' },
			},
			{
				filterType: 'select',
				placeholder: 'All status',
				paramsKey: 'status',
				options: STATUS_OPTIONS,
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
				colProps: { xs: 12, sm: 6, md: 'auto' },
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					setKeyword('');
					setStatus('');
					cancelMatch(TRANSACTION_LIST_API);
					setTimeout(() => tableFilterRef.current?.handleSearch?.(), 300);
				},
				colProps: { xs: 12, sm: 6, md: 'auto' },
			},
			{
				filterType: 'custom',
				colProps: {
					xs: 24,
					sm: 24,
					md: 24,
					lg: 24,
					xl: 24,
					flex: '1 1 200px',
					style: { display: 'flex', justifyContent: 'flex-end' },
				},
				children: (
					<div className="flex justify-end w-full">
						<Button
							type="primary"
							size="large"
							icon={<PlusOutlined />}
							onClick={() => setNewTransactionOpen(true)}
						>
							New Transaction
						</Button>
					</div>
				),
			},
		],
		[keyword, status, cancelMatch],
	);

	const columns = useMemo(
		() => [
			{
				title: 'ID',
				dataIndex: 'id',
				key: 'id',
				sorter: true,
				width: 140,
			},
			{
				title: 'Property',
				dataIndex: 'property',
				key: 'property',
				sorter: true,
				render: (_, record) => {
					const text = record?.property?.address?.address1 ?? '—';
					return (
						<Typography.Link
							onClick={(e) => {
								e.preventDefault();
								message.info('Property detail is not wired yet.');
							}}
						>
							{text}
						</Typography.Link>
					);
				},
			},
			{
				title: 'City',
				dataIndex: 'city',
				key: 'city',
				sorter: true,
				width: 120,
				render: (_, record) => record?.property?.address?.city ?? '—',
			},
			{
				title: 'State',
				dataIndex: 'state',
				key: 'state',
				width: 80,
				sorter: true,
				render: (_, record) => record?.property?.address?.state ?? '—',
			},
			{
				title: 'Type',
				dataIndex: 'type',
				key: 'type',
				width: 120,
				sorter: true,
				render: (_, record) => record.type ?? '—',
			},
			{
				title: 'Status',
				dataIndex: 'status',
				key: 'status',
				width: 120,
				sorter: true,
				render: (v) => (v ? statusTag(v) : '—'),
			},
			{
				title: 'COE Date',
				dataIndex: 'coe_date',
				key: 'coe_date',
				width: 120,
				sorter: true,
				render: (_, record) => {
					const raw = record.coe_date ?? record.coeDate;
					const formatted = raw ? formatDate(raw) : '';
					return formatted || '—';
				},
			},
			{
				title: 'Files',
				key: 'files',
				width: 80,
				align: 'center',
				render: () => (
					<Button
						type="text"
						icon={<FolderOpenOutlined className="text-base-warning" />}
						onClick={() => message.info('File list is not wired yet.')}
					/>
				),
			},
			{
				title: 'Action',
				key: 'action',
				width: 160,
				fixed: 'right',
				render: (_, record) => (
					<Space size="small">
						<Button
							size="small"
							onClick={() =>
								message.info(`Edit is a placeholder for transaction ${record.id}.`)
							}
						>
							Edit
						</Button>
						<Button
							size="small"
							danger
							type="default"
							onClick={() =>
								message.info(
									`Delete is a placeholder for transaction ${record.id}.`,
								)
							}
						>
							Delete
						</Button>
					</Space>
				),
			},
		],
		[],
	);

	return (
		<>
			<FilterTable
				ref={tableFilterRef}
				rowKey="id"
				columns={columns}
				filterList={filterList}
				filterName="Transaction filters"
				getData={getData}
				scroll={{ x: 'max-content' }}
				filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden flex flex-col min-h-0"
			/>
			<NewTransactionModal
				open={newTransactionOpen}
				onClose={() => setNewTransactionOpen(false)}
				onSaved={() => {
					cancelMatch(TRANSACTION_LIST_API);
					tableFilterRef.current?.handleSearch?.();
				}}
			/>
		</>
	);
};

export default MyTransactionsPage;
