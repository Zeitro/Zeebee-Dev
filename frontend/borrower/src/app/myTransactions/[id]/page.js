'use client';

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Button, Modal, Select, Tabs, message } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import {
	createRequest,
	displayValue,
	eventBus,
	formatDate,
	formatTime,
} from '@repo/shared-methods/utils';
import { useGlobalStore } from '@/store/global';
import DocumentsTab from './components/DocumentsTab';
import DetailsTab from './components/DetailsTab';
import LoansTab from './components/LoansTab';
import LogTab from './components/LogTab';
import { buildMockTransactionDetail } from './components/mockTransactionDetail';
import {
	TRANSACTION_STATUS_OPTIONS,
	formatTransactionId,
	getTransactionTitle,
	resolveTransactionDetailTab,
} from './components/utils';
import styles from './index.module.scss';

const TransactionDetailsPage = () => {
	const { id } = useParams();
	const router = useRouter();
	const searchParams = useSearchParams();
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState({});
	const tabFromUrl = resolveTransactionDetailTab(searchParams.get('tab'));
	const [statusSaving, setStatusSaving] = useState(false);
	const { changeSidebarCollapsed } = useGlobalStore();
	const api = `/transaction/${id}`;

	const getDetails = useCallback(async () => {
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data);
			eventBus.emit('changeCrumbs', { 2: getTransactionTitle(data) });
		} catch {
			const mock = buildMockTransactionDetail(id);
			setRowData(mock);
			eventBus.emit('changeCrumbs', { 2: getTransactionTitle(mock) });
		} finally {
			setLoading(false);
		}
	}, [api, id, request]);

	const handleStatusChange = async (status) => {
		try {
			setStatusSaving(true);
			await request(api, {
				method: 'PUT',
				body: { ...rowData, status },
			});
			setRowData((prev) => ({ ...prev, status }));
			message.success('Status updated');
		} catch {
			setRowData((prev) => ({ ...prev, status }));
		} finally {
			setStatusSaving(false);
		}
	};

	const handleDelete = () => {
		Modal.confirm({
			title: 'Are you sure delete this transaction?',
			content: formatTransactionId(rowData),
			okText: 'Delete',
			okButtonProps: { danger: true },
			onOk: async () => {
				try {
					await request(api, { method: 'DELETE' });
					message.success('Transaction deleted');
					router.push('/myTransactions');
				} catch {
					message.info('Delete is not wired to the API yet.');
				}
			},
		});
	};

	const handleTabChange = (key) => {
		const params = new URLSearchParams(searchParams.toString());
		if (key === 'documents') {
			params.delete('tab');
		} else {
			params.set('tab', key);
		}
		const query = params.toString();
		router.replace(query ? `/myTransactions/${id}?${query}` : `/myTransactions/${id}`, {
			scroll: false,
		});
	};

	useEffect(() => {
		changeSidebarCollapsed(true);
		const timer = setTimeout(() => {
			getDetails();
		}, 0);
		return () => clearTimeout(timer);
	}, [changeSidebarCollapsed, getDetails]);

	const updatedRaw = rowData?.updated_at ?? rowData?.updatedAt;
	const updatedDate = updatedRaw ? formatDate(updatedRaw) : displayValue(null);
	const updatedTime = updatedRaw ? formatTime(updatedRaw) : null;
	const coeRaw = rowData?.coe_date ?? rowData?.coeDate;
	const status = rowData?.status;

	const tabItems = useMemo(
		() => [
			{
				key: 'documents',
				label: 'Document',
				children: <DocumentsTab rowData={rowData} />,
			},
			{
				key: 'details',
				label: 'Details',
				children: (
					<div className="px-6 py-6">
						<DetailsTab rowData={rowData} transactionId={id} onRefresh={getDetails} />
					</div>
				),
			},
			{
				key: 'loans',
				label: 'Loan',
				children: <LoansTab rowData={rowData} loading={loading} onRefresh={getDetails} />,
			},
			{
				key: 'log',
				label: 'Log',
				children: <LogTab transactionId={id} />,
			},
		],
		[getDetails, id, loading, rowData],
	);

	return (
		<div className="flex h-full flex-1 flex-col overflow-hidden bg-white">
			<div className="flex shrink-0 flex-col gap-4 px-6 pb-4 pt-6">
				<div className="flex flex-wrap items-start justify-between gap-4">
					<div className="flex min-w-0 flex-col gap-2">
						<div className="flex flex-wrap items-center gap-8">
							<h1 className="m-0 text-xl font-medium leading-7 text-foreground">
								{getTransactionTitle(rowData)}
							</h1>
							<div className="flex items-center gap-2">
								<span className="text-sm font-medium leading-[22px] text-base-gray">
									Status
								</span>
								<Select
									className="min-w-[120px]"
									value={status}
									placeholder="Status"
									options={TRANSACTION_STATUS_OPTIONS}
									loading={statusSaving}
									onChange={handleStatusChange}
									allowClear={false}
								/>
							</div>
						</div>
						<div className="flex flex-wrap gap-6 text-sm font-medium leading-[22px]">
							<span className="text-base-gray">
								ID{' '}
								<span className="text-foreground">
									{formatTransactionId(rowData)}
								</span>
							</span>
							<span className="text-base-gray">
								COE{' '}
								<span className="text-foreground">
									{displayValue(coeRaw, formatDate)}
								</span>
							</span>
							<span className="text-base-gray">
								Updated{' '}
								<span className="text-foreground">
									{updatedDate}
									{updatedTime ? ` ${updatedTime}` : ''}
								</span>
							</span>
						</div>
					</div>
					<Button danger icon={<DeleteOutlined />} onClick={handleDelete}>
						Delete
					</Button>
				</div>
			</div>

			<Tabs
				activeKey={tabFromUrl}
				onChange={handleTabChange}
				items={tabItems}
				className={styles.transactionDetailTabs}
				destroyOnHidden={false}
			/>
		</div>
	);
};

export default TransactionDetailsPage;
