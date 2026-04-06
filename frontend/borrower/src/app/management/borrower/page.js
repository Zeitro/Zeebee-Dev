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
import useBorrower from './components/useBorrower';
import AddModal from '@/components/SelectCard/Add/borrower';
import DistanceSvg from '@repo/shared-methods/assets/images/distance.svg';
const BorrowerManagement = () => {
	const { cancel } = createRequest();
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const { columns, handleDelete, loading } = useBorrower({
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
						Add borrower
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
			api="/borrower_applicant/page"
			onRow={() => ({
				onClick: () => {
					// handleSelectedRow(row, stage);
				},
			})}
			typeKey="card"
			itemClassName="w-[420px]"
			cardTitle="Borrowers"
			renderItem={(row) => {
				return (
					<>
						<div className="flex items-center justify-between">
							<div>
								<p className="font-medium text-base text-black-dark pb-2">
									{row?.full_name}
								</p>
								<div className="flex items-center gap-2">
									<span className="text-sm text-base-gray">Credit score</span>
									<div className="border border-[#B7EB8F] rounded-sm px-2 text-xs text-green-medium bg-[#F6FFED] w-max">
										{row?.est_credit_score}
									</div>
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
									title="Are you sure delete this borrower?"
									content={row.full_name}
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
						<div className="pt-4 mt-4 border-t border-base-border flex items-center justify-between">
							<span className="text-base font-medium text-base-gray">SSN</span>
							<span className="text-base font-medium">
								{formatTINOrSSN(row?.ssn)}
							</span>
						</div>
						<div className="pt-4 flex items-center justify-between">
							<span className="text-base font-medium text-base-gray">Phone</span>
							<span className="text-base font-medium">{formatPhone(row?.phone)}</span>
						</div>
						<div className="py-4 border-b border-base-border flex items-center justify-between">
							<span className="text-base font-medium text-base-gray">Email</span>
							<span className="text-base font-medium">{row?.email}</span>
						</div>
						<div className="text-base font-medium  text-base-gray pt-4">
							Primary address
						</div>
						<div className="flex items-center gap-2 pt-3">
							<DistanceSvg />
							<span className="text-sm/6 text-black-dark">
								{assemblyAddress(row?.address, true)}
							</span>
						</div>
					</>
				);
			}}
		/>
	);
};
export default BorrowerManagement;
