import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
import MainModal from './Main';
import { createRequest, formatPhone } from '@repo/shared-methods/utils';
import { useRef, useState } from 'react';
import { message, Space } from 'antd';
import { getSelectOptions } from '@repo/shared-methods/api';
import { useUserStore } from '@repo/shared-methods/store';
import { userStatusColor } from '@repo/shared-methods/enum';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
const useUser = ({ getNewMainModal }) => {
	const { request, cancel } = createRequest();
	const { userInfo } = useUserStore();
	const [loading, setLoading] = useState(false);
	const [statusLoading, setStatusLoading] = useState(false);
	const [roles, setRoles] = useState([]);
	const [rolesLoading, setRolesLoading] = useState(false);
	const [userStatusOptions, setUserStatusOptions] = useState([]);
	const [sendLoading, setSendLoading] = useState(false);
	//----start user 列表所需字段
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [role, setRole] = useState(null);
	const [status, setStatus] = useState(null);
	const filterList = [
		{
			filterType: 'input',
			placeholder: 'Search by any field',
			paramsKey: 'keyword',
			value: searchValue,
			onChange: ({ target: { value } }) => setSearchValue(value),
		},
		{
			filterType: 'select',
			label: 'Role',
			paramsKey: 'role_id',
			placeholder: 'Role',
			value: role,
			onChange: (value) => setRole(value),
			options: roles,
			fieldNames: {
				label: 'name',
				value: 'id',
			},
		},
		{
			filterType: 'select',
			label: 'Status',
			paramsKey: 'status',
			placeholder: 'Status',
			value: status,
			onChange: (value) => setStatus(value),
			options: userStatusOptions,
		},
		{
			filterType: 'button',
			onClick: () => {
				cancel('/organization_user');
				tableFilterRef?.current?.handleSearch();
			},
		},
		{
			filterType: 'button',
			type: 'default',
			text: 'Reset',
			icon: <ReloadOutlined />,
			onClick: () => {
				cancel('/organization_user');
				setSearchValue('');
				setTimeout(() => {
					tableFilterRef?.current?.handleSearch();
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
					getContent={(setOpen) => {
						const mainProps = {
							setOpen,
							handleSearch: () => tableFilterRef?.current?.handleSearch(),
							roles,
						};
						return getNewMainModal(mainProps) || <MainModal {...mainProps} />;
					}}
					type="primary"
					icon={<PlusOutlined />}
					size="large"
				>
					Add user
				</BtnPopUpModal>
			),
		},
	];
	// ----end

	const handleSearch = (flag) => {
		if (flag === 'list') {
			tableFilterRef?.current?.getList();
			return;
		}
		tableFilterRef?.current?.handleSearch();
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
			handleSearch('list');
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
			width: 240,
			align: 'left',
			fixed: 'right',
			render: (_, record) => {
				const isDisabled = record?.is_owner;
				return (
					<Space>
						<BtnPopUpModal
							getContent={(setOpen) => {
								const mainProps = {
									setOpen,
									handleSearch: () => tableFilterRef?.current?.handleSearch(),
									rowData: record,
									roles,
								};
								return getNewMainModal(mainProps) || <MainModal {...mainProps} />;
							}}
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
								title="Are you sure you want to disable this user?"
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
								title="Are you sure you want to enable this user?"
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
	];
	return {
		columns,
		getRolesOptions,
		roles,
		rolesLoading,
		getUserStatusOptions,
		userStatusOptions,
		// user列表所需字段
		filterList,
		tableFilterRef,
	};
};
export default useUser;
