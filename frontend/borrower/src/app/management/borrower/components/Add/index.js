import React, { useState } from 'react';
import { NewModal } from '@repo/shared-methods';
import { Col, Form, Input, InputNumber, Row, Select } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { InputPhone, SelectTINType } from '@repo/shared-methods';
import { emailRegex, phoneRegex, validateRule, tinOrSSNRegex } from '@repo/shared-methods/utils';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@/hooks/useGetOptions';
let name = 'borrower';
const AddModal = ({ setOpen }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const { borrowerTypeOptions } = useGetOptions([
		{
			key: 'applicant_type',
			valueName: 'borrowerTypeOptions',
		},
	]);
	const handleSave = () => {};
	return (
		<NewModal
			width={600}
			title="Add borrowing entity"
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				className=""
				size="large"
				disabled={loading}
				scrollToFirstError
			>
				<Row>
					<Col {...layoutCol}>
						<Form.Item
							label="Entity type"
							name={[name, 'type']}
							rules={[{ required: true }]}
						>
							<Select
								options={borrowerTypeOptions}
								placeholder="Select borrower type"
								onChange={(v) => {
									form.setFieldValue(
										[name, 'own_percent'],
										v === 'Individual' ? 100 : undefined,
									);
									form.setFieldValue([name, 'tin'], {
										tin_type: undefined,
										tin: '',
									});
								}}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Entity name"
							name={[name, 'name']}
							rules={[{ required: true }]}
						>
							<Input placeholder="Enter entity name" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="TIN"
							name={[name, 'tin']}
							rules={[
								{
									required: true,
								},
								{
									validator: (_, value) =>
										validateRule(
											value?.tin,
											tinOrSSNRegex,
											'Please enter valid TIN!',
										),
								},
							]}
						>
							<SelectTINType isIndividual={false} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="%Own"
							name={[name, 'own_percent']}
							rules={[{ required: true }]}
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
			</Form>
		</NewModal>
	);
};
export default AddModal;
