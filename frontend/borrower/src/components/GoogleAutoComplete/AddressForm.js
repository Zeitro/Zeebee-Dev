'use client';
import React, { useEffect } from 'react';
import { Form, Input, Col, Select } from 'antd';
import GoogleAutoComplete from './index';
import { set } from 'lodash';
import { useAddressStore } from '@/store/address';
import { useDisUpdateEffect } from '@/hooks/common';
const layoutCol = {
	xs: 24,
	sm: 24,
};
const getPath = (name, listName) => {
	const listNameArr = Array.isArray(listName) ? listName : [listName];
	let pathArr = Array.isArray(name) ? name : [name];
	if (listName) {
		pathArr = [...listNameArr, ...pathArr];
	}
	return pathArr;
};
/**
 * 地址自动完成（包括address1/address2/city/state/zip/unit）
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
	disabled,
	isShowAddress2 = true,
	isShowZip = true,
	listName = '', // 此组件使用在Form.List内时传此字段,且Form.List上的name属性值与此listName字段值一致；方便用于设置form.setFieldsValue
}) => {
	const { getStateList, getCityList, stateOptions, cityOptions } = useAddressStore();
	const stateValue = Form.useWatch(getPath(stateName, listName), form);
	const handlePlaceSelected = (place) => {
		if (!form) return;
		const newValue = form.getFieldsValue() || {};
		set(newValue, getPath(cityName, listName), place.city);
		set(newValue, getPath(stateName, listName), place.state);
		if (isShowZip) {
			set(newValue, getPath(zipName, listName), place.zip);
		}
		if (isShowAddress2) {
			set(newValue, getPath(addressName2, listName), place.unit);
		}
		form.setFieldsValue(newValue);
	};

	useEffect(() => {
		getStateList();
	}, []);

	useDisUpdateEffect(() => {
		if (stateValue) {
			getCityList(stateValue);
		}
	}, [stateValue]);

	return (
		<>
			<Col {...layoutCol} {...addressCol}>
				<Form.Item
					name={addressName}
					label={addressLabel}
					rules={[{ required: true, message: 'Please enter your address' }]}
				>
					<GoogleAutoComplete
						onPlaceSelected={handlePlaceSelected}
						placeholder={`Enter ${addressLabel}`}
						disabled={disabled}
					/>
				</Form.Item>
			</Col>
			{isShowAddress2 && (
				<Col {...layoutCol} {...address2Col}>
					<Form.Item name={addressName2} label={addressLabel2}>
						<Input placeholder={`Enter ${addressLabel2}`} disabled={disabled} />
					</Form.Item>
				</Col>
			)}
			<Col {...layoutCol} {...stateCol}>
				<Form.Item
					name={stateName}
					label={stateLabel}
					rules={[{ required: true, message: 'Please select state!' }]}
				>
					<Select
						options={stateOptions}
						optionFilterProp="label"
						placeholder="Select state"
						showSearch
						disabled={disabled}
					/>
				</Form.Item>
			</Col>
			<Col md={12} {...layoutCol} {...cityCol}>
				<Form.Item
					name={cityName}
					label={cityLabel}
					rules={[{ required: true, message: 'Please select city!' }]}
				>
					<Select
						options={cityOptions}
						optionFilterProp="label"
						placeholder="Select city"
						showSearch
						disabled={disabled}
					/>
				</Form.Item>
			</Col>
			{isShowZip && (
				<Col md={12} {...layoutCol} {...zipCol}>
					<Form.Item
						name={zipName}
						label={zipLabel}
						rules={[{ required: true, message: 'Please enter zip code!' }]}
					>
						<Input placeholder="Enter zip" disabled={disabled} />
					</Form.Item>
				</Col>
			)}
		</>
	);
};

export default AddressForm;
