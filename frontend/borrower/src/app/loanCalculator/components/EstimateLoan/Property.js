import React from 'react';
// import styles from './index.module.scss';
import { Row, Col, Form, Radio, Select, Input, InputNumber } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { useGetOptions } from '@/hooks/useGetOptions';
import { layoutCol, titleCol, yesNoOptions } from '@/enum/common';
let name = 'loan_property';
const PropertyInfo = ({ gutter = 24, form, addressProps, isRequestProperty, loan_type }) => {
	const { propertyTypeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);

	return (
		<Row gutter={gutter}>
			<Col {...titleCol} className="!text-3xl/10 font-medium pb-8">
				Property information
			</Col>
			<AddressForm
				{...addressProps}
				addressName={[name, 'property_base', 'address', 'address1']}
				addressName2={[name, 'property_base', 'address', 'address2']}
				cityName={[name, 'property_base', 'address', 'city']}
				stateName={[name, 'property_base', 'address', 'state']}
				zipName={[name, 'property_base', 'address', 'zip_code']}
				form={form}
				stateCol={{ ...titleCol }}
				cityCol={{ ...layoutCol, lg: 12, md: 24 }}
				zipCol={{ ...layoutCol, lg: 12, md: 24 }}
				disabled={isRequestProperty}
			/>

			<Col {...titleCol}>
				<Form.Item
					label="Property type"
					name={[name, 'property_type']}
					rules={[{ required: true }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select property type" />
				</Form.Item>
			</Col>
			{loan_type === 'Ground up construction' && (
				<Col {...titleCol}>
					<Form.Item
						label="Permit ready"
						name="has_permit_ready"
						rules={[{ required: true }]}
					>
						<Radio.Group options={yesNoOptions} style={{ width: '100%' }} />
					</Form.Item>
				</Col>
			)}
		</Row>
	);
};
export default PropertyInfo;
