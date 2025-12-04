'use client';
import { Radio } from 'antd';
import React, { useState } from 'react';
import TableLoan from './components/ListView';
import PipelinePendingLoan from './components/Pipeline';
import NewHeader from '@/components/NewHeader';
const Loans = () => {
	const [type, setType] = useState('list');

	return (
		<>
			<NewHeader isShowCrumbs>
				<Radio.Group
					value={type}
					buttonStyle="solid"
					onChange={(e) => setType(e.target.value)}
				>
					<Radio.Button value="list">List view</Radio.Button>
					<Radio.Button value="pipeline">Pipeline view</Radio.Button>
				</Radio.Group>
			</NewHeader>
			{type === 'list' ? <TableLoan /> : <PipelinePendingLoan />}
		</>
	);
};
export default Loans;
