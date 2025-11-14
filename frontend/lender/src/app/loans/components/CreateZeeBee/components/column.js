'use client';
import { AddEditBtn } from './AddEditBtn';
import { DeleteBtn, FilterField } from '@repo/shared-methods';
import { formatAmount, formatPhone } from '@repo/shared-methods/utils';
export const useColumnByZeeBeeLoan = (handleAdd = () => {}, handleDelete = () => {}) => {
	const beforeCommon = [
		{
			title: 'Property address',
			dataIndex: ['address', 'address1'],
			key: 'address',
			ellipsis: true,
			fixed: 'left',
			minWidth: 200,
		},
		{
			title: 'City',
			dataIndex: ['address', 'city'],
			key: 'city',
			ellipsis: true,
			width: 120,
		},
		{
			title: 'State',
			dataIndex: ['address', 'state'],
			key: 'state',
			ellipsis: true,
			width: 70,
		},
		{
			title: 'Owner ship on title',
			dataIndex: 'ownership_on_title',
			key: 'ownership_on_title',
			ellipsis: true,
			width: 160,
		},
		{
			title: 'Property type',
			dataIndex: 'property_type',
			key: 'property_type',
			ellipsis: true,
			width: 160,
		},
		{
			title: 'Purchase price',
			dataIndex: 'purchase_amount',
			key: 'purchase_amount',
			ellipsis: true,
			width: 140,
			render: (v) => formatAmount(v),
		},
		{
			title: 'Purchase date',
			dataIndex: 'purchase_date',
			key: 'purchase_date',
			ellipsis: true,
			width: 140,
		},
	];

	const afterCommon = [
		{
			title: 'Sales price',
			dataIndex: 'sold_amount',
			key: 'sold_amount',
			ellipsis: true,
			width: 140,
			render: (v) => formatAmount(v),
		},
		{
			title: 'Date sold',
			dataIndex: 'sold_date',
			key: 'sold_date',
			ellipsis: true,
			width: 140,
		},
		{
			title: 'Profit/Loss',
			dataIndex: 'profit',
			key: 'profit',
			ellipsis: true,
			width: 140,
			render: (v) => formatAmount(v),
		},
	];

	const getAction = (type, title, deleteTitle) => ({
		title: 'Action',
		key: 'operation',
		fixed: 'right',
		ellipsis: true,
		width: 130,
		render: (v) => {
			return (
				<FilterField
					gutter={0}
					list={[
						{
							filterType: 'button',
							children: (
								<AddEditBtn
									type="link"
									size="small"
									rowData={v}
									icon={null}
									title={title}
									commonType={type}
									onOk={handleAdd}
								>
									Edit
								</AddEditBtn>
							),
						},
						{
							filterType: 'button',
							children: (
								<DeleteBtn
									title={deleteTitle}
									onOk={() => handleDelete(v.id, type)}
								>
									Delete
								</DeleteBtn>
							),
						},
					]}
				/>
			);
		},
	});
	const transactionHistoryColumns = [
		...beforeCommon,
		{
			title: 'Rehab budget/Construction cost',
			dataIndex: 'budget_or_cost',
			key: 'budget_or_cost',
			ellipsis: true,
			width: 240,
		},
		...afterCommon,
		getAction(
			'transactionHistory',
			'Edit transaction history',
			'Are you sure you want to delete this transaction history?',
		),
	];

	const realEstateColumns = [
		...beforeCommon,
		{
			title: 'Tenant occupied',
			dataIndex: 'is_tenant_occupied',
			key: 'is_tenant_occupied',
			ellipsis: true,
			width: 140,
		},
		...afterCommon,
		getAction(
			'realEstate',
			'Edit real estate',
			'Are you sure you want to delete this real state?',
		),
	];

	const thirdPartyServicesColumns = [
		{
			title: 'Service type',
			dataIndex: 'type',
			key: 'type',
			ellipsis: true,
			minWidth: 160,
		},
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			ellipsis: true,
			minWidth: 160,
		},
		{
			title: 'Company',
			dataIndex: 'company',
			key: 'company',
			ellipsis: true,
			minWidth: 160,
		},
		{
			title: 'Phone',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
			minWidth: 160,
			render: (v) => formatPhone(v),
		},
		{
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
			minWidth: 160,
		},
		getAction(
			'thirdPartyService',
			'Edit third party service',
			'Are you sure you want to delete this third party service?',
		),
	];
	return { transactionHistoryColumns, realEstateColumns, thirdPartyServicesColumns };
};
