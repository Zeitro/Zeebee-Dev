import React, { useState } from 'react';
import { Row, Col, Form, Select, Input, message, Button, Space } from 'antd';
import {
	emailRegex,
	phoneRegex,
	validateRule,
	createRequest,
	copyToClipboard,
} from '@repo/shared-methods/utils';
import { layoutCol, titleCol } from '@repo/shared-methods/enum';
import { InputNameCombination, InputPhone, NewModal } from '@repo/shared-methods';
import ConfirmPassword from '@repo/shared-methods/componentsUi/ConfirmPassword';
import { useUserStore } from '@/store/user';
const { TextArea } = Input;
// 添加/编辑用户
const MainModal = ({ setOpen, rowData, handleSearch, roles }) => {
	const { request } = createRequest();
	const { getUserInfo, userInfo } = useUserStore();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);

	const handleSave = async () => {
		try {
			// 此处admin端与lender和borrower端不同的一点是：admin创建不需要给用户发送短信邀请确认邮件，默认就是active状态
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/organization_user';
			let method = 'POST';
			let msg = 'User creation successful.';
			if (rowData?.id) {
				url = `/organization_user/${rowData?.id}`;
				method = 'PUT';
				msg = 'User update successful.';
			}
			const res = await request(url, {
				method,
				body: { ...value },
			});
			handleSearch?.(res?.id);
			if (rowData?.id === userInfo?.organization_user_id) {
				await getUserInfo();
				// await getOrganization();
			}
			message.open({ type: 'success', content: msg });
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width={600}
			title={`${rowData?.id ? 'Edit' : 'Add'} user`}
			okText={rowData?.id ? 'Save change' : 'Add'}
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
			footer={[
				<Button key={'cancel'} onClick={() => setOpen(false)}>
					Cancel
				</Button>,
				<Space key={'ok'}>
					{!rowData?.id && (
						<Button
							key={'cancel'}
							type="primary"
							ghost
							onClick={() => {
								const v = form.getFieldsValue();
								copyToClipboard(`Email: ${v.email}\nPassword: ${v.password}`);
							}}
						>
							Copy account
						</Button>
					)}
					<Button
						key={'save'}
						type="primary"
						onClick={handleSave}
						disabled={confirmLoading}
						loading={confirmLoading}
					>
						{rowData?.id ? 'Save change' : 'Add'}
					</Button>
				</Space>,
			]}
		>
			<Form
				form={form}
				name="user-form"
				layout="vertical"
				autoComplete="off"
				className={[
					' overflow-hidden overflow-y-auto',
					rowData?.id ? 'h-[64vh]' : 'h-[56vh]',
				].join(' ')}
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
				initialValues={{ ...rowData, role_ids: rowData?.roles.map((i) => i.id) }}
			>
				<Row gutter={16}>
					<Col {...titleCol} className="!text-xl font-medium text-base-gray pb-4">
						User information
					</Col>
					{/* 添加用户时：borrower和lender端不显示 ，只有admin端添加用户需要传name相关*/}
					<InputNameCombination
						layoutCol={{ ...layoutCol, md: 12 }}
						firstName={['first_name']}
						lastName={['last_name']}
						hideMiddleName={true}
					/>
					<Col {...titleCol}>
						<Form.Item
							label="Email"
							name={['email']}
							rules={[
								{
									required: true,
								},
								{
									validator: (_, value) =>
										validateRule(value, emailRegex, 'Invalid email format'),
								},
							]}
						>
							<Input placeholder="Enter email" disabled={rowData?.id} />
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item
							label="Phone"
							name={['phone']}
							rules={[
								{
									validator: (_, value) =>
										validateRule(value, phoneRegex, 'Invalid phone number'),
								},
							]}
						>
							<InputPhone />
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item label="Remark" name={['remark']}>
							<TextArea
								placeholder="Enter remark"
								autoSize={{ minRows: 3, maxRows: 5 }}
							/>
						</Form.Item>
					</Col>
					<ConfirmPassword
						renderCol={(e) => <Col {...titleCol}>{e}</Col>}
						required={!rowData?.id}
						passwordFormProps={{
							extra: "You can omit the password when editing; it will not affect the user's password usage.",
						}}
						confirmPasswordFormProps={{
							extra: "You can omit the confirm password when editing; it will not affect the user's password usage.",
						}}
					/>
					<Col {...titleCol} className="!text-xl font-medium text-base-gray pb-4">
						User role
					</Col>
					<Col {...titleCol}>
						<Form.Item label="Role" name={['role_ids']} rules={[{ required: true }]}>
							<Select
								options={roles}
								placeholder="Select role"
								mode="multiple"
								fieldNames={{
									label: 'name',
									value: 'id',
								}}
							/>
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};
export default MainModal;
