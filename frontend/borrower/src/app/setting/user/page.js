'use client';
import CommonSettingUser from '@repo/shared-methods/componentsUi/SettingUser';
import { PERMISSIONS } from '@config/permissions';
const SettingUser = () => <CommonSettingUser PERMISSIONS={PERMISSIONS} />;
export default SettingUser;
// import { useRef, useState } from 'react';
// import { createRequest } from '@repo/shared-methods/utils';
// import { ReloadOutlined, PlusOutlined } from '@ant-design/icons';
// import { BtnPopUpModal, FilterTableCard } from '@repo/shared-methods';
// import useContent from './components/useContent';
// import MainModal from './components/Main';
// import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

// const SettingUser = () => {
// 	const { cancel } = createRequest();
// 	const tableFilterRef = useRef();
// 	const [searchValue, setSearchValue] = useState('');
// 	const [role, setRole] = useState(null);
// 	const [status, setStatus] = useState(null);
// 	const { columns, getRolesOptions, roles, getUserStatusOptions, userStatusOptions } = useContent(
// 		{
// 			handleSearch: () => tableFilterRef?.current?.handleSearch(),
// 		},
// 	);

// 	const filterList = [
// 		{
// 			filterType: 'input',
// 			placeholder: 'Search by any field',
// 			paramsKey: 'keyword',
// 			value: searchValue,
// 			onChange: ({ target: { value } }) => setSearchValue(value),
// 		},
// 		{
// 			filterType: 'select',
// 			label: 'Role',
// 			paramsKey: 'role_id',
// 			placeholder: 'Role',
// 			value: role,
// 			onChange: (value) => setRole(value),
// 			options: roles,
// 			fieldNames: {
// 				label: 'name',
// 				value: 'id',
// 			},
// 		},
// 		{
// 			filterType: 'select',
// 			label: 'Status',
// 			paramsKey: 'status',
// 			placeholder: 'Status',
// 			value: status,
// 			onChange: (value) => setStatus(value),
// 			options: userStatusOptions,
// 		},
// 		{
// 			filterType: 'button',
// 			onClick: () => {
// 				cancel('/organization_user');
// 				tableFilterRef?.current?.handleSearch();
// 			},
// 		},
// 		{
// 			filterType: 'button',
// 			type: 'default',
// 			text: 'Reset',
// 			icon: <ReloadOutlined />,
// 			onClick: () => {
// 				cancel('/organization_user');
// 				setSearchValue('');
// 				setTimeout(() => {
// 					tableFilterRef?.current?.handleSearch();
// 				}, 300);
// 			},
// 			colProps: {
// 				className: 'flex-1',
// 			},
// 		},
// 		{
// 			filterType: 'button',
// 			children: (
// 				<BtnPopUpModal
// 					getContent={(setOpen) => (
// 						<MainModal
// 							setOpen={setOpen}
// 							handleSearch={() => tableFilterRef?.current?.handleSearch()}
// 							roles={roles}
// 						/>
// 					)}
// 					type="primary"
// 					icon={<PlusOutlined />}
// 					size="large"
// 				>
// 					Add user
// 				</BtnPopUpModal>
// 			),
// 		},
// 	];

// 	useDeferredSyncEffect({
// 		deferred: () => {
// 			getRolesOptions();
// 			getUserStatusOptions();
// 		},
// 		deps: [],
// 		isUseMicrotask: true,
// 	});
// 	return (
// 		<FilterTableCard
// 			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
// 			ref={tableFilterRef}
// 			columns={columns}
// 			scroll={{ x: 1200 }}
// 			filterList={filterList}
// 			hasFilterField={true}
// 			// pagination={pagination}
// 			api="/organization_user/page"
// 			onRow={() => ({
// 				onClick: () => {
// 					// handleSelectedRow(row, stage);
// 				},
// 			})}
// 		/>
// 	);
// };
// export default SettingUser;
