import React from 'react';
import DocumentList from '../DocumentList';
import {
	borrowerUploadColumns,
	LenderUploadColumns,
	lenderClosingColumns,
} from '../documentColumns';
import { Button, List, Space } from 'antd';
import HeaderBtn from '../HeaderBtn';
import { AddDocument } from '../AddDocument';
const Processing = () => {
	const lenderList = [
		{
			type: 'Closing docs',
			status: 'Pending',
		},
		{
			type: 'title docs',
			status: 'Pending',
		},
		{
			type: 'formal LOI',
			status: 'Pending',
		},
	];
	return (
		<div>
			<HeaderBtn>Closing</HeaderBtn>
			<div className="px-[34px] py-6">
				<AddDocument />
				<DocumentList columns={lenderClosingColumns()} dataSource={[...lenderList]} />
			</div>
		</div>
	);
};
export default Processing;
