'use client';
import React, {useState} from 'react';
import {Button, Checkbox, Col, Form, Input, Row} from 'antd';
import {useRouter, useSearchParams} from 'next/navigation';
import {useUserStore} from '@repo/shared-methods/store';
import {emailRegex, validateRule} from '@repo/shared-methods/utils';
import styles from '../index.module.scss';
import ConfirmPassWord from '../components/ConfirmPassword';
import {throttle} from 'lodash';

const Register = () => {
	const [form] = Form.useForm();
	const router = useRouter();
	const searchParams = useSearchParams();
	const newEmail = Form.useWatch('email', form);
	const { registerAndSendEmail } = useUserStore();
	const [loading, setLoading] = useState(false);

	// 注册
	const handleCreateAccount = throttle(async () => {
		try {
			const values = await form.validateFields();
			if (!values.agreement) return;
			delete values.agreement;
			setLoading(true);
			await registerAndSendEmail(values);
			setLoading(false);
			router.replace(`/login/sendRegisterTips?email=${newEmail}`);
		} catch (error) {
			setLoading(false);
		}
	}, 1000);

	// 服务条款
	const handleTermsService = (e) => {
		e.preventDefault();
		e.stopPropagation();
	};

	// 隐私政策
	const handlePrivacyPolicy = (e) => {
		e.preventDefault();
		e.stopPropagation();
	};

	return (
		<div className="h-full bg-[#f5f5f5] flex items-center justify-center overflow-hidden">
			<Form
				className={[
					styles.loginForm,
					'w-[400px] rounded-lg bg-background !px-6 !py-6 max-h-[95%] overflow-y-auto',
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
				initialValues={{ email: searchParams.get('email') }}
			>
				<div className="text-3xl/10 font-medium">Create your account</div>
				<Row gutter={16} className="pt-8">
					<Col sm={12} xs={24}>
						<Form.Item
							label="First name"
							name="first_name"
							rules={[
								{
									required: true,
									message: 'Please enter first name!',
								},
							]}
						>
							<Input placeholder="First name" allowClear />
						</Form.Item>
					</Col>
					<Col sm={12} xs={24}>
						<Form.Item
							label="Last name"
							name="last_name"
							rules={[
								{
									required: true,
									message: 'Please enter last name!',
								},
							]}
							dependencies={['last_name']}
						>
							<Input placeholder="Last name" allowClear />
						</Form.Item>
					</Col>
				</Row>
				<Form.Item
					label="Organization name"
					name="organization_name"
					rules={[
						{
							required: true,
						},
					]}
				>
					<Input placeholder="Organization name" allowClear minLength={2} />
				</Form.Item>
				<Form.Item
					label="Email address"
					name="email"
					rules={[
						{
							required: true,
							message: 'Please enter email',
						},
						{
							validator: (_, value) =>
								validateRule(
									value,
									emailRegex,
									'Please enter a valid email address',
								),
						},
					]}
				>
					<Input placeholder="Email" className="w-full" allowClear />
				</Form.Item>
				<ConfirmPassWord form={form} />
				<Form.Item
					name="agreement"
					rules={[
						{
							required: true,
							pattern: /true/,
							message:
								'You must accept the Terms of Service and Privacy Policy before creating an account.',
						},
					]}
					valuePropName="checked"
					label={null}
					className="pt-6"
				>
					<Checkbox size="small">
						I agree to the
						<span
							className="text-base-primary cursor-pointer"
							onClick={handleTermsService}
						>
							&nbsp;Terms of Service&nbsp;
						</span>
						and
						<span
							className="text-base-primary cursor-pointer"
							onClick={handlePrivacyPolicy}
						>
							&nbsp;Privacy Policy&nbsp;
						</span>
					</Checkbox>
				</Form.Item>
				<div className="flex align-center justify-center gap-4">
					<Button
						size="middle"
						onClick={() => router.replace('/login')}
						className="w-[130px]"
					>
						Cancel
					</Button>
					<Button
						size="middle"
						type="primary"
						loading={loading}
						onClick={handleCreateAccount}
					>
						Create account
					</Button>
				</div>
				<div className="w-full text-center pt-6 text-black-dark">
					Already have an account?
					<span
						className="text-base-primary cursor-pointer"
						onClick={() => router.replace('/login')}
					>
						&nbsp;Sign in
					</span>
				</div>
			</Form>
		</div>
	);
};
export default Register;
