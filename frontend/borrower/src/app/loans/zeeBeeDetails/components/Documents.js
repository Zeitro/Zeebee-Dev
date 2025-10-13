import React, { useEffect, useState } from 'react';
import DocumentList from './DocumentList';
import { documentsColumns } from './DocumentList/documentColumns';
const Documents = ({ rowData }) => {
	const [lenderList, setLenderList] = useState([]);

	useEffect(() => {
		console.log(rowData?.loan_documents);

		if (rowData?.loan_documents?.length) {
			let lList = rowData?.loan_documents?.filter((item) =>
				['Borrower', 'Extend'].includes(item.category),
			);
			setLenderList(lList);
		}
	}, [rowData]);
	return (
		<div className="border border-base-border px-6 pb-6 rounded-2xl">
			<DocumentList columns={documentsColumns()} dataSource={[...lenderList]} />
		</div>
	);
};
export default Documents;
