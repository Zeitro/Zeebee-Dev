import React, { useState } from 'react';
import { Row, Col, Form, Select, Input } from 'antd';
import { emailRegex, phoneRegex, validateRule, createRequest } from '@repo/shared-methods/utils';
import { layoutCol, titleCol } from '@/enum/common';
import { InputNameCombination, InputPhone, NewModal } from '@repo/shared-methods';
const { TextArea } = Input;
const MainModal = ({ setOpen, rowData, handleSearch, roles }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/organization_user';
			let method = 'POST';
			if (rowData?.id) {
				url = `/organization_user/${rowData?.id}`;
				method = 'PUT';
			}
			const res = await request(url, {
				method,
				body: { ...value },
			});
			handleSearch?.(res?.id);
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
					{rowData?.id && (
						<InputNameCombination
							layoutCol={{ ...layoutCol, md: 12 }}
							firstName={['first_name']}
							lastName={['last_name']}
							hideMiddleName={true}
						/>
					)}
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
