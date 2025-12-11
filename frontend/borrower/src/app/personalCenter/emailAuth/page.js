'use client';

import { Button, Form, Input } from 'antd';
import { useState } from 'react';
const { TextArea } = Input;
const EmailAuth = () => {
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
			<Form.Item label="SMTP(Authorized)" name={['authorized']} rules={[{ required: true }]}>
				<TextArea placeholder="Enter SMTP" autoSize={{ minRows: 3, maxRows: 5 }} />
			</Form.Item>
			<Button type="primary" htmlType="submit" loading={loading} size="middle">
				Authorize
			</Button>
		</Form>
	);
};
export default EmailAuth;
