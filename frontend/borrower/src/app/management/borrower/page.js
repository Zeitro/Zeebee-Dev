'use client';
import React, { useMemo, useRef, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { ReloadOutlined } from '@ant-design/icons';
import { BtnPopUpModal, FilterTable } from '@repo/shared-methods';
import useBorrower from './components/useBorrower';
import AddModal from './components/Add';
import NewHeader from '@/components/NewHeader';
const BorrowerManagement = () => {
	const { cancel } = createRequest();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const { columns } = useBorrower({
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
			},
		],
		[searchValue, cancel],
	);
	return (
		<>
			<NewHeader isShowCrumbs={true}>
				<BtnPopUpModal
					getContent={(setOpen) => (
						<AddModal
							setOpen={setOpen}
							handleSearch={() => tableFilterRef?.current?.handleSearch()}
						/>
					)}
					type="primary"
					icon={''}
					size="large"
				>
					Add borrower
				</BtnPopUpModal>
			</NewHeader>
			<FilterTable
				filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
				ref={tableFilterRef}
				columns={columns}
				scroll={{ x: 1200 }}
				filterList={filterList}
				hasFilterField={true}
				api="/borrower_applicant/page"
				onRow={(row) => ({
					onClick: () => {
						// handleSelectedRow(row, stage);
					},
				})}
			/>
		</>
	);
};
export default BorrowerManagement;
