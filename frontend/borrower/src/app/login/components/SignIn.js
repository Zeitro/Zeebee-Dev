'use client';
import React, { useEffect, useState } from 'react';
import { Button, Checkbox, Col, Form, Input, Row } from 'antd';
import Password from 'antd/es/input/Password';
import { useRouter, usePathname } from 'next/navigation';
import { useUserStore } from '@/store/user';
import { setCookie, clearAllCookie } from '@/lib/cookie';
import { localStorageUtil } from '@/lib/local';
import { validateRule, emailRegex } from '@/lib/format';
import { MailOutlined, LockOutlined } from '@ant-design/icons';
const SignIn = ({ form, setPageType }) => {
	const router = useRouter();
	const { getUserInfo, login, clearUserInfo } = useUserStore();
	const [apiError, setApiError] = useState('');
	const [loading, setLoading] = useState(false);
	const onFinish = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			const access = await login(values, { isHideTips: true });
			await getUserInfo();
			setCookie('token', access?.access_token, access?.access_token_expire * 1000);
			setLoading(false);
			router.push('/');
		} catch (error) {
			setApiError(error.message);
			setLoading(false);
		}
	};
	return (
		<>
			<div className="text-[30px]/[40px] font-medium pb-6">Sign in to ZeeBee</div>
			<Form.Item
				label="Email"
				name="email"
				rules={[
					{
						required: true,
						message: 'Please enter email',
					},
					{
						validator: (_, value) =>
							validateRule(value, emailRegex, 'Please enter a valid email address'),
					},
				]}
			>
				<Input
					placeholder="Email"
					className="w-full"
					prefix={
						<MailOutlined
							className="text-sm"
							style={{ color: 'rgba(0, 0, 0, 0.25)' }}
						/>
					}
					allowClear
				/>
			</Form.Item>
			<Form.Item
				label="Password"
				name="password"
				rules={[
					{
						required: true,
						message: 'Please enter password',
					},
				]}
			>
				<Password
					placeholder="Password"
					className="w-full"
					prefix={
						<LockOutlined
							className="text-sm"
							style={{ color: 'rgba(0, 0, 0, 0.25)' }}
						/>
					}
					allowClear
				/>
			</Form.Item>
			<div className="text-xs/5 text-base-error">{apiError}</div>
			<div className="flex items-center justify-between">
				<Form.Item name="remember" valuePropName="checked" label={null} className="!mb-4">
					<Checkbox size="small">Remember me</Checkbox>
				</Form.Item>
				<Form.Item label={null} className="!mb-4">
					<Button
						type="link"
						size="small"
						onClick={() => setPageType('sendResetPasswordEmail')}
					>
						Forgot password?
					</Button>
				</Form.Item>
			</div>

			<Form.Item label={null}>
				<Button type="primary" onClick={onFinish} className="w-full" loading={loading}>
					Login
				</Button>
			</Form.Item>
			<div className="flex items-center justify-end">
				Don’t have an account?
				<Button type="link" size="small" onClick={() => router.push('/login/register')}>
					Create one!
				</Button>
			</div>
		</>
	);
};

export default SignIn;
