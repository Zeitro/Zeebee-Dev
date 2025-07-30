'use client';
import FilterField from '@/components/FilterField';
import Image from 'next/image';
import React, { useEffect, useMemo, useState } from 'react';
import styles from './index.module.scss';
import { sort_asc, sort_des } from '@/enum/properties';
import { Button, Col, Row } from 'antd';
const PropertiesItem = ({data}) => {
	//    const [type, setType] = useState('');

	const list1 = [
		{
			label: 'Total loan',
			value: 'totalLoan',
			render: (v) => `$ 980,000`,
		},
		{
			label: 'Outstanding loan',
			value: 'outstandingLoan',
			render: (v) => `$ 980,000`,
		},
		{
			label: 'Avg. rate',
			value: 'avgRate',
			render: (v) => `8.92%`,
		},
		{
			label: 'Total monthly payment',
			value: 'totalMonthlyPayment',
			render: (v) => `$ 980,000`,
		},
	];
	const list2 = [
		{
			label: 'Borrower',
			value: 'borrower',
		},
		{
			label: 'Lender',
			value: 'lender',
		},
		{
			label: 'Nearest maturity date',
			value: 'nearestMaturityDate',
			render: () => `11-01-2025`,
		},
		{
			render: () => <Button>View detials</Button>,
		},
	];

	const getLabelValue = (item, index) => (
		<Col lg={6} sm={12} xs={24} key={index}>
			{item?.label && <p className='text-base-label pb-1'>{item.label}</p>}
			<div className='text-base/[24px] font-medium flex items-center'>{item?.render?.() || item.value}</div>
		</Col>
	);
	return (
		<div className={[styles.propertiesItem].join(' ')}>
			<div className={styles.header}>
				<p className={styles.address}>512 Ocean Breeze Ave, Santa Monica, CA</p>
				<div className="flex items-center shrink-0 gap-2">
					<div className={styles.active}>2 Active Loans</div>
					<div className={styles.pending}>2 Pending Loans</div>
				</div>
			</div>
			<Row className={styles.info}>{list1.map(getLabelValue)}</Row>
			<Row className={styles.info}>{list2.map(getLabelValue)}</Row>
		</div>
	);
};
export default PropertiesItem;
