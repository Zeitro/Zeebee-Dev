'use client';
import React, { useState } from 'react';
import NewModal from '@/components/NewModal';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { Col, Form, Input, Radio, Row } from 'antd';
import { useGetOptions } from '@/hooks/useGetOptions';
import InputNumAmount from '@/components/InputNumAmount';
import CommonFormDate from '@/components/CommonDate';
import { yesNoOptions } from '@/enum/common';
import { emailRegex } from '@/lib/format';
const AddThirdService = ({ title, rowData, onOk, onCancel }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const layoutCol = {
		md: 24,
	};
	const { businessStrategyOptions, getSelectFormDom } = useGetOptions([
		// {
		// 	key: 'business_strategy',
		// 	valueName: 'businessStrategyOptions',
		// },
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
					await onOk(values, rowData?.id ? 'Edit' : 'Add');
					setConfirmLoading(false);
					onCancel();
				} catch (error) {
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
					{/* {getSelectFormDom({
						label: 'Service type',
						name: 'property_type',
						placeholder: 'Select service type',
						valueName: 'propertyTypeOptions',
						requiredMessage: 'Please select service type!',
					})} */}
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
						<Form.Item
							label="Company"
							name="company"
							rules={[
								{
									required: true,
									message: 'Please enter service company',
								},
							]}
						>
							<Input placeholder="Enter service company" />
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
				</Row>
			</Form>
		</NewModal>
	);
};
export default AddThirdService;
