'use client';
import { titleCol } from '@/enum/common';
import { InfoCircleOutlined, MailTwoTone, PlusOutlined } from '@ant-design/icons';
import {
	FieldTips,
	InputNumAmount,
	InputPhone,
	UploadFileButton,
	OrgProfile,
} from '@repo/shared-methods';
import {
	capitalizeFirstLetter,
	emailRegex,
	phoneRegex,
	validateRule,
} from '@repo/shared-methods/utils';
import { Button, Col, Form, Input, InputNumber, Row } from 'antd';
import { useState } from 'react';
const OrganizationProfile = () => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);

	return (
		<div className=" py-6 px-8 flex-1 overflow-hidden w-full overflow-y-auto">
			<OrgProfile />
			{/* <Form
				form={form}
				name="org-profile-form"
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
				// initialValues={{ ...rowData, role_ids: rowData?.roles.map((i) => i.id) }}
			>
				<Row gutter={16}>
					<Col {...titleCol} className="!text-xl font-medium text-base-gray pb-4">
						Organization profile
					</Col>
					<Col {...titleCol}>
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
					<Col {...titleCol}>
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
					</Col>
					<Col {...titleCol}>
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
					</Col>
					<Col {...titleCol}>
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
					</Col>
					<Col {...titleCol}>
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
					</Col>
					<Col {...titleCol}>
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
					</Col>
					<Col {...titleCol} className="!flex">
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
					</Col>
					<Col {...titleCol} className="!flex">
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
					</Col>
					<Col {...titleCol} className="!flex">
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
					</Col>
					<Col {...titleCol} className="!flex">
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
					</Col>

					<Col className="!flex justify-end w-full">
						<Button type="primary" loading={confirmLoading}>
							Save
						</Button>
					</Col>
				</Row>
			</Form> */}
		</div>
	);
};
export default OrganizationProfile;
