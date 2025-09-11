'use client';
import React, { useEffect, useState } from 'react';
import { Button, Checkbox, Col, Form, Input, Row } from 'antd';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';
import { useUserStore } from '@/store/user';
import { setCookie, clearAllCookie } from '@/lib/cookie';
import { localStorageUtil } from '@/lib/local';
import { validateRule, emailRegex } from '@/lib/format';
import { MailOutlined, LockOutlined } from '@ant-design/icons';
import styles from '../index.module.scss';
import ConfirmPassWord from '../components/ConfirmPassword';

const ResetPassword = () => {
	const [form] = Form.useForm();
	const router = useRouter();
	const { code, email } = useSearchParams();
	const { getUserInfo, login, clearUserInfo } = useUserStore();
	const [loading, setLoading] = useState(false);

	// 重置密码并登录
	const handleResetPassword = (values) => {
		if (!code || !email) return;
		console.log('Received values of form: ', values);
	};

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
					Create a new secure password for {email}
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
