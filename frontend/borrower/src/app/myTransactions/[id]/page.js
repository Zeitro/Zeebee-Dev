'use client';

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Tabs, message } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import {
	createRequest,
	displayValue,
	eventBus,
	formatDate,
	formatTime,
} from '@repo/shared-methods/utils';
import { DeleteBtn } from '@repo/shared-methods/components/NewModal/BtnPopUpModal';
import { useGetOptions } from '@repo/shared-methods/hooks';
import TransactionStatusSelect from './components/TransactionStatusSelect';
import DocumentsTab from './components/DocumentsTab';
import DetailsTab from './components/DetailsTab';
import LoansTab from './components/LoansTab';
import LogTab from './components/LogTab';
import TransactionDetailSkeleton from './components/TransactionDetailSkeleton';
import styles from './index.module.scss';
import { PermissionSlot } from '@repo/shared-methods';
import { PERMISSIONS } from '@config/permissions';

const TransactionDetailsPage = () => {
	const { id } = useParams();
	const router = useRouter();
	const searchParams = useSearchParams();
	const isSearchParamsEditing = useMemo(
		() => searchParams.get('isEditing') === 'true',
		[searchParams],
	);
	const [messageApi, contextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [isInitialLoading, setIsInitialLoading] = useState(true);
	const [deleteLoading, setDeleteLoading] = useState(false);
	const [rowData, setRowData] = useState({});
	const [tabFromUrl, setTabFromUrl] = useState(searchParams.get('tab') || 'documents');
	// const tabFromUrl = resolveTransactionDetailTab(searchParams.get('tab'));
	const { transactionStatusOptions } = useGetOptions([
		{ key: 'borrower_property_transaction_status', valueName: 'transactionStatusOptions' },
	]);
	const api = `/borrower_property_transaction/${id}`;

	const getDetails = useCallback(async () => {
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data);
			eventBus.emit('changeCrumbs', { 2: `${data.type} transaction` });
		} catch {
			setRowData({});
		} finally {
			setLoading(false);
			setIsInitialLoading(false);
		}
	}, [api, request]);

	const handleDelete = async () => {
		try {
			setDeleteLoading(true);
			await request(`/borrower_property_transaction/${id}`, { method: 'DELETE' });
			messageApi.success('Delete transaction successfully');
			router.push('/myTransactions');
		} finally {
			setDeleteLoading(false);
		}
	};

	const handleTabChange = (key) => {
		const params = new URLSearchParams(searchParams.toString());
		if (key === 'documents') {
			params.delete('tab');
			setTabFromUrl('documents');
		} else {
			params.set('tab', key);
			setTabFromUrl(key);
		}
		// 仅 Details 用父级 rowData，切回时刷新交易详情；
		// Loan / Log 靠 destroyOnHidden remount 自行请求，不必调 getDetails
		if (key === 'details') {
			getDetails();
		}
		const query = params.toString();
		router.replace(query ? `/myTransactions/${id}?${query}` : `/myTransactions/${id}`, {
			scroll: false,
		});
	};

	useEffect(() => {
		const timer = setTimeout(() => {
			getDetails();
		}, 0);
		return () => clearTimeout(timer);
	}, [getDetails]);

	/**
	 * Documents 用 destroyOnHidden:false 常驻；其余 tab 隐藏时销毁，切回 remount 拿最新数据。
	 */
	const tabItems = useMemo(
		() => [
			{
				key: 'documents',
				label: 'Documents',
				destroyOnHidden: false,
				children: <DocumentsTab rowData={rowData} onRefresh={getDetails} />,
			},
			{
				key: 'details',
				label: 'Details',
				destroyOnHidden: true,
				children: (
					<div className="px-6 py-6">
						<DetailsTab
							rowData={rowData}
							transactionId={id}
							onRefresh={getDetails}
							isSearchParamsEditing={isSearchParamsEditing}
						/>
					</div>
				),
			},
			{
				key: 'loans',
				label: 'Loan',
				destroyOnHidden: true,
				children: <LoansTab rowData={rowData} loading={loading} onRefresh={getDetails} />,
			},
			{
				key: 'log',
				label: 'Log',
				destroyOnHidden: true,
				children: <LogTab transactionId={id} />,
			},
		],
		[getDetails, id, loading, rowData, isSearchParamsEditing],
	);

	if (isInitialLoading) {
		return <TransactionDetailSkeleton />;
	}

	return (
		<div className="flex h-full flex-1 flex-col overflow-hidden bg-white">
			{contextHolder}
			<div className="flex shrink-0 flex-col gap-4 px-6 pb-4 pt-6">
				<div className="flex flex-wrap items-start justify-between gap-4">
					<div className="flex min-w-0 flex-col gap-2 flex-1">
						<div className="flex flex-wrap items-center gap-8">
							<h1 className="m-0 text-xl font-medium leading-7 text-foreground">
								{rowData.type} transaction
							</h1>
							<div className="flex items-center gap-2">
								<span className="text-sm font-medium leading-[22px] text-base-gray">
									Status
								</span>
								<TransactionStatusSelect
									value={rowData?.status}
									options={transactionStatusOptions}
									id={id}
									onSuccess={getDetails}
								/>
							</div>
						</div>
						<div className="flex flex-wrap gap-6 text-sm font-medium leading-[22px]">
							<span className="text-base-gray">
								ID
								<span className="text-foreground pl-1">
									{displayValue(rowData?.code)}
								</span>
							</span>
							<span className="text-base-gray">
								COE
								<span className="text-foreground pl-1">
									{displayValue(rowData?.coe_date, formatDate)}
								</span>
							</span>
							<span className="text-base-gray">
								Updated
								<span className="text-foreground pl-1">
									{displayValue(rowData?.updated_at, formatTime)}
								</span>
							</span>
						</div>
					</div>
					<PermissionSlot code={PERMISSIONS?.TRANSACTION?.DELETE}>
						<DeleteBtn
							type="default"
							size="middle"
							icon={<DeleteOutlined />}
							content={`ID: ${rowData?.code}`}
							onOk={handleDelete}
							loading={deleteLoading}
						/>
					</PermissionSlot>
				</div>
			</div>

			<Tabs
				activeKey={tabFromUrl}
				onChange={handleTabChange}
				items={tabItems}
				className={styles.transactionDetailTabs}
			/>
		</div>
	);
};

export default TransactionDetailsPage;
