'use client';

import { Button, Form, Input } from 'antd';
import { useState } from 'react';
import ConfirmPassword from '@/app/login/components/ConfirmPassword';
import { passwordRegex } from '@repo/shared-methods/utils';
const { TextArea } = Input;
const ChangePassword = () => {
	const [form] = Form.useForm();
	const [loading, setLoading] = useState(false);
	const onFinish = async (values) => {
		console.log(values);
	};
	return (
		<Form
			form={form}
			name="user-form"
			layout="vertical"
			autoComplete="off"
			className="!px-6 !py-6"
			size="large"
			disabled={loading}
			scrollToFirstError
			onFinish={onFinish}
		>
			<Form.Item
				label="Old Password"
				name="old_password"
				rules={[
					{
						required: true,
						message: 'Please enter old password',
					},
					{
						pattern: passwordRegex,
						message: 'Incorrect format!',
					},
				]}
			>
				<Input.Password placeholder="Old password" className="w-full" allowClear />
			</Form.Item>
			<ConfirmPassword />
			<Button type="primary" htmlType="submit" loading={loading} size="middle">
				Change
			</Button>
		</Form>
	);
};
export default ChangePassword;
