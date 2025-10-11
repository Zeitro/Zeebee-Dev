'use client';
import React, { useEffect, useState } from 'react';
import DocumentList from '../DocumentList';
import { borrowerUploadColumns, lenderUploadColumns } from '../documentColumns';
import { message } from 'antd';
import HeaderBtn from '../HeaderBtn';
import { AddDocument } from '../AddDocument';
import { SendRequestBtn } from '../SendRequestModal';
import { useSearchParams } from 'next/navigation';
import { createRequest } from '@/lib/fetch';
import eventBus from '@/lib/bus';
const Processing = ({
	isHide,
	title = 'Processing',
	isHaveRejectBtn = false,
	disabledOperate,
	rowData,
}) => {
	const searchParams = useSearchParams();
	const loan_id = Number(searchParams.get('id'));
	const [borrowerList, setBorrowerList] = useState([]);
	const [lenderList, setLenderList] = useState([]);
	const [deleteLoading, setDeleteLoading] = useState(false);
	const { request } = createRequest();

	useEffect(() => {
		if (rowData?.loan_documents?.length) {
			let bList = [],
				lList = [];
			rowData?.loan_documents?.map((item) => {
				if (['Borrower', 'Extend'].includes(item.category)) {
					bList.push(item);
				}
				if (item.category === 'Lender') {
					lList.push(item);
				}
			});
			setBorrowerList(bList);
			setLenderList(lList);
		}
	}, [rowData]);

	const handleDelete = async (record) => {
		try {
			setDeleteLoading(true);
			await request(`/loan_document/${record.id}/request`, {
				method: 'DELETE',
			});
			message.open({
				type: 'success',
				content: 'Delete Document request successfully',
			});
			eventBus.emit('getLoanDetail');
			setDeleteLoading(false);
		} catch (error) {
			setDeleteLoading(false);
		}
	};

	const handleClickOpt = (key, record) => {
		switch (key) {
			case 'delete':
				handleDelete(record);
				break;

			default:
				break;
		}
	};

	return (
		<div>
			<HeaderBtn isHaveRejectBtn={isHaveRejectBtn} disabledOperate={disabledOperate}>
				{title}
			</HeaderBtn>
			<div className="px-[34px] py-6">
				<p className="text-xl font-medium text-black-dark pb-6">Document</p>
				{!isHide && (
					<div className="flex justify-between items-center">
						<AddDocument disabled={disabledOperate} loan_id={loan_id} />
						<SendRequestBtn
							list={
								borrowerList.filter((item) => item.status === 'Pending request') ||
								[]
							}
							size="middle"
							disabled={disabledOperate}
						>
							Send multiple request
						</SendRequestBtn>
					</div>
				)}
				<DocumentList
					columns={borrowerUploadColumns(isHide, handleClickOpt, deleteLoading)}
					dataSource={[...borrowerList]}
					isHide={isHide}
				/>
				<DocumentList
					columns={lenderUploadColumns(isHide)}
					dataSource={[...lenderList]}
					isHide={isHide}
				/>
			</div>
		</div>
	);
};
export default Processing;
