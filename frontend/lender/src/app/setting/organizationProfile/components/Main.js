import { NewModal } from '@repo/shared-methods';
import { emailRegex, validateRule, createRequest } from '@repo/shared-methods/utils';
import { Form, Input, message } from 'antd';
import { useState } from 'react';
const MainModal = ({ rowData, setOpen, getData }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setConfirmLoading(true);
			await request('/organization_email/save_email', {
				method: 'POST',
				body: {
					...values,
					auth_strategy: 'Organization',
					purpose: 'Contact',
				},
			});
			setConfirmLoading(true);
			message.destroy();
			message.open({
				content: `Contact Email ${rowData?.id ? 'edited' : 'added'} successfully`,
				type: 'success',
			});
			getData?.();
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};
	return (
		<NewModal
			width={500}
			title={`${rowData?.id ? 'Edit' : 'Add'} Contact Email`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			{/* <Form
				form={form}
				layout="vertical"
				autoComplete="off"
				className="h-[130px] overflow-hidden overflow-y-auto"
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
				initialValues={{ email: rowData?.email }}
			>
				<div className="text-base-error pb-3">
					Only Gmail and google verified email account are supported.
				</div>
				<Form.Item
					label="Email address"
					name={['email']}
					rules={[
						{
							required: true,
							message: 'Please enter email',
						},
						{
							validator: (_, value) =>
								validateRule(value, emailRegex, 'Invalid email format'),
						},
					]}
				>
					<Input placeholder="Enter email" allowClear />
				</Form.Item>
			</Form> */}
		</NewModal>
	);
};
export default MainModal;
