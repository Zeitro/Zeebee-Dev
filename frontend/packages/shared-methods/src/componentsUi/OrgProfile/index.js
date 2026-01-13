'use client';
import { PlusOutlined } from '@ant-design/icons';
import {
	ConfirmPassWord,
	InputNameCombination,
	InputNumAmount,
	InputPhone,
	UploadFileButton,
} from '@repo/shared-methods';
import {
	capitalizeFirstLetter,
	emailRegex,
	passwordRegex,
	phoneRegex,
	validateRule,
	createRequest,
} from '@repo/shared-methods/utils';
import { Button, Col, Form, Input, InputNumber, message, Row, Space } from 'antd';
import { useState, forwardRef, useImperativeHandle } from 'react';
const { TextArea } = Input;
const OrganizationProfile = (
	{
		api,
		rowData,
		successCallback,
		className,
		leftCol = { lg: 8, sm: 24, xl: 6 },
		rightCol = { lg: 16, sm: 24, xl: 12 },
		isShowAccountInput,
	},
	ref,
) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const to = <div className="flex items-center justify-center px-4 flex-shrink-0">to</div>;
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
			<Row gutter={24}>
				<Col {...leftCol} className="!flex flex-col items-center">
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
					{!isShowAccountInput && (
						<>
							<div className="text-2xl font-medium text-black-dark">
								Organization name
							</div>
							<div className="text-xl font-medium text-[rgba(0,0,0,0.65)] pt-3">
								abc@example.com
							</div>
						</>
					)}
				</Col>
				<Col {...rightCol}>
					{isShowAccountInput && (
						<>
							<Row gutter={16}>
								<InputNameCombination
									layoutCol={{ span: 12 }}
									firstName={['first_name']}
									lastName={['last_name']}
									hideMiddleName={true}
								/>
							</Row>
							<Form.Item
								label="Organization name"
								name={['org_name']}
								rules={[
									{
										required: true,
									},
								]}
							>
								<Input allowClear placeholder="Organization name" />
							</Form.Item>
							<Form.Item
								label="Email account"
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
								extra="Email account used to log in to the ZeeBee lender system."
							>
								<Input placeholder="Email account" allowClear />
							</Form.Item>
							{/* <Form.Item
						label="Password"
						name={['password']}
						rules={[
							{
								required: true,
							},
							{
								validator: (_, value) =>
									validateRule(value, passwordRegex, 'Invalid password format'),
							},
						]}
					>
						<Input.Password placeholder="Password" />
					</Form.Item>
					<Form.Item
						label="Confirm password"
						name={['password']}
						rules={[
							{
								required: true,
							},
							{
								validator: (_, value) =>
									validateRule(value, passwordRegex, 'Invalid password format'),
							},
						]}
					>
						<Input.Password placeholder="Password" />
					</Form.Item> */}
							<ConfirmPassWord isShowRule={false} />
						</>
					)}

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
						<Input placeholder="Contact name" allowClear />
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
						<Input placeholder="Contact email" allowClear />
					</Form.Item>
					<Form.Item
						label="Contact phone"
						name={['phone']}
						rules={[
							{
								required: true,
							},
							{
								validator: (_, value) =>
									validateRule(value, phoneRegex, 'Invalid phone number'),
							},
						]}
					>
						<InputPhone placeholder="Contact phone" />
					</Form.Item>
					<Form.Item label="Organization tagline" name={['tagline']}>
						<TextArea
							placeholder="Organization tagline"
							autoSize={{ minRows: 2, maxRows: 2 }}
						/>
					</Form.Item>
					<Form.Item label="Website link" name={['link']}>
						<Input addonBefore="https://" placeholder="Website link" allowClear />
					</Form.Item>
					<div className="!flex">
						<Form.Item label="Loan amount (US $)" name={['amount']} className="flex-1">
							<InputNumAmount placeholder="Min.loan amount" className="!w-full" />
						</Form.Item>
						{to}
						<Form.Item label=" " name={['amount_max']} className="flex-1">
							<InputNumAmount placeholder="Max.loan amount" className="!w-full" />
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item label="Loan term (Months)" name={['term']} className="flex-1">
							<InputNumber
								placeholder="Min.loan term"
								className="!w-full"
								controls={false}
								// suffix="Months"
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['term_max']} className="flex-1">
							<InputNumber
								placeholder="Max.loan term"
								className="!w-full"
								controls={false}
								// suffix="Months"
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item label="Loan rate" name={['rate']} className="flex-1">
							<InputNumber
								placeholder="Min.loan rate"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['rate_max']} className="flex-1">
							<InputNumber
								placeholder="Max.loan rate"
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
								placeholder="Min.loan points"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={['point_max']} className="flex-1">
							<InputNumber
								placeholder="Max.loan points"
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
