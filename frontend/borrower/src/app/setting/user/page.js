'use client';
import React, { useMemo, useRef, useState } from 'react';
import {
	assemblyAddress,
	createRequest,
	formatPhone,
	formatTINOrSSN,
} from '@repo/shared-methods/utils';
import { ReloadOutlined, PlusOutlined } from '@ant-design/icons';
import { BtnPopUpModal, DeleteBtn, FilterTableCard } from '@repo/shared-methods';
import useContent from './components/useContent';
import MainModal from './components/Main';
const SettingUser = () => {
	const { cancel } = createRequest();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [role, setRole] = useState(null);
	const [status, setStatus] = useState(null);
	const { columns, handleDelete, loading } = useContent({
		handleSearch: () => tableFilterRef?.current?.handleSearch(),
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
			filterType: 'select',
			label: 'Role',
			paramsKey: 'role',
			placeholder: 'Role',
			value: role,
			onChange: (value) => setRole(value),
			options: [],
			// treeData: transformArrayStructure({
			// 	sourceArray: loanPurposeOptions,
			// 	outputKey: 'title',
			// }).map((i) => ({ ...i, key: i.value })),
		},
		{
			filterType: 'select',
			label: 'Status',
			paramsKey: 'status',
			placeholder: 'Status',
			value: status,
			onChange: (value) => setStatus(value),
			options: [],
			// treeData: transformArrayStructure({
			// 	sourceArray: loanPurposeOptions,
			// 	outputKey: 'title',
			// }).map((i) => ({ ...i, key: i.value })),
		},
		{
			filterType: 'button',
			onClick: () => {
				cancel('/property/page');
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
			children: (
				<BtnPopUpModal
					getContent={(setOpen) => (
						<MainModal
							setOpen={setOpen}
							handleSearch={() => tableFilterRef?.current?.handleSearch()}
						/>
					)}
					type="primary"
					icon={<PlusOutlined />}
					size="large"
				>
					Add user
				</BtnPopUpModal>
			),
		},
	];
	return (
		<FilterTableCard
			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
			ref={tableFilterRef}
			columns={columns}
			scroll={{ x: 1200 }}
			filterList={filterList}
			hasFilterField={true}
			// pagination={pagination}
			api="/borrower_applicant/page"
			onRow={() => ({
				onClick: () => {
					// handleSelectedRow(row, stage);
				},
			})}
		/>
	);
};
export default SettingUser;
