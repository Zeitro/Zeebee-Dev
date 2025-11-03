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
			title: 'Entity type',
			dataIndex: 'address',
			key: 'address',
			ellipsis: true,
		},
		{
			title: 'EIN',
			dataIndex: 'ein',
			key: 'ein',
			ellipsis: true,
		},
		{
			title: 'Own percentage',
			dataIndex: 'state',
			key: 'own',
			ellipsis: true,
		},
		{
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
		},
		{
			title: 'Phone number',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
		},
		{
			title: 'Contact name',
			dataIndex: 'contact_name',
			key: 'contact_name',
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
