import React, { useMemo, useState } from 'react';
import { Checkbox, List } from 'antd';
import CollapseList from './CollapseList';
import HeaderBtn from './HeaderBtn';
const QuotingPage = () => {
	const [list, setList] = useState([
		{
			label: 'Soft quote loan amount,rate, term and COE date with borrower',
			key: 1,
		},
	]);
	const isApprove = useMemo(() => list.every((item) => item.checked), [list]);
	return (
		<div>
			<HeaderBtn disabled={!isApprove}>Quoting</HeaderBtn>
			<div className="px-6">
				<p className="text-black-dark text-xl font-medium pt-6">CheckList</p>
				<CollapseList list={list} hideCollapsible={true} setList={setList} />
			</div>
		</div>
	);
};
export default QuotingPage;
