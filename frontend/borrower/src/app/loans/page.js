'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import styles from './loans.module.scss';
import PendingPage from './components/Pending';
import ActiveLoans from './components/Active';
import ArchivedLoans from './components/Archived';
import Drafts from './components/Drafts';
import { useGetOptions } from '@/hooks/useGetOptions';
const api = '/loan/page';
const Loans = () => {
	const [tabKey, setTabKey] = useState('pending');
	useGetOptions([
		{
			key: 'loan_purpose',
			valueName: 'loanPurposeOptions',
		},
		{
			key: 'loan_type',
			valueName: 'loanTypeOptions',
		},
		{
			key: 'loan_pending_status',
			valueName: 'loanPendingStatus',
		},
		{
			type: 'lender',
		},
	]);
	return (
		<Tabs
			activeKey={tabKey}
			items={[
				{
					key: 'pending',
					label: 'Pending Loans',
					children: <PendingPage filterWrapClassName="px-4 pb-4" api={api} />,
				},
				{
					key: 'active',
					label: 'Active Loans',
					children: <ActiveLoans filterWrapClassName="px-4 pb-4" api={api} />,
				},
				{
					key: 'archived',
					label: 'Archived Loans',
					children: <ArchivedLoans filterWrapClassName="px-4 pb-4" api={api} />,
				},
				{
					key: 'drafts',
					label: 'Drafts',
					children: <Drafts filterWrapClassName="px-4 pb-4" api={api} />,
				},
			]}
			onChange={(v) => setTabKey(v)}
			className={`h-full ${styles.loans}`}
		/>
	);
};
export default Loans;
