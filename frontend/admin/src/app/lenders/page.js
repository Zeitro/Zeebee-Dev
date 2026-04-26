'use client';
import { useEffect, useMemo, useRef, useState } from 'react';
import { Button, InputNumber, message, Modal, Space, Image, Tag } from 'antd';
import { BtnPopUpModal, DeleteBtn, FilterTable, PermissionSlot } from '@repo/shared-methods';
import { createRequest, formatPhone, formatTime } from '@repo/shared-methods/utils';
import { PlusOutlined, ReloadOutlined, SearchOutlined, SettingOutlined } from '@ant-design/icons';
import MainModal from './components/Main';
import DiscoverSettingsDrawer from './components/DiscoverSettingsDrawer';
import DiscoverPreviewModal from './components/DiscoverPreviewModal';
import { userStatusColor } from '@repo/shared-methods/enum';
import { useAuth } from '@repo/shared-methods/hooks';
import { PERMISSIONS } from '@config/permissions';
const api = '/organization/page';

const OrderInput = ({ value, onCommit }) => {
	const [draft, setDraft] = useState(value ?? 2);
	useEffect(() => {
		setDraft(value ?? 2);
	}, [value]);

	const commit = () => {
		if (draft == null || draft === value) return;
		onCommit(draft);
	};

	return (
		<InputNumber
			size="small"
			min={0}
			value={draft}
			onChange={(val) => setDraft(val)}
			onBlur={commit}
			onPressEnter={commit}
			style={{ width: 64 }}
		/>
	);
};

const Lenders = () => {
	const tableFilterRef = useRef();
	const { request } = createRequest();
	const requestRef = useRef(request);
	useEffect(() => {
		requestRef.current = request;
	});

	const [searchValue, setSearchValue] = useState('');
	const [loading, setLoading] = useState(false);
	const [statusLoading, setStatusLoading] = useState(false);
	const [discoverDrawerOpen, setDiscoverDrawerOpen] = useState(false);
	const [discoverPreviewOpen, setDiscoverPreviewOpen] = useState(false);
	const { hasPermission, getActiveByPermission } = useAuth();

	const openDiscoverPreview = () => {
		Modal.confirm({
			title: 'Search lenders with ChatGPT?',
			content:
				'This can take several minutes. You will review new vs existing lenders before anything is saved. Results are unverified research leads.',
			okText: 'Continue',
			onOk: () => setDiscoverPreviewOpen(true),
		});
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
			{
				filterType: 'button',
				children: (
					<PermissionSlot code={PERMISSIONS.LENDERS.DISCOVER_SETTINGS}>
						<Button
							icon={<SettingOutlined />}
							size="large"
							onClick={() => setDiscoverDrawerOpen(true)}
						>
							Search settings
						</Button>
					</PermissionSlot>
				),
				hidden: !hasPermission(PERMISSIONS.LENDERS.DISCOVER_SETTINGS),
			},
			{
				filterType: 'button',
				children: (
					<PermissionSlot code={PERMISSIONS.LENDERS.DISCOVER_RUN}>
						<Button
							type="primary"
							icon={<SearchOutlined />}
							size="large"
							onClick={openDiscoverPreview}
						>
							Search lenders
						</Button>
					</PermissionSlot>
				),
				hidden: !hasPermission(PERMISSIONS.LENDERS.DISCOVER_RUN),
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

	const handleUpdateOrder = async (id, order) => {
		try {
			await request(`/organization/${id}/order`, {
				method: 'PATCH',
				body: { order },
			});
			tableFilterRef?.current?.getList?.();
		} catch {
			message.error('Failed to update order');
		}
	};

	return (
		<div className="flex-1 overflow-hidden w-full flex flex-col py-4 px-4 gap-3">
			<div className="flex-1 min-h-0 flex items-center justify-center">
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
									loading="lazy"
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
							title: 'Source',
							dataIndex: 'source',
							key: 'source',
							width: 108,
							render: (v) => {
								if (v === 'Search') return <Tag color="blue">Search</Tag>;
								if (v === 'Manual') return <Tag>Manual</Tag>;
								return v || '—';
							},
						},
						{
							title: 'Order',
							dataIndex: 'order',
							key: 'order',
							width: 90,
							render: (v, record) =>
								hasPermission(PERMISSIONS.LENDERS.EDIT) ? (
									<OrderInput
										value={v}
										onCommit={(val) => handleUpdateOrder(record.id, val)}
									/>
								) : (
									(v ?? '—')
								),
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
													onOk={() =>
														handleChangeStatus(record.id, 'Active')
													}
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
			<DiscoverSettingsDrawer
				open={discoverDrawerOpen}
				onClose={() => setDiscoverDrawerOpen(false)}
			/>
			<DiscoverPreviewModal
				open={discoverPreviewOpen}
				onClose={() => setDiscoverPreviewOpen(false)}
				onCommitted={() => {
					message.success('Selected lenders added as Not activated.');
					tableFilterRef?.current?.handleSearch?.();
				}}
			/>
		</div>
	);
};
export default Lenders;
