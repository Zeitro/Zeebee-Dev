import React, { useState } from 'react';
import { NewModal } from '@repo/shared-methods';
import { Form, Input, message } from 'antd';
import { eventBus, createRequest } from '@repo/shared-methods/utils';
const { TextArea } = Input;
const EditCommentModal = ({ open, setOpen, record = {} }) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleOK = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			await request(`/loan_document/${record.id}/comment`, {
				method: 'PATCH',
				body: {
					...values,
				},
			});
			message.open({
				type: 'success',
				content: 'Update comment successfully',
			});
			eventBus.emit('getLoanDetail');
			setLoading(false);
			setOpen(false);
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<NewModal
			title="Edit comment"
			width={500}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
			okText="Save"
			confirmLoading={loading}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				className="h-[300px] overflow-hidden overflow-y-auto"
				initialValues={{
					comment: record.comment,
				}}
			>
				<div className="pb-4">
					<label>Document type</label>
					<div className="font-medium text-black-dark">{record.file_type}</div>
				</div>
				<Form.Item label="Comment" name="comment" rules={[{ required: true }]}>
					<TextArea placeholder="Add a comment" autoSize={{ minRows: 3, maxRows: 5 }} />
				</Form.Item>
			</Form>
		</NewModal>
	);
};

export default EditCommentModal;
export const EditCommentModalBtn = ({ children, isStopPropagation, record }) => {
	const [open, setOpen] = useState(false);
	return (
		<>
			<div
				onClick={(e) => {
					if (isStopPropagation) {
						e.stopPropagation();
						e.preventDefault();
					}
					setOpen(true);
				}}
			>
				{children || 'Edit comment'}
			</div>

			{open && <EditCommentModal open={open} setOpen={setOpen} record={record} />}
		</>
	);
};
