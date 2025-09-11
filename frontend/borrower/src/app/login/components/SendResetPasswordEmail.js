'use client';
import React, { useEffect, useState } from 'react';
import { Button, Form, Input } from 'antd';
import { useRouter, usePathname } from 'next/navigation';
import { useUserStore } from '@/store/user';
import { validateRule, emailRegex } from '@/lib/format';
import { MailOutlined, LockOutlined } from '@ant-design/icons';
import Image from 'next/image';
const SendResetPasswordEmail = ({ form, setPageType }) => {
	const router = useRouter();
	const { getUserInfo, login, sendResetPasswordEmail } = useUserStore();
	const [apiError, setApiError] = useState('');
	const [loading, setLoading] = useState(false);
	const onResetPassword = async () => {
		try {
			setApiError('');
			// const values = await form.validateFields();
			setLoading(true);
			// const access = await sendResetPasswordEmail(values, { isHideTips: true });
			setLoading(false);
			setPageType('resetPasswordEmailTips');
		} catch (error) {
			setApiError(error.message);
			setLoading(false);
		}
	};
	return (
		<>
			<div className="text-[30px]/[40px] font-medium pb-6">
				Enter your email to reset password
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
				help={apiError}
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
			<Form.Item label={null}>
				<Button
					type="primary"
					onClick={onResetPassword}
					className="w-full"
					loading={loading}
				>
					Reset password
				</Button>
			</Form.Item>

			<Form.Item label={null}>
				<Button
					type="link"
					onClick={() => setPageType('signIn')}
					className="w-full"
					loading={loading}
				>
					Cancel
				</Button>
			</Form.Item>
		</>
	);
};

export default SendResetPasswordEmail;

export const ResetPasswordEmailTips = ({ setPageType }) => {
	return (
		<div className="flex flex-col items-center">
			<Image src="/images/email.svg" width={48} height={48} alt="email" />
			<div className="text-base pt-6 pb-8">
				If an account exists with this email address, you will receive password reset
				instructions within the next few minutes. Please check your inbox and spam folder.
			</div>
			<Button
				type="link"
				onClick={() => {
					setPageType('signIn');
				}}
			>
				Back to login
			</Button>
		</div>
	);
};
