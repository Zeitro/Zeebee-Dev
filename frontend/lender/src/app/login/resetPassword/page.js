'use client';
import React, { useState } from 'react';
import { Button, Form } from 'antd';
import { useRouter, useSearchParams } from 'next/navigation';
import { useUserStore } from '@repo/shared-methods/store';
import styles from '../index.module.scss';
import ConfirmPassWord from '../components/ConfirmPassword';
import { throttle } from 'lodash';
import useContent from '../components/useContent';
const ResetPassword = () => {
	const [form] = Form.useForm();
	const router = useRouter();
	const searchParams = useSearchParams();
	const { resetPassword } = useUserStore();
	const { loginByToken } = useContent();
	const [loading, setLoading] = useState(false);

	// 重置密码并登录
	const handleResetPassword = throttle(async () => {
		try {
			const code = searchParams.get('code');
			const email = searchParams.get('email');
			if (!code || !email) return;
			const values = await form.validateFields();
			setLoading(true);
			const access = await resetPassword({
				email,
				code,
				...values,
			});
			loginByToken(access, setLoading);
		} catch {
			setLoading(false);
		}
	}, 1000);

	return (
		<div className="h-full bg-[#f5f5f5] flex items-center justify-center">
			<Form
				className={[
					styles.loginForm,
					'w-[400px] rounded-lg bg-background !px-6 !py-6',
				].join(' ')}
				name="resetPasswordForm"
				size="large"
				layout="vertical"
				autoComplete="off"
				form={form}
				style={{
					boxShadow:
						'0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
				}}
			>
				<div className="text-3xl/10 font-medium">Reset your password</div>
				<p style={{ color: 'rgba(0, 0, 0, 0.65)' }} className="pt-4 pb-8">
					Create a new secure password for {searchParams.get('email')}
				</p>
				<ConfirmPassWord
					form={form}
					passwordLabel="New password"
					confirmPasswordLabel="Confirm new password"
					passwordPlaceholder="New password"
					confirmPasswordPlaceholder="Re-enter new password"
				/>
				<div className="flex align-center justify-center gap-4 pt-8">
					<Button size="middle" onClick={() => router.replace('/login')}>
						Cancel
					</Button>
					<Button
						size="middle"
						type="primary"
						loading={loading}
						onClick={handleResetPassword}
					>
						Reset password
					</Button>
				</div>
			</Form>
		</div>
	);
};
export default ResetPassword;
