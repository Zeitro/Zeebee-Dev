'use client';
import { Col, Form, InputNumber, Row, Select } from 'antd';
import React, { useState } from 'react';
import styles from './index.module.scss';
import { useGetOptions } from '@/hooks/useGetOptions';
const CreateProperty = ({ stepBtn, gutter, ...colProps }) => {
	const { propertyTypeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	return (
		<Row className={styles.create} gutter={gutter}>
			<Col md={24} {...(colProps || {})}>
				<Form.Item
					label="Property type"
					name={['property', 'property_type']}
					rules={[{ required: true, message: 'Please select Property type!' }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select Property type" />
				</Form.Item>
			</Col>
			<Col md={24} {...(colProps || {})}>
				<Form.Item
					label="Number of unit"
					name={['property', 'unit']}
					rules={[{ required: true, message: 'Please input!' }]}
				>
					<InputNumber
						style={{ width: '100%' }}
						min={1}
						max={10}
						controls={false}
						placeholder="Enter initial loan amount"
					/>
				</Form.Item>
			</Col>

			{stepBtn}
		</Row>
	);
};
export default CreateProperty;
