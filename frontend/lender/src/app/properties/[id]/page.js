'use client';
import Image from 'next/image';
import React, { useEffect, useMemo, useState } from 'react';
import { propertiesFieldArr } from '@/enum/properties';
import { Col, Row } from 'antd';
import PendingPage from '../../loans/components/Pending';
import ActiveLoans from '../../loans/components/Active';
import ArchivedLoans from '../../loans/components/Archived';
import { useParams } from 'next/navigation';
import { createRequest } from '@/lib/fetch';
import { assemblyAddress } from '@/lib/format';
import { useGetOptions } from '@/hooks/useGetOptions';
import eventBus from '@/lib/bus';
import { useGlobalStore } from '@/store/global';
import FieldTips from '@/components/FieldTips';
const PropertiesDetails = () => {
	const { id } = useParams();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState({});
	const { changeSiberCollapsed } = useGlobalStore();
	const api = `/property/${id}`;

	const getDetails = async () => {
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data);
			eventBus.emit('changeCrumbs', { 2: assemblyAddress(data?.address) });
			setLoading(false);
		} catch (error) {
			setLoading(false);
		}
	};

	const list = useMemo(() => {
		return [
			{
				title: 'Pending loan',
				total: rowData?.pending_loans?.length,
				content: (
					<PendingPage
						hasFilterField={false}
						hasPagination={false}
						dataSource={rowData?.pending_loans}
						isNotAddress={true}
						loading={loading}
						scrollX={2000}
						successCallback={getDetails}
					/>
				),
			},
			{
				title: 'Active loan',
				total: rowData?.active_loans?.length,
				content: (
					<ActiveLoans
						hasFilterField={false}
						hasPagination={false}
						dataSource={rowData?.active_loans}
						isNotAddress={true}
						loading={loading}
						scrollX={2100}
						successCallback={getDetails}
					/>
				),
			},
			{
				title: 'Archived loan',
				total: rowData?.archived_loans?.length,
				content: (
					<ArchivedLoans
						hasFilterField={false}
						hasPagination={false}
						dataSource={rowData?.archived_loans}
						isNotAddress={true}
						loading={loading}
						scrollX={1000}
						successCallback={getDetails}
					/>
				),
			},
		];
	}, [loading, rowData]);

	useGetOptions([
		{
			key: 'loan_pending_status',
			valueName: 'loanPendingStatus',
		},
	]);

	useEffect(() => {
		changeSiberCollapsed(true);
		getDetails();
	}, []);
	return (
		<div className={['h-full px-8 py-6 overflow-hidden overflow-y-auto'].join(' ')}>
			<div className="flex items-center pb-6">
				<Image src="/images/home.svg" width={32} height={32} alt="" />
				<p className="text-[24px]/8 font-medium">{assemblyAddress(rowData?.address)}</p>
			</div>
			<Row className="border border-base-disabled rounded-lg px-4 py-4">
				{propertiesFieldArr.map((item, index) => (
					<Col lg={6} sm={12} xs={24} key={index}>
						{item?.label && <div className="text-base-label pb-1">{item.label}</div>}
						<div className="text-base/[24px] font-medium flex items-center">
							{item?.render?.(rowData?.[item.value])}
						</div>
					</Col>
				))}
			</Row>
			{list.map((item, index) => {
				return (
					<div key={index}>
						<div className="pt-6 pb-2 mb-4 text-xl/7 font-medium border-b border-[#91CAFF]">
							{item.title}
						</div>
						{item.content}
					</div>
				);
			})}
		</div>
	);
};
export default PropertiesDetails;
