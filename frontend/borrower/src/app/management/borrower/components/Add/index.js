import React, { useState } from 'react';
import { Row, Col, Form, Select, Input, InputNumber, Button } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import {
	emailRegex,
	phoneRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
	validateRule,
	tinOrSSNRegex,
} from '@repo/shared-methods/utils';
import { MinusCircleOutlined, PlusCircleOutlined } from '@ant-design/icons';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@/hooks/useGetOptions';
import {
	CommonFormDate,
	InputNameCombination,
	InputPhone,
	InputTIN,
	NewModal,
} from '@repo/shared-methods';
let name = 'borrower';
const AddModal = ({ setOpen }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const { maritalStatusOptions } = useGetOptions([
		{
			key: 'marital_status',
			valueName: 'maritalStatusOptions',
		},
	]);
	const handleSave = () => {};
	return (
		<NewModal
			width={600}
			title="Add borrower"
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				className="h-[60vh] overflow-hidden overflow-y-auto"
				size="large"
				disabled={loading}
				scrollToFirstError
			>
				<Row
					// className={styles.createZeeBee}
					gutter={16}
				>
					<InputNameCombination
						layoutCol={{ ...layoutCol, md: 8 }}
						firstName={['first_name']}
						middleName={['mid_name']}
						lastName={['last_name']}
						config={
							{
								// firstNameRequired: !index,
								// lastNameRequired: !index,
							}
						}
						listName="applicants"
						form={form}
						className="!ml-0 !mr-0"
					/>
					<Col {...layoutCol}>
						<Form.Item
							label="SSN"
							name={['ssn']}
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
							name={['dob']}
							rules={[{ required: true, message: 'Please select date of birth!' }]}
						>
							<CommonFormDate placeholder="select" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Phone"
							name={['phone']}
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
							label="Marital status"
							name={['marital_status']}
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
							name={['email']}
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
							label="Estimated Credit Score"
							name={['est_credit_score']}
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
							name={['driver_license_num']}
							rules={[{ required: true, message: 'Please Enter!' }]}
						>
							<Input placeholder="Enter driver’s license/ID number" />
						</Form.Item>
					</Col>
					<AddressForm
						addressName={['address', 'address1']}
						addressName2={['address', 'address2']}
						cityName={['address', 'city']}
						stateName={['address', 'state']}
						zipName={['address', 'zip_code']}
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
				</Row>
			</Form>
		</NewModal>
	);
};
export default AddModal;
