'use client';
import React, { useCallback, useState } from 'react';
import { Button, Form } from 'antd';
import { useRouter, useSearchParams } from 'next/navigation';
import { useUserStore } from '@repo/shared-methods/store';
import { syncCredentialAfterPasswordReset } from '@repo/shared-methods/hooks';
import styles from '../index.module.scss';
import ConfirmPassWord from '../components/ConfirmPassword';
import useData from '../components/useData';

const ResetPassword = () => {
	const [form] = Form.useForm();
	const { getUserPermissions } = useData();
	const router = useRouter();
	const searchParams = useSearchParams();
	const { resetPassword } = useUserStore();
	const [loading, setLoading] = useState(false);

	const email = searchParams.get('email');
	const code = searchParams.get('code');

	const onFinish = useCallback(
		async (values) => {
			if (!code || !email) return;
			try {
				setLoading(true);
				await syncCredentialAfterPasswordReset(email, values.password);
				const access = await resetPassword({
					email,
					code,
					...values,
				});
				await getUserPermissions(access, setLoading);
			} catch {
				setLoading(false);
			}
		},
		[code, email, getUserPermissions, resetPassword],
	);

	return (
		<div className="h-full bg-[#f5f5f5] flex items-center justify-center">
			<Form
				className={[
					styles.loginForm,
					'w-[400px] rounded-lg bg-background !px-6 !py-6',
				].join(' ')}
				name="reset-password-form"
				size="large"
				layout="vertical"
				autoComplete="on"
				form={form}
				onFinish={onFinish}
				style={{
					boxShadow:
						'0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
				}}
			>
				<div className="text-3xl/10 font-medium">Reset your password</div>
				<p style={{ color: 'rgba(0, 0, 0, 0.65)' }} className="pt-4 pb-8">
					Create a new secure password for {email}
				</p>
				{email ? (
					<input
						type="text"
						name="username"
						autoComplete="username"
						value={email}
						readOnly
						tabIndex={-1}
						aria-hidden="true"
						className="absolute h-0 w-0 overflow-hidden opacity-0 pointer-events-none"
					/>
				) : null}
				<ConfirmPassWord
					form={form}
					passwordLabel="New password"
					confirmPasswordLabel="Confirm new password"
					passwordPlaceholder="New password"
					confirmPasswordPlaceholder="Re-enter new password"
					passwordAutoComplete="new-password"
					confirmPasswordAutoComplete="new-password"
				/>
				<div className="flex align-center justify-center gap-4 pt-8">
					<Button size="middle" onClick={() => router.replace('/login')}>
						Cancel
					</Button>
					<Button size="middle" type="primary" htmlType="submit" loading={loading}>
						Reset password
					</Button>
				</div>
			</Form>
		</div>
	);
};
export default ResetPassword;
