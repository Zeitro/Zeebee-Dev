'use client';
import React, { useEffect, useState } from 'react';
import NewTable from '@/components/NewTable';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
import { Form, message } from 'antd';
import _ from 'lodash';
let paramsKey = 'loan_third_services';
const ThirdPartyService = ({ form, rowData, setIsChange }) => {
	const [list, setList] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();

	useEffect(() => {
		if (rowData?.[paramsKey]?.length) {
			setList(rowData?.[paramsKey] || []);
		}
	}, [rowData]);

	const commonAdd = (arr, setFn, values, key) => {
		let newArr = [];
		if (key === 'Edit') {
			const index = arr.findIndex((i) => i.id === values.id);
			arr.splice(index, 1, values);
			newArr = [...arr];
		} else {
			newArr = [...arr, { ...values, id: `${arr.length + 1}Custom` }]; // 加Custom为了区分自定义的id，草稿时后端会返回id并且更新草稿需要传id
		}
		setFn(newArr);
		const noIdArr = _.cloneDeep(newArr).map((i) => {
			if (typeof i.id === 'string') {
				delete i.id;
			}
			return i;
		});
		form.setFieldsValue({ [paramsKey]: noIdArr });
		setIsChange(true);
		message.destroy();
		messageApi.open({
			type: 'success',
			content: `${key}ed successfully!`,
		});
	};

	const commonDelete = (arr, setFn, id) => {
		const newArr = arr.filter((i) => i.id !== id);
		setFn([...newArr]);
		const noIdArr = _.cloneDeep(newArr).map((i) => {
			if (typeof i.id === 'string') {
				delete i.id;
			}
			return i;
		});
		form.setFieldsValue({ [paramsKey]: noIdArr });
		setIsChange(true);
		message.destroy();
		messageApi.open({
			type: 'success',
			content: 'Deleted successfully!',
		});
	};

	const handleAdd = (values, v) => commonAdd(list, setList, values, v);

	const handleDelete = (id) => commonDelete(list, setList, id);

	const { thirdPartyServicesColumns } = useColumnByZeeBeeLoan(handleAdd, handleDelete);

	return (
		<div className="gap-y-8 flex flex-col w-full">
			<div className="text-2xl font-medium">Third party services</div>
			{list?.length > 0 && (
				<NewTable
					dataSource={list}
					columns={thirdPartyServicesColumns}
					pagination={false}
					rowKey="id"
				/>
			)}
			<AddEditBtn onOk={handleAdd} commonType="thirdPartyService">
				Add third party service
			</AddEditBtn>
			<div className="pb-6"></div>
			{contextHolder}
			<Form.Item name={paramsKey} noStyle></Form.Item>
		</div>
	);
};

export default ThirdPartyService;
