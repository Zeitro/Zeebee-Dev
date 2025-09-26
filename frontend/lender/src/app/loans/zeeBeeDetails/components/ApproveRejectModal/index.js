import React from 'react';
import { BtnPopUpModal } from '@/components/NewModal/BtnPopUpModal';
import NewModal from '@/components/NewModal';
import { Form, Input } from 'antd';
export const ApproveModal = ({ key, title, ...props }) => {
	const [form] = Form.useForm();
	const handleOk = async () => {
		const values = await form.validateFields();
		console.log(values);
	};
	return (
		<BtnPopUpModal
			type="primary"
			{...props}
			getContent={(setOpen) => (
				<NewModal
					title="Approval"
					onOk={handleOk}
					onCancel={() => setOpen(false)}
					okText="Approve"
				>
					<div className="pb-6">
						Approval will move this loan to processing stage.Please ensure that all
						items in the underwriting checklist have been validated before officially
						approving the loan application.you can add the approval comments below
					</div>
					<Form form={form} layout="vertical" autoComplete="off" size="large">
						<Form.Item
							label="Comment"
							name="comment"
							rules={[
								{
									required: true,
									message: 'Please input your comments!',
								},
							]}
						>
							<Input.TextArea rows={5} />
						</Form.Item>
					</Form>
				</NewModal>
			)}
		>
			Approve
		</BtnPopUpModal>
	);
};

export const RejectedModal = ({ key, title, ...props }) => {
	const [form] = Form.useForm();
	const handleOk = async () => {
		const values = await form.validateFields();
		console.log(values);
	};
	return (
		<BtnPopUpModal
			type="primary"
			danger={true}
			{...props}
			getContent={(setOpen) => (
				<NewModal
					title="Reject"
					onOk={handleOk}
					onCancel={() => setOpen(false)}
					okText="Reject"
				>
					<div className="pb-6">
						reject will move this loan to rejected loan.Please ensure that all items in
						the underwriting checklist have been validated before officially approving
						the loan application. Rejected comment is required
					</div>
					<Form form={form} layout="vertical" autoComplete="off" size="large">
						<Form.Item
							label="Comment"
							name="comment"
							rules={[
								{
									required: true,
									message: 'Please input your comments!',
								},
							]}
						>
							<Input.TextArea rows={5} />
						</Form.Item>
					</Form>
				</NewModal>
			)}
		>
			Reject
		</BtnPopUpModal>
	);
};
