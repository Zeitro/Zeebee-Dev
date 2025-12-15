import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
import MainModal from './Main';
import { createRequest, formatPhone } from '@repo/shared-methods/utils';
import { useState } from 'react';
import { message, Space } from 'antd';
import { getSelectOptions } from '@repo/shared-methods/api';
import { useUserStore } from '@/store/user';
const useUser = ({ handleSearch }) => {
	const { request } = createRequest();
	const { userInfo } = useUserStore();
	const [loading, setLoading] = useState(false);
	const [statusLoading, setStatusLoading] = useState(false);
	const [roles, setRoles] = useState([]);
	const [rolesLoading, setRolesLoading] = useState(false);
	const [userStatusOptions, setUserStatusOptions] = useState([]);
	const [sendLoading, setSendLoading] = useState(false);

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

	// 删除用户
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/organization_user/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			handleSearch();
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

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
			handleSearch();
			message.open({
				type: 'success',
				content: 'User status updated successfully！',
			});
		} catch (error) {
			console.log(error);
			setStatusLoading(false);
		}
	};

	// 发送邮件邀请用户激活并登录系统
	const handleSendEmail = async (id) => {
		try {
			setSendLoading(true);
			await request(`/organization_user/${id}/send_email`, {
				method: 'POST',
			});
			setSendLoading(false);
			// handleSearch();
			message.open({
				type: 'success',
				content: 'Email sent successfully.',
			});
		} catch (error) {
			console.log(error);
			setSendLoading(false);
		}
	};
	// 获取角色下拉列表
	const getRolesOptions = async () => {
		try {
			setRolesLoading(true);
			const { list } = await request('/role/all');
			setRoles(list);
			setRolesLoading(false);
		} catch (error) {
			console.log(error);
			setRolesLoading(false);
		}
	};

	// 获取用户状态下拉列表
	const getUserStatusOptions = async () => {
		const { options } = await getSelectOptions('organization_user_status');
		setUserStatusOptions(options);
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
			dataIndex: 'full_name',
			key: 'full_name',
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
			dataIndex: 'roles',
			key: 'roles',
			ellipsis: true,
			render: (v) => (
				<Space>
					{v?.map((item, index) => (
						<div
							key={index}
							className="text-xs/5 px-2 border border-[#91D5FF] bg-[#E6F7FF] text-base-primary"
						>
							{item?.name}
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
					className={'text-xs/5 px-2 border rounded-sm w-max'}
					style={{
						borderColor: statusColor[v]?.border,
						backgroundColor: statusColor[v]?.bg,
						color: statusColor[v]?.color,
					}}
				>
					{v}
				</div>
			),
		},
		{
			title: 'Action',
			dataIndex: 'operation',
			width: 220,
			align: 'left',
			fixed: 'right',
			render: (_, record) => {
				const isDisabled = record?.is_owner;
				return (
					<Space>
						<BtnPopUpModal
							getContent={(setOpen) => (
								<MainModal
									setOpen={setOpen}
									handleSearch={handleSearch}
									rowData={record}
									roles={roles}
								/>
							)}
							size="small"
							icon={null}
							isStopPropagation={false}
							// disabled={isDisabled}
						>
							Edit
						</BtnPopUpModal>
						{record?.status === 'Not activated' && (
							<BtnPopUpModal
								config={{
									title: `Are you sure you want to invite this user to join the organization(${userInfo?.organization_name})?`,
									content: record.email,
									okText: 'Send Invitation',
									onOk: () => handleSendEmail(record.id),
									okButtonProps: { danger: false },
									confirmLoading: sendLoading,
								}}
								modalType="confirm"
								isStopPropagation={false}
								// disabled={isDisabled}
							>
								Send Invitation
							</BtnPopUpModal>
						)}
						{record?.status === 'Active' && (
							<DeleteBtn
								title="Are you sure you want to deactivate this user?"
								content={record.email}
								onOk={() => handleChangeStatus(record.id, 'Disabled')}
								loading={statusLoading}
								isStopPropagation={false}
								okText="Deactivate"
								danger={false}
								disabled={isDisabled}
							>
								Disable
							</DeleteBtn>
						)}
						{record?.status === 'Disabled' && (
							<DeleteBtn
								title="Are you sure you want to reactivate this user?"
								content={record.email}
								onOk={() => handleChangeStatus(record.id, 'Active')}
								loading={statusLoading}
								isStopPropagation={false}
								okText="Recover"
								danger={false}
								disabled={isDisabled}
								okButtonProps={{ danger: false }}
							>
								Recover
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
	];
	return {
		columns,
		getRolesOptions,
		roles,
		rolesLoading,
		getUserStatusOptions,
		userStatusOptions,
	};
};
export default useUser;
