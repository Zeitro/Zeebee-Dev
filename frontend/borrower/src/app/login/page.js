'use client';
import React, {useState} from 'react';
import {Form} from 'antd';
import styles from './index.module.scss';
import SignIn from './components/SignIn';
import SendResetPasswordEmail, {ResetPasswordEmailTips,} from './components/SendResetPasswordEmail';

const Login = () => {
	const [form] = Form.useForm();
	const [pageType, setPageType] = useState('signIn');

	const getContent = () => {
		switch (pageType) {
			case 'signIn':
				return <SignIn form={form} setPageType={setPageType} />;
			case 'sendResetPasswordEmail':
				return <SendResetPasswordEmail form={form} setPageType={setPageType} />;
			case 'resetPasswordEmailTips':
				return <ResetPasswordEmailTips setPageType={setPageType} />;
			default:
				break;
		}
	};

	return (
		<div className="h-full w-full flex items-center justify-center">
			<div className="h-full flex-1 bg-[#f5f5f5] pt-[7.5%]">
				<Form
					className={[
						styles.loginForm,
						'w-[440px] min-h-[484px] rounded-lg bg-background !px-12 !py-12 !m-auto',
					].join(' ')}
					name="basic"
					size="large"
					layout="vertical"
					autoComplete="off"
					form={form}
					style={{
						boxShadow:
							'0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
					}}
				>
					{getContent()}
				</Form>
			</div>
			<div className="h-full flex-1 bg-base-primary">
				<p className="text-[38px]/[46px] font-medium text-background pl-[9%] pt-[11%]">
					The simplest way to manage your loan
				</p>
			</div>
		</div>
	);
};
export default Login;
