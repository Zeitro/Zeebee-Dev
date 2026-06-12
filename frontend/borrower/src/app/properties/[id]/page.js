'use client';

import Image from 'next/image';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Modal, message, Tabs } from 'antd';
import { useRouter } from 'next/navigation';
import { useTabRoute } from '@/components/TabWorkspace/useTabRoute';
import AddModal from '@repo/shared-methods/componentsUi/SelectCard/Add/property';
import { useProperty } from '@repo/shared-methods/componentsUi/SelectCard/Item/useCommon';
import { useAuth } from '@repo/shared-methods/hooks';
import { assemblyAddress, createRequest, displayValue, eventBus } from '@repo/shared-methods/utils';
import NewTransactionModal from '@/app/myTransactions/components/NewTransactionModal';
import { PERMISSIONS } from '@config/permissions';
import DetailsTab from './components/DetailsTab';
import TransactionsTab from './components/TransactionsTab';
import LoansTab from './components/LoansTab';
import styles from './index.module.scss';

const PropertiesDetails = () => {
	const { params } = useTabRoute();
	const { id } = params;
	const router = useRouter();
	const { hasPermission } = useAuth();
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState({});
	const [activeTab, setActiveTab] = useState('details');
	const [editOpen, setEditOpen] = useState(false);
	const [transactionModalOpen, setTransactionModalOpen] = useState(false);
	const refreshTransactionsListRef = useRef(() => {});
	const api = `/property/${id}/summary`;

	const { handleDelete, loading: deleteLoading } = useProperty({
		handleSearch: () => router.push('/properties'),
	});

	const getDetails = useCallback(async () => {
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data);
			eventBus.emit('changeCrumbs', { 2: assemblyAddress(data?.address) });
		} catch {
			// keep previous rowData on error
		} finally {
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [api]);

	const actionItems = useMemo(() => {
		const items = [];
		if (hasPermission(PERMISSIONS?.PROPERTY?.EDIT)) {
			items.push({
				key: 'edit',
				label: 'Edit',
				onClick: () => setEditOpen(true),
			});
		}
		items.push({
			key: 'newTransaction',
			label: 'New Transaction',
			onClick: () => setTransactionModalOpen(true),
		});
		if (hasPermission(PERMISSIONS?.PROPERTY?.DELETE)) {
			items.push({
				key: 'delete',
				label: 'Delete',
				danger: true,
				onClick: () => {
					if (rowData?.loan_count) {
						message.destroy();
						message.open({
							type: 'warning',
							content:
								'This property has an associated loan, so it cannot be deleted.',
						});
						return;
					}
					Modal.confirm({
						title: 'Are you sure delete this Property?',
						content: assemblyAddress(rowData?.address),
						okText: 'Delete',
						okButtonProps: { danger: true },
						onOk: () => handleDelete(rowData?.id),
						confirmLoading: deleteLoading,
					});
				},
			});
		}
		return items;
	}, [deleteLoading, handleDelete, hasPermission, rowData]);

	useEffect(() => {
		const timer = setTimeout(() => {
			getDetails();
		}, 0);
		return () => clearTimeout(timer);
	}, [getDetails]);

	const tabItems = useMemo(() => {
		let items = [
			{
				key: 'details',
				label: 'Details',
				children: (
					<div className="px-6 py-6">
						<DetailsTab
							rowData={rowData}
							propertyId={id}
							canEdit={hasPermission(PERMISSIONS?.PROPERTY?.EDIT)}
							onRefresh={getDetails}
						/>
					</div>
				),
			},

			{
				key: 'loans',
				label: 'Loans',
				children: (
					<LoansTab
						rowData={rowData}
						loading={loading}
						onRefresh={getDetails}
						className="px-4 pb-4"
					/>
				),
			},
		];
		if (hasPermission(PERMISSIONS?.TRANSACTION?.PAGE)) {
			items.splice(1, 0, {
				key: 'transactions',
				label: 'Transactions',
				children: (
					<TransactionsTab
						propertyId={id}
						onNewTransaction={() => setTransactionModalOpen(true)}
						onRegisterRefresh={(fn) => {
							refreshTransactionsListRef.current = fn ?? (() => {});
						}}
					/>
				),
			});
		}
		return items;
	}, [getDetails, hasPermission, id, loading, rowData]);

	return (
		<div className="h-full flex flex-col overflow-hidden flex-1 bg-white">
			<div className="flex items-center justify-between gap-4 shrink-0 px-6 pt-6 pb-4">
				<div className="flex items-center gap-2 min-w-0">
					<Image src="/images/home.svg" width={32} height={32} alt="" />
					<p className="text-xl leading-7 font-medium truncate">
						{displayValue(assemblyAddress(rowData?.address))}
					</p>
				</div>
				{/* {actionItems.length > 0 && (
					<Dropdown menu={{ items: actionItems }} trigger={['click']}>
						<Button>Actions</Button>
					</Dropdown>
				)} */}
			</div>
			<Tabs
				activeKey={activeTab}
				onChange={setActiveTab}
				items={tabItems}
				className={styles.propertyDetailTabs}
				destroyOnHidden={false}
			/>
			{editOpen && (
				<AddModal setOpen={setEditOpen} rowData={rowData} handleSearch={getDetails} />
			)}
			<NewTransactionModal
				open={transactionModalOpen}
				initialProperty={rowData?.id ? rowData : null}
				onClose={() => setTransactionModalOpen(false)}
				onSaved={() => {
					setTransactionModalOpen(false);
					void refreshTransactionsListRef.current?.();
				}}
			/>
		</div>
	);
};

export default PropertiesDetails;
