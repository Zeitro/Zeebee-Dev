'use client';

import React, { useMemo } from 'react';
import { useGetOptions } from '@repo/shared-methods/hooks';
import LoansPage from '@/app/loans/components/ListView/LoansPage';

const LOAN_SECTIONS = [
	{ title: 'Pending loan', key: 'pending_loans', stage: 'Pending' },
	{ title: 'Active loan', key: 'active_loans', stage: 'Active' },
	{ title: 'Archived loan', key: 'archived_loans', stage: 'Archived' },
];

const LoansTab = ({ rowData, loading, onRefresh }) => {
	useGetOptions([
		{
			key: 'loan_pending_status',
			valueName: 'loanPendingStatus',
		},
	]);

	const sections = useMemo(
		() =>
			LOAN_SECTIONS.map((section) => ({
				...section,
				dataSource: rowData?.[section.key],
			})),
		[rowData],
	);

	return (
		<div className="[&_.ant-table-column-title]:!text-base-gray [&_th]:!text-base-gray">
			{sections.map((section) => (
				<div key={section.key} className="border border-base-border rounded-lg mb-4 pb-1">
					<div className="py-2 px-4 text-base/6 font-medium border-b border-base-border">
						{section.title} ({section.dataSource.length})
					</div>
					<LoansPage
						hasFilterField={false}
						hasPagination={false}
						dataSource={section.dataSource || []}
						isNotAddress={true}
						loading={loading}
						successCallback={onRefresh}
						stage={section.stage}
					/>
				</div>
			))}
		</div>
	);
};

export default LoansTab;
