import React from 'react';
import { Checkbox, Col, Form, Radio, Row } from 'antd';
import styles from './index.module.scss';
import { titleCol } from '@/enum/common';
import NestedCheckboxes from '@/components/NestedCheckboxes';
import { raceOptions, ethnicityOptions } from '@/enum/loan';
import { useGetOptions } from '@/hooks/useGetOptions';
const Demographics = ({ gutter = 24, form }) => {
	const { businessStrategyOptions, getSelectFormDom } = useGetOptions([
		// {
		// 	key: 'business_strategy',
		// 	valueName: 'businessStrategyOptions',
		// },
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
					name="race"
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
					name: 'applicant_sex',
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
