'use client';
import React, { useRef, useState } from 'react';
import { ReloadOutlined, PlusOutlined } from '@ant-design/icons';
import { BtnPopUpModal, FilterTableCard } from '@repo/shared-methods';
import useContent from './components/useContent';
import MainModal from './components/Main';
import { useAuth } from '@repo/shared-methods/hooks';
const SettingRole = ({ api, systemType, PERMISSIONS }) => {
	const { hasPermission } = useAuth();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const { columns } = useContent({
		handleSearch: () => tableFilterRef?.current?.handleSearch(),
		systemType,
		PERMISSIONS,
	});

	const filterList = [
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
				// cancel('/property/page');
				tableFilterRef?.current?.handleSearch();
			},
		},
		{
			filterType: 'button',
			type: 'default',
			text: 'Reset',
			icon: <ReloadOutlined />,
			onClick: () => {
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
			hidden: !hasPermission(PERMISSIONS.SETTING_ROLE.ADD),
			children: (
				<BtnPopUpModal
					getContent={(setOpen) => (
						<MainModal
							setOpen={setOpen}
							handleSearch={() => tableFilterRef?.current?.handleSearch()}
							systemType={systemType}
						/>
					)}
					type="primary"
					icon={<PlusOutlined />}
					size="large"
				>
					Add role
				</BtnPopUpModal>
			),
		},
	];
	return (
		<FilterTableCard
			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
			ref={tableFilterRef}
			columns={columns}
			// scroll={{ x: 1200 }}
			filterList={filterList}
			hasFilterField={true}
			// pagination={pagination}
			api={api}
			onRow={() => ({
				onClick: () => {
					// handleSelectedRow(row, stage);
				},
			})}
		/>
	);
};
export default SettingRole;
