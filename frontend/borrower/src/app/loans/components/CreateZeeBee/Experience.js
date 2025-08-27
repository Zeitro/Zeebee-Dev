'use client';
import React, { useState } from 'react';
import NewTable from '@/components/NewTable';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
import { message } from 'antd';
const Experience = () => {
	const [historyList, setHistoryList] = useState([]);
	const [realEstateList, setRealEstateList] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();

	const commonAdd = (arr, setFn, values, key) => {
		if (key === 'Edit') {
			const index = arr.findIndex((i) => i.id === values.id);
			arr.splice(index, 1, values);
			setFn([...arr]);
		} else {
			setFn([...arr, { ...values, id: arr.length + 1 }]);
		}
		message.destroy();
		messageApi.open({
			type: 'success',
			content: `${key}ed successfully!`,
		});
	};

	const commonDelete = (arr, setFn, id) => {
		const newArr = arr.filter((i) => i.id !== id);
		setFn([...newArr]);
		message.destroy();
		messageApi.open({
			type: 'success',
			content: 'Deleted successfully!',
		});
	};

	const handleAdd = (values, v, type) => {
		if (type === 'transactionHistory') {
			commonAdd(historyList, setHistoryList, values, v);
			return;
		}
		if (type === 'realEstate') commonAdd(realEstateList, setRealEstateList, values, v);
	};

	const handleDelete = (id, type) => {
		if (type === 'transactionHistory') {
			commonDelete(historyList, setHistoryList, id);
			return;
		}
		if (type === 'realEstate') commonDelete(realEstateList, setRealEstateList, id);
	};

	const { transactionHistoryColumns, realEstateColumns } = useColumnByZeeBeeLoan(
		handleAdd,
		handleDelete,
	);

	return (
		<div className="gap-y-8 flex flex-col w-full">
			<div className="text-2xl font-medium">
				Transaction history
				<p className="text-sm text-base-label">
					(Flips, Ground Up, Value Add Properties) Note: Please exclude all homes that
					were a primary residence
				</p>
			</div>
			{historyList?.length > 0 && (
				<NewTable
					dataSource={historyList}
					columns={transactionHistoryColumns}
					pagination={false}
					rowKey="id"
				/>
			)}
			<AddEditBtn onOk={handleAdd} commonType="transactionHistory">
				Add transaction history
			</AddEditBtn>
			<div className="text-2xl font-medium pt-6 border-t border-base-border">
				Real estate owned
				<p className="text-sm text-base-label">Properties held as Rentals</p>
			</div>
			{realEstateList?.length > 0 && (
				<NewTable
					dataSource={realEstateList}
					columns={realEstateColumns}
					pagination={false}
					rowKey="id"
				/>
			)}
			<AddEditBtn onOk={handleAdd} commonType="realEstate">
				Add real estate
			</AddEditBtn>
			<div className="pb-6"></div>
			{contextHolder}
		</div>
	);
};

export default Experience;
