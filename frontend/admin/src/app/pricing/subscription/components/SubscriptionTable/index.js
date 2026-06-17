'use client';

import React, { useEffect, useMemo, useRef, useState } from 'react';
import { FilterTable } from '@repo/shared-methods';
import {
	createRequest,
	displayValue,
	eventBus,
	formatAmount,
	formatDate,
	transformArrayStructure,
} from '@repo/shared-methods/utils';
import { ReloadOutlined } from '@ant-design/icons';
import ChangeStatus from '../ChangeStatus';
import EditableExpireDate from '../EditableExpireDate';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { useOptionsStore } from '@repo/shared-methods/store';

const api = '/organization_subscription/page';

/**
 * 单端（Borrower / Lender / Admin）组织订阅列表；type 写入 filter.type。
 *
 * @param {object} props
 * @param {'Borrower'|'Lender'|'Admin'} [props.type]
 * @param {string} [props.eventBusPrefix] 与父级 Add 弹窗联动刷新
 * @param {string} [props.filterWrapClassName]
 */
const SubscriptionTable = ({
	type = 'Borrower',
	eventBusPrefix,
	filterWrapClassName = 'px-4 pb-4',
}) => {
	const tableFilterRef = useRef();
	const { request } = useMemo(() => createRequest(), []);
	const [searchValue, setSearchValue] = useState('');
	const [status, setStatus] = useState([]);
	useGetOptions([
		{ key: 'organization_subscription_status', valueName: 'subscriptionStatusOptions' },
	]);
	const { subscriptionStatusOptions } = useOptionsStore();

	useEffect(() => {
		if (!type || !eventBusPrefix) return;
		const searchList = () => tableFilterRef?.current?.handleSearch?.();
		eventBus.on(`${eventBusPrefix}${type}`, searchList);
		return () => eventBus.off(`${eventBusPrefix}${type}`, searchList);
	}, [type, eventBusPrefix]);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by ID, Plan name or organization name',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
				style: { width: 400 },
			},
			{
				filterType: 'treeSelect',
				label: 'Status',
				placeholder: 'Status',
				paramsKey: 'status',
				value: status,
				treeData: transformArrayStructure({
					sourceArray: subscriptionStatusOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
				onChange: (v) => setStatus(v || []),
			},
			{
				filterType: 'button',
				onClick: () => {
					tableFilterRef?.current?.handleSearch?.();
				},
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					setSearchValue('');
					setStatus([]);
					tableFilterRef?.current?.handleSearch?.({
						filter: {
							keyword: '',
							status: [],
						},
					});
				},
			},
		],
		[searchValue, status, subscriptionStatusOptions],
	);

	return (
		<FilterTable
			key={type}
			ref={tableFilterRef}
			filterList={filterList}
			filterWrapClassName={filterWrapClassName}
			scroll={{ x: 1300 }}
			columns={[
				{
					title: 'ID',
					dataIndex: 'code',
					key: 'code',
					ellipsis: true,
					minWidth: 220,
					render: (v) => displayValue(v),
				},
				{
					title: 'Plan name',
					dataIndex: 'pricing_plan_name',
					key: 'pricing_plan_name',
					ellipsis: true,
					minWidth: 140,
					render: (v) => displayValue(v),
				},
				{
					title: 'Organization name',
					dataIndex: 'organization_name',
					key: 'organization_name',
					ellipsis: true,
					minWidth: 160,
					render: (v) => displayValue(v),
				},
				{
					title: 'Status',
					dataIndex: 'status',
					key: 'status',
					width: 150,
					render: (v, record) => (
						<ChangeStatus
							data={record}
							successCallback={() => tableFilterRef?.current?.getList?.()}
						/>
					),
				},
				{
					title: 'Expire date',
					dataIndex: 'expire_date',
					key: 'expire_date',
					width: 140,
					sorter: true,
					render: (v, record) => (
						<EditableExpireDate
							record={record}
							successCallback={() => tableFilterRef?.current?.getList?.()}
						/>
					),
				},
				{
					title: (
						<div>
							<span>Pricing plan policy type</span>
							<span className="text-xs text-gray-500"> (Actual/Original)</span>
						</div>
					),
					dataIndex: 'pricing_plan_policy_type',
					key: 'pricing_plan_policy_type',
					width: 300,
					render: (v, record) => (
						<div>
							<span>{v}</span>
							<span className="text-xs text-gray-500 pl-1">
								({formatAmount(record.actual_amount)} /
								{formatAmount(record.original_amount)})
							</span>
						</div>
					),
				},
				{
					title: 'Auto renew',
					dataIndex: 'is_auto_renew',
					key: 'is_auto_renew',
					width: 120,
					render: (v) => displayValue(v === true ? 'Yes' : v === false ? 'No' : null),
				},
				{
					title: 'Created date',
					dataIndex: 'created_at',
					key: 'created_at',
					width: 160,
					sorter: true,
					defaultSortOrder: 'descend',
					render: (v) => displayValue(v, formatDate),
				},
			]}
			getData={(params) => {
				const { sorter, ...rest } = params || {};
				const filter = {
					...(params?.filter || {}),
					type,
					status: params?.filter?.status ?? status ?? [],
				};
				const body = {
					...rest,
					filter,
					sorter: sorter?.column ? sorter : { column: '', sort: '' },
				};
				return request(api, { method: 'POST', body });
			}}
		/>
	);
};

export default SubscriptionTable;
