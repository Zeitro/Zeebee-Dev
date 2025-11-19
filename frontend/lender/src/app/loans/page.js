'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import TableLoan from './components/ListView';
import PipelinePendingLoan from './components/Pipeline';

const Loans = () => {
	const [tabKey, setTabKey] = useState('List');

	return <>{tabKey === 'List' ? <TableLoan /> : <PipelinePendingLoan />}</>;
};
export default Loans;
