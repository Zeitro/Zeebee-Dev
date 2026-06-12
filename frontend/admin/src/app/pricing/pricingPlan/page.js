'use client';

import React, { useMemo, useRef, useState } from 'react';
import { message, Space } from 'antd';
import { BtnPopUpModal, DeleteBtn, FilterTable } from '@repo/shared-methods';
import { displayValue, formatAmount } from '@repo/shared-methods/utils';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import PlanFormModal from './components/PlanFormModal';
import { getMockPricingPlanList, planStatusColor, planStatusOptions } from './constants';

// 接入后端后改为真实路径，并在 getData 中调用 request
// const api = '/pricing_plan/page';

const PricingPlan = () => {
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [status, setStatus] = useState(null);
	const [deleteLoading, setDeleteLoading] = useState(false);
	const [planModalOpen, setPlanModalOpen] = useState(false);
	const [planRowData, setPlanRowData] = useState({});

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
					setStatus(null);
					tableFilterRef?.current?.handleSearch?.({
						filter: {
							keyword: '',
							status: null,
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
							<PlanFormModal
								setOpen={setOpen}
								successCallback={() => tableFilterRef?.current?.handleSearch?.()}
							/>
						)}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add plan
					</BtnPopUpModal>
				),
			},
		],
		[searchValue, status],
	);

	const handleDelete = async (record) => {
		try {
			setDeleteLoading(true);
			// TODO: 接入 DELETE /pricing_plan/:id
			await Promise.resolve();
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

	const handlePermission = (record) => {
		// TODO: 打开权限配置
		message.info(`Permission for ${record.name}`);
	};

	return (
		<div className="flex-1 overflow-hidden w-full flex items-center justify-center py-4 px-4">
			<FilterTable
				ref={tableFilterRef}
				filterList={filterList}
				scroll={{ x: 1400 }}
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
						title: 'Description',
						dataIndex: 'description',
						key: 'description',
						ellipsis: true,
						minWidth: 240,
						render: (v) => displayValue(v),
					},
					{
						title: 'Monthly price',
						dataIndex: 'monthly_price',
						key: 'monthly_price',
						width: 130,
						render: (v) => displayValue(v, formatAmount),
					},
					{
						title: 'Yearly price',
						dataIndex: 'yearly_price',
						key: 'yearly_price',
						width: 120,
						render: (v) => displayValue(v, formatAmount),
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
						width: 100,
						render: (v) => {
							const colors = planStatusColor[v];
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
						title: 'Permission',
						key: 'permission',
						width: 110,
						render: (_, record) => (
							<span
								className="text-base-primary cursor-pointer"
								onClick={() => handlePermission(record)}
							>
								Permission
							</span>
						),
					},
					{
						title: 'Action',
						key: 'action',
						fixed: 'right',
						width: 120,
						render: (_, record) => (
							<Space>
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
					// 有 api 时: createRequest + request(api, { method: 'POST', body: { ...params, filter: params?.filter || {} } })
					return getMockPricingPlanList(params);
				}}
			/>
			{planModalOpen && (
				<PlanFormModal
					rowData={planRowData}
					setOpen={(value) => {
						setPlanModalOpen(value);
						if (!value) setPlanRowData({});
					}}
					successCallback={() => tableFilterRef?.current?.handleSearch?.()}
				/>
			)}
		</div>
	);
};

export default PricingPlan;
