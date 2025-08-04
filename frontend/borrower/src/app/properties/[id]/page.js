'use client';
import Image from 'next/image';
import React, { useEffect, useMemo, useState } from 'react';
import { propertiesFieldArr } from '@/enum/properties';
import { Col, Row } from 'antd';
import PendingPage from '../../loans/components/Pending';
import ActiveLoans from '../../loans/components/Active';
import ArchivedLoans from '../../loans/components/Archived';
const PropertiesDetails = () => {
	const list = [
		{
			title: 'Pending loan',
			content: <PendingPage hasFilterField={false} hasPagination={false}/>,
		},
		{
			title: 'Active loan',
			content: <ActiveLoans hasFilterField={false} hasPagination={false}/>,
		},
		{
			title: 'Archived loan',
			content: <ArchivedLoans hasFilterField={false} hasPagination={false}/>,
		},
	];
	return (
		<div className={['h-full px-8 py-6 overflow-hidden overflow-y-auto'].join(' ')}>
			<div className="flex items-center pb-6">
				<Image src={require('@/assets/images/home.svg')} alt="" className="w-8 h-8" />
				<p className="text-[24px]/8 font-medium">1396 Cedar Hollow Rd, Austin, TX</p>
			</div>
			<Row className="border border-base-disabled rounded-lg px-4 py-4">
				{propertiesFieldArr.map((item, index) => (
					<Col lg={6} sm={12} xs={24} key={index}>
						{item?.label && <p className="text-base-label pb-1">{item.label}</p>}
						<div className="text-base/[24px] font-medium flex items-center">
							{item?.render?.() || item.value}
						</div>
					</Col>
				))}
			</Row>
			{list.map((item, index) => {
				return (
					<div key={index}>
						<div className='pt-6 pb-2 mb-4 text-xl/7 font-medium border-b border-[#91CAFF]'>{item.title}</div>
						{item.content}
					</div>
				);
			})}
		</div>
	);
};
export default PropertiesDetails;
