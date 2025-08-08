'use client';
import React, { useEffect } from 'react';
import { Button, Checkbox, Form, Input } from 'antd';
import Item from 'antd/es/form/FormItem';
import Password from 'antd/es/input/Password';
import { useRouter } from 'next/navigation';
import { useUserStore } from '@/store/user';
import { setCookie, clearAllCookie } from '@/lib/cookie';
const Login = () => {
	const [form] = Form.useForm();
	const router = useRouter();
	const { getUserInfo, login, clearInfo } = useUserStore();

	useEffect(() => {
		clearInfo();
		clearAllCookie();
	}, []);

	const onFinish = async () => {
		const values = await form.validateFields();
		const access = await login(values);
		await getUserInfo();
		setCookie('token', access?.access_token, access?.access_token_expire*1000);
		router.push('/');
	};

	return (
		<div className="h-full w-full flex items-center justify-center flex-col">
			<div>Zeebee</div>
			<Form
				name="basic"
				labelCol={{ span: 8 }}
				wrapperCol={{ span: 16 }}
				style={{ maxWidth: 600 }}
				initialValues={{ remember: true }}
				autoComplete="off"
				form={form}
				// className="w-[200px]"
			>
				<Item
					label="Email"
					name="email"
					rules={[{ required: true, message: 'Please input your email!' }]}
				>
					<Input />
				</Item>

				<Item
					label="Password"
					name="password"
					rules={[{ required: true, message: 'Please input your password!' }]}
				>
					<Password />
				</Item>

				<Item name="remember" valuePropName="checked" label={null}>
					<Checkbox>Remember me</Checkbox>
				</Item>

				<Item label={null}>
					<Button type="primary" onClick={onFinish}>
						Submit
					</Button>
				</Item>
			</Form>
		</div>
	);
};
export default Login;
