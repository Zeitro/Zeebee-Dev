import { BtnPopUpModal, DeleteBtn, TableOperation } from '@repo/shared-methods';
import MainModal from './Main';
import AssignModal from './AssignModal';
import { createRequest, formatTime } from '@repo/shared-methods/utils';
import { useState } from 'react';
import { userStatusColor } from '@repo/shared-methods/enum';
import { message, Tooltip } from 'antd';

const useContent = ({ handleSearch = () => {}, systemType }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/role/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			message.open({
				type: 'success',
				content: 'Role delete successful.',
			});
			handleSearch?.();
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};
	const columns = [
		{
			title: 'Role name',
			dataIndex: 'name',
			key: 'name',
			ellipsis: true,
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
			title: 'Create Time',
			dataIndex: 'created_at',
			key: 'created_at',
			ellipsis: true,
			width: 180,
			render: (v) => formatTime(v) || '--',
		},
		{
			title: 'Description',
			dataIndex: 'description',
			key: 'description',
			width: 200,
			ellipsis: true,
			render: (v) => v || '--',
		},
		{
			title: 'Action',
			dataIndex: 'operation',
			width: 240,
			align: 'left',
			fixed: 'right',
			render: (_, record) => (
				<TableOperation
					record={record}
					getItems={(editItem, deleteItem) => {
						editItem.label = (
							<BtnPopUpModal
								getContent={(setOpen) => (
									<MainModal
										setOpen={setOpen}
										handleSearch={handleSearch}
										rowData={record}
										systemType={systemType}
									/>
								)}
								type="primary"
								size="large"
								isButton={false}
								isStopPropagation={false}
							>
								{editItem.icon}&nbsp;Edit
							</BtnPopUpModal>
						);
						deleteItem.label = (
							<DeleteBtn
								type="default"
								title="Are you sure delete this role?"
								content={record.name}
								isButton={false}
								onOk={() => handleDelete(record.id)}
								loading={loading}
								isStopPropagation={false}
							>
								{deleteItem?.icon}&nbsp;Delete
							</DeleteBtn>
						);
						editItem.icon = null;
						deleteItem.icon = null;
						return [
							editItem,
							deleteItem,
							{
								label: (
									<BtnPopUpModal
										getContent={(setOpen) => (
											<AssignModal
												setOpen={setOpen}
												handleSearch={handleSearch}
												rowData={record}
											/>
										)}
										type="primary"
										size="large"
										isButton={false}
										icon={editItem.icon}
										isStopPropagation={false}
									>
										Assign
									</BtnPopUpModal>
								),
								value: 'assign',
							},
						];
					}}
				/>
			),
		},
	];

	const assignColumns = [
		{
			title: 'Name',
			dataIndex: 'full_name',
			key: 'full_name',
			ellipsis: true,
			width: 160,
		},
		{
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
			width: 200,
		},
		{
			title: 'Phone',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
			width: 200,
		},
		{
			title: 'Roles',
			dataIndex: 'roles',
			key: 'roles',
			ellipsis: true,
			render: (v) => {
				const roles = v?.map?.((item) => item.name).join(', ') || [];
				return <Tooltip title={roles}>{roles}</Tooltip>;
			},
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
			title: 'Create Time',
			dataIndex: 'created_at',
			key: 'created_at',
			ellipsis: true,
			render: (v) => formatTime(v) || '--',
			width: 180,
		},
	];
	return { columns, assignColumns };
};
export default useContent;
