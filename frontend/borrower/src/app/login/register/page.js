'use client';
import React, { useMemo, useState } from 'react';
import { Button, Checkbox, Col, Form, Input, Row } from 'antd';
import { useRouter, useSearchParams } from 'next/navigation';
import { useUserStore } from '@/store/user';
import { validateRule, emailRegex, capitalizeFirstLetter } from '@repo/shared-methods/utils';
import styles from '../index.module.scss';
import ConfirmPassWord from '../components/ConfirmPassword';
import { throttle } from 'lodash';
const Register = () => {
	const [form] = Form.useForm();
	const router = useRouter();
	const searchParams = useSearchParams();
	const newEmail = Form.useWatch('email', form);
	const { registerAndSendEmail } = useUserStore();
	const [loading, setLoading] = useState(false);

	// 从官网pre-approved页面跳转过来
	const perApprovedInfo = useMemo(() => {
		const request_id = searchParams.get('request_id');
		const lender_id = searchParams.get('lender_id');
		return {
			request_id,
			lender_id,
			isPerApproved: request_id && lender_id,
		};
	}, [searchParams]);

	// 注册
	const handleCreateAccount = throttle(async () => {
		try {
			const values = await form.validateFields();
			if (!values.agreement) return;
			delete values.agreement;
			setLoading(true);
			if (perApprovedInfo.isPerApproved) {
				values.pre_approval_request_id = perApprovedInfo.request_id;
				values.target_lender_company_id = Number(perApprovedInfo.lender_id);
			}
			await registerAndSendEmail(values);
			setLoading(false);
			router.replace(`/login/sendRegisterTips?email=${newEmail}`);
		} catch {
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
					perApprovedInfo.isPerApproved ? 'w-[500px]' : 'w-[400px]',
					'rounded-lg bg-background !px-6 !py-6 max-h-[95%] overflow-y-auto',
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
				<div className="text-3xl/10 font-medium">
					{perApprovedInfo.isPerApproved ? (
						<>
							Create an account to processed
							<p className="text-sm/[22px] font-normal">
								To accept your loan terms, you’ll need to create a ZeeBee account.
							</p>
						</>
					) : (
						'Create your account'
					)}
				</div>
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
							getValueFromEvent={(v) => {
								return capitalizeFirstLetter(v.target.value);
							}}
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
							getValueFromEvent={(v) => {
								return capitalizeFirstLetter(v.target.value);
							}}
						>
							<Input placeholder="Last name" allowClear />
						</Form.Item>
					</Col>
				</Row>
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
					{!perApprovedInfo.isPerApproved && (
						<Button
							size="middle"
							onClick={() => router.replace('/login')}
							className="w-[130px]"
						>
							Cancel
						</Button>
					)}
					<Button
						size="middle"
						type="primary"
						loading={loading}
						onClick={handleCreateAccount}
					>
						{perApprovedInfo.isPerApproved
							? 'Create  account & Accept Loan Terms'
							: 'Create account'}
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
