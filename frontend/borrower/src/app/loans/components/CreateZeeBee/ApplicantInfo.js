import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Radio, Select, Input, InputNumber } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import InputPhone from '@/components/InputPhone';
import SelectTINType from '@/components/InputTIN/SelectTINType';
import { emailRegex, phoneRegex } from '@/lib/format';
const ApplicantInfo = ({ gutter, form, addressProps }) => {
	const borrowerType = Form.useWatch('borrower_type', form);
	const layoutCol = {
		xs: 24,
		sm: 24,
	};
	const nameLabel = borrowerType === 'Individual' ? 'Individual name' : 'Entity name';
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...layoutCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Applicant information
			</Col>
			<Col {...layoutCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Applicant type
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Borrower type"
					name="borrower_type"
					rules={[{ required: true, message: 'Please select borrower type!' }]}
				>
					<Select
						options={[
							{ label: 'Individual', value: 'Individual' },
							{ label: 'LLC', value: 'LLC' },
						]}
						placeholder="Select borrower type"
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label={nameLabel}
					name="name"
					rules={[
						{
							required: true,
							message: `Please enter ${nameLabel}!`,
						},
					]}
				>
					<Input placeholder={`Enter ${nameLabel}`} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="TIN"
					name="tin"
					rules={[
						{
							required: true,
							message: 'Please enter TIN!',
						},
					]}
				>
					<SelectTINType isIndividual={borrowerType === 'Individual'} />
				</Form.Item>
			</Col>
			{borrowerType !== 'Individual' && (
				<Col {...layoutCol}>
					<Form.Item
						label="%Own"
						name="own_percent"
						rules={[
							{
								required: true,
								message: 'Please enter %Own!',
							},
						]}
					>
						<InputNumber
							style={{ width: '100%' }}
							min={1}
							max={100}
							controls={false}
							placeholder="Enter %Own"
							suffix="%"
						/>
					</Form.Item>
				</Col>
			)}
			<Col {...layoutCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Mailing information
			</Col>
			<AddressForm
				{...addressProps}
				addressName={['loan_property', 'property_base', 'address', 'address1']}
				addressName2={['loan_property', 'property_base', 'address', 'address2']}
				cityName={['loan_property', 'property_base', 'address', 'city']}
				stateName={['loan_property', 'property_base', 'address', 'state']}
				zipName={['loan_property', 'property_base', 'address', 'zip_code']}
				addressLabel="Mailing Address 1"
				addressLabel2="Mailing Address 2"
				form={form}
			/>
			<Col {...layoutCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Contact information
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
					label="Phone number"
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
					label="Contact name"
					name="contact_name"
					rules={[
						{
							required: true,
							message: 'Please enter contact name',
						},
					]}
				>
					<Input placeholder="Enter contact name" />
				</Form.Item>
			</Col>
		</Row>
	);
};
export default ApplicantInfo;
