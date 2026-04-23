import React, { useState } from 'react';
import { BtnPopUpModal, NewModal } from '@repo/shared-methods';
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { Form, Input, message } from 'antd';
import { createRequest, eventBus } from '@repo/shared-methods/utils';

export const ApproveRejectModal = ({
	pending_status,
	id,
	action = 'Approve',
	setNextStep = () => {},
	...props
}) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	let url = `/loan/${id}/review_pending_status`;
	const isApprove = action === 'Approve';
	const [loading, setLoading] = useState(false);
	const handleOk = async (setOpen) => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			const res = await request(url, {
				method: 'POST',
				body: {
					action,
					current_pending_status: pending_status,
					...values,
				},
			});
			message.open({
				type: 'success',
				content: `${isApprove ? 'Approved' : 'Rejected'} successfully`,
			});
			eventBus.emit('getLoanDetail');
			setLoading(false);
			setOpen(false);
			setNextStep();
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<BtnPopUpModal
			type="primary"
			danger={!isApprove}
			ghost={!isApprove}
			icon={isApprove ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
			{...props}
			getContent={(setOpen) => (
				<NewModal
					title={action}
					onOk={() => handleOk(setOpen)}
					onCancel={() => setOpen(false)}
					okText={action}
					confirmLoading={loading}
				>
					<div className="pb-6">
						{isApprove
							? `Approval will move this loan to processing stage.Please ensure that all
						items in the underwriting checklist have been validated before officially
						approving the loan application.you can add the approval comments below`
							: `Reject will move this loan to processing.Please ensure that all items in
						the underwriting checklist have been validated before officially approving
						the loan application. Rejected comment is required`}
					</div>
					<Form form={form} layout="vertical" autoComplete="off" size="large">
						<Form.Item
							label="Comment"
							name="comment"
							rules={[
								{
									required: !isApprove,
									message: 'Please input your comment!',
								},
							]}
						>
							<Input.TextArea rows={5} />
						</Form.Item>
					</Form>
				</NewModal>
			)}
		>
			{action}
		</BtnPopUpModal>
	);
};

export const RejectedModal = ({ pending_status, id, title, ...props }) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	let url = `loan/${id}/review_pending_status`;
	const [loading, setLoading] = useState(false);
	const handleOk = async (setOpen) => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			const res = await request(url, {
				method: 'POST',
				body: {
					action: 'Approve',
					current_pending_status: pending_status,
					...values,
				},
			});
			console.log(values);
			message.open({
				type: 'success',
				content: 'Approved successfully',
			});
			setLoading(false);
			setOpen(false);
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<BtnPopUpModal
			type="primary"
			danger={true}
			ghost={true}
			icon={<CloseCircleOutlined />}
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
