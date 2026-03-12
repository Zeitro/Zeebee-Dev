'use client';
import React, { useMemo, useRef, useState } from 'react';
import { message, Space, Button, Alert } from 'antd';
import { DeleteBtn, FilterTable, PermissionSlot } from '@repo/shared-methods';
import { createRequest, formatTime } from '@repo/shared-methods/utils';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import MainDrawer from './components/Main';
import { userStatusColor } from '@repo/shared-methods/enum';
import { useAuth } from '@repo/shared-methods/hooks';
import { PERMISSIONS } from '@config/permissions';
const api = '/base_template/page';

const EditorTemplate = () => {
	const { hasPermission } = useAuth();
	const tableFilterRef = useRef();
	const { request } = createRequest();
	const [searchValue, setSearchValue] = useState('');
	const [type, setType] = useState(null);
	const [status, setStatus] = useState(null);
	const [loading, setLoading] = useState(false);
	const [open, setOpen] = useState(false);
	const [rowData, setRowData] = useState({});
	const [statusLoading, setStatusLoading] = useState(false);
	const [updateLoading, setUpdateLoading] = useState(false);

	const statusOptions = [
		{
			label: 'Active',
			value: 'Active',
		},
		{
			label: 'Disabled',
			value: 'Disabled',
		},
	];

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by any field',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
			},
			// {
			// 	filterType: 'select',
			// 	placeholder: 'System',
			// 	paramsKey: 'type',
			// 	value: type,
			// 	options: systemOptions,
			// 	onChange: (v) => setType(v),
			// },
			{
				filterType: 'select',
				placeholder: 'Status',
				paramsKey: 'status',
				value: status,
				options: statusOptions,
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
					setType(null);
					setTimeout(() => {
						tableFilterRef?.current?.handleSearch?.();
					}, 300);
				},
				colProps: {
					className: 'flex-1',
				},
			},
			hasPermission(PERMISSIONS.TEMPLATE.ADD, {
				filterType: 'button',
				text: 'Add template',
				icon: <PlusOutlined />,
				onClick: () => {
					setRowData({});
					setOpen(true);
				},
			}),
		],
		[searchValue, type, status, hasPermission],
	);

	// 删除模板
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/base_template/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			tableFilterRef?.current?.handleSearch?.();
			message.open({
				type: 'success',
				content: 'Deleted successfully！',
			});
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	// 改变状态
	const handleChangeStatus = async (id, status = 'Active') => {
		try {
			setStatusLoading(true);
			await request(`/base_template/${id}/status`, {
				method: 'PATCH',
				body: {
					status,
				},
			});
			setStatusLoading(false);
			tableFilterRef?.current?.getList?.();
			message.open({
				type: 'success',
				content: 'Template status updated successfully！',
			});
		} catch (error) {
			console.log(error);
			setStatusLoading(false);
		}
	};

	const handleUpdate = async (id) => {
		try {
			setUpdateLoading(true);
			await request(`/base_template/${id}/reset`, {
				method: 'POST',
			});
			setUpdateLoading(false);
			tableFilterRef?.current?.getList?.();
			message.open({
				type: 'success',
				content: 'Template force an updated successfully！',
			});
		} catch (error) {
			console.log(error);
			setUpdateLoading(false);
		}
	};

	return (
		<div className="flex-1 overflow-hidden w-full flex items-center justify-center py-4 px-4">
			<FilterTable
				ref={tableFilterRef}
				columns={[
					{
						title: 'Purpose',
						dataIndex: 'purpose',
						key: 'purpose',
						ellipsis: true,
						width: 240,
					},
					{
						title: 'Template name',
						dataIndex: 'name',
						key: 'name',
						ellipsis: true,
						render: (v) => v || '--',
						width: 240,
					},
					// {
					// 	title: 'Type',
					// 	dataIndex: 'type',
					// 	key: 'type',
					// 	ellipsis: true,
					// 	width: 120,
					// 	render: (v) => v || '--',
					// },
					{
						title: 'Description',
						dataIndex: 'description',
						key: 'description',
						ellipsis: true,
						render: (v) => v || '--',
					},
					{
						title: 'Created time',
						dataIndex: 'created_at',
						key: 'created_at',
						ellipsis: true,
						sorter: true,
						defaultSortOrder: 'descend',
						render: (v) => formatTime(v) || '--',
						width: 180,
					},
					{
						title: 'Status',
						dataIndex: 'status',
						key: 'status',
						ellipsis: true,
						fixed: 'right',
						width: 100,
						render: (v) => (
							<div
								className={'text-xs/5 px-2 border rounded-sm w-max'}
								style={{
									borderColor: userStatusColor[v]?.border,
									backgroundColor: userStatusColor[v]?.bg,
									color: userStatusColor[v]?.color,
								}}
							>
								{v}
							</div>
						),
					},
					{
						title: 'Action',
						dataIndex: 'operation',
						width: 260,
						align: 'left',
						fixed: 'right',
						render: (_, record) => {
							return (
								<Space>
									<PermissionSlot code={PERMISSIONS.TEMPLATE.EDIT}>
										<Button
											size="small"
											type="link"
											onClick={() => {
												setRowData(record);
												setOpen(true);
											}}
										>
											Edit
										</Button>
									</PermissionSlot>
									<PermissionSlot code={PERMISSIONS.TEMPLATE.CHANGE_STATUS}>
										{record?.status === 'Active' && (
											<DeleteBtn
												title="Are you sure you want to disable this template?"
												content={record.name}
												onOk={() =>
													handleChangeStatus(record.id, 'Disabled')
												}
												loading={statusLoading}
												isStopPropagation={false}
												okText="Disable"
												danger={false}
											>
												Disable
											</DeleteBtn>
										)}
										{record?.status === 'Disabled' && (
											<DeleteBtn
												title="Are you sure you want to enable this template?"
												content={record.name}
												onOk={() => handleChangeStatus(record.id, 'Active')}
												loading={statusLoading}
												isStopPropagation={false}
												okText="Enable"
												danger={false}
												okButtonProps={{ danger: false }}
											>
												Enable
											</DeleteBtn>
										)}
									</PermissionSlot>
									<PermissionSlot code={PERMISSIONS.TEMPLATE.UPDATE}>
										<DeleteBtn
											title="Are you sure you want to force a template update?"
											content={
												<div>
													<Alert
														message="Users will use the latest template."
														type="warning"
														className="!mb-4"
													/>
													<span>{record.name}</span>
												</div>
											}
											onOk={() => handleUpdate(record.id)}
											loading={updateLoading}
											isStopPropagation={false}
											okText="Update"
											danger={false}
											okButtonProps={{ danger: false }}
										>
											Update
										</DeleteBtn>
									</PermissionSlot>
									<PermissionSlot code={PERMISSIONS.TEMPLATE.DELETE}>
										<DeleteBtn
											title="Are you sure delete this editor template?"
											content={record.name}
											// isButton={false}
											onOk={() => handleDelete(record.id)}
											loading={loading}
											isStopPropagation={false}
										>
											Delete
										</DeleteBtn>
									</PermissionSlot>
								</Space>
							);
						},
					},
				]}
				filterList={filterList}
				scroll={{ x: 1300 }}
				// pagination={true}

				getData={(params) => {
					if (!api) return;
					const filter = {
						...(params?.filter || {}),
					};
					const payload = {
						method: 'POST',
						body: {
							...params,
							filter,
						},
					};
					return request(api, payload);
				}}
			/>
			<MainDrawer
				rowData={rowData}
				open={open}
				setOpen={(value) => {
					setOpen(value);
					setRowData({});
				}}
				successCallback={async () => {
					if (rowData?.id) {
						tableFilterRef?.current?.getList?.();
						return;
					}
					tableFilterRef?.current?.handleSearch?.();
				}}
			/>
		</div>
	);
};
export default EditorTemplate;
