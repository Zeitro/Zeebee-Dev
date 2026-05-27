'use client';

import React, { useMemo, useState } from 'react';
import { Spin, Tabs } from 'antd';
import DocumentsTab from '@/app/myTransactions/[id]/components/DocumentsTab';
import DetailsTab from '@/app/myTransactions/[id]/components/DetailsTab';
import LoansTab from '@/app/myTransactions/[id]/components/LoansTab';
import LogTab from '@/app/myTransactions/[id]/components/LogTab';
import { useTransactionDetail } from './useTransactionDetail';
import styles from './TransactionCard.module.scss';

const TransactionCardExpandedContent = ({ transactionId }) => {
	const [activeTab, setActiveTab] = useState('documents');
	const { rowData, loading, getDetails } = useTransactionDetail(transactionId, true);

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
						<DetailsTab
							rowData={rowData}
							transactionId={transactionId}
							onRefresh={getDetails}
						/>
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
				children: <LogTab transactionId={transactionId} />,
			},
		],
		[getDetails, loading, rowData, transactionId],
	);

	return (
		<div
			className="border-t border-base-border bg-background"
			onClick={(e) => e.stopPropagation()}
			onKeyDown={(e) => e.stopPropagation()}
		>
			<Spin spinning={loading}>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={tabItems}
					className={styles.transactionCardTabs}
					destroyOnHidden={false}
				/>
			</Spin>
		</div>
	);
};

export default TransactionCardExpandedContent;
