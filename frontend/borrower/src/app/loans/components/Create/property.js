'use client';
import { Col, Form, Input, InputNumber, Row, Select } from 'antd';
import styles from './index.module.scss';
import { useGetOptions } from '@repo/shared-methods/hooks';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { loan_pending } from '@/enum/common';
import { InputNumAmount, CommonFormDate } from '@repo/shared-methods';
import SelectCard from '@/components/SelectCard';
const layoutCol = {
	xs: 24,
	sm: 24,
};
const CreateProperty = ({
	stage,
	stepBtn,
	gutter,
	addressProps = {},
	isRequestProperty,
	unitInfo,
	className,
	selectCardCol = { xl: 6, sm: 24, xs: 24 },
	btnCol = { xl: 6, sm: 24, xs: 24 },
	sliceIndex = 3,
	...colProps
}) => {
	const form = Form.useFormInstance();
	const borrower_property_id = Form.useWatch(['loan_property', 'borrower_property_id'], form);
	const { propertyTypeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	return (
		<Row className={[styles.create, className].join(' ')} gutter={gutter}>
			<SelectCard
				className="mb-8"
				changeFormValue={(values) => {
					const property = {
						loan_property: { ...values, borrower_property_id: values.id },
					};
					form.setFieldsValue(property);
				}}
				name={['loan_property', 'borrower_property_id']}
				id={borrower_property_id}
				disabled={isRequestProperty}
				sliceIndex={sliceIndex}
				{...selectCardCol}
				btnCol={btnCol}
			/>
			<Col {...layoutCol} {...(colProps || {})}>
				<Form.Item
					label="Property type"
					name={['loan_property', 'property_type']}
					rules={[{ required: true, message: 'Please select Property type!' }]}
				>
					<Select
						options={propertyTypeOptions}
						placeholder="Select Property type"
						// disabled={isRequestProperty}
						disabled={true}
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol} {...(colProps || {})}>
				<Form.Item
					label="Number of units"
					name={['loan_property', 'unit']}
					rules={[{ required: true, message: 'Please input!' }]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={unitInfo.minUnit}
						max={unitInfo.maxUnit}
						controls={false}
						placeholder="Enter number of units"
						disabled={true}
					/>
				</Form.Item>
			</Col>
			<AddressForm
				{...addressProps}
				addressName={['loan_property', 'address', 'address1']}
				addressName2={['loan_property', 'address', 'address2']}
				cityName={['loan_property', 'address', 'city']}
				stateName={['loan_property', 'address', 'state']}
				zipName={['loan_property', 'address', 'zip_code']}
				form={form}
				disabled={true}
			/>
			{stage === loan_pending && (
				<>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Purchase date" name={['loan_property', 'purchase_date']}>
							<CommonFormDate />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Owner" name={['loan_property', 'owner']}>
							<Input placeholder="Enter owner" />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Purchase price"
							name={['loan_property', 'purchase_amount']}
						>
							<InputNumAmount placeholder="Enter purchase price" />
						</Form.Item>
					</Col>
				</>
			)}
			<Col {...layoutCol} {...(colProps || {})}>
				<Form.Item
					label="Estimated after repair value (ARV)" // 修复后预计估值
					name={['loan_property', 'arv']}
					rules={[
						{
							required: true,
							message: 'Please enter estimated ARV!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter estimated ARV" />
				</Form.Item>
			</Col>
			{stepBtn}
		</Row>
	);
};
export default CreateProperty;
