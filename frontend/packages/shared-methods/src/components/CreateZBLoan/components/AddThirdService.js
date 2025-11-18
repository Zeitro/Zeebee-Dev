'use client';
import React, { useState } from 'react';
import { Col, Form, Input, Row } from 'antd';
import { useGetOptions } from '@/hooks/useGetOptions';
import { emailRegex, phoneRegex, validateRule } from '../../../lib/format';
import { InputPhone } from '../../InputPhone';
import { NewModal } from '../../NewModal';
const AddThirdService = ({ title, rowData, onOk, onCancel }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const layoutCol = {
		md: 24,
	};
	const { getSelectFormDom } = useGetOptions([
		{
			key: 'loan_third_service_type',
			valueName: 'serviceTypeOptions',
		},
	]);
	return (
		<NewModal
			width={572}
			title={title}
			onCancel={onCancel}
			okText={`${rowData?.id ? 'Edit' : 'Add'} Service`}
			confirmLoading={confirmLoading}
			onOk={async () => {
				try {
					setConfirmLoading(true);
					const values = await form.validateFields();
					if (rowData?.id) values.id = rowData?.id;
					await onOk(values, rowData?.id ? 'Edit' : 'Add');
					setConfirmLoading(false);
					onCancel();
				} catch {
					setConfirmLoading(false);
				}
			}}
		>
			<Form
				form={form}
				name="addThirdServiceForm"
				layout="vertical"
				autoComplete="off"
				className="w-full"
				size="large"
				disabled={confirmLoading}
				initialValues={{
					...rowData,
				}}
			>
				<Row gutter={24}>
					{getSelectFormDom({
						label: 'Service type',
						name: 'type',
						placeholder: 'Select service type',
						valueName: 'serviceTypeOptions',
						requiredMessage: 'Please select service type!',
					})}
					<Col {...layoutCol}>
						<Form.Item
							label="Name"
							name="name"
							rules={[
								{
									required: true,
									message: 'Please enter service name',
								},
							]}
						>
							<Input placeholder="Enter service name" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item label="Company" name="company">
							<Input placeholder="Enter service company" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Phone"
							name="phone"
							rules={[
								{
									required: true,
									message: 'Please enter phone',
								},
								{
									validator: (_, value) => validateRule(value, phoneRegex),
								},
							]}
						>
							<InputPhone placeholder="Enter phone" />
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
									validator: (_, value) => validateRule(value, emailRegex),
								},
							]}
						>
							<Input placeholder="Enter email" />
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};
export default AddThirdService;
