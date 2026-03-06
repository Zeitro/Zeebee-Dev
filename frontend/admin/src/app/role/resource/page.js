'use client';
import { Button, Tabs } from 'antd';
import React, { useState } from 'react';
import TreeTable from '../components/TreeTable';
import ResourceMain from '../components/ResourceMain';
import { PlusOutlined } from '@ant-design/icons';
const Resource = ({ className }) => {
	const [tabKey, setTabKey] = useState('Borrower');
	const [open, setOpen] = useState(false);
	const eventBusPrefix = 'searchResource-';
	const api = '/resource/tree';
	const menuKey = 'Resource';
	return (
		<div className="relative flex h-full flex-col">
			<Tabs
				activeKey={tabKey}
				items={[
					{
						key: 'Borrower',
						label: 'Borrower',
						children: (
							<TreeTable
								filterWrapClassName="px-4 pb-4"
								type="Borrower"
								eventBusPrefix={eventBusPrefix}
								menuKey={menuKey}
								api={api}
							/>
						),
					},
					{
						key: 'Lender',
						label: 'Lender',
						children: (
							<TreeTable
								filterWrapClassName="px-4 pb-4"
								type="Lender"
								eventBusPrefix={eventBusPrefix}
								menuKey={menuKey}
								api={api}
							/>
						),
					},
					{
						key: 'Admin',
						label: 'Admin',
						children: (
							<TreeTable
								filterWrapClassName="px-4 pb-4"
								type="Admin"
								eventBusPrefix={eventBusPrefix}
								menuKey={menuKey}
								api={api}
							/>
						),
					},
				]}
				onChange={(v) => setTabKey(v)}
				className={`flex-1 overflow-hidden !pt-4 !px-4 [&_.ant-tabs-content]:h-full [&_.ant-tabs-tabpane]:h-full ${className}`}
				destroyOnHidden
				size="large"
			/>
			<Button
				type="primary"
				onClick={() => {
					setOpen(true);
				}}
				className="!absolute right-4 top-6"
				size="large"
				icon={<PlusOutlined />}
			>
				Add
			</Button>
			{open && (
				<ResourceMain
					setOpen={setOpen}
					rowData={{ type: tabKey }}
					eventBusPrefix={eventBusPrefix}
					api={api}
				/>
			)}
		</div>
	);
};
export default Resource;
