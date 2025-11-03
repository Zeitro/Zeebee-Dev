import React from 'react';
import { TableOperation } from '@repo/shared-methods';
const useBorrowingEntity = () => {
	const columns = [
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			ellipsis: true,
		},

		{
			title: 'SSN',
			dataIndex: 'ssn',
			key: 'ssn',
			ellipsis: true,
		},
		{
			title: 'Phone number',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
		},
		{
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
		},
		{
			title: 'Estimated credit score',
			dataIndex: 'address',
			key: 'address',
			ellipsis: true,
		},
		{
			title: 'Primary address',
			dataIndex: 'state',
			key: 'own',
			ellipsis: true,
		},

		{
			title: '',
			dataIndex: 'operation',
			width: 50,
			align: 'left',
			fixed: 'right',
			render: (_, record) => (
				<TableOperation
					record={record}
					// getItems={(editItem, deleteItem, detailsItem) => {
					// 	const arr = [];

					// 	return arr;
					// }}
					onSelect={(key) => {}}
				/>
			),
		},
	];
	return { columns };
};
export default useBorrowingEntity;
