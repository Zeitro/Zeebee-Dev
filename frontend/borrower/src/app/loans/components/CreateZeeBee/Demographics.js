import React from 'react';
import { Checkbox, Col, Form, Radio, Row } from 'antd';
import styles from './index.module.scss';
import { layoutCol, titleCol } from '@/enum/common';
import { yesNoOptions } from '@/enum/common';
import NestedCheckboxes from '@/components/NestedCheckboxes';
import { raceOptions, ethnicityOptions } from '@/enum/loan';
const Demographics = ({ gutter = 24, form }) => {
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Borrower demographic information
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Ethnicity (Check one or more)"
					name="ethnicity"
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<NestedCheckboxes options={ethnicityOptions} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Race (Check one or more)"
					name="race"
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<NestedCheckboxes options={raceOptions} />
				</Form.Item>
			</Col>
		</Row>
	);
};

export default Demographics;
