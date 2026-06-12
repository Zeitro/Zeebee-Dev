'use client';
import React, { useCallback, useEffect, useState } from 'react';
import { Button, Checkbox, Form, Input, message, Modal } from 'antd';
import Password from 'antd/es/input/Password';
import { useRememberLoginEmail } from '@repo/shared-methods/hooks';
import { useUserStore, usePermissionStore } from '@repo/shared-methods/store';
import {
	validateRule,
	emailRegex,
	storeBrowserPasswordCredential,
} from '@repo/shared-methods/utils';
import { MailOutlined, LockOutlined } from '@ant-design/icons';
import { initConfigFn } from '@repo/shared-methods';
import useData from './useData';

const SignIn = ({ form, setPageType, registerFinish }) => {
	const { getUserPermissions } = useData();
	const { resendVerifyEmail } = useUserStore();
	const { login } = usePermissionStore();
	const [apiError, setApiError] = useState('');
	const [loading, setLoading] = useState(false);
	const [verifyLoading, setVerifyLoading] = useState(false);

	const [modal, contextHolder] = Modal.useModal();
	const email = Form.useWatch('email', form);
	const remember = Form.useWatch('remember', form);
	const { persistRememberedEmail, shouldAutofill } = useRememberLoginEmail(form);
	const [autofillGuardReleased, setAutofillGuardReleased] = useState(false);

	// 保存邮箱
	useEffect(() => {
		persistRememberedEmail(email, remember);
	}, [email, persistRememberedEmail, remember]);

	// 重新验证邮箱（针对已注册未验证用户）
	const reVerify = useCallback(async () => {
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
	}, [email, resendVerifyEmail, setVerifyLoading]);

	const onFinish = useCallback(
		async (values) => {
			try {
				setLoading(true);
				const access = await login(values, { isHideTips: true });
				if (values.remember) {
					await storeBrowserPasswordCredential(values.email, values.password);
				}
				persistRememberedEmail(values.email, values.remember);
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
									Click the link in the email to activate your account. The link
									will expire in 24 hours.
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
		},
		[email, getUserPermissions, login, modal, persistRememberedEmail, reVerify, verifyLoading],
	);

	useEffect(() => {
		registerFinish?.(onFinish);
		return () => registerFinish?.(null);
	}, [registerFinish, onFinish]);

	return (
		<>
			{contextHolder}
			<div className="text-[30px]/[40px] font-medium pb-6">
				Sign in to Zeebee Admin Console
			</div>
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
					id="email"
					name="email"
					autoComplete={shouldAutofill ? 'username' : 'off'}
					readOnly={!shouldAutofill && !autofillGuardReleased}
					onFocus={() => !shouldAutofill && setAutofillGuardReleased(true)}
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
					id="password"
					name="password"
					autoComplete="current-password"
					readOnly={!shouldAutofill && !autofillGuardReleased}
					onFocus={() => !shouldAutofill && setAutofillGuardReleased(true)}
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
				<Button type="primary" htmlType="submit" className="w-full" loading={loading}>
					Login
				</Button>
			</Form.Item>
			{/* <div className="flex items-center justify-end">
				Don’t have an account?
				<Button type="link" size="small" onClick={() => router.push('/login/register')}>
					Create one!
				</Button>
			</div> */}
		</>
	);
};

export default SignIn;
