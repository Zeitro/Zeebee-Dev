import React from 'react';
import DocumentList from '../DocumentList';
import { borrowerUploadColumns, lenderUploadColumns } from '../documentColumns';
import { Button, List, Space } from 'antd';
import HeaderBtn from '../HeaderBtn';
import { AddDocument } from '../AddDocument';
import { SendRequestBtn } from '../SendRequestModal';
const Processing = ({ isHide, title = 'Processing', isHaveRejectBtn = false }) => {
	const borrowerList = [
		{
			type: 'Valid ID of borrower',
			status: 'Pending',
			id: 1,
			children: [
				{
					title: 'ssss',
					id: 2,
				},
				{
					title: '3',
					id: 3,
				},
			],
		},
		{
			type: 'Purchase agreement',
			status: 'Pending',
		},
		{
			type: 'Evidence of Insurance',
			status: 'Pending',
		},
		{
			type: 'Preliminary Title report/escrow contact',
			status: 'Pending',
		},
		{
			type: 'RE experience',
			status: 'Pending',
		},
		{
			type: 'Bank statements for most recent two months',
			status: 'Pending',
		},
	];
	const lenderList = [
		{
			type: 'Borrower‘s credit score',
			status: 'Pending',
		},
		{
			type: 'Borrower’s default/fraud/risk/view/criminal records and RE experience',
			status: 'Pending',
		},
		{
			type: 'property data/historical property ownership',
			status: 'Pending',
		},
		{
			type: 'Title report',
			status: 'Pending',
		},
	];
	return (
		<div>
			<HeaderBtn isHaveRejectBtn={isHaveRejectBtn}>{title}</HeaderBtn>
			<div className="px-[34px] py-6">
				<p className="text-xl font-medium text-black-dark pb-6">Document</p>
				{!isHide && (
					<div className="flex justify-between items-center">
						<AddDocument />
						<SendRequestBtn
							list={borrowerList.map((item) => ({ document_type: item.type }))}
							size="middle"
						>
							Send multiple request
						</SendRequestBtn>
					</div>
				)}
				<DocumentList
					columns={borrowerUploadColumns(isHide)}
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
