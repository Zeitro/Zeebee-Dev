'use client';
import React, { useState } from 'react';
import NewTable from '@/components/NewTable';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
import { message } from 'antd';
const ThirdPartyService = () => {
	const [list, setList] = useState([]);
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
