'use client';
import React, { useState } from 'react';
import { Button, Form, Input } from 'antd';
import { useUserStore } from '@repo/shared-methods/store';
import { emailRegex } from '@repo/shared-methods/utils';
import { MailOutlined } from '@ant-design/icons';
import Image from 'next/image';
const SendResetPasswordEmail = ({ form, setPageType }) => {
	const { sendResetPasswordEmail } = useUserStore();
	const [apiError, setApiError] = useState('');
	const [loading, setLoading] = useState(false);
	const onResetPasswordEmail = async () => {
		try {
			setApiError('');
			const values = await form.validateFields();
			setLoading(true);
			await sendResetPasswordEmail(values.email, { isHideTips: true });
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
						validator: (_, value) => {
							if (value && emailRegex.test(value)) {
								return Promise.resolve();
							}
							setApiError('Please enter a valid email address');
							return Promise.reject('enter a valid email');
						},
					},
				]}
				help={apiError}
				validateStatus={apiError ? 'error' : ''}
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
					onChange={() => setApiError('')}
					allowClear
					onKeyDown={(e) => {
						if (e.key === 'Enter') {
							onResetPasswordEmail();
						}
					}}
				/>
			</Form.Item>
			<Form.Item label={null}>
				<Button
					type="primary"
					onClick={onResetPasswordEmail}
					className="w-full"
					loading={loading}
				>
					Reset password
				</Button>
			</Form.Item>

			<Form.Item label={null}>
				<Button type="link" onClick={() => setPageType('signIn')} className="w-full">
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
