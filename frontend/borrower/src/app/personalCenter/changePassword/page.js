'use client';

import { Button, Form, Input, message } from 'antd';
import { useState } from 'react';
import { ConfirmPassWord } from '@repo/shared-methods';
import { syncCredentialAfterPasswordReset } from '@repo/shared-methods/hooks';
import { passwordRegex, setPendingLoginCredential } from '@repo/shared-methods/utils';
import { debounce } from 'lodash';
import { useUserStore } from '@/store/user';
import { useRouter } from 'next/navigation';
const ChangePassword = () => {
	const router = useRouter();
	const { changePassword, userInfo } = useUserStore();
	const [form] = Form.useForm();
	const [loading, setLoading] = useState(false);
	const [disabled, setDisabled] = useState(true);
	const userEmail = userInfo?.email;

	const onFinish = async (values) => {
		try {
			setLoading(true);
			const email = userEmail;
			const newPassword = values.new_password;
			if (email && newPassword) {
				await syncCredentialAfterPasswordReset(email, newPassword);
			}
			await changePassword(values);
			if (email && newPassword) {
				setPendingLoginCredential({ email, password: newPassword });
			}
			setLoading(false);
			message.open({
				type: 'success',
				content: 'You’ve successfully changed your password. Please log in again!',
			});
			router.push('/login');
		} catch {
			setLoading(false);
		}
	};
	return (
		<div className="w-full h-full flex items-center flex-col pt-12">
			<div className="text-3xl/10 font-medium text-black-dark pb-8">Change password</div>
			<Form
				form={form}
				name="change-password-form"
				layout="vertical"
				autoComplete="on"
				className="w-[355px]"
				size="large"
				disabled={loading}
				scrollToFirstError
				onFinish={onFinish}
				onFieldsChange={debounce((changedFields, allFields) => {
					const isSuccess = allFields.every(
						(item) => item.validated && !item.errors?.length,
					);
					setDisabled(!isSuccess);
				}, 500)}
			>
				{userEmail ? (
					<input
						type="text"
						name="username"
						autoComplete="username"
						value={userEmail}
						readOnly
						tabIndex={-1}
						aria-hidden="true"
						className="absolute h-0 w-0 overflow-hidden opacity-0 pointer-events-none"
					/>
				) : null}
				<Form.Item
					label="Old Password"
					name="old_password"
					rules={[
						{
							required: true,
							message: 'Please enter old password',
						},
						{
							pattern: passwordRegex,
							message: 'Incorrect format!',
						},
					]}
				>
					<Input.Password
						name="old_password"
						autoComplete="current-password"
						placeholder="Old password"
						className="w-full"
						allowClear
					/>
				</Form.Item>
				<ConfirmPassWord
					passwordPlaceholder="New password"
					passwordName="new_password"
					passwordAutoComplete="new-password"
					confirmPasswordAutoComplete="new-password"
				/>
				<div className="pt-8 flex justify-center">
					<Button
						type="primary"
						htmlType="submit"
						loading={loading}
						size="middle"
						disabled={disabled}
					>
						Change password
					</Button>
				</div>
			</Form>
		</div>
	);
};
export default ChangePassword;
