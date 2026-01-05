'use client';
import { Radio } from 'antd';
import React, { useState } from 'react';
import TableLoan from './components/ListView';
import PipelinePendingLoan from './components/Pipeline';
import { NewHeader } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
const Loans = () => {
	const { cancelAll } = createRequest();
	const [type, setType] = useState('list');

	return (
		<>
			<NewHeader isShowCrumbs>
				<Radio.Group
					value={type}
					buttonStyle="solid"
					onChange={(e) => {
						cancelAll();
						setType(e.target.value);
					}}
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
