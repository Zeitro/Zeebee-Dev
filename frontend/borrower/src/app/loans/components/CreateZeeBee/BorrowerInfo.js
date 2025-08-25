import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Radio, Select, Input, InputNumber, DatePicker, Button } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import InputPhone from '@/components/InputPhone';
import SelectTINType from '@/components/InputTIN/SelectTINType';
import InputTIN from '@/components/InputTIN';
import {
	emailRegex,
	phoneRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
} from '@/lib/format';
import { formatStr } from '@/lib/date';
import { MinusCircleOutlined, PlusCircleOutlined } from '@ant-design/icons';
import { layoutCol, titleCol } from '@/enum/common';
const ApplicantInfo = ({ gutter = 24, form, addressProps }) => {
	const borrowerType = Form.useWatch('borrower_type', form);
	return (
		<div className={styles.createZeeBee}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Borrower information
			</Col>
			<Form.List name="borrowers">
				{(fields, { add, remove }) => {
					return fields?.map((field, index) => (
						<Row
							className={styles.createZeeBee}
							gutter={gutter}
							key={field.key}
							fieldkey={field.fieldKey}
						>
							<Col
								{...titleCol}
								className="!text-xl font-medium text-[#000000a6] pb-6 !flex items-center justify-between"
							>
								Borrower {index + 1}
								{fields.length > 1 && (
									<Button
										type="link"
										size="small"
										icon={<MinusCircleOutlined />}
										onClick={() => remove(field.name)}
									>
										Delete borrower
									</Button>
								)}
							</Col>
							<Col {...layoutCol} md={8}>
								<Form.Item
									label="First name"
									name={[field.name, 'first_name']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<Input placeholder="Enter first name" />
								</Form.Item>
							</Col>
							<Col {...layoutCol} md={8}>
								<Form.Item
									label="Middle name"
									name={[field.name, 'middle_name']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<Input placeholder="Enter middle name" />
								</Form.Item>
							</Col>
							<Col {...layoutCol} md={8}>
								<Form.Item
									label="Last name"
									name={[field.name, 'last_name']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<Input placeholder="Enter last name" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="SSN"
									name={[field.name, 'ssn']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<InputTIN type="SSN" placeholder="Enter SSN" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Date of birth"
									name="date_of_birth"
									rules={[
										{ required: true, message: 'Please select date of birth!' },
									]}
								>
									<DatePicker
										format={formatStr}
										placeholder="select"
										style={{ width: '100%' }}
									/>
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Phone"
									name="phone"
									rules={[
										{
											required: true,
											message: 'Please enter phone number',
										},
										{
											validator: (_, value) => {
												if (value && !phoneRegex.test(value)) {
													return Promise.reject('Invalid phone number');
												}
												return Promise.resolve();
											},
										},
									]}
								>
									<InputPhone />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Marital status"
									name="marital_status"
									rules={[
										{
											required: true,
											message: 'Please select marital status!',
										},
									]}
								>
									<Select options={[]} placeholder="Select marital status" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Email"
									name="email"
									rules={[
										{
											required: true,
											message: 'Please enter email',
										},
										{
											validator: (_, value) => {
												if (value && !emailRegex.test(value)) {
													return Promise.reject('Invalid email format');
												}
												return Promise.resolve();
											},
										},
									]}
								>
									<Input placeholder="Enter email" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Est. Credit Score"
									name="credit_score"
									rules={[
										{
											required: true,
											message: 'Please enter est. credit score!',
										},
										{
											validator: (_, value) => {
												if (value && !creditScoreRegex.test(value)) {
													return Promise.reject('Invalid Credit Score');
												}
												return Promise.resolve();
											},
										},
									]}
								>
									<InputNumber
										style={{ width: '100%' }}
										min={minCreditScore}
										max={maxCreditScore}
										controls={false}
										placeholder="Enter est. credit score"
									/>
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Driver’s License/ID number"
									name={[field.name, 'id']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<Input placeholder="Enter driver’s license/ID number" />
								</Form.Item>
							</Col>
							<AddressForm
								{...addressProps}
								addressName={[
									'loan_property',
									'property_base',
									'address',
									'address1',
								]}
								addressName2={[
									'loan_property',
									'property_base',
									'address',
									'address2',
								]}
								cityName={['loan_property', 'property_base', 'address', 'city']}
								stateName={['loan_property', 'property_base', 'address', 'state']}
								zipName={['loan_property', 'property_base', 'address', 'zip_code']}
								addressLabel="Primary Address 1"
								addressLabel2="Primary Address 2"
								form={form}
								// addressCol={{ ...layoutCol, md: 12 }}
								// address2Col={{ ...layoutCol, md: 12 }}
								stateCol={{ ...layoutCol, md: 8 }}
								cityCol={{ ...layoutCol, md: 8 }}
								zipCol={{ ...layoutCol, md: 8 }}
							/>
							{fields.length === index + 1 && (
								<Col {...layoutCol}>
									<Button
										type="link"
										size="small"
										icon={<PlusCircleOutlined />}
										onClick={() => add()}
									>
										Add borrower
									</Button>
								</Col>
							)}
						</Row>
					));
				}}
			</Form.List>
		</div>
	);
};
export default ApplicantInfo;
