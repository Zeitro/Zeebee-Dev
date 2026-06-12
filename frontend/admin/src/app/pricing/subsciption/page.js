'use client';

import React, { useMemo, useRef, useState } from 'react';
import { message, Popconfirm, Space } from 'antd';
import { BtnPopUpModal, CommonFormDate, FilterTable } from '@repo/shared-methods';
import { displayValue, formatDate } from '@repo/shared-methods/utils';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import SubscriptionFormModal from './components/SubscriptionFormModal';
import {
	formatSubscriptionType,
	getMockSubscriptionList,
	planNameOptions,
	SUBSCRIPTION_STATUS,
	subscriptionStatusColor,
	subscriptionStatusOptions,
	subscriptionTypeOptions,
	toggleMockSubscriptionStatus,
} from './constants';

// 接入后端后改为真实路径，并在 getData 中调用 request
// const api = '/subscription/page';

const Subscription = () => {
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [planName, setPlanName] = useState(null);
	const [subscriptionType, setSubscriptionType] = useState(null);
	const [status, setStatus] = useState(null);
	const [startDateRange, setStartDateRange] = useState(null);
	const [toggleLoading, setToggleLoading] = useState(false);
	const [subscriptionModalOpen, setSubscriptionModalOpen] = useState(false);
	const [subscriptionRowData, setSubscriptionRowData] = useState({});

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by subscription no. or customer name',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
			},
			{
				filterType: 'select',
				placeholder: 'Subscription type',
				paramsKey: 'subscription_type',
				value: subscriptionType,
				options: subscriptionTypeOptions,
				onChange: (v) => setSubscriptionType(v),
			},
			{
				filterType: 'select',
				placeholder: 'Status',
				paramsKey: 'status',
				value: status,
				options: subscriptionStatusOptions,
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
					setPlanName(null);
					setSubscriptionType(null);
					setStatus(null);
					setStartDateRange(null);
					tableFilterRef?.current?.handleSearch?.({
						filter: {
							keyword: '',
							plan_name: null,
							subscription_type: null,
							status: null,
							start_date: null,
						},
					});
				},
				colProps: {
					className: 'flex-1',
				},
			},
			{
				filterType: 'button',
				children: (
					<BtnPopUpModal
						getContent={(setOpen) => (
							<SubscriptionFormModal
								setOpen={setOpen}
								successCallback={() => tableFilterRef?.current?.handleSearch?.()}
							/>
						)}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add subscription
					</BtnPopUpModal>
				),
			},
		],
		[searchValue, planName, subscriptionType, status, startDateRange],
	);

	const handleEdit = (record) => {
		setSubscriptionRowData(record);
		setSubscriptionModalOpen(true);
	};

	const handleToggleStatus = async (record, nextStatus) => {
		try {
			setToggleLoading(true);
			// TODO: 接入 PATCH /subscription/:id/status
			await Promise.resolve();
			toggleMockSubscriptionStatus(record.id, nextStatus);
			tableFilterRef?.current?.handleSearch?.();
			message.success(
				`Subscription ${nextStatus === SUBSCRIPTION_STATUS.SUSPENDED ? 'suspended' : 'activated'} successfully.`,
			);
		} catch (error) {
			console.log(error);
		} finally {
			setToggleLoading(false);
		}
	};

	const canSuspend = (recordStatus) =>
		[
			SUBSCRIPTION_STATUS.ACTIVE,
			SUBSCRIPTION_STATUS.TRIALING,
			SUBSCRIPTION_STATUS.EXPIRING,
		].includes(recordStatus);

	const canActivate = (recordStatus) => recordStatus === SUBSCRIPTION_STATUS.SUSPENDED;

	return (
		<div className="flex-1 overflow-hidden w-full flex items-center justify-center py-4 px-4">
			<FilterTable
				ref={tableFilterRef}
				filterList={filterList}
				scroll={{ x: 1400 }}
				columns={[
					{
						title: 'Subscription no.',
						dataIndex: 'subscription_no',
						key: 'subscription_no',
						ellipsis: true,
						minWidth: 140,
						render: (v) => displayValue(v),
					},
					{
						title: 'Customer name',
						dataIndex: 'customer_name',
						key: 'customer_name',
						ellipsis: true,
						minWidth: 160,
						render: (v) => displayValue(v),
					},
					{
						title: 'Plan name',
						dataIndex: 'plan_name',
						key: 'plan_name',
						width: 120,
						render: (v) => displayValue(v),
					},
					{
						title: 'Subscription type',
						dataIndex: 'subscription_type',
						key: 'subscription_type',
						width: 140,
						render: (v) => displayValue(v, formatSubscriptionType),
					},
					{
						title: 'Start date',
						dataIndex: 'start_date',
						key: 'start_date',
						width: 120,
						render: (v) => displayValue(v, formatDate),
					},
					{
						title: 'Expiry date',
						dataIndex: 'expiry_date',
						key: 'expiry_date',
						width: 120,
						render: (v) => displayValue(v, formatDate),
					},
					{
						title: 'Status',
						dataIndex: 'status',
						key: 'status',
						width: 110,
						render: (v) => {
							const colors = subscriptionStatusColor[v];
							if (!colors) return displayValue(v);
							return (
								<div
									className="text-xs/5 px-2 border rounded-sm w-max capitalize"
									style={{
										borderColor: colors.border,
										backgroundColor: colors.bg,
										color: colors.color,
									}}
								>
									{v}
								</div>
							);
						},
					},
					{
						title: 'Action',
						key: 'action',
						fixed: 'right',
						width: 140,
						render: (_, record) => (
							<Space>
								<span
									className="text-base-primary cursor-pointer"
									onClick={() => handleEdit(record)}
								>
									Edit
								</span>
								{canSuspend(record.status) && (
									<Popconfirm
										title="Are you sure you want to suspend this subscription?"
										onConfirm={() =>
											handleToggleStatus(
												record,
												SUBSCRIPTION_STATUS.SUSPENDED,
											)
										}
										okButtonProps={{ loading: toggleLoading }}
									>
										<span className="text-base-error cursor-pointer">
											Suspend
										</span>
									</Popconfirm>
								)}
								{canActivate(record.status) && (
									<Popconfirm
										title="Are you sure you want to activate this subscription?"
										onConfirm={() =>
											handleToggleStatus(record, SUBSCRIPTION_STATUS.ACTIVE)
										}
										okButtonProps={{ loading: toggleLoading }}
									>
										<span className="text-base-primary cursor-pointer">
											Active
										</span>
									</Popconfirm>
								)}
							</Space>
						),
					},
				]}
				getData={(params) => {
					// 有 api 时: createRequest + request(api, { method: 'POST', body: { ...params, filter: params?.filter || {} } })
					return getMockSubscriptionList(params);
				}}
			/>
			{subscriptionModalOpen && (
				<SubscriptionFormModal
					rowData={subscriptionRowData}
					setOpen={(value) => {
						setSubscriptionModalOpen(value);
						if (!value) setSubscriptionRowData({});
					}}
					successCallback={() => tableFilterRef?.current?.handleSearch?.()}
				/>
			)}
		</div>
	);
};

export default Subscription;
