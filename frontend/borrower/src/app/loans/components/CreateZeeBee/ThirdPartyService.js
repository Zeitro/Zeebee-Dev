'use client';
import React, { useState } from 'react';
import NewTable from '@/components/NewTable';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
import { message } from 'antd';
import _ from 'lodash';
let paramsKey = 'loan_third_services';
const ThirdPartyService = ({ form }) => {
	const [list, setList] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();
	const commonAdd = (arr, setFn, values, key) => {
		let newArr = [];
		if (key === 'Edit') {
			const index = arr.findIndex((i) => i.id === values.id);
			arr.splice(index, 1, values);
			newArr = [...arr];
		} else {
			newArr = [...arr, { ...values, id: arr.length + 1 }];
		}
		setFn(newArr);
		const noIdArr = _.cloneDeep(newArr).map((i) => {
			delete i.id;
			return i;
		});
		form.setFieldsValue({ [paramsKey]: noIdArr });
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
			delete i.id;
			return i;
		});
		form.setFieldsValue({ [paramsKey]: noIdArr });
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
		</div>
	);
};

export default ThirdPartyService;
