'use client';
import React from 'react';
import { Button } from 'antd';
import Image from 'next/image';
import { useRouter, useSearchParams } from 'next/navigation';
const SendRegisterTips = () => {
	const router = useRouter();
	const searchParams = useSearchParams();

	return (
		<div className="h-full bg-[#f5f5f5] flex items-center justify-center text-sm">
			<div
				className="w-[400px] rounded-lg px-6 py-6 flex flex-col items-center gap-8 text-center bg-background"
				style={{
					boxShadow:
						'0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
				}}
			>
				<Image src="/images/email.svg" width={48} height={48} alt="email" />
				<div className="text-[30px]/[40px] font-medium text-black-dark">
					Check your email
				</div>
				<div className="text-black-dark">
					We’ve sent a verification link to
					<div className="text-base-primary">{searchParams.get('email')}</div>
				</div>
				<div>
					Click the link in the email to activate your account. The link will expire in 24
					hours.
				</div>
				<Button
					type="link"
					onClick={() => {
						router.replace('/login');
					}}
				>
					Back to Log in
				</Button>
			</div>
		</div>
	);
};
export default SendRegisterTips;
