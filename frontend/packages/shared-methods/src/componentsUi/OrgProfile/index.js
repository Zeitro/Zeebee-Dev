'use client';
import { PlusOutlined } from '@ant-design/icons';
import { InputNumAmount, InputPhone, UploadFileButton } from '@repo/shared-methods';
import {
	capitalizeFirstLetter,
	emailRegex,
	phoneRegex,
	validateRule,
	createRequest,
} from '@repo/shared-methods/utils';
import { Button, Col, Form, Input, InputNumber, message, Row } from 'antd';
import { useState, forwardRef, useImperativeHandle } from 'react';
const OrganizationProfile = ({ api, rowData, successCallback, className }, ref) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const to = <div className="flex items-center justify-center px-2 flex-shrink-0">to</div>;
	const uploadButton = (
		<button style={{ border: 0, background: 'none' }} type="button">
			<PlusOutlined />
			<div style={{ marginTop: 8 }}>Upload</div>
		</button>
	);

	const handleSaveOrgProfile = async () => {
		try {
			if (!api) return;
			const values = await form.validateFields();
			setConfirmLoading(true);
			if (rowData?.id) {
				values.id = rowData.id;
			}
			await request(api, {
				method: 'POST',
				body: {
					...values,
				},
			});
			setConfirmLoading(false);
			message.destroy();
			message.open({
				content: `Organization Profile ${rowData?.id ? 'edited' : 'added'} successfully.`,
				type: 'success',
			});
			successCallback?.();
		} catch (error) {
			setConfirmLoading(false);
		}
	};

	useImperativeHandle(ref, () => ({
		handleSaveOrgProfile,
	}));
	return (
		<Form
			form={form}
			name="org-profile-form"
			layout="vertical"
			autoComplete="off"
			size="large"
			disabled={confirmLoading}
			scrollToFirstError
			initialValues={{ ...rowData }}
			className={className}
		>
			<Row gutter={16}>
				<Col md={6} sm={24}>
					<Form.Item
						label={null}
						name={['avatar']}
						rules={[
							{
								required: true,
							},
						]}
					>
						<UploadFileButton
							accept=".png,.jpg,.jpeg"
							showUploadList={false}
							listType="picture-card"
						>
							{uploadButton}
						</UploadFileButton>
					</Form.Item>
				</Col>
				<Col md={18} sm={24}>
					<Form.Item
						label="Organization name"
						name={['org_name']}
						rules={[
							{
								required: true,
							},
						]}
					>
						<Input disabled />
					</Form.Item>
					<Form.Item
						label="Contact name"
						name={['name']}
						rules={[
							{
								required: true,
							},
						]}
						getValueFromEvent={(e) => {
							return capitalizeFirstLetter(e.target.value);
						}}
					>
						<Input placeholder="Enter contact name" />
					</Form.Item>
					<Form.Item
						label="Contact email"
						name={['email']}
						rules={[
							{
								required: true,
							},
							{
								validator: (_, value) =>
									validateRule(value, emailRegex, 'Invalid email format'),
							},
						]}
					>
						<Input placeholder="Enter email" />
					</Form.Item>
					<Form.Item
						label="Contact phone"
						name={['phone']}
						rules={[
							{
								validator: (_, value) =>
									validateRule(value, phoneRegex, 'Invalid phone number'),
							},
						]}
					>
						<InputPhone />
					</Form.Item>
					<Form.Item
						label="Website link"
						name={['link']}
						rules={[
							{
								required: true,
							},
						]}
					>
						<Input placeholder="Enter website link" />
					</Form.Item>
					<div className="!flex">
						<Form.Item label="Loan amount" name={['amount']} className="flex-1">
							<InputNumAmount
								placeholder="Enter min loan amount"
								className="!w-full"
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['amount_max']} className="flex-1">
							<InputNumAmount
								placeholder="Enter min loan amount"
								className="!w-full"
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item label="Loan term" name={['term']} className="flex-1">
							<InputNumber
								placeholder="Enter min loan term"
								className="!w-full"
								controls={false}
								suffix="Months"
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['term_max']} className="flex-1">
							<InputNumber
								placeholder="Enter max loan term"
								className="!w-full"
								controls={false}
								suffix="Months"
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item label="Loan rates" name={['rate']} className="flex-1">
							<InputNumber
								placeholder="Enter min loan amount"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['rate_max']} className="flex-1">
							<InputNumber
								placeholder="Enter max loan amount"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item label="Loan points" name={['point']} className="flex-1">
							<InputNumber
								placeholder="Enter min loan amount"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['point_max']} className="flex-1">
							<InputNumber
								placeholder="Enter max loan amount"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
					</div>
				</Col>
			</Row>
		</Form>
	);
};
OrganizationProfile.displayName = 'OrgProfile';
export const OrgProfile = forwardRef(OrganizationProfile);
