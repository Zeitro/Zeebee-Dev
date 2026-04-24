'use client';

import { Button, Form, Input, message } from 'antd';
import { useState } from 'react';
import { ConfirmPassWord } from '@repo/shared-methods';
import { localStorageUtil, passwordRegex } from '@repo/shared-methods/utils';
import { debounce } from 'lodash';
import { useUserStore } from '@repo/shared-methods/store';
import { useRouter } from 'next/navigation';
const ChangePassword = () => {
	const router = useRouter();
	const { changePassword } = useUserStore();
	const [form] = Form.useForm();
	const [loading, setLoading] = useState(false);
	const [disabled, setDisabled] = useState(true);
	const onFinish = async (values) => {
		try {
			setLoading(true);
			await changePassword(values);
			setLoading(false);
			message.open({
				type: 'success',
				content: 'You’ve successfully changed your password. Please log in again!',
			});
			const account = localStorageUtil.get('account');
			if (account) {
				account.password = values.new_password;
				localStorageUtil.set('account', account);
			}
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
				name="user-form"
				layout="vertical"
				autoComplete="off"
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
					<Input.Password placeholder="Old password" className="w-full" allowClear />
				</Form.Item>
				<ConfirmPassWord passwordPlaceholder="New password" passwordName="new_password" />
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
