'use client';
import React, { useMemo } from 'react';
import { Form, Input } from 'antd';
import { passwordRegex } from '@repo/shared-methods/utils';
import { CheckCircleOutlined } from '@ant-design/icons';

const ConfirmPassWord = ({
	form,
	passwordLabel = 'Password',
	confirmPasswordLabel = 'Confirm password',
	passwordPlaceholder = 'Password',
	confirmPasswordPlaceholder = 'Re-enter password',
}) => {
	const password = Form.useWatch('password', form);
	const passwordRules = useMemo(() => {
		let arr = [
			{
				label: 'Upper and lowercase letters',
				pattern: /(?=.*[a-z])(?=.*[A-Z])/,
			},
			{
				label: 'At least one number',
				pattern: /[0-9]{1,}/,
			},
			{
				label: '8 or more characters',
				pattern: /^\S{8,}$/,
			},
			{
				label: 'Password match',
				pattern: passwordRegex,
			},
		];
		const newArr = arr.map((item) => {
			item.status = '';
			if (password) {
				if (item.pattern.test(password)) {
					item.status = 'success';
				} else {
					item.status = 'error';
				}
			}
			return item;
		});
		return newArr;
	}, [password]);

	return (
		<>
			<Form.Item
				label={passwordLabel}
				name="password"
				rules={[
					{
						required: true,
						message: 'Please enter password',
					},
					{
						pattern: passwordRegex,
						message: 'Incorrect format!',
					},
				]}
			>
				<Input.Password placeholder={passwordPlaceholder} className="w-full" allowClear />
			</Form.Item>
			<Form.Item
				label={confirmPasswordLabel}
				name="confirm_password"
				rules={[
					{
						required: true,
						message: 'Please enter confirm password',
					},
					{
						validator: (_, value) => {
							if (value && value !== password) {
								return Promise.reject(new Error('The two passwords do not match!'));
							}
							return Promise.resolve();
						},
					},
				]}
				dependencies={['password']}
			>
				<Input.Password
					placeholder={confirmPasswordPlaceholder}
					className="w-full"
					allowClear
				/>
			</Form.Item>
			<div className="text-base text-black-dark">Your password must have:</div>
			{passwordRules.map((item, index) => (
				<div
					key={index}
					style={{
						color: item.status
							? item.status === 'success'
								? 'var(--base-success)'
								: 'var(--base-error)'
							: 'rgba(0, 0, 0, 0.65',
					}}
					className="pt-4 flex items-center"
				>
					<CheckCircleOutlined className="mr-1" />
					{item.label}
				</div>
			))}
		</>
	);
};

export default ConfirmPassWord;
