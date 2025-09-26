'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import styles from './loans.module.scss';
import PendingPage from './components/Pending';
import ActiveLoans from './components/Active';
import ArchivedLoans from './components/Archived';
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
			key: 'loan_stage',
			valueName: 'loanStatusOptions',
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
					children: (
						<PendingPage filterWrapClassName="px-4 pb-4" api={api} key="Pending" />
					),
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
					key: 'rejected',
					label: 'Rejected Loans',
					children: (
						<PendingPage
							filterWrapClassName="px-4 pb-4"
							api={api}
							stage="Rejected"
							key="Rejected"
							scrollX={2200}
						/>
					),
				},
			]}
			onChange={(v) => setTabKey(v)}
			className={`h-full ${styles.loans}`}
		/>
	);
};
export default Loans;
