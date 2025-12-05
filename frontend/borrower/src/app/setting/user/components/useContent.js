import { BtnPopUpModal, DeleteBtn, TableOperation } from '@repo/shared-methods';
import MainModal from './Main';
import { createRequest, formatPhone } from '@repo/shared-methods/utils';
import { useState } from 'react';
import { Button, Space } from 'antd';

const useUser = ({ handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const statusColor = {
		Active: {
			border: '#B7EB8F',
			bg: '#F6FFED',
			color: '#52C41A',
		},
		'Not activated': {
			border: '#FFD591',
			bg: '#FFF7E6',
			color: '#FA8C16',
		},
		Disabled: {
			border: 'var(--base-border)',
			bg: '#FAFAFA',
			color: 'rgba(0, 0, 0, 0.85)',
		},
	};
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
			title: 'Email',
			dataIndex: 'email',
			key: 'email',
			ellipsis: true,
		},
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			ellipsis: true,
		},
		{
			title: 'Phone',
			dataIndex: 'phone',
			key: 'phone',
			ellipsis: true,
			render: (v) => formatPhone(v) || '--',
		},
		{
			title: 'Role',
			dataIndex: 'role',
			key: 'role',
			ellipsis: true,
			render: (v) => (
				<Space>
					{v?.map((item, index) => (
						<div
							key={index}
							className="text-xs/5 px-2 border border-[#91D5FF] bg-[#E6F7FF]"
						>
							{item}
						</div>
					))}
				</Space>
			),
		},
		{
			title: 'Status',
			dataIndex: 'status',
			key: 'status',
			ellipsis: true,
			render: (v) => (
				<div
					className={[
						'text-xs/5 px-2 border, rounded-lg',
						`border-[${statusColor[v]?.border}]`,
						`bg-[${statusColor[v]?.bg}]`,
						`text-[${statusColor[v]?.color}]`,
					].join(' ')}
				>
					{v}
				</div>
			),
		},
		{
			title: 'Action',
			dataIndex: 'operation',
			width: 200,
			align: 'left',
			fixed: 'right',
			render: (_, record) => (
				<Space>
					<BtnPopUpModal
						getContent={(setOpen) => (
							<MainModal
								setOpen={setOpen}
								handleSearch={handleSearch}
								rowData={record}
							/>
						)}
						size="small"
						// isButton={false}
						icon={null}
						isStopPropagation={false}
						// className="text-primary"
					>
						Edit
					</BtnPopUpModal>
					{record?.status === 'Not activated' && (
						<Button type="link" size="small">
							Send Email
						</Button>
					)}
					{record?.status === 'Active' && (
						<DeleteBtn
							title="Are you sure you want to deactivate this user?"
							content={record.full_name}
							// isButton={false}
							onOk={() => handleDelete(record.id)}
							loading={loading}
							isStopPropagation={false}
							okText="Deactivate"
						>
							Disabled
						</DeleteBtn>
					)}
					{record?.status === 'Disabled' && (
						<DeleteBtn
							title="Are you sure you want to reactivate this user?"
							content={record.full_name}
							// isButton={false}
							onOk={() => handleDelete(record.id)}
							loading={loading}
							isStopPropagation={false}
							okText="Reactivate"
						>
							Reactivate
						</DeleteBtn>
					)}
					<DeleteBtn
						title="Are you sure delete this user?"
						content={record.full_name}
						// isButton={false}
						onOk={() => handleDelete(record.id)}
						loading={loading}
						isStopPropagation={false}
					>
						Delete
					</DeleteBtn>
				</Space>
			),
		},
	];
	return { columns };
};
export default useUser;
