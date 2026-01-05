import { BtnPopUpModal, DeleteBtn, TableOperation } from '@repo/shared-methods';
import MainModal from './Main';
import AssignModal from './AssignModal';
import { createRequest } from '@repo/shared-methods/utils';
import { useState } from 'react';

const useContent = (handleSearch = () => {}) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_applicant/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			handleSearch?.();
		} catch (error) {
			console.log(error);
			setLoading(false);
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
			title: 'Create Time',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
		},
		{
			title: '',
			dataIndex: 'operation',
			width: 50,
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
									/>
								)}
								type="primary"
								size="large"
								isButton={false}
								icon={editItem.icon}
								isStopPropagation={false}
							>
								Edit
							</BtnPopUpModal>
						);
						deleteItem.label = (
							<DeleteBtn
								type="default"
								title="Are you sure delete this user?"
								content={record.full_name}
								isButton={false}
								onOk={() => handleDelete(record.id)}
								loading={loading}
								isStopPropagation={false}
							>
								Delete
							</DeleteBtn>
						);

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
			dataIndex: 'name',
			key: 'name',
			ellipsis: true,
		},
		{
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
		},
		{
			title: 'Phone',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
		},
		{
			title: 'Role',
			dataIndex: 'role',
			key: 'role',
			ellipsis: true,
		},
		{
			title: 'Status',
			dataIndex: 'status',
			key: 'status',
			ellipsis: true,
		},
		{
			title: 'Create Time',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
		},
	];
	return { columns, assignColumns };
};
export default useContent;
