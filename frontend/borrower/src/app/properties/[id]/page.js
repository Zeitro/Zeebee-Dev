'use client';
import Image from 'next/image';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { propertiesFieldArr, propertyList } from '@/enum/properties';
import PendingPage from '../../loans/components/Pending';
import ActiveLoans from '../../loans/components/Active';
import ArchivedLoans from '../../loans/components/Archived';
import { useParams } from 'next/navigation';
import { createRequest, assemblyAddress, eventBus } from '@repo/shared-methods/utils';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { useGlobalStore } from '@/store/global';
const PropertiesDetails = () => {
	const { id } = useParams();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState({});
	const { changeSidebarCollapsed } = useGlobalStore();
	const api = `/property/${id}/summary`;

	const getDetails = useCallback(async () => {
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data);
			eventBus.emit('changeCrumbs', { 2: assemblyAddress(data?.address) });
			setLoading(false);
		} catch (error) {
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [api]);

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
	}, [
		getDetails,
		loading,
		rowData?.active_loans,
		rowData?.archived_loans,
		rowData?.pending_loans,
	]);

	useGetOptions([
		{
			key: 'loan_pending_status',
			valueName: 'loanPendingStatus',
		},
	]);

	useEffect(() => {
		changeSidebarCollapsed(true);

		// macrotask 调度，能让 React 完成当前渲染周期后再执行 setState
		const timer = setTimeout(async () => {
			await getDetails();
		}, 0);

		return () => {
			clearTimeout(timer);
		};
	}, [changeSidebarCollapsed, getDetails]);

	const getItem = (item, index) => {
		return (
			<div key={index} className={[item?.className || 'w-1/2', 'pb-4'].join(' ')}>
				<p className="text-base-gray text-sm pb-1">{item.label}</p>
				<span className="text-black-dark text-base font-medium">
					{(item?.render ? item?.render?.(rowData[item.value]) : rowData[item.value]) ||
						'--'}
				</span>
			</div>
		);
	};

	return (
		<div className={['h-full px-8 py-6 overflow-hidden overflow-y-auto flex-1'].join(' ')}>
			<div className="flex items-center pb-6">
				<Image src="/images/home.svg" width={32} height={32} alt="" />
				<p className="text-[24px]/8 font-medium">{assemblyAddress(rowData?.address)}</p>
			</div>
			<div className="flex flex-wrap gap-6 flex-col md:flex-row">
				<div className="w-full flex-1 border border-base-border rounded-lg pt-4 px-4 flex flex-wrap">
					{propertyList.map(getItem)}
				</div>
				<div className="w-full  flex-1 border border-base-border rounded-lg pt-4 px-4 flex flex-wrap">
					{propertiesFieldArr.map(getItem)}
				</div>
			</div>
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
