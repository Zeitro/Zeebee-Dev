import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Select, InputNumber } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { useGetOptions } from '@/hooks/useGetOptions';
import { layoutCol, titleCol } from '@/enum/common';
import { livingAreaProps } from '@/lib/format';
let name = 'loan_property';
const PropertyInfo = ({ gutter = 24, form, addressProps, isRequestProperty, unitInfo }) => {
	const { loanPurposeOptions, loanTypeOptions, propertyTypeOptions } = useGetOptions([
		{
			key: 'loan_purpose',
			valueName: 'loanPurposeOptions',
		},
		{
			key: 'loan_type',
			valueName: 'loanTypeOptions',
		},
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);

	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Property information
			</Col>
			<AddressForm
				{...addressProps}
				addressName={[name, 'property_base', 'address', 'address1']}
				addressName2={[name, 'property_base', 'address', 'address2']}
				cityName={[name, 'property_base', 'address', 'city']}
				stateName={[name, 'property_base', 'address', 'state']}
				zipName={[name, 'property_base', 'address', 'zip_code']}
				addressLabel="Subject property address 1"
				addressLabel2="Subject property address 2"
				form={form}
				stateCol={{ ...layoutCol, md: 8 }}
				cityCol={{ ...layoutCol, md: 8 }}
				zipCol={{ ...layoutCol, md: 8 }}
				disabled={isRequestProperty}
			/>
			<Col {...layoutCol}>
				<Form.Item
					label="Loan purpose"
					name="loan_purpose"
					rules={[{ required: true, message: 'Please select loan purpose!' }]}
				>
					<Select options={loanPurposeOptions} placeholder="Select loan purpose" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Loan type"
					name="loan_type"
					rules={[{ required: true, message: 'Please select loan type!' }]}
				>
					<Select options={loanTypeOptions} placeholder="Select loan type" />
				</Form.Item>
			</Col>
			<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Current status
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Property type"
					name={[name, 'property_type']}
					rules={[{ required: true, message: 'Please select property type!' }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select property type" />
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Living area sqft"
					name={[name, 'living_area_sqft']}
					rules={[
						{
							required: true,
							message: 'Please enter Living area sqft!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						{...livingAreaProps}
						placeholder="Enter Living area sqft"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of units"
					name={[name, 'property_base', 'unit']}
					rules={[
						{
							required: true,
							message: 'Please enter Number of units!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={unitInfo?.unit?.minUnit}
						max={unitInfo?.unit?.maxUnit}
						controls={false}
						placeholder="Enter number of units"
						disabled={isRequestProperty || unitInfo?.unit?.disabled}
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of bedrooms"
					name={[name, 'bedrooms']}
					rules={[
						{
							required: true,
							message: 'Please enter Number of bedrooms!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={100}
						controls={false}
						placeholder="Enter Number of bedrooms"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of bathrooms"
					name={[name, 'bathrooms']}
					rules={[
						{
							required: true,
							message: 'Please enter Number of bathrooms!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={100}
						controls={false}
						placeholder="Enter Number of bathrooms"
					/>
				</Form.Item>
			</Col>
			<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Future status
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Property type"
					name={[name, 'future_property_type']}
					rules={[{ required: true, message: 'Please select property type!' }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select property type" />
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Living area sqft"
					name={[name, 'future_living_area_sqft']}
					rules={[
						{
							required: true,
							message: 'Please enter Living area sqft!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						{...livingAreaProps}
						placeholder="Enter Living area sqft"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of units"
					name={[name, 'future_unit']}
					rules={[
						{
							required: true,
							message: 'Please enter Number of units!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						controls={false}
						min={unitInfo?.future_unit?.minUnit}
						max={unitInfo?.future_unit?.maxUnit}
						placeholder="Enter number of units"
						disabled={unitInfo?.future_unit?.disabled}
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of bedrooms"
					name={[name, 'future_bedrooms']}
					rules={[
						{
							required: true,
							message: 'Please enter Number of bedrooms!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={100}
						controls={false}
						placeholder="Enter Number of bedrooms"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of bathrooms"
					name={[name, 'future_bathrooms']}
					rules={[
						{
							required: true,
							message: 'Please enter Number of bathrooms!',
						},
					]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={100}
						controls={false}
						placeholder="Enter Number of bathrooms"
					/>
				</Form.Item>
			</Col>
		</Row>
	);
};
export default PropertyInfo;
