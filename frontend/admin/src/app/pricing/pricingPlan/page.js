'use client';

import React, { useState } from 'react';
import { Button, Tabs } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useAuth } from '@repo/shared-methods/hooks';
import { eventBus } from '@repo/shared-methods/utils';
import { PERMISSIONS } from '@config/permissions';
import PlanFormDrawer from './components/PlanFormDrawer';
import PricingPlanTable from './components/PricingPlanTable';

const eventBusPrefix = 'searchPricingPlan-';

const PricingPlan = () => {
	const { hasPermission } = useAuth();
	const [tabKey, setTabKey] = useState('Borrower');
	const [open, setOpen] = useState(false);
	const canAdd = hasPermission(PERMISSIONS.PRICING_PLAN.ADD);

	return (
		<div className="flex-1 overflow-hidden relative flex flex-col">
			<Tabs
				activeKey={tabKey}
				items={[
					{
						key: 'Borrower',
						label: 'Borrower',
						children: (
							<PricingPlanTable
								filterWrapClassName="px-4 pb-4"
								type="Borrower"
								eventBusPrefix={eventBusPrefix}
							/>
						),
					},
					{
						key: 'Lender',
						label: 'Lender',
						children: (
							<PricingPlanTable
								filterWrapClassName="px-4 pb-4"
								type="Lender"
								eventBusPrefix={eventBusPrefix}
							/>
						),
					},
					{
						key: 'Admin',
						label: 'Admin',
						children: (
							<PricingPlanTable
								filterWrapClassName="px-4 pb-4"
								type="Admin"
								eventBusPrefix={eventBusPrefix}
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
			{canAdd && (
				<Button
					type="primary"
					onClick={() => setOpen(true)}
					className="!absolute right-4 top-6"
					size="large"
					icon={<PlusOutlined />}
				>
					Add plan
				</Button>
			)}
			{open && canAdd && (
				<PlanFormDrawer
					setOpen={setOpen}
					rowData={{ type: tabKey }}
					successCallback={() => {
						eventBus.emit(`${eventBusPrefix}${tabKey}`);
					}}
				/>
			)}
		</div>
	);
};

export default PricingPlan;
