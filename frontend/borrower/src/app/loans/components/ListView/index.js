'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import styles from './index.module.scss';
import LoansPage from './LoansPage';
import { useGetOptions } from '@repo/shared-methods/hooks';
const api = '/loan/page';

const ListView = ({ className, options = [] }) => {
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
		...options,
	]);
	return (
		<Tabs
			activeKey={tabKey}
			items={[
				{
					key: 'pending',
					label: 'Pending',
					children: (
						<LoansPage filterWrapClassName="px-4 pb-4" stage="Pending" api={api} />
					),
				},
				{
					key: 'active',
					label: 'Active',
					children: (
						<LoansPage filterWrapClassName="px-4 pb-4" stage="Active" api={api} />
					),
				},
				{
					key: 'archived',
					label: 'Archived',
					children: (
						<LoansPage filterWrapClassName="px-4 pb-4" stage="Archived" api={api} />
					),
				},
				{
					key: 'Draft',
					label: 'Drafts',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Draft"
							key="Draft"
							api={api}
						/>
					),
				},
			]}
			onChange={(v) => setTabKey(v)}
			className={`flex-1 overflow-hidden ${styles.loans} ${className}`}
			// destroyOnHidden
		/>
	);
};
export default ListView;
