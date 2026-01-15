'use client';
import React, { useMemo, useState } from 'react';
import DocumentList from '../DocumentList';
import { borrowerUploadColumns, lenderUploadColumns } from '../documentColumns';
import { message } from 'antd';
import HeaderBtn from '../HeaderBtn';
import { AddDocument } from '../AddDocument';
import { SendRequestBtn } from '../SendRequestModal';
import { useSearchParams } from 'next/navigation';
import { eventBus, createRequest } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const Processing = ({
	isHide,
	title = 'Processing',
	isHaveRejectBtn = false,
	disabledOperate,
	rowData,
	id,
	setNextStep,
	setRejectStep,
}) => {
	const searchParams = useSearchParams();
	const loan_id = Number(searchParams.get('id'));
	const [borrowerList, setBorrowerList] = useState([]);
	const [lenderList, setLenderList] = useState([]);
	const [deleteLoading, setDeleteLoading] = useState(false);
	const [loading, setLoading] = useState(true);
	const { request } = createRequest();
	const sendList = useMemo(
		() => borrowerList.filter((item) => item.status === 'Pending request') || [],
		[borrowerList],
	);

	useDeferredSyncEffect({
		notPassedCallback: () => setLoading(false),
		predicate: () => rowData?.loan_documents?.length,
		deferred: () => {
			let bList = [],
				lList = [];
			setLoading(true);
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
			setLoading(false);
		},
		deps: [rowData],
		delay: 0,
	});

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
			<HeaderBtn
				isHaveRejectBtn={isHaveRejectBtn}
				disabledOperate={disabledOperate}
				id={id}
				pending_status={title}
				setNextStep={setNextStep}
				setRejectStep={setRejectStep}
			>
				{title}
			</HeaderBtn>
			<div className="px-[34px] py-6">
				<p className="text-xl font-medium text-black-dark pb-6">Document</p>
				{!isHide && (
					<div className="flex justify-between items-center">
						<AddDocument
							// disabled={disabledOperate}
							loan_id={loan_id}
							category="Extend"
						/>
						<SendRequestBtn list={sendList} size="middle" disabled={!sendList.length}>
							Send multiple request
						</SendRequestBtn>
					</div>
				)}
				<DocumentList
					columns={borrowerUploadColumns(isHide, handleClickOpt, deleteLoading)}
					dataSource={[...borrowerList]}
					isHide={isHide}
					loading={loading}
				/>
				<DocumentList
					columns={lenderUploadColumns(isHide)}
					dataSource={[...lenderList]}
					isHide={isHide}
					loading={loading}
				/>
			</div>
		</div>
	);
};
export default Processing;
