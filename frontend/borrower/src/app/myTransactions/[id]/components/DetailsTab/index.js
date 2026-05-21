'use client';

import React from 'react';
import { EnvironmentOutlined, SwapOutlined } from '@ant-design/icons';
import DetailField from '@/app/properties/[id]/components/DetailsTab/DetailField';
import DetailSection from '@/app/properties/[id]/components/DetailsTab/DetailSection';
import { PROPERTY_FIELDS, TRANSACTION_FIELDS } from './constants';

function renderFields(fields, data) {
	return fields.map((field) => (
		<DetailField key={field.label} label={field.label} value={field.getValue(data)} />
	));
}

const DetailsTab = ({ rowData }) => {
	return (
		<div className="rounded-lg border border-base-border bg-background p-4">
			<DetailSection icon={<SwapOutlined />} title="Transaction">
				<div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{renderFields(TRANSACTION_FIELDS, rowData)}
				</div>
			</DetailSection>

			<DetailSection icon={<EnvironmentOutlined />} title="Property" bordered={false}>
				<div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{renderFields(PROPERTY_FIELDS, rowData)}
				</div>
			</DetailSection>
		</div>
	);
};

export default DetailsTab;
