import React, { useState } from 'react';
import DocumentList from './DocumentList';
import { documentsColumns } from './DocumentList/documentColumns';
import { NoData } from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const Documents = ({ rowData }) => {
	const [lenderList, setLenderList] = useState([]);

	useDeferredSyncEffect({
		predicate: () => rowData?.loan_documents?.length,
		deferred: () => {
			let lList = rowData?.loan_documents?.filter(
				(item) =>
					['Borrower', 'Extend'].includes(item.category) &&
					item.status !== 'Pending request',
			);
			setLenderList(lList);
		},
		deps: [rowData],
		delay: 0,
	});
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
