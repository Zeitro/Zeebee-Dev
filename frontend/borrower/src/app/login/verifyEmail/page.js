'use client';
import React, { useState } from 'react';
import { Button } from 'antd';
import { useUserStore } from '@/store/user';
import { CheckCircleTwoTone } from '@ant-design/icons';
import { useRouter, useSearchParams } from 'next/navigation';
import { setCookie } from '@/lib/cookie';

// 注册成功后，邮箱打开此页面，触发验证并登录
const VerifyEmail = () => {
	const router = useRouter();
	const searchParams = useSearchParams();
	const { verifyEmail, getUserInfo } = useUserStore();
	const [loading, setLoading] = useState(false);
	const handleActivateAccount = async () => {
		try {
			setLoading(true);
			const email = searchParams.get('email');
			const code = searchParams.get('code');
			const access = await verifyEmail({ email, code });
			await getUserInfo();
			setCookie('token', access?.access_token, access?.access_token_expire * 1000);
			setLoading(false);
			router.replace('/');
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<div className="h-full bg-[#f5f5f5] flex items-center justify-center text-sm">
			<div
				className="w-[400px] rounded-lg px-6 py-6 flex flex-col items-center gap-8 text-center bg-background"
				style={{
					boxShadow:
						'0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
				}}
			>
				<CheckCircleTwoTone />
				<div className="text-[30px]/[40px] font-medium text-black-dark">
					Email Verified Successfully!
				</div>
				<div>
					Your email address has been confirmed. Click the button below to activate your
					account and start using cur services.
				</div>
				<Button type="primary" onClick={handleActivateAccount} loading={loading}>
					Activate account
				</Button>
			</div>
		</div>
	);
};
export default VerifyEmail;
