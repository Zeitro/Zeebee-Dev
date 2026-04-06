'use client';
import React, { useMemo, useRef, useState } from 'react';
import { createRequest, formatPhone, formatTINOrSSN } from '@repo/shared-methods/utils';
import { ReloadOutlined, PlusOutlined } from '@ant-design/icons';
import { BtnPopUpModal, FilterTableCard, DeleteBtn } from '@repo/shared-methods';
import useBorrowingEntity from './components/useBorrowingEntity';
import AddModal from '@/components/SelectCard/Add/borrowingEntity';
import PersonSvg from '@repo/shared-methods/assets/images/person.svg';
import EmailSvg from '@repo/shared-methods/assets/images/email.svg';
import PhoneSvg from '@repo/shared-methods/assets/images/phone.svg';
const BorrowingEntity = () => {
	const { cancel } = createRequest();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	// const [pagination, setPagination] = useState(getPagination());
	const { columns, handleDelete, loading } = useBorrowingEntity({
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
		<FilterTableCard
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
			typeKey="card"
			cardTitle="Borrowing entities"
			itemClassName="w-[420px]"
			renderItem={(row) => {
				return (
					<>
						<div className="flex items-center justify-between">
							<div>
								<p className="font-medium text-base text-black-dark pb-2">
									{row?.name}
								</p>
								<div className="border border-[#91D5FF] rounded-sm px-2 text-xs text-blue-medium bg-[#E6F7FF] w-max">
									{row?.type}
								</div>
							</div>
							<div className="flex items-center">
								<BtnPopUpModal
									getContent={(setOpen) => (
										<AddModal
											setOpen={setOpen}
											handleSearch={() =>
												tableFilterRef?.current?.handleSearch()
											}
											rowData={row}
										/>
									)}
									type="link"
									size="large"
									// isButton={false}
									// icon={<EditOutlined />}
									className="!px-2"
								>
									Edit
								</BtnPopUpModal>
								<DeleteBtn
									type="link"
									title="Are you sure delete this borrowing entity?"
									content={row.name}
									// isButton={false}
									onOk={() => handleDelete(row.id)}
									loading={loading}
									danger
									className="!px-2"
								>
									Delete
								</DeleteBtn>
							</div>
						</div>
						<div className="py-4 border-t border-b border-base-border my-4 flex items-center justify-between">
							<span className="text-base font-medium text-base-gray">EIN</span>
							<span className="text-base font-medium">
								{formatTINOrSSN(row?.tin)}
							</span>
						</div>
						<div className="text-base font-medium  text-base-gray">
							Contact information
						</div>
						<div className="flex items-center gap-4 pt-3">
							<PersonSvg />
							<span className="text-sm/6">{row?.contact_name}</span>
						</div>
						<div className="flex items-center gap-4 pt-3">
							<EmailSvg />
							<span className="text-sm/6">{row?.contact_email}</span>
						</div>
						<div className="flex items-center gap-4 pt-3">
							<PhoneSvg />
							<span className="text-sm/6">{formatPhone(row?.contact_phone)}</span>
						</div>
					</>
				);
			}}
		/>
	);
};
export default BorrowingEntity;
