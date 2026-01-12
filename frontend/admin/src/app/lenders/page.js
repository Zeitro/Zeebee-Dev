'use client';
import React, { useMemo, useRef, useState } from 'react';
import { message, Space } from 'antd';
import { BtnPopUpModal, DeleteBtn, FilterTable, NewModal, OrgProfile } from '@repo/shared-methods';
import { createRequest, formatPhone } from '@repo/shared-methods/utils';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
const api = '';
const Lenders = () => {
	const tableFilterRef = useRef();
	const mainRef = useRef();
	const { request } = createRequest();
	const [searchValue, setSearchValue] = useState('');
	const [loading, setLoading] = useState(false);
	const [statusLoading, setStatusLoading] = useState(false);

	const mainModal = (setOpen, rowData) => {
		return (
			<NewModal
				width={800}
				title={rowData?.id ? 'Edit Lender' : 'Add Lender'}
				onCancel={() => setOpen(false)}
				onOk={() => mainRef?.current?.handleSaveOrgProfile?.()}
			>
				<OrgProfile
					ref={mainRef}
					rowData={rowData}
					className="h-[70vh] overflow-hidden overflow-y-auto"
					successCallback={() => tableFilterRef?.current?.handleSearch?.()}
				/>
			</NewModal>
		);
	};
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
					setTimeout(() => {
						tableFilterRef?.current?.handleSearch?.();
					}, 300);
				},
				colProps: {
					className: 'flex-1',
				},
			},
			{
				filterType: 'button',
				children: (
					<BtnPopUpModal
						getContent={(setOpen) => mainModal(setOpen, {})}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add Lender
					</BtnPopUpModal>
				),
			},
		],
		[searchValue],
	);
	// 改变状态
	const handleChangeStatus = async (id, status = 'Active') => {
		try {
			setStatusLoading(true);
			await request(`/organization_user/${id}/status`, {
				method: 'PATCH',
				body: {
					status,
				},
			});
			setStatusLoading(false);
			// handleSearch();
			message.open({
				type: 'success',
				content: 'Lender status updated successfully！',
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
			await request(`/organization_user/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			// handleSearch();
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	return (
		<div className="h-full w-full flex items-center justify-center py-4 px-4">
			<FilterTable
				ref={tableFilterRef}
				columns={[
					{
						title: 'Lender name',
						dataIndex: 'name',
						key: 'name',
						ellipsis: true,
					},
					{
						title: 'Email account',
						dataIndex: 'email_account',
						key: 'email_account',
						ellipsis: true,
					},
					{
						title: 'Contact name',
						dataIndex: 'contact_name',
						key: 'contact_name',
						ellipsis: true,
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
					},
					{
						title: 'Status',
						dataIndex: 'status',
						key: 'status',
						ellipsis: true,
					},
					{
						title: 'Action',
						dataIndex: 'operation',
						width: 240,
						align: 'left',
						fixed: 'right',
						render: (_, record) => {
							const isDisabled = record?.is_owner;
							return (
								<Space>
									<BtnPopUpModal
										getContent={(setOpen) => mainModal(setOpen, record)}
										size="small"
										icon={null}
										isStopPropagation={false}
										// disabled={isDisabled}
									>
										Edit
									</BtnPopUpModal>

									{record?.status === 'Active' && (
										<DeleteBtn
											title="Are you sure you want to disable this lender?"
											content={record.email}
											onOk={() => handleChangeStatus(record.id, 'Disabled')}
											loading={statusLoading}
											isStopPropagation={false}
											okText="Disable"
											danger={false}
											disabled={isDisabled}
										>
											Disable
										</DeleteBtn>
									)}
									{record?.status === 'Disabled' && (
										<DeleteBtn
											title="Are you sure you want to enable this lender?"
											content={record.email}
											onOk={() => handleChangeStatus(record.id, 'Active')}
											loading={statusLoading}
											isStopPropagation={false}
											okText="Enable"
											danger={false}
											disabled={isDisabled}
											okButtonProps={{ danger: false }}
										>
											Enable
										</DeleteBtn>
									)}
									<DeleteBtn
										title="Are you sure delete this user?"
										content={record.email}
										// isButton={false}
										onOk={() => handleDelete(record.id)}
										loading={loading}
										isStopPropagation={false}
										disabled={isDisabled}
									>
										Delete
									</DeleteBtn>
								</Space>
							);
						},
					},
				]}
				filterList={filterList}
				// scroll={{ x: scrollX }}
				// pagination={true}

				getData={(params) => {
					if (!api) return;
					const filter = {
						...(params?.filter || {}),
					};
					const payload = {
						method: 'POST',
						...params,
						filter,
					};

					return request(api, payload);
				}}
			/>
		</div>
	);
};
export default Lenders;
