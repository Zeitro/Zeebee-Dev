import { BtnPopUpModal, DeleteBtn, PermissionSlot } from '@repo/shared-methods';
import { userStatusColor, iconData } from '@repo/shared-methods/enum';
import { Button, message, Space } from 'antd';
import { useState, useRef } from 'react';
import ApiMainModal from './ApiMain';
import ResourceMain from './ResourceMain';
import { createRequest, eventBus, copyToClipboard } from '@repo/shared-methods/utils';
import { useAuth } from '@repo/shared-methods/hooks';
import { PERMISSIONS } from '@config/permissions';
import { CopyOutlined } from '@ant-design/icons';
const useColumn = ({ eventBusPrefix, api }) => {
	const { getActiveByPermission } = useAuth();
	const { request } = createRequest();
	const tableFilterRef = useRef(null);
	const [deleteLoading, setDeleteLoading] = useState(false);

	const handleApiDelete = async (record) => {
		try {
			setDeleteLoading(true);
			await request(`/api/${record?.id}`, {
				method: 'DELETE',
				body: {
					type: record?.type,
				},
			});
			message.open({
				type: 'success',
				content: 'Delete API successfully',
			});
			eventBus.emit(`${eventBusPrefix}${record?.type}`);
			setDeleteLoading(false);
		} catch (error) {
			setDeleteLoading(false);
			throw error;
		}
	};

	const handleResourceDelete = async (record) => {
		try {
			setDeleteLoading(true);
			await request(`/resource/${record?.id}`, {
				method: 'DELETE',
				body: {
					type: record?.type,
				},
			});
			message.open({
				type: 'success',
				content: 'Delete resource successfully',
			});
			eventBus.emit(`${eventBusPrefix}${record?.type}`);
			setDeleteLoading(false);
		} catch (error) {
			setDeleteLoading(false);
			throw error;
		}
	};

	const handleCopyCode = (record) => {
		let str = JSON.stringify([record]);
		copyToClipboard(str);
		message.open({
			type: 'success',
			content: 'Copy code successfully',
		});
	};

	const apiColumns = [
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			width: 260,
			ellipsis: true,
		},
		{
			title: 'Method',
			dataIndex: 'method',
			key: 'method',
			width: 100,
			ellipsis: true,
			render: (v) => v || '--',
		},
		{
			title: 'Path',
			dataIndex: 'path',
			key: 'path',
			width: 300,
			ellipsis: true,
			render: (v) => v || '--',
		},
		{
			title: 'Status',
			dataIndex: 'status',
			key: 'status',
			width: 120,
			ellipsis: true,
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
			title: 'Auth required',
			dataIndex: 'auth_required',
			key: 'auth_required',
			width: 120,
			ellipsis: true,
			align: 'center',
			render: (text) => (text ? 'Yes' : 'No'),
		},
		{
			title: 'Permission required',
			dataIndex: 'permission_required',
			key: 'permission_required',
			width: 160,
			ellipsis: true,
			align: 'center',
			render: (text) => (text ? 'Yes' : 'No'),
		},
		{
			title: 'Description',
			dataIndex: 'description',
			key: 'description',
			ellipsis: true,
			render: (v) => v || '--',
		},
		getActiveByPermission({
			buttonOptions: [
				{
					permission: PERMISSIONS.API.EDIT,
					width: 50,
				},
				{
					permission: PERMISSIONS.API.DELETE,
					width: 70,
				},
			],
			render: (_, record) => {
				return (
					<Space>
						<PermissionSlot code={PERMISSIONS.API.EDIT}>
							<BtnPopUpModal
								getContent={(setOpen) => (
									<ApiMainModal
										setOpen={setOpen}
										handleSearch={() => tableFilterRef?.current?.handleSearch()}
										rowData={record}
										eventBusPrefix={eventBusPrefix}
									/>
								)}
								size="small"
								icon={null}
								isStopPropagation={false}
								// disabled={isDisabled}
							>
								Edit
							</BtnPopUpModal>
						</PermissionSlot>

						<PermissionSlot code={PERMISSIONS.API.DELETE}>
							<DeleteBtn
								title="Are you sure delete this api?"
								content={record?.path || record.name}
								// isButton={false}
								onOk={() => handleApiDelete(record)}
								loading={deleteLoading}
								isStopPropagation={false}
							>
								Delete
							</DeleteBtn>
						</PermissionSlot>
					</Space>
				);
			},
		}),
	];

	const resourceColumns = [
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			width: 260,
			ellipsis: true,
		},
		{
			title: 'Resource Type',
			dataIndex: 'resource_type',
			key: 'resource_type',
			width: 130,
			ellipsis: true,
			align: 'center',
		},
		{
			title: 'Tag',
			dataIndex: 'tag',
			key: 'tag',
			width: 300,
			ellipsis: true,
		},
		{
			title: 'Path',
			dataIndex: 'path',
			key: 'path',
			width: 200,
			ellipsis: true,
			render: (v) => v || '--',
		},

		{
			title: 'Status',
			dataIndex: 'status',
			key: 'status',
			width: 100,
			ellipsis: true,
			align: 'center',
			render: (v) => (
				<span
					className={'text-xs/5 px-2 py-[2px] border rounded-sm w-max'}
					style={{
						borderColor: userStatusColor[v]?.border,
						backgroundColor: userStatusColor[v]?.bg,
						color: userStatusColor[v]?.color,
					}}
				>
					{v}
				</span>
			),
		},
		{
			title: 'Show menu',
			dataIndex: 'is_show',
			key: 'is_show',
			width: 120,
			ellipsis: true,
			align: 'center',
			render: (text) => (text ? 'Yes' : 'No'),
		},
		{
			title: 'Is public',
			dataIndex: 'is_public',
			key: 'is_public',
			width: 120,
			ellipsis: true,
			align: 'center',
			render: (text) => (text ? 'Yes' : 'No'),
		},
		{
			title: 'Icon',
			dataIndex: 'icon',
			key: 'icon',
			width: 120,
			align: 'center',
			ellipsis: true,
			render: (v) =>
				iconData[v] ? (
					<div className="flex justify-center items-center">{iconData[v]}</div>
				) : (
					v || '--'
				),
		},
		{
			title: 'Description',
			dataIndex: 'description',
			key: 'description',
			ellipsis: true,
			render: (v) => v || '--',
			minWidth: 300,
		},
		getActiveByPermission({
			buttonOptions: [
				{
					permission: PERMISSIONS.RESOURCE.EDIT,
					width: 50,
				},
				{
					width: 50,
				},
				{
					permission: PERMISSIONS.RESOURCE.DELETE,
					width: 70,
				},
			],
			render: (_, record) => {
				return (
					<Space>
						<PermissionSlot code={PERMISSIONS.RESOURCE.EDIT}>
							<BtnPopUpModal
								getContent={(setOpen) => (
									<ResourceMain
										setOpen={setOpen}
										rowData={record}
										eventBusPrefix={eventBusPrefix}
										api={api}
									/>
								)}
								size="small"
								icon={null}
								isStopPropagation={false}
								// disabled={isDisabled}
							>
								Edit
							</BtnPopUpModal>
						</PermissionSlot>
						<Button
							type="link"
							icon={<CopyOutlined />}
							onClick={() => handleCopyCode(record)}
						/>
						<PermissionSlot code={PERMISSIONS.RESOURCE.DELETE}>
							<DeleteBtn
								title="Are you sure delete this resource?"
								content={record?.path || record.name}
								// isButton={false}
								onOk={() => handleResourceDelete(record)}
								loading={deleteLoading}
								isStopPropagation={false}
							>
								Delete
							</DeleteBtn>
						</PermissionSlot>
					</Space>
				);
			},
		}),
	];

	const getColumnsByMenuKey = (menuKey) => {
		switch (menuKey) {
			case 'Api':
				return apiColumns;
			case 'Resource':
				return resourceColumns;
			default:
				return [];
		}
	};
	return { getColumnsByMenuKey, tableFilterRef };
};
export default useColumn;
