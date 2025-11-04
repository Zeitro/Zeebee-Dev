'use client';
import React, { useCallback, useMemo, useRef, useState } from 'react';
import { createRequest, getPagination, transformArrayStructure } from '@repo/shared-methods/utils';
import { ReloadOutlined } from '@ant-design/icons';
import { BtnPopUpModal, FilterField, FilterTable } from '@repo/shared-methods';
import useBorrower from './components/useBorrower';
import AddModal from './components/Add';
const BorrowerManagement = () => {
	const { request, cancel } = createRequest();
	const filterFieldRef = useRef();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [pagination, setPagination] = useState(getPagination());
	const { columns } = useBorrower();
	const getList = useCallback(() => {
		// request('');
	}, [cancel, getPagination]);

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
					cancel('/property/page');
					setPagination((v) => {
						v.current = 1;
						return { ...v };
					});
					getList();
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
						setPagination((v) => {
							v.current = 1;
							return { ...v };
						});
						getList?.();
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
						getContent={(setOpen) => <AddModal setOpen={setOpen} />}
						type="primary"
						icon={''}
						size="large"
					>
						Add
					</BtnPopUpModal>
				),
			},
		],
		[searchValue, cancel, getList],
	);
	return (
		<FilterTable
			filterWrapClassName="px-8 pt-10 pb-8"
			ref={tableFilterRef}
			columns={columns}
			scroll={{ x: 1400 }}
			filterList={filterList}
			hasFilterField={true}
			pagination={pagination}
			onRow={(row) => ({
				onClick: () => {
					// handleSelectedRow(row, stage);
				},
			})}
			getData={(params) => {
				const payload = {
					method: 'POST',
					body: params,
				};
				// return request('api', payload);
				return {};
			}}
		/>
	);
};
export default BorrowerManagement;
