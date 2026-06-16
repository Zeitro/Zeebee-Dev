'use client';

import React, { useEffect, useMemo, useRef, useState } from 'react';
import { message, Space } from 'antd';
import { DeleteBtn, FilterTable } from '@repo/shared-methods';
import {
	createRequest,
	displayValue,
	EMPTY_DISPLAY,
	eventBus,
	formatAmount,
} from '@repo/shared-methods/utils';
import { ReloadOutlined } from '@ant-design/icons';
import PlanFormDrawer from '../PlanFormDrawer';
import MenuPermissionDrawer from '../MenuPermissionDrawer';
import { planStatusOptions, planTypeOptions } from '../../constants';
import ChangeStatus from '../ChangeStatus';
const api = '/pricing_plan/page';

/**
 * 从 policies 中取指定计费周期，展示 actual_amount/original_amount。
 *
 * @param {Array<{ type: string, actual_amount?: number, original_amount?: number }>} policies
 * @param {'Monthly'|'Yearly'} policyType API 字段 policies[].type
 * @returns {string}
 */
const renderPolicyPrice = (policies, policyType) => {
	const policy = policies?.find((item) => item.type === policyType);
	if (!policy) return displayValue(null);
	const actual = displayValue(policy.actual_amount, formatAmount);
	const original = displayValue(policy.original_amount, formatAmount);
	if (actual === EMPTY_DISPLAY && original === EMPTY_DISPLAY) return EMPTY_DISPLAY;
	if (original === EMPTY_DISPLAY) return actual;
	if (actual === EMPTY_DISPLAY) return original;
	return `${actual} / ${original}`;
};

/**
 * 单端（Borrower / Lender / Admin）定价计划列表；type 写入 filter 与新增 body。
 *
 * @param {object} props
 * @param {'Borrower'|'Lender'|'Admin'} [props.type]
 * @param {string} [props.eventBusPrefix] 与父级 Add 弹窗联动刷新，如 searchPricingPlan-
 * @param {string} [props.filterWrapClassName]
 */
