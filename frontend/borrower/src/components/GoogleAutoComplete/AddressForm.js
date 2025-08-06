'use client';
import React from 'react';
import { Form, Input, Col, Select } from 'antd';
import GoogleAutoComplete from './index';
import { set } from 'lodash';
/**
 * 地址自动完成（包括address1/address2/city/state/zip）
 */
const AddressForm = ({
	form,
	addressName = 'address',
	addressLabel = 'Street address 1',
	addressName2 = 'address',
	addressLabel2 = 'Street address 2',
	cityName = 'city',
	cityLabel = 'City',
	stateName = 'state',
	stateLabel = 'State',
	zipName = 'zip_code',
	zipLabel = 'Zip',
	addressCol = {},
	address2Col = {},
	cityCol = {},
	stateCol = {},
	zipCol = {},
}) => {
	const handlePlaceSelected = (place) => {
		if (!form) return;
		const newValue = form.getFieldsValue() || {};
		set(newValue, cityName, place.city);
		set(newValue, stateName, place.state);
		set(newValue, zipName, place.zip);
		set(newValue, addressName2, place.unit);
		form.setFieldsValue(newValue);
	};

	return (
		<>
			<Col md={24} {...addressCol}>
				<Form.Item
					name={addressName}
					label={addressLabel}
					rules={[{ required: true, message: 'Please enter your address' }]}
				>
					<GoogleAutoComplete
						onPlaceSelected={handlePlaceSelected}
						placeholder={`Enter ${addressLabel}`}
					/>
				</Form.Item>
			</Col>
			<Col md={24} {...address2Col}>
				<Form.Item name={addressName2} label={addressLabel2}>
					<Input placeholder={`Enter ${addressLabel2}`} />
				</Form.Item>
			</Col>
			<Col md={24} {...cityCol} >
				<Form.Item name={cityName} label={cityLabel} rules={[{ required: true, message: 'Please enter city!' }]}>
					<Input placeholder="Enter city" />
				</Form.Item>
			</Col>
			<Col md={12} {...stateCol}>
				<Form.Item
					name={stateName}
					label={stateLabel}
					rules={[{ required: true, message: 'Please enter state!' }]}
				>
					<Select options={[]} placeholder="Enter state" />
				</Form.Item>
			</Col>
			<Col
				md={12}
				{...zipCol}
			>
				<Form.Item name={zipName} label={zipLabel} rules={[{ required: true, message: 'Please enter zip code!' }]}>
					<Input placeholder="Enter zip" />
				</Form.Item>
			</Col>
		</>
	);
};

export default AddressForm;
