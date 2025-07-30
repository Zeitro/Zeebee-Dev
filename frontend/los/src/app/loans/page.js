'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import NewTable from '../../components/NewTable';
import { activeColumns, archivedColumns, pendingColumns } from './components/column';
import styles from './loans.module.scss';
import PendingPage from './components/Pending';
import ActiveLoans from './components/Active';
import ArchivedLoans from './components/Archived';

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
			children: <ActiveLoans />,
		},
		{
			key: 'archived',
			label: 'Archived Loans',
			children: <ArchivedLoans />,
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
