import React, { useEffect, useState } from 'react';
import DocumentList from '../DocumentList';
import {
	borrowerUploadColumns,
	LenderUploadColumns,
	lenderClosingColumns,
} from '../documentColumns';
import { Button, List, Space } from 'antd';
import HeaderBtn from '../HeaderBtn';
import { AddDocument } from '../AddDocument';
const Closing = ({ disabledOperate, rowData }) => {
	const [lenderList, setLenderList] = useState([]);

	useEffect(() => {
		if (rowData?.loan_documents?.length) {
			let lList = rowData?.loan_documents?.filter(
				(item) => item.category === 'LenderClosing',
			);
			setLenderList(lList);
		}
	}, [rowData]);
	return (
		<div>
			<HeaderBtn disabledOperate={disabledOperate}>Closing</HeaderBtn>
			<div className="px-[34px] py-6">
				<AddDocument disabled={disabledOperate} />
				<DocumentList columns={lenderClosingColumns()} dataSource={[...lenderList]} />
			</div>
		</div>
	);
};
export default Closing;
