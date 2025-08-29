import React from 'react';
import { Col, Form, Radio, Row } from 'antd';
import styles from './index.module.scss';
import { titleCol } from '@/enum/common';
import { yesNoOptions } from '@/enum/common';
const Declaration = ({ gutter, form }) => {
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				About this property and your money for this loan
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Will you occupy the property as your primary residence?"
					name={['loan_demographic', 'property', 'as_primary_resident']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="If this is a Purchase Transaction. Do you have a family relationship or business affiliation with the seller of the property?"
					name={['loan_demographic', 'property', 'has_transaction_relationship']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Are you borrowing any money for this real estate transaction (e.g. money for your closing costs or down payment) or obtaining any money from another party such as the seller or realtor, that you have not disclosed on this loan application? lf Yes, what is the amount of this money?"
					name={['loan_demographic', 'property', 'has_borrowed_money']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Have you or you be applying for a mortgage loan on another property (not the property securing this loan) on or before closing this transaction that is not disclosed on this loan application?"
					name={['loan_demographic', 'property', 'has_another_mortgage_loan']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Will this property be subject to a lien that could take priority over the first mortgage lien. such as clean energy lien paid through your property taxes (e.g., the Property Assessed Clean Energy Program)?"
					name={['loan_demographic', 'property', 'has_priority_over_1st']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Are there any outstanding judgments against you?"
					name={['loan_demographic', 'property', 'has_outstanding_judgments']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Are you a party to a lawsuit in which you potentially have any personal financial liability?"
					name={['loan_demographic', 'property', 'has_a_lawsuit']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Have you conveyed title to any property in lieu of foreclosure in the past 7 years?"
					name={['loan_demographic', 'property', 'has_conveyed_title']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Within the past 7 years, have you completed a pre-foreclosure sale or short sale, whereby the property was sold to a third party and the Lender agreed to accept less than the outstanding mortgage balance due?"
					name={['loan_demographic', 'property', 'has_a_pre_foreclosure_sale']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Have you had property foreclosed upon in the last 7 years?"
					name={['loan_demographic', 'property', 'has_foreclosed_upon']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Have you declared bankruptcy within the past 7 years?If YES, identify the type(s) of bankruptcy:"
					name={['loan_demographic', 'property', 'has_a_bankruptcy']}
					rules={[
						{
							required: true,
							message: 'Please select!',
						},
					]}
				>
					<Radio.Group options={yesNoOptions} className={styles.yesNo2} />
				</Form.Item>
			</Col>
		</Row>
	);
};

export default Declaration;
