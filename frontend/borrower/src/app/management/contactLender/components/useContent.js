import React from 'react';
import { BtnPopUpModal, DeleteBtn, TableOperation } from '@repo/shared-methods';
import AddModal from '@/components/SelectCard/Add/borrower';
import {
	assemblyAddress,
	createRequest,
	formatPhone,
	formatTINOrSSN,
} from '@repo/shared-methods/utils';
import { useState } from 'react';
const useContent = ({ handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_applicant/${id}`, {
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
			dataIndex: 'full_name',
			key: 'full_name',
			ellipsis: true,
		},

		{
			title: 'SSN',
			dataIndex: 'ssn',
			key: 'ssn',
			ellipsis: true,
			width: 160,
			render: (v) => formatTINOrSSN(v) || '--',
		},
		{
			title: 'Phone number',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
			width: 160,
			render: (v) => formatPhone(v) || '--',
		},
		{
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
		},
		{
			title: 'Estimated credit score',
			dataIndex: 'est_credit_score',
			key: 'est_credit_score',
			ellipsis: true,
			width: 200,
		},
		{
			title: 'Primary address',
			dataIndex: ['address', 'address1'],
			key: 'address1',
			ellipsis: true,
			render: (v, { address }) => assemblyAddress(address) || '--',
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
									title="Are you sure delete this borrower?"
									content={record.full_name}
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
export default useContent;
