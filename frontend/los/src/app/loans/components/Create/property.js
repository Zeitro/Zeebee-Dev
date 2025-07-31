'use client';
import { Button, Col, DatePicker, Form, InputNumber, Radio, Row, Select, Space, Steps, Tabs } from 'antd';
import React, { useState } from 'react';
import {
	loanStatusOptions,
	loanPendingStatus,
	loanPurposeOptions,
	loanTypeOptions,
} from '@/enum/loan';
import styles from './index.module.scss';
import { usePropertyStore } from '@/store/property';
import { useRouter } from 'next/navigation';
const CreateProperty = ({stepBtn }) => {
	const { propertyTypeOptions } = usePropertyStore();
	const router = useRouter();
	return (
		<Row className={styles.create}>
			<Col md={24}>
				<Form.Item
					label="Property type"
					name="propertyType"
					rules={[{ required: true, message: 'Please select Property type!' }]}
				>
					<Select options={propertyTypeOptions} placeholder="Select Property type" />
				</Form.Item>
			</Col>
			<Col md={24}>
				<Form.Item
					label="Number of unit"
					name="unit"
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
