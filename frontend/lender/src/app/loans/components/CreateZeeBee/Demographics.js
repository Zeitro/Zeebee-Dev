import React from 'react';
import {Col, Form, Row} from 'antd';
import styles from './index.module.scss';
import {titleCol} from '@/enum/common';
import {NestedCheckboxes} from '@repo/shared-methods';
import {ethnicityOptions, raceOptions} from '@/enum/loan';
import {useGetOptions} from '@repo/shared-methods/hooks';

const Demographics = ({ gutter = 24, form }) => {
	const { getSelectFormDom } = useGetOptions([
		{
			key: 'applicant_sex',
			valueName: 'sexOptions',
		},
	]);
	const layoutCol = {
		...titleCol,
		xl: 12,
	};
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Borrower demographic information
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Ethnicity (Check one or more)"
					name={['loan_demographic', 'ethnicity']}
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
			<Col {...layoutCol}></Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Race (Check one or more)"
					name={['loan_demographic', 'race']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<NestedCheckboxes options={raceOptions} form={form} name="race" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}></Col>
			{getSelectFormDom(
				{
					label: 'Sex',
					name: ['loan_demographic', 'sex'],
					placeholder: 'Select Sex',
					valueName: 'sexOptions',
					requiredMessage: 'Please select sex!',
				},
				layoutCol,
			)}
		</Row>
	);
};

export default Demographics;
