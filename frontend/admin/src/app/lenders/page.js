'use client';
import React, { useMemo, useRef, useState } from 'react';
import { message, Space, Image } from 'antd';
import { BtnPopUpModal, DeleteBtn, FilterTable, PermissionSlot } from '@repo/shared-methods';
import { createRequest, formatPhone, formatTime } from '@repo/shared-methods/utils';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import MainModal from './components/Main';
import { userStatusColor } from '@repo/shared-methods/enum';
import { useAuth } from '@repo/shared-methods/hooks';
import { PERMISSIONS } from '@config/permissions';
const api = '/organization/page';
const Lenders = () => {
	const tableFilterRef = useRef();
	const { request } = createRequest();
	const [searchValue, setSearchValue] = useState('');
	const [loading, setLoading] = useState(false);
	const [statusLoading, setStatusLoading] = useState(false);
	const { hasPermission, getActiveByPermission } = useAuth();

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by any field',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
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
					tableFilterRef?.current?.handleSearch?.({
						filter: {
							keyword: '',
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
							<MainModal
								setOpen={setOpen}
								successCallback={async () =>
									tableFilterRef?.current?.handleSearch?.()
								}
							/>
						)}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add Lender org
					</BtnPopUpModal>
				),
				hidden: !hasPermission(PERMISSIONS.LENDERS.ADD),
			},
		],
		[searchValue, hasPermission],
	);
	// 改变状态
	const handleChangeStatus = async (id, status = 'Active') => {
		try {
			setStatusLoading(true);
			await request(`/organization/${id}/status`, {
				method: 'PATCH',
				body: {
					status,
				},
			});
			setStatusLoading(false);
			tableFilterRef?.current?.getList?.();
			message.open({
				type: 'success',
				content: 'Status updated successfully！',
			});
		} catch (error) {
			console.log(error);
			setStatusLoading(false);
		}
	};

	// 删除用户
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/organization/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			tableFilterRef?.current?.handleSearch?.();
			message.open({
				type: 'success',
				content: 'deleted successfully！',
			});
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	return (
		<div className="flex-1 overflow-hidden w-full flex items-center justify-center py-4 px-4">
			<FilterTable
				ref={tableFilterRef}
				columns={[
					{
						title: 'Logo',
						dataIndex: 'logo',
						key: 'logo',
						ellipsis: true,
						render: (v) => (
							<Image
								src={v || null}
								width={40}
								height={40}
								alt=""
								preview={false}
								fallback="/images/noOrgLogo.svg"
								className="object-contain rounded-lg"
							/>
						),
					},
					{
						title: 'Organization name',
						dataIndex: 'name',
						key: 'name',
						ellipsis: true,
						render: (v) => v || '--',
					},
					{
						title: 'Email account',
						dataIndex: 'owner_email',
						key: 'owner_email',
						ellipsis: true,
						render: (v) => v || '--',
					},
					{
						title: 'Contact name',
						dataIndex: 'contact_name',
						key: 'contact_name',
						ellipsis: true,
						render: (v) => v || '--',
					},
					{
						title: 'Contact phone',
						dataIndex: 'contact_phone',
						key: 'contact_phone',
						ellipsis: true,
						render: (v) => formatPhone(v) || '--',
					},
					{
						title: 'Contact email',
						dataIndex: 'contact_email',
						key: 'contact_email',
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
					},
					{
						title: 'Status',
						dataIndex: 'status',
						key: 'status',
						ellipsis: true,
						fixed: 'right',
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
					getActiveByPermission({
						permissionList: [
							PERMISSIONS.LENDERS.EDIT,
							PERMISSIONS.LENDERS.CHANGE_STATUS,
							PERMISSIONS.LENDERS.DELETE,
						],
						render: (_, record) => {
							return (
								<Space>
									<PermissionSlot code={PERMISSIONS.LENDERS.EDIT}>
										<BtnPopUpModal
											getContent={(setOpen) => (
												<MainModal
													rowData={record}
													setOpen={setOpen}
													successCallback={async () =>
														tableFilterRef?.current?.getList?.()
													}
												/>
											)}
											size="small"
											icon={null}
											isStopPropagation={false}
										>
											Edit
										</BtnPopUpModal>
									</PermissionSlot>

									<PermissionSlot code={PERMISSIONS.LENDERS.CHANGE_STATUS}>
										{record?.status === 'Active' && (
											<DeleteBtn
												title="Are you sure you want to disable this lender organization?"
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
												title="Are you sure you want to enable this lender organization?"
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
									<PermissionSlot code={PERMISSIONS.LENDERS.DELETE}>
										<DeleteBtn
											title="Are you sure delete this lender organization?"
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
					}),
				]}
				filterList={filterList}
				// scroll={{ x: scrollX }}
				// pagination={true}

				getData={(params) => {
					if (!api) return;
					const filter = {
						...(params?.filter || {}),
						type: 'Lender',
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
		</div>
	);
};
export default Lenders;
