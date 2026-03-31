'use client';
import { Tabs } from 'antd';
import React, { useEffect, useState } from 'react';
import styles from './index.module.scss';
import LoansPage from './LoansPage';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { sessionStorageUtil } from '@repo/shared-methods/utils';
const api = '/loan/page';

const ListView = ({ className, options = [] }) => {
	const [tabKey, setTabKey] = useState(null);
	const [searchQuery, setSearchQuery] = useState({});

	// 跳转到loans页面所带的隐藏参数（因刷新后不带这些参数，所以remove掉）
	useEffect(() => {
		const getUrlQuery = () => {
			const query = sessionStorageUtil.get('LOANS-FILTER');
			setSearchQuery(query);
			setTabKey(query?.stage || 'Pending');
			sessionStorageUtil.remove('LOANS-FILTER');
		};
		getUrlQuery();
	}, []);
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
					key: 'Pending',
					label: 'Pending',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Pending"
							api={api}
							searchQuery={searchQuery}
						/>
					),
				},
				{
					key: 'Active',
					label: 'Active',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Active"
							api={api}
							searchQuery={searchQuery}
						/>
					),
				},
				{
					key: 'Archived',
					label: 'Archived',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Archived"
							api={api}
							searchQuery={searchQuery}
						/>
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
							searchQuery={searchQuery}
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
