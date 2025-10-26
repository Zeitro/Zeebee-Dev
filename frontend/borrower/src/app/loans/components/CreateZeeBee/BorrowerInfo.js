import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Select, Input, InputNumber, Button } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import InputPhone from '@/components/InputPhone';
import InputTIN from '@/components/InputTIN';
import {
	emailRegex,
	phoneRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
	validateRule,
	tinOrSSNRegex,
} from '@/lib/format';
import { MinusCircleOutlined, PlusCircleOutlined } from '@ant-design/icons';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@/hooks/useGetOptions';
import { CommonFormDate } from '@repo/shared-methods';
import InputName from '@/components/InputName';
const ApplicantInfo = ({ gutter = 24, form, addressProps }) => {
	const { maritalStatusOptions } = useGetOptions([
		{
			key: 'marital_status',
			valueName: 'maritalStatusOptions',
		},
	]);
	return (
		<div className={styles.createZeeBee}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Borrower information
			</Col>
			<Form.List name="applicants">
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
							<InputName
								layoutCol={{ ...layoutCol, md: 8 }}
								firstName={[field.name, 'first_name']}
								middleName={[field.name, 'mid_name']}
								lastName={[field.name, 'last_name']}
								config={{
									firstNameRequired: !index,
									lastNameRequired: !index,
								}}
								listName="applicants"
								form={form}
								className="!ml-0 !mr-0"
							/>
							{/* <Col {...layoutCol} md={8}>
								<Form.Item
									label="First name"
									name={[field.name, 'first_name']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<Input placeholder="Enter first name" />
								</Form.Item>
							</Col>
							<Col {...layoutCol} md={8}>
								<Form.Item label="Middle name" name={[field.name, 'mid_name']}>
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
							</Col> */}
							<Col {...layoutCol}>
								<Form.Item
									label="SSN"
									name={[field.name, 'ssn']}
									rules={[
										{ required: true, message: 'Please Enter!' },
										{
											pattern: tinOrSSNRegex,
											message: 'Please enter valid SSN!',
										},
									]}
								>
									<InputTIN type="SSN" placeholder="Enter SSN" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Date of birth"
									name={[field.name, 'dob']}
									rules={[
										{ required: true, message: 'Please select date of birth!' },
									]}
								>
									<CommonFormDate placeholder="select" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Phone"
									name={[field.name, 'phone']}
									rules={[
										{
											required: true,
											message: 'Please enter phone number',
										},
										{
											validator: (_, value) =>
												validateRule(
													value,
													phoneRegex,
													'Invalid phone number',
												),
										},
									]}
								>
									<InputPhone />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Marital status"
									name={[field.name, 'marital_status']}
									rules={[
										{
											required: true,
											message: 'Please select marital status!',
										},
									]}
								>
									<Select
										options={maritalStatusOptions}
										placeholder="Select marital status"
									/>
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Email"
									name={[field.name, 'email']}
									rules={[
										{
											required: true,
											message: 'Please enter email',
										},
										{
											validator: (_, value) =>
												validateRule(
													value,
													emailRegex,
													'Invalid email format',
												),
										},
									]}
								>
									<Input placeholder="Enter email" />
								</Form.Item>
							</Col>
							<Col {...layoutCol}>
								<Form.Item
									label="Estimated Credit Score"
									name={[field.name, 'est_credit_score']}
									rules={[
										{
											required: true,
											message: 'Please enter est. credit score!',
										},
										{
											validator: (_, value) =>
												validateRule(
													value,
													creditScoreRegex,
													'Invalid Credit Score',
												),
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
									name={[field.name, 'driver_license_num']}
									rules={[{ required: true, message: 'Please Enter!' }]}
								>
									<Input placeholder="Enter driver’s license/ID number" />
								</Form.Item>
							</Col>
							<AddressForm
								{...addressProps}
								addressName={[field.name, 'address', 'address1']}
								addressName2={[field.name, 'address', 'address2']}
								cityName={[field.name, 'address', 'city']}
								stateName={[field.name, 'address', 'state']}
								zipName={[field.name, 'address', 'zip_code']}
								addressLabel="Primary Address 1"
								addressLabel2="Primary Address 2"
								form={form}
								// addressCol={{ ...layoutCol, md: 12 }}
								// address2Col={{ ...layoutCol, md: 12 }}
								stateCol={{ ...layoutCol, md: 8 }}
								cityCol={{ ...layoutCol, md: 8 }}
								zipCol={{ ...layoutCol, md: 8 }}
								listName="applicants"
							/>
							{fields.length === index + 1 && (
								<Col {...layoutCol}>
									<Button
										type="link"
										size="small"
										icon={<PlusCircleOutlined />}
										onClick={() =>
											add({ type: 'Individual', is_primary: false })
										}
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
