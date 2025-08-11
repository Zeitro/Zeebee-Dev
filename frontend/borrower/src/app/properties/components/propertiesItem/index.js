'use client';
import React from 'react';
import styles from './index.module.scss';
import { propertiesFieldArr } from '@/enum/properties';
import { Button, Col, Row } from 'antd';
import { useRouter } from 'next/navigation';
import { formatDate } from '@/lib/date';
import { assemblyAddress } from '@/lib/common';
const PropertiesItem = ({ data }) => {
	const router = useRouter();
	const list2 = [
		{
			label: 'Borrower',
			value: 'applicants',
			render: (v) => v?.join(','),
		},
		{
			label: 'Lender',
			value: 'lenders',
		},
		{
			label: 'Nearest maturity date',
			value: 'nearest_maturity_date',
			render: (v) => formatDate(v),
		},
		{
			render: () => (
				<Button onClick={() => router.push(`/properties/${data?.id}`)}>View details</Button>
			),
		},
	];

	const getLabelValue = (item, index) => (
		<Col lg={6} sm={12} xs={24} key={index}>
			{item?.label && <p className="text-base-label pb-1">{item.label}</p>}
			<div className="text-base/[24px] font-medium flex items-center">
				{item?.render?.(data[item.value])}
			</div>
		</Col>
	);
	return (
		<div className={[styles.propertiesItem].join(' ')}>
			<div className={styles.header}>
				<p className={styles.address}>{assemblyAddress(data?.address)}</p>
				<div className="flex items-center shrink-0 gap-2">
					{!!data?.active_loan_count && (
						<div className={styles.active}>
							{data?.active_loan_count || 0} Active Loans
						</div>
					)}
					{!!data?.pending_loan_count && (
						<div className={styles.pending}>
							{data?.pending_loan_count || 0} Pending Loans
						</div>
					)}
				</div>
			</div>
			<Row className={styles.info}>{propertiesFieldArr.map(getLabelValue)}</Row>
			<Row className={styles.info}>{list2.map(getLabelValue)}</Row>
		</div>
	);
};
export default PropertiesItem;
