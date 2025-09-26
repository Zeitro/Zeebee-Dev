import React, { useState } from 'react';
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
	return (
		<div>
			<HeaderBtn>Quoting</HeaderBtn>
			<div className="px-6">
				<CollapseList list={list} />
			</div>
		</div>
	);
};
export default QuotingPage;
