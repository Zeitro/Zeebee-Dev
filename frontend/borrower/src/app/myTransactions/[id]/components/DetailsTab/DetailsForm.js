'use client';

import React from 'react';
import { Form, Input, Select } from 'antd';
import { CommonFormDate, InputNumAmount } from '@repo/shared-methods';
import { useGetOptions } from '@repo/shared-methods/hooks';

const DetailsForm = () => {
	const { propertyTypeOptions, businessStrategyOptions } = useGetOptions([
		{ key: 'property_type', valueName: 'propertyTypeOptions' },
		{ key: 'business_strategy', valueName: 'businessStrategyOptions' },
	]);

	return (
		<div className="flex w-full flex-col gap-6">
			<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-3">
				<Form.Item label="COE" name="coe_date">
					<CommonFormDate size="middle" />
				</Form.Item>
				<Form.Item label="Purchase price" name="purchase_price">
					<InputNumAmount placeholder="Enter purchase price" size="middle" />
				</Form.Item>
				<Form.Item label="EMD" name="emd">
					<InputNumAmount placeholder="Enter EMD" size="middle" />
				</Form.Item>
			</div>
			<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-3">
				<Form.Item label="Property type" name="property_type">
					<Select
						options={propertyTypeOptions}
						placeholder="Select property type"
						allowClear
					/>
				</Form.Item>
				<Form.Item label="PM" name="pm">
					<Input placeholder="Enter PM" />
				</Form.Item>
				<Form.Item label="Strategy" name="business_strategy">
					<Select
						options={businessStrategyOptions}
						placeholder="Select strategy"
						allowClear
					/>
				</Form.Item>
			</div>
			<Form.Item label="TC Notes" name="tc_notes" className="mb-0">
				<Input.TextArea placeholder="Enter TC notes" rows={5} />
			</Form.Item>
		</div>
	);
};

export default DetailsForm;
