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
			// setSearchQuery(query); 暂时取消获取dashboard跳转带的参数
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
			key: 'loan_stage',
			valueName: 'loanStatusOptions',
		},
		{
			type: 'lender',
		},
		...options,
	]);

	const loansProps = {
		filterWrapClassName: 'px-4 pb-4',
		api,
		searchQuery,
		setSearchQuery,
	};

	return (
		<Tabs
			activeKey={tabKey}
			items={[
				{
					key: 'Pending',
					label: 'Pending',
					children: <LoansPage stage="Pending" {...loansProps} />,
				},
				{
					key: 'Active',
					label: 'Active',
					children: <LoansPage stage="Active" {...loansProps} />,
				},
				{
					key: 'Archived',
					label: 'Archived',
					children: (
						<LoansPage
							filterWrapClassName="px-4 pb-4"
							stage="Archived"
							{...loansProps}
						/>
					),
				},
				{
					key: 'Draft',
					label: 'Drafts',
					children: <LoansPage stage="Draft" key="Draft" {...loansProps} />,
				},
			]}
			onChange={(v) => setTabKey(v)}
			className={`flex-1 overflow-hidden ${styles.loans} ${className}`}
			destroyOnHidden
		/>
	);
};
export default ListView;
