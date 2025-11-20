import React, { useEffect, useState } from 'react';
import DocumentList from './DocumentList';
import { documentsColumns } from './DocumentList/documentColumns';
import NoData from '@/components/NoData';
const Documents = ({ rowData }) => {
	const [lenderList, setLenderList] = useState([]);

	useEffect(() => {
		if (rowData?.loan_documents?.length) {
			let lList = rowData?.loan_documents?.filter(
				(item) =>
					['Borrower', 'Extend'].includes(item.category) &&
					item.status !== 'Pending request',
			);
			setLenderList(lList);
		}
	}, [rowData]);
	return (
		<div className="border border-base-border px-6 pb-6 rounded-2xl">
			{lenderList.length ? (
				<DocumentList
					columns={documentsColumns()}
					dataSource={[...lenderList]}
					isHide={true}
				/>
			) : (
				<NoData />
			)}
		</div>
	);
};
export default Documents;
