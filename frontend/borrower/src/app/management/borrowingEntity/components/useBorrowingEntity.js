import React from 'react';
import { BtnPopUpModal, DeleteBtn, TableOperation } from '@repo/shared-methods';
import { createRequest, formatPhone } from '@repo/shared-methods/utils';
import AddModal from '@/components/SelectCard/Add/borrowingEntity';
import { useState } from 'react';
const useBorrowingEntity = ({ handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_entity/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			handleSearch();
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};
	const columns = [
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			ellipsis: true,
		},
		{
			title: 'Entity type',
			dataIndex: 'type',
			key: 'type',
			ellipsis: true,
		},
		{
			title: 'EIN',
			dataIndex: 'tin',
			key: 'tin',
			ellipsis: true,
		},
		{
			title: 'Email',
			dataIndex: 'contact_email',
			key: 'contact_email',
			ellipsis: true,
		},
		{
			title: 'Phone number',
			dataIndex: 'contact_phone',
			key: 'contact_phone',
			ellipsis: true,
			render: (v) => formatPhone(v) || '--',
		},
		{
			title: 'Contact name',
			dataIndex: 'contact_name',
			key: 'contact_name',
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
						const edit = {
							key: editItem.key,
							label: (
								<BtnPopUpModal
									getContent={(setOpen) => (
										<AddModal
											setOpen={setOpen}
											handleSearch={handleSearch}
											rowData={record}
										/>
									)}
									type="primary"
									size="large"
									isButton={false}
									icon={editItem.icon}
								>
									Edit
								</BtnPopUpModal>
							),
						};
						const deleteInfo = {
							key: deleteItem.key,
							label: (
								<DeleteBtn
									type="default"
									title="Are you sure delete this borrowing entity?"
									content={record.name}
									isButton={false}
									onOk={() => handleDelete(record.id)}
									loading={loading}
								>
									Delete
								</DeleteBtn>
							),
						};
						return [edit, deleteInfo];
					}}
				/>
			),
		},
	];
	return { columns, handleDelete, loading };
};
export default useBorrowingEntity;
