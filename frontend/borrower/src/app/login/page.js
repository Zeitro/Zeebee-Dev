'use client';
import React from 'react';
import { Button, Checkbox, Form, Input } from 'antd';
import Item from 'antd/es/form/FormItem';
import Password from 'antd/es/input/Password';
import { useRouter } from 'next/navigation';
import {createRequest} from '@/lib/fetch';
import {setLocalStorage} from '@/lib/local';
import {useUserStore} from '@/store/user';
const Login = () => {
    const [form] = Form.useForm();
    const router = useRouter();
	const {request} = createRequest();
	const {getUserInfo} = useUserStore();

	const onFinish = async () => {
        const values = await form.validateFields();
		const data = await request('/user/login',{
			method:'POST',
			body:values
		});
		await getUserInfo();
		setLocalStorage('access',data)
        router.push('/');
	};

	return (
		<div className='h-full w-full flex items-center justify-center flex-col'>
			<div>Zeebee</div>
			<Form
				name="basic"
				labelCol={{ span: 8 }}
				wrapperCol={{ span: 16 }}
				style={{ maxWidth: 600 }}
				initialValues={{ remember: true }}
				autoComplete="off"
                form={form}
                // className="w-[200px]"
			>
				<Item
					label="Email"
					name="email"
					rules={[{ required: true, message: 'Please input your email!' }]}
				>
					<Input />
				</Item>

				<Item
					label="Password"
					name="password"
					rules={[{ required: true, message: 'Please input your password!' }]}
				>
					<Password />
				</Item>

				<Item name="remember" valuePropName="checked" label={null}>
					<Checkbox>Remember me</Checkbox>
				</Item>

				<Item label={null}>
					<Button type="primary" onClick={onFinish}>
						Submit
					</Button>
				</Item>
			</Form>
		</div>
	);
};
export default Login;
