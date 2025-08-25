'use client';
import React, { useState } from 'react';
import NewTable from '@/components/NewTable';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
const Experience = () => {
	const [historyList, setHistoryList] = useState([{ id: 1 }]);
	const [realEstateList, setRealEstateList] = useState([{ id: 2 }]);
	const { transactionHistoryColumns, realEstateColumns } = useColumnByZeeBeeLoan();

	return (
		<div className="gap-y-8 flex flex-col w-full">
			<div className="text-2xl font-medium">Transaction history</div>
			{historyList?.length > 0 && (
				<NewTable
					dataSource={historyList}
					columns={transactionHistoryColumns}
					pagination={false}
					rowKey="id"
				/>
			)}
			<AddEditBtn>Add transaction history</AddEditBtn>
			<div className="text-2xl font-medium pt-6 border-t border-base-border">
				Real estate owned
			</div>
			{realEstateList?.length > 0 && (
				<NewTable
					dataSource={realEstateList}
					columns={realEstateColumns}
					pagination={false}
					rowKey="id"
				/>
			)}
			<AddEditBtn>Add real estate</AddEditBtn>
			<div className="pb-6"></div>
		</div>
	);
};

export default Experience;
