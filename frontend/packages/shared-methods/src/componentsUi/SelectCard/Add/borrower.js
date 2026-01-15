import React, { useState, useCallback } from 'react';
import { Row, Col, Form, Select, Input, InputNumber, Spin } from 'antd';
// import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import {
	emailRegex,
	phoneRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
	validateRule,
	tinOrSSNRegex,
	createRequest,
} from '@repo/shared-methods/utils';
import { layoutCol } from '@repo/shared-methods/enum';
import { useGetOptions } from '@repo/shared-methods/hooks';
import {
	CommonFormDate,
	InputNameCombination,
	InputPhone,
	InputTIN,
	NewModal,
	AddressForm,
} from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const AddModal = ({ setOpen, rowData, handleSearch }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const { maritalStatusOptions } = useGetOptions([
		{
			key: 'marital_status',
			valueName: 'maritalStatusOptions',
		},
	]);

	const handleDetails = useCallback(async () => {
		try {
			setLoading(true);
			const res = await request(`/borrower_applicant/${rowData?.id}`, {
				method: 'GET',
			});
			form.setFieldsValue(res);
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [form, rowData]);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/borrower_applicant';
			let method = 'POST';
			if (rowData?.id) {
				url = `/borrower_applicant/${rowData?.id}`;
				method = 'PUT';
			}
			const res = await request(url, {
				method,
				body: { ...value },
			});
			handleSearch?.(res?.id);
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	useDeferredSyncEffect({
		predicate: () => !!rowData?.id,
		deferred: () => {
			handleDetails();
		},
		deps: [handleDetails, rowData?.id],
		isUseMicrotask: true,
	});
	return (
		<NewModal
			width={600}
			title={`${rowData?.id ? 'Edit' : 'Add'} borrower`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Spin spinning={loading}>
				<Form
					form={form}
					layout="vertical"
					autoComplete="off"
					className="h-[60vh] overflow-hidden overflow-y-auto"
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
				>
					<Row gutter={16}>
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
						/>
					</Row>
				</Form>
			</Spin>
		</NewModal>
	);
};
export default AddModal;
