import React from 'react';
import DocumentList from '../DocumentList';
import {
	borrowerUploadColumns,
	LenderUploadColumns,
	lenderClosingColumns,
} from '../documentColumns';
import { Button, List, Space } from 'antd';
import HeaderBtn from '../HeaderBtn';
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
			<Button>Submit</Button>
			<DocumentList columns={lenderClosingColumns()} dataSource={[...lenderList]} />
		</div>
	);
};
export default Processing;
