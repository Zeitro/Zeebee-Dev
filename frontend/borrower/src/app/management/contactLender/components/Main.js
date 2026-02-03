import React, { useState, useCallback } from 'react';
import { Row, Col, Form, Input, Spin } from 'antd';
import { emailRegex, phoneRegex, validateRule, createRequest } from '@repo/shared-methods/utils';
import { layoutCol } from '@/enum/common';
import { InputPhone, NewModal } from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const MainModal = ({ setOpen, rowData, handleSearch }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);

	const handleDetails = useCallback(async () => {
		try {
			setLoading(true);
			const res = await request(`/borrower_applicant/${rowData?.id}`, {
				method: 'GET',
			});
			form.setFieldsValue(res);
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [form, rowData]);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/borrower_applicant';
			let method = 'POST';
			if (rowData?.id) {
				url = `/borrower_applicant/${rowData?.id}`;
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

	useDeferredSyncEffect({
		predicate: () => !!rowData?.id,
		deferred: () => {
			handleDetails();
		},
		deps: [handleDetails, rowData?.id],
		isUseMicrotask: true,
	});
	return (
		<NewModal
			width={600}
			title={`${rowData?.id ? 'Edit' : 'Add'} contact lender`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Spin spinning={loading}>
				<Form
					form={form}
					layout="vertical"
					autoComplete="off"
					className="h-[200px] overflow-hidden overflow-y-auto"
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
				>
					<Row gutter={16}>
						<Col {...layoutCol}>
							<Form.Item
								label="Lender name"
								name={['name']}
								rules={[
									{
										required: true,
									},
								]}
							>
								<Input placeholder="Enter lender name" />
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
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
								<Input placeholder="Enter email" />
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
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
						<Col {...layoutCol}>
							<Form.Item label="Contact name" name={['contact_name']}>
								<Input placeholder="Enter contact name" />
							</Form.Item>
						</Col>
					</Row>
				</Form>
			</Spin>
		</NewModal>
	);
};
export default MainModal;