const PricingPlanTable = ({
	type = 'Borrower',
	eventBusPrefix,
	filterWrapClassName = 'px-4 pb-4',
}) => {
	const tableFilterRef = useRef();
	const { request } = useMemo(() => createRequest(), []);
	const [searchValue, setSearchValue] = useState('');
	const [planType, setPlanType] = useState(null);
	const [status, setStatus] = useState(null);
	const [deleteLoading, setDeleteLoading] = useState(false);
	const [planModalOpen, setPlanModalOpen] = useState(false);
	const [planRowData, setPlanRowData] = useState({});
	const [addMenuModalOpen, setAddMenuModalOpen] = useState(false);
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
				placeholder: 'Search by plan name',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
			},
			{
				filterType: 'select',
				placeholder: 'Plan type',
				paramsKey: 'plan_type',
				value: planType,
				options: planTypeOptions,
				onChange: (v) => setPlanType(v),
			},
			{
				filterType: 'select',
				placeholder: 'Status',
				paramsKey: 'status',
				value: status,
				options: planStatusOptions,
				onChange: (v) => setStatus(v),
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
					setPlanType(null);
					setStatus(null);
					tableFilterRef?.current?.handleSearch?.({
						filter: {
							keyword: '',
							plan_type: null,
							status: null,
						},
					});
				},
				colProps: {
					className: 'flex-1',
				},
			},
		],
		[searchValue, planType, status],
	);

	const handleDelete = async (record) => {
		try {
			setDeleteLoading(true);
			await request(`/pricing_plan/${record.id}`, { method: 'DELETE' });
			setDeleteLoading(false);
			tableFilterRef?.current?.handleSearch?.();
			message.success(`Deleted ${record.name} successfully!`);
		} catch (error) {
			console.log(error);
			setDeleteLoading(false);
		}
	};

	const handleEdit = (record) => {
		setPlanRowData(record);
		setPlanModalOpen(true);
	};

	const handleAddMenu = (record) => {
		setPlanRowData(record);
		setAddMenuModalOpen(true);
	};

	return (
		<>
			<FilterTable
				key={type}
				ref={tableFilterRef}
				filterList={filterList}
				filterWrapClassName={filterWrapClassName}
				scroll={{ x: 1700 }}
				columns={[
					{
						title: 'Plan name',
						dataIndex: 'name',
						key: 'name',
						ellipsis: true,
						minWidth: 120,
						render: (v) => displayValue(v),
					},
					{
						title: 'Plan subtitle',
						dataIndex: 'subtitle',
						key: 'subtitle',
						ellipsis: true,
						minWidth: 180,
						render: (v) => displayValue(v),
					},
					{
						title: 'Plan type',
						dataIndex: 'plan_type',
						key: 'plan_type',
						ellipsis: true,
						width: 120,
						render: (v) => v,
					},
					{
						title: (
							<div>
								<span>Monthly price</span>
								<span className="text-xs text-gray-500"> (Actual/Original)</span>
							</div>
						),
						dataIndex: 'policies',
						key: 'monthly_price',
						width: 220,
						render: (policies) => renderPolicyPrice(policies, 'Monthly'),
					},
					{
						title: (
							<div>
								<span>Quarterly price</span>
								<span className="text-xs text-gray-500"> (Actual/Original)</span>
							</div>
						),
						dataIndex: 'policies',
						key: 'quarterly_price',
						width: 220,
						render: (policies) => renderPolicyPrice(policies, 'Quarterly'),
					},
					{
						title: (
							<div>
								<span>Yearly price</span>
								<span className="text-xs text-gray-500"> (Actual/Original)</span>
							</div>
						),
						dataIndex: 'policies',
						key: 'yearly_price',
						width: 220,
						render: (policies) => renderPolicyPrice(policies, 'Yearly'),
					},
					{
						title: 'Trial days',
						dataIndex: 'trial_days',
						key: 'trial_days',
						width: 100,
						render: (v) => displayValue(v),
					},
					{
						title: 'Status',
						dataIndex: 'status',
						key: 'status',
						width: 150,
						render: (v, record) => {
							return (
								<ChangeStatus
									data={record}
									successCallback={() => tableFilterRef?.current?.getList?.()}
								/>
							);
						},
					},
					{
						title: 'Action',
						key: 'action',
						fixed: 'right',
						width: 260,
						render: (_, record) => (
							<Space>
								<span
									className="text-base-primary cursor-pointer pr-2"
									onClick={() => handleAddMenu(record)}
								>
									Menu permission
								</span>
								<span
									className="text-base-primary cursor-pointer"
									onClick={() => handleEdit(record)}
								>
									Edit
								</span>
								<DeleteBtn
									title="Are you sure you want to delete this pricing plan?"
									content={record.name}
									onOk={() => handleDelete(record)}
									loading={deleteLoading}
									isStopPropagation={false}
								>
									Delete
								</DeleteBtn>
							</Space>
						),
					},
				]}
				getData={(params) => {
					const { sorter, ...rest } = params || {};
					const filter = {
						...(params?.filter || {}),
						type,
					};
					const body = {
						...rest,
						filter,
						sort: sorter?.column ? sorter : { column: '', sort: '' },
					};
					return request(api, { method: 'POST', body });
				}}
			/>
			{planModalOpen && (
				<PlanFormDrawer
					rowData={planRowData}
					setOpen={(value) => {
						setPlanModalOpen(value);
						if (!value) setPlanRowData({});
					}}
					successCallback={() => tableFilterRef?.current?.handleSearch?.()}
				/>
			)}
			<MenuPermissionDrawer
				open={addMenuModalOpen}
				setOpen={(value) => {
					setAddMenuModalOpen(value);
					if (!value) setPlanRowData({});
				}}
				rowData={planRowData}
				systemType={type}
				successCallback={() => tableFilterRef?.current?.handleSearch?.()}
			/>
		</>
	);
};

export default PricingPlanTable;
