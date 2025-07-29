'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import NewTable from '../../components/NewTable';
import { activeColumns, archivedColumns, pendingColumns } from './components/column';
import styles from './loans.module.scss';
import PendingPage from './components/LoanFilterTable/pending';

const Loans = () => {
	const [tabKey, setTabKey] = useState('pending');
	const items = [
		{
			key: 'pending',
			label: 'Pending Loans',
			children: <PendingPage />,
		},
		{
			key: 'active',
			label: 'Active Loans',
			children: (
				<NewTable dataSource={[]} columns={activeColumns} scroll={{ x: 3500, y: '100%' }} />
			),
		},
		{
			key: 'archived',
			label: 'Archived Loans',
			children: (
				<NewTable
					dataSource={[]}
					columns={archivedColumns}
					scroll={{ x: 2000, y: '100%' }}
				/>
			),
		},
	];
	return (
		<Tabs
			activeKey={tabKey}
			items={items}
			onChange={(v) => setTabKey(v)}
			className={`h-full ${styles.loans}`}
		/>
	);
};
export default Loans;
