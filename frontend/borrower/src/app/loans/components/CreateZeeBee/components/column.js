'use client';
import FilterField from '@/components/FilterField';
import { AddEditBtn } from './AddEditBtn';
import { DeleteBtn } from '@/components/NewModal/BtnPopUpModal';
export const useColumnByZeeBeeLoan = () => {
	const beforeCommon = [
		{
			title: 'Property address',
			dataIndex: ['loan_property', 'property_base', 'address', 'address1'],
			key: 'address',
			ellipsis: true,
			sorter: true,
			sortDirections: ['ascend', 'descend'],
			defaultSortOrder: 'descend',
			fixed: 'left',
			minWidth: 200,
		},
		{
			title: 'City',
			dataIndex: 'city',
			key: 'city',
			ellipsis: true,
			width: 120,
		},
		{
			title: 'State',
			dataIndex: 'state',
			key: 'state',
			ellipsis: true,
			width: 70,
		},
		{
			title: 'Owner ship on title',
			dataIndex: 'owner_ship_on_title',
			key: 'owner_ship_on_title',
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
			dataIndex: 'purchase_price',
			key: 'purchase_price',
			ellipsis: true,
			width: 140,
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
			dataIndex: 'sales_price',
			key: 'sales_price',
			ellipsis: true,
			width: 140,
		},
		{
			title: 'Date sold',
			dataIndex: 'date_sold',
			key: 'date_sold',
			ellipsis: true,
			width: 140,
		},
		{
			title: 'Profit/Loan',
			dataIndex: 'profit_loan',
			key: 'profit_loan',
			ellipsis: true,
			width: 140,
		},
	];

	const getAction = (title, deleteTitle) => ({
		title: 'Action',
		dataIndex: 'action',
		key: 'action',
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
								>
									Edit
								</AddEditBtn>
							),
						},
						{
							filterType: 'button',
							children: <DeleteBtn title={deleteTitle}>Delete</DeleteBtn>,
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
			dataIndex: 'rehab_budget_construction_cost',
			key: 'rehab_budget_construction_cost',
			ellipsis: true,
			width: 240,
		},
		...afterCommon,
		getAction(
			'Edit transaction history',
			'Are you sure you want to delete this transaction history?',
		),
	];

	const realEstateColumns = [
		...beforeCommon,
		{
			title: 'Tenant occupied',
			dataIndex: 'tenant_occupied',
			key: 'tenant_occupied',
			ellipsis: true,
			width: 140,
		},
		...afterCommon,
		getAction('Edit real estate', 'Are you sure you want to delete this real state?'),
	];

	const thirdPartyServicesColumns = [
		{
			title: 'Service type',
			dataIndex: 'service_type',
			key: 'service_type',
			ellipsis: true,
		},
		{
			title: 'Name',
			dataIndex: 'Name',
			key: 'Name',
			ellipsis: true,
		},
		{
			title: 'Company',
			dataIndex: 'company',
			key: 'company',
			ellipsis: true,
		},
		{
			title: 'Phone',
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
		getAction(
			'Edit third party service',
			'Are you sure you want to delete this third party service?',
		),
	];
	return { transactionHistoryColumns, realEstateColumns, thirdPartyServicesColumns };
};
