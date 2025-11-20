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
	]);
	return (
		<Tabs
			activeKey={tabKey}
			items={[
				{
					key: 'pending',
					label: 'Pending',
					children: (
						<PendingPage filterWrapClassName="px-4 pb-4" api={api} key="Pending" />
					),
				},
				{
					key: 'active',
					label: 'Active',
					children: <ActiveLoans filterWrapClassName="px-4 pb-4" api={api} />,
				},
				{
					key: 'archived',
					label: 'Archived',
					children: <ArchivedLoans filterWrapClassName="px-4 pb-4" api={api} />,
				},
				{
					key: 'rejected',
					label: 'Rejected',
					children: (
						<PendingPage
							filterWrapClassName="px-4 pb-4"
							api={api}
							stage="Rejected"
							key="Rejected"
							scrollX={2400}
						/>
					),
				},
				{
					key: 'Draft',
					label: 'Pre approval',
					children: (
						<PendingPage
							filterWrapClassName="px-4 pb-4"
							api={api}
							stage="Draft"
							key="Draft"
							scrollX={2250}
						/>
					),
				},
			]}
			onChange={(v) => setTabKey(v)}
			className={`h-full ${styles.loans}`}
			destroyOnHidden
		/>
	);
};
export default Loans;
