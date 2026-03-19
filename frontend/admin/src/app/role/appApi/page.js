'use client';
import { Button, Tabs } from 'antd';
import React, { useState } from 'react';
import TreeTable from '../components/TreeTable';
import ApiMainModal from '../components/ApiMain';
import { PlusOutlined } from '@ant-design/icons';
import { PERMISSIONS } from '@config/permissions';
import { PermissionSlot } from '@repo/shared-methods';
const AppApi = () => {
	const [tabKey, setTabKey] = useState('Borrower');
	const [open, setOpen] = useState(false);
	const eventBusPrefix = 'searchApi-';
	const menuKey = 'Api';
	const api = '/api/tree';
	return (
		<div className="flex-1 overflow-hidden relative flex h-full flex-col">
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
								scrollX={1400}
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
								scrollX={1400}
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
								scrollX={1400}
							/>
						),
					},
				]}
				onChange={(v) => setTabKey(v)}
				className={
					'flex-1 overflow-hidden !pt-4 !px-4 [&_.ant-tabs-content]:h-full [&_.ant-tabs-tabpane]:h-full'
				}
				destroyOnHidden
				size="large"
			/>
			<PermissionSlot code={PERMISSIONS.API.ADD}>
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
			</PermissionSlot>
			{open && (
				<ApiMainModal
					setOpen={setOpen}
					rowData={{ type: tabKey }}
					eventBusPrefix={eventBusPrefix}
				/>
			)}
		</div>
	);
};
export default AppApi;
