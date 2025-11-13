'use client';
import React, { useMemo, useRef, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { ReloadOutlined, PlusOutlined } from '@ant-design/icons';
import { BtnPopUpModal, FilterTable } from '@repo/shared-methods';
import useBorrowingEntity from './components/useBorrowingEntity';
import AddModal from '@/components/SelectCard/Add/borrowingEntity';
const BorrowingEntity = () => {
	const { cancel } = createRequest();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	// const [pagination, setPagination] = useState(getPagination());
	const { columns } = useBorrowingEntity({
		handleSearch: () => tableFilterRef?.current?.handleSearch(),
	});

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
					cancel('/borrower_entity/page');
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
							<AddModal
								setOpen={setOpen}
								handleSearch={() => tableFilterRef?.current?.handleSearch()}
							/>
						)}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add borrowing entity
					</BtnPopUpModal>
				),
			},
		],
		[searchValue, cancel],
	);
	return (
		<FilterTable
			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
			ref={tableFilterRef}
			columns={columns}
			scroll={{ x: 1200 }}
			filterList={filterList}
			hasFilterField={true}
			// pagination={pagination}
			api="/borrower_entity/page"
			onRow={() => ({
				onClick: () => {
					// handleSelectedRow(row, stage);
				},
			})}
		/>
	);
};
export default BorrowingEntity;
