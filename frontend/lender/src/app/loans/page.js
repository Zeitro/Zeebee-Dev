'use client';
import { Tabs } from 'antd';
import React, { useState } from 'react';
import TableLoan from './components/ListView';
import PipelinePendingLoan from './components/Pipeline';

const Loans = () => {
	const [tabKey, setTabKey] = useState('List');

	return (
		<Tabs
			activeKey={tabKey}
			items={[
				{
					key: 'List',
					children: <TableLoan />,
				},
				{
					key: 'Pipeline',
					children: <PipelinePendingLoan />,
				},
			]}
			onChange={(v) => setTabKey(v)}
			// className={`h-full ${styles.loans}`}
			renderTabBar={() => <div> </div>}
			destroyOnHidden
		/>
	);
};
export default Loans;
