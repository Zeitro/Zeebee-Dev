'use client';

import React, { useMemo } from 'react';
import { Form, Select, InputNumber } from 'antd';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { changeUnitByPropertyType, livingAreaProps } from '@repo/shared-methods/utils';

const PropertyDetailsForm = ({ rowData }) => {
	const { propertyTypeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	const form = Form.useFormInstance();
	const propertyType = Form.useWatch('property_type', form);
	const unit = Form.useWatch('unit', form);

	const unitInfo = useMemo(
		() =>
			changeUnitByPropertyType(propertyType ?? rowData?.property_type, unit ?? rowData?.unit),
		[propertyType, rowData?.property_type, unit, rowData?.unit],
	);

	return (
		<div className="grid w-full grid-cols-1 gap-x-6 sm:grid-cols-2 lg:grid-cols-3">
			<Form.Item
				label="Property type"
				name="property_type"
				rules={[{ required: true, message: 'Please select property type!' }]}
			>
				<Select
					options={propertyTypeOptions}
					placeholder="Select property type"
					onChange={(value) => {
						const info = changeUnitByPropertyType(value);
						form.setFieldValue('unit', info.unit);
					}}
				/>
			</Form.Item>
			<Form.Item
				label="Living area"
				name="living_area_sqft"
				rules={[{ required: true, message: 'Please enter Living area sqft!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					{...livingAreaProps}
					placeholder="Enter Living area sqft"
					suffix="sqft"
				/>
			</Form.Item>
			<Form.Item
				label="Units"
				name="unit"
				rules={[{ required: true, message: 'Please enter Number of units!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					min={unitInfo?.minUnit}
					max={unitInfo?.maxUnit}
					controls={false}
					placeholder="Enter number of units"
					disabled={unitInfo?.disabled}
				/>
			</Form.Item>
			<Form.Item
				label="Lot area"
				name="lot_area_sqft"
				rules={[{ required: true, message: 'Please enter Lot area sqft!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					{...livingAreaProps}
					placeholder="Enter Lot area sqft"
					suffix="sqft"
				/>
			</Form.Item>
			<Form.Item
				label="Beds"
				name="bedrooms"
				rules={[{ required: true, message: 'Please enter Number of bedrooms!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					min={0}
					max={100}
					controls={false}
					placeholder="Enter Number of bedrooms"
				/>
			</Form.Item>
			<Form.Item
				label="Baths"
				name="bathrooms"
				rules={[{ required: true, message: 'Please enter Number of bathrooms!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					min={0}
					max={100}
					controls={false}
					placeholder="Enter Number of bathrooms"
				/>
			</Form.Item>
		</div>
	);
};

export default PropertyDetailsForm;
