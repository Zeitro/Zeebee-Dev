'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import styles from './index.module.scss';
import LoansPage from './LoansPage';
import { useGetOptions } from '@repo/shared-methods/hooks';
const ListView = ({ className }) => {
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
					children: <LoansPage filterWrapClassName="px-4 pb-4" stage="Pending" />,
				},
				{
					key: 'active',
					label: 'Active',
					children: (
						<LoansPage filterWrapClassName="px-4 pb-4" stage="Active" scrollX={2860} />
					),
				},
				{
					key: 'archived',
					label: 'Archived',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Archived"
							scrollX={1880}
						/>
					),
				},
				{
					key: 'rejected',
					label: 'Rejected',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Rejected"
							key="Rejected"
							scrollX={2610}
						/>
					),
				},
				{
					key: 'Draft',
					label: 'Pre approval',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Draft"
							key="Draft"
							scrollX={2460}
						/>
					),
				},
			]}
			onChange={(v) => setTabKey(v)}
			className={`flex-1 overflow-hidden ${styles.loans} ${className}`}
			destroyOnHidden
		/>
	);
};
export default ListView;
