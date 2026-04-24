import React from 'react';
// import styles from './index.module.scss';
import { Col, Form, InputNumber, Radio, Row, Select } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { layoutCol, titleCol, yesNoOptions } from '@/enum/common';

let name = 'property';
const PropertyInfo = ({ gutter = 24, unitInfo }) => {
	const form = Form.useFormInstance();
	const loan_type = Form.useWatch(['loan', 'loan_type'], form);
	const { propertyTypeOptions, loanTypeOptions, loanPurposeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
		{
			key: 'loan_type',
			valueName: 'loanTypeOptions',
		},
		{
			key: 'loan_purpose',
			valueName: 'loanPurposeOptions',
		},
	]);

	return (
		<Row gutter={gutter}>
			<Col {...titleCol} className="!text-3xl/10 font-medium pb-8">
				Property information
			</Col>
			<AddressForm
				addressName={[name, 'address', 'address1']}
				addressName2={[name, 'address', 'address2']}
				cityName={[name, 'address', 'city']}
				stateName={[name, 'address', 'state']}
				zipName={[name, 'address', 'zip_code']}
				form={form}
				stateCol={{ ...titleCol }}
				cityCol={{ ...layoutCol, lg: 12, md: 24 }}
				zipCol={{ ...layoutCol, lg: 12, md: 24 }}
			/>

			<Col {...titleCol}>
				<Form.Item
					label="Loan purpose"
					name={['loan', 'loan_purpose']}
					rules={[{ required: true }]}
				>
					<Select options={loanPurposeOptions} placeholder="Loan purpose" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Loan type"
					name={['loan', 'loan_type']}
					rules={[{ required: true }]}
				>
					<Select
						options={loanTypeOptions.filter((item) => item.value !== 'Rental')}
						placeholder="Loan type"
					/>
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Property type"
					name={[name, 'property_type']}
					rules={[{ required: true }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select property type" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Number of units"
					name={[name, 'unit']}
					rules={[{ required: true }]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={unitInfo.minUnit}
						max={unitInfo.maxUnit}
						controls={false}
						placeholder="Enter number of units"
						disabled={unitInfo.disabled}
					/>
				</Form.Item>
			</Col>
			{loan_type === 'Ground up construction' && (
				<Col {...titleCol}>
					<Form.Item
						label="Permit ready"
						name={[name, 'is_permit_ready']}
						rules={[{ required: true }]}
					>
						<Radio.Group block options={yesNoOptions} style={{ width: '100%' }} />
					</Form.Item>
				</Col>
			)}
		</Row>
	);
};
export default PropertyInfo;
