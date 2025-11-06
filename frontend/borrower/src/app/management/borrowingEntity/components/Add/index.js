import React, { useState } from 'react';
import { NewModal } from '@repo/shared-methods';
import { Col, Form, Input, InputNumber, Row, Select, Spin } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { InputPhone, SelectFormTINType } from '@repo/shared-methods';
import {
	emailRegex,
	phoneRegex,
	validateRule,
	createRequest,
	capitalizeFirstLetter,
} from '@repo/shared-methods/utils';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@/hooks/useGetOptions';
import { useCallback } from 'react';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const AddModal = ({ setOpen, rowData, handleSearch }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const { borrowerTypeOptions } = useGetOptions([
		{
			key: 'applicant_type',
			valueName: 'borrowerTypeOptions',
		},
	]);

	const handleDetails = useCallback(async () => {
		try {
			setLoading(true);
			const res = await request(`/borrower_entity/${rowData?.id}`, {
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
			let url = '/borrower_entity';
			let method = 'POST';
			if (rowData?.id) {
				url = `borrower_entity/${rowData?.id}`;
				method = 'PUT';
			}
			await request(url, {
				method,
				body: { ...value },
			});
			handleSearch?.();
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
			width={700}
			title={`${rowData?.id ? 'Edit' : 'Add'} borrowing entity`}
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
						<Col {...layoutCol}>
							<Form.Item
								label="Entity type"
								name={['type']}
								rules={[{ required: true }]}
							>
								<Select
									options={borrowerTypeOptions}
									placeholder="Select borrower type"
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Entity name"
								name={['name']}
								rules={[{ required: true }]}
							>
								<Input placeholder="Enter entity name" />
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<SelectFormTINType isIndividual={false} />
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="%Own"
								name={['own_percent']}
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
							addressName={['mailing_address', 'address1']}
							addressName2={['mailing_address', 'address2']}
							cityName={['mailing_address', 'city']}
							stateName={['mailing_address', 'state']}
							zipName={['mailing_address', 'zip_code']}
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
								name={['contact_email']}
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
								name={['contact_phone']}
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
								name={['contact_name']}
								rules={[
									{
										required: true,
										message: 'Please enter contact name',
									},
								]}
								getValueFromEvent={(e) => {
									return capitalizeFirstLetter(e.target.value);
								}}
							>
								<Input placeholder="Enter contact name" />
							</Form.Item>
						</Col>
					</Row>
				</Form>
			</Spin>
		</NewModal>
	);
};
export default AddModal;
