import React from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Select, Input } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { InputPhone, SelectFormTINType } from '@repo/shared-methods';
import { emailRegex, phoneRegex, validateRule } from '@repo/shared-methods/utils';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@repo/shared-methods/hooks';
import SelectCard from '@/components/SelectCard';
const name = 'loan_borrower';
const ApplicantInfo = ({
	gutter = 24,
	addressProps,
	disabled,
	isHideTitle,
	isShowSelectCard = true,
}) => {
	const form = Form.useFormInstance();
	const borrowerType = Form.useWatch([name, 'type'], form);
	const borrower_entity_id = Form.useWatch([name, 'borrower_entity_id'], form);
	const nameLabel = borrowerType === 'Individual' ? 'Individual name' : 'Entity name';
	const { borrowerTypeOptions } = useGetOptions([
		{
			key: 'applicant_type',
			valueName: 'borrowerTypeOptions',
		},
	]);
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			{!isHideTitle && (
				<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
					Applicant information
				</Col>
			)}
			{isShowSelectCard && (
				<SelectCard
					type="borrowingEntity"
					addTitle="Add borrowing entity"
					title="Choose a borrowing entity"
					api="/borrower_entity/all"
					className="mb-8"
					changeFormValue={(values) => {
						const borrower_entity = {
							[name]: {
								...values,
								borrower_entity_id: values.id,
							},
						};
						form.setFieldsValue(borrower_entity);
					}}
					name={[name, 'borrower_entity_id']}
					id={borrower_entity_id}
					disabled={disabled}
				/>
			)}
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
						disabled
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
					<Input placeholder={`Enter ${nameLabel}`} disabled />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				{/* <Form.Item
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
				> */}
				{/* </Form.Item> */}
				<SelectFormTINType
					isIndividual={borrowerType === 'Individual'}
					disabled
					typeName={[name, 'tin_type']}
					tinName={[name, 'tin']}
				/>
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
				disabled={true}
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
					<Input placeholder="Enter email" disabled />
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
					<InputPhone disabled={true} />
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
					<Input placeholder="Enter contact name" disabled={true} />
				</Form.Item>
			</Col>
		</Row>
	);
};
export default ApplicantInfo;
