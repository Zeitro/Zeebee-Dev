import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Select, Input, InputNumber } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import InputPhone from '@/components/InputPhone';
import SelectTINType from '@/components/InputTIN/SelectTINType';
import { emailRegex, phoneRegex, validateRule, tinOrSSNRegex } from '@/lib/format';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@/hooks/useGetOptions';
const name = 'loan_borrower';
const ApplicantInfo = ({ gutter = 24, form, addressProps }) => {
	const borrowerType = Form.useWatch([name, 'type'], form);
	const nameLabel = borrowerType === 'Individual' ? 'Individual name' : 'Entity name';
	const { borrowerTypeOptions } = useGetOptions([
		{
			key: 'applicant_type',
			valueName: 'borrowerTypeOptions',
		},
	]);
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Applicant information
			</Col>
			<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Applicant type
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Borrower type"
					name={[name, 'type']}
					rules={[{ required: true, message: 'Please select borrower type!' }]}
				>
					<Select
						options={borrowerTypeOptions}
						placeholder="Select borrower type"
						onChange={(v) => {
							form.setFieldValue(
								[name, 'own_percent'],
								v === 'Individual' ? 100 : undefined,
							);
							form.setFieldValue([name, 'tin'], { tin_type: undefined, tin: '' });
						}}
					/>
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label={nameLabel}
					name={[name, 'name']}
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
					name={[name, 'tin']}
					rules={[
						{
							required: true,
							message: 'Please enter TIN!',
						},
						{
							validator: (_, value) =>
								validateRule(value?.tin, tinOrSSNRegex, 'Please enter valid TIN!'),
						},
					]}
				>
					<SelectTINType isIndividual={borrowerType === 'Individual'} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="%Own"
					name={[name, 'own_percent']}
					rules={[
						{
							required: true,
							message: 'Please enter %Own!',
						},
					]}
					hidden={borrowerType === 'Individual'}
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
			<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Mailing information
			</Col>
			<AddressForm
				{...addressProps}
				addressName={[name, 'mailing_address', 'address1']}
				addressName2={[name, 'mailing_address', 'address2']}
				cityName={[name, 'mailing_address', 'city']}
				stateName={[name, 'mailing_address', 'state']}
				zipName={[name, 'mailing_address', 'zip_code']}
				addressLabel="Mailing Address 1"
				addressLabel2="Mailing Address 2"
				form={form}
				// addressCol={layoutCol}
				// address2Col={layoutCol}
				stateCol={{ xs: 24, sm: 24, md: 8 }}
				cityCol={{ xs: 24, sm: 24, md: 8 }}
				zipCol={{ xs: 24, sm: 24, md: 8 }}
			/>
			<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Contact information
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Email"
					name={[name, 'contact_email']}
					rules={[
						{
							required: true,
							message: 'Please enter email',
						},
						{
							validator: (_, value) =>
								validateRule(value, emailRegex, 'Invalid email format'),
						},
					]}
				>
					<Input placeholder="Enter email" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Phone number"
					name={[name, 'contact_phone']}
					rules={[
						{
							required: true,
							message: 'Please enter phone number',
						},
						{
							validator: (_, value) =>
								validateRule(value, phoneRegex, 'Invalid phone number'),
						},
					]}
				>
					<InputPhone />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Contact name"
					name={[name, 'contact_name']}
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
