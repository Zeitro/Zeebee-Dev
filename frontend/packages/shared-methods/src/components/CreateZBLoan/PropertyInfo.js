'use client';
import React, { useMemo, useState } from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Select, InputNumber, Radio } from 'antd';
import AddressForm from '../GoogleAutoComplete/AddressForm';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { layoutCol, titleCol } from '@repo/shared-methods/enum';
import { livingAreaProps } from '../../lib/format';
// import SelectCard from '@/components/SelectCard/index';
import { SelectCard } from '@repo/shared-methods';
let name = 'loan_property';

const SELECT_MODE = {
	property: 'property',
	transaction: 'transaction',
};

export const PropertyInfo = ({
	gutter = 24,
	addressProps,
	isRequestProperty,
	unitInfo,
	isShowSelectCard = true,
}) => {
	const form = Form.useFormInstance();
	const borrower_property_id = Form.useWatch([name, 'borrower_property_id'], form);
	const borrower_property_transaction_id = Form.useWatch(
		['borrower_property_transaction_id'],
		form,
	);
	const [manualSelectMode, setManualSelectMode] = useState(null);
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

	const resolvedSelectMode =
		manualSelectMode ||
		(borrower_property_transaction_id ? SELECT_MODE.transaction : SELECT_MODE.property);
	const isTransactionMode = resolvedSelectMode === SELECT_MODE.transaction;

	const selectCardProps = useMemo(() => {
		if (isTransactionMode) {
			return {
				type: 'transaction',
				api: '/borrower_property_transaction/all',
				viewMoreTitle: 'View all transactions',
				formFieldName: ['borrower_property_transaction_id'],
				selectedId: borrower_property_transaction_id,
				isHideAddBtn: true,
				title: '',
			};
		}
		return {
			type: 'property',
			api: '/property/all',
			title: '',
			addTitle: 'Add a property',
			formFieldName: [name, 'borrower_property_id'],
			selectedId: borrower_property_id,
		};
	}, [isTransactionMode, borrower_property_id, borrower_property_transaction_id]);

	const handleSelectModeChange = (mode) => {
		setManualSelectMode(mode);
	};
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Property information
			</Col>
			{isShowSelectCard && (
				<Col span={24}>
					<Form.Item
						noStyle
						name={[name, 'borrower_property_id']}
						rules={[{ required: true }]}
					></Form.Item>
					<Form.Item
						noStyle
						name={['borrower_property_transaction_id']}
						rules={[{ required: true }]}
					></Form.Item>
					<Radio.Group
						value={resolvedSelectMode}
						onChange={(e) => handleSelectModeChange(e.target.value)}
						disabled={isRequestProperty}
						size="large"
					>
						<Radio
							value={SELECT_MODE.property}
							className="[&_.ant-radio-label]:!font-medium [&_.ant-radio-label]:text-base !text-base-text"
						>
							Choose a property
						</Radio>
						<span className="text-base text-base-gray select-none px-2">or</span>
						<Radio
							value={SELECT_MODE.transaction}
							className="[&_.ant-radio-label]:!font-medium [&_.ant-radio-label]:text-base !text-base-text"
						>
							Choose a transaction
						</Radio>
					</Radio.Group>
					<SelectCard
						key={resolvedSelectMode}
						className="mb-8"
						type={selectCardProps.type}
						api={selectCardProps.api}
						title={selectCardProps.title}
						addTitle={selectCardProps.addTitle}
						name={selectCardProps.formFieldName}
						id={selectCardProps.selectedId}
						isHideAddBtn={selectCardProps.isHideAddBtn}
						changeFormValue={(values) => {
							if (isTransactionMode) {
								form.setFieldsValue({
									borrower_property_transaction_id: values.id,
									[name]: {
										borrower_property_id: values.borrower_property_id,
										...values?.borrower_property,
									},
								});
								return;
							}
							form.setFieldsValue({
								borrower_property_transaction_id: null,
								[name]: { ...values, borrower_property_id: values.id },
							});
						}}
						disabled={isRequestProperty}
					/>
				</Col>
			)}
			<AddressForm
				{...addressProps}
				addressName={[name, 'address', 'address1']}
				addressName2={[name, 'address', 'address2']}
				cityName={[name, 'address', 'city']}
				stateName={[name, 'address', 'state']}
				zipName={[name, 'address', 'zip_code']}
				addressLabel="Subject property address 1"
				addressLabel2="Subject property address 2"
				form={form}
				stateCol={{ ...layoutCol, md: 8 }}
				cityCol={{ ...layoutCol, md: 8 }}
				zipCol={{ ...layoutCol, md: 8 }}
				disabled={true}
				isRequest={false}
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
					<Select
						options={propertyTypeOptions}
						placeholder="Select property type"
						disabled
					/>
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
						disabled
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} md={12}>
				<Form.Item
					label="Number of units"
					name={[name, 'unit']}
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
						placeholder="Enter number of units"
						disabled={true}
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
						disabled
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
						disabled
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
