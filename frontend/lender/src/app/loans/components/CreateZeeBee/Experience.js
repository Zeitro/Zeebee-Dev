'use client';
import React, { useEffect, useState } from 'react';
import { NewTable } from '@repo/shared-methods';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
import { Form, message } from 'antd';
import _ from 'lodash';
const paramsKeyInfo = {
	transactionHistory: 'transaction_experiences',
	realEstate: 'real_estate_experiences',
};
const paramsTypeInfo = {
	transactionHistory: 'Transaction',
	realEstate: 'Real Estate',
};
const Experience = ({ form, rowData, setIsChange, changeFieldsValue = () => {} }) => {
	const [historyList, setHistoryList] = useState([]);
	const [realEstateList, setRealEstateList] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();

	useEffect(() => {
		if (rowData?.transaction_experiences?.length) {
			setHistoryList(rowData.transaction_experiences || []);
		}
		if (rowData?.real_estate_experiences?.length) {
			setRealEstateList(rowData.real_estate_experiences || []);
		}
	}, [rowData]);

	const commonAdd = (arr, setFn, values, key, type) => {
		let newArr = [];
		if (key === 'Edit') {
			const index = arr.findIndex((i) => i.id === values.id);
			arr.splice(index, 1, values);
			newArr = [...arr];
		} else {
			// 因为后端的transactionHistory和realEstate的数据在一张表里面，所以要区分type，is_tenant_occupied和budget_or_cost没有值时传默认值
			newArr = [
				{
					is_tenant_occupied: false,
					budget_or_cost: 0,
					...values,
					id: `${arr.length + 1}Custom`, // 加Custom为了区分自定义的id，草稿时后端会返回id并且更新草稿需要传id
					type: paramsTypeInfo[type],
				},
				...arr,
			];
		}
		setFn(newArr);
		const noIdArr = _.cloneDeep(newArr).map((i) => {
			if (typeof i.id === 'string') {
				delete i.id;
			}
			return i;
		});
		form.setFieldsValue({ [paramsKeyInfo[type]]: noIdArr });
		changeFieldsValue(noIdArr, paramsKeyInfo[type]);
		setIsChange(true);
		message.destroy();
		messageApi.open({
			type: 'success',
			content: `${key}ed successfully!`,
		});
	};

	const commonDelete = (arr, setFn, id, type) => {
		const newArr = arr.filter((i) => i.id !== id);
		setFn([...newArr]);
		const noIdArr = _.cloneDeep(newArr).map((i) => {
			if (typeof i.id === 'string') {
				delete i.id;
			}
			return i;
		});
		let paramsKey =
			type === 'transactionHistory' ? 'transaction_experiences' : 'real_estate_experiences';
		form.setFieldsValue({ [paramsKey]: noIdArr });
		changeFieldsValue(noIdArr, paramsKey);
		setIsChange(true);
		message.destroy();
		messageApi.open({
			type: 'success',
			content: 'Deleted successfully!',
		});
	};

	const handleAdd = (values, v, type) => {
		if (type === 'transactionHistory') {
			commonAdd(historyList, setHistoryList, values, v, type);
			return;
		}
		if (type === 'realEstate') commonAdd(realEstateList, setRealEstateList, values, v, type);
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
			<Form.Item name="transaction_experiences" noStyle></Form.Item>
			<Form.Item name="real_estate_experiences" noStyle></Form.Item>
		</div>
	);
};

export default Experience;
