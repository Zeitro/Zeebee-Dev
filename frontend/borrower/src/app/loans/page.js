'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import styles from './loans.module.scss';
import PendingPage from './components/Pending';
import ActiveLoans from './components/Active';
import ArchivedLoans from './components/Archived';
const Loans = () => {
	const [tabKey, setTabKey] = useState('pending');
    
	return (
		<Tabs
			activeKey={tabKey}
			items={[
		{
			key: 'pending',
			label: 'Pending Loans',
			children: <PendingPage filterWrapClassName="px-4 pb-4"/>,
		},
		{
			key: 'active',
			label: 'Active Loans',
			children: <ActiveLoans filterWrapClassName="px-4 pb-4"/>,
		},
		{
			key: 'archived',
			label: 'Archived Loans',
			children: <ArchivedLoans filterWrapClassName="px-4 pb-4"/>,
		},
	]}
			onChange={(v) => setTabKey(v)}
			className={`h-full ${styles.loans}`}
		/>
	);
};
export default Loans;
