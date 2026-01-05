'use client';
import React, { useEffect, useState } from 'react';
import DocumentList from '../DocumentList';
import { lenderClosingColumns } from '../documentColumns';
import HeaderBtn from '../HeaderBtn';
import { AddDocument } from '../AddDocument';
import { useSearchParams } from 'next/navigation';
const Closing = ({ disabledOperate, rowData, id }) => {
	const [lenderList, setLenderList] = useState([]);
	const searchParams = useSearchParams();
	const loan_id = Number(searchParams.get('id'));
	useEffect(() => {
		if (rowData?.loan_documents?.length) {
			let lList = rowData?.loan_documents?.filter((item) =>
				['LenderClosingExtend', 'LenderClosing'].includes(item.category),
			);
			setLenderList(lList);
		}
	}, [rowData]);
	return (
		<div>
			<HeaderBtn disabledOperate={disabledOperate} id={id} pending_status="Closing">
				Closing
			</HeaderBtn>
			<div className="px-[34px] py-6">
				<AddDocument
					// disabled={disabledOperate}
					category="LenderClosingExtend"
					loan_id={loan_id}
					isHideAddRequest={true}
				/>
				<DocumentList columns={lenderClosingColumns()} dataSource={[...lenderList]} />
			</div>
		</div>
	);
};
export default Closing;
