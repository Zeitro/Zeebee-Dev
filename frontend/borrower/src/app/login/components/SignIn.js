'use client';
import React, {useEffect, useState} from 'react';
import {Button, Checkbox, Form, Input, message, Modal} from 'antd';
import Password from 'antd/es/input/Password';
import {useRouter} from 'next/navigation';
import {usePermissionStore, useUserStore} from '@repo/shared-methods/store';
import {emailRegex, localStorageUtil, validateRule} from '@repo/shared-methods/utils';
import {LockOutlined, MailOutlined} from '@ant-design/icons';
import {initConfigFn} from '@repo/shared-methods';
import useData from './useData';

const SignIn = ({ form, setPageType }) => {
	const router = useRouter();
	const { getUserPermissions } = useData();
	const { resendVerifyEmail } = useUserStore();
	const { login } = usePermissionStore();
	const [apiError, setApiError] = useState('');
	const [loading, setLoading] = useState(false);
	const [verifyLoading, setVerifyLoading] = useState(false);

	const [modal, contextHolder] = Modal.useModal();
	const email = Form.useWatch('email', form);
	const password = Form.useWatch('password', form);
	const remember = Form.useWatch('remember', form);

	useEffect(() => {
		const account = localStorageUtil.get('account');
		if (account?.email || account?.password) {
			form.setFieldsValue({
				email: account.email,
				password: account.password,
				remember: true,
			});
		}
	}, [form]);

	// 保存登录信息
	useEffect(() => {
		if (remember) {
			localStorageUtil.set('account', { email, password });
		} else {
			localStorageUtil.remove('account');
		}
	}, [email, password, remember]);

	// 重新验证邮箱（针对已注册未验证用户）
	const reVerify = async () => {
		try {
			setVerifyLoading(true);
			await resendVerifyEmail(email);
			setVerifyLoading(false);
			message.open({
				type: 'success',
				content: 'Email sent successfully! Please check your email for verification!',
			});
		} catch {
			setVerifyLoading(false);
		}
	};
	const onFinish = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			const access = await login(values, { isHideTips: true });
			// await getUserInfo();
			// setCookie('token', access?.access_token, access?.access_token_expire * 1000);
			// setLoading(false);
			// router.push('/');
			await getUserPermissions(access, setLoading);
		} catch (error) {
			if (error.code === 401001) {
				const initValue = initConfigFn({
					title: 'Check your email',
					content: (
						<div>
							<div className="text-black-dark pb-6">
								We will sent a verification link to
								<p className="text-base-primary">{email}</p>
							</div>
							<div className="text-gray-dark">
								Click the link in the email to activate your account. The link will
								expire in 24 hours.
							</div>
						</div>
					),
					okText: 'Send email',
					onOk: reVerify,
					confirmLoading: verifyLoading,
				});
				modal.confirm(initValue);
			}
			setApiError(error.message);
			setLoading(false);
		}
	};
	return (
		<>
			{contextHolder}
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
					onKeyDown={(e) => {
						if (e.key === 'Enter') {
							onFinish();
						}
					}}
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
