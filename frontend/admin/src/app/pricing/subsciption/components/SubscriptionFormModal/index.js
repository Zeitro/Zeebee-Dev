'use client';

import React, { useEffect, useState } from 'react';
import { Row, Col, Form, Input, Select, message } from 'antd';
import { CommonFormDate, NewModal } from '@repo/shared-methods';
import { layoutCol } from '@repo/shared-methods/enum';
import {
	planNameOptions,
	SUBSCRIPTION_STATUS,
	subscriptionStatusOptions,
	subscriptionTypeOptions,
	upsertMockSubscription,
} from '../../constants';

/**
 * Add / Edit subscription 弹窗表单。
 * Edit 模式仅可改 status、plan_name、subscription_type、expiry_date（与产品说明一致）。
 *
 * @param {object} props
 * @param {(open: boolean) => void} props.setOpen 控制弹窗显隐
 * @param {object} [props.rowData] 有 id 时为编辑
 * @param {() => void} [props.successCallback] 保存成功后刷新列表
 */
const SubscriptionFormModal = ({ setOpen, rowData = {}, successCallback }) => {
	const isEdit = Boolean(rowData?.id);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);

	useEffect(() => {
		if (!isEdit) return;
		form.setFieldsValue({
			subscription_no: rowData.subscription_no,
			status: rowData.status,
			customer_name: rowData.customer_name,
			plan_name: rowData.plan_name,
			subscription_type: rowData.subscription_type,
			start_date: rowData.start_date,
			expiry_date: rowData.expiry_date,
		});
	}, [form, isEdit, rowData]);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();

			// 接入后端后启用真实请求
			// let api = '/subscription';
			// let method = 'POST';
			// if (isEdit) {
			// 	api = `/subscription/${rowData.id}`;
			// 	method = 'PUT';
			// }
			// await request(api, { method, body: values });
			upsertMockSubscription(values, rowData?.id);

			message.success(`Subscription ${isEdit ? 'updated' : 'added'} successfully.`);
			successCallback?.();
			form.resetFields();
			setOpen(false);
		} catch {
			// 表单校验失败或请求异常
		} finally {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width={720}
			title={`${isEdit ? 'Edit' : 'Add'} Subscription`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				name="subscription-form"
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
				initialValues={{
					status: SUBSCRIPTION_STATUS.TRIALING,
				}}
			>
				<Row gutter={16}>
					<Col {...layoutCol}>
						<Form.Item
							label="Subscription no."
							name="subscription_no"
							rules={[{ required: true, message: 'Please enter subscription no.' }]}
						>
							<Input placeholder="Enter subscription no." disabled={isEdit} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Status"
							name="status"
							rules={[{ required: true, message: 'Please select status' }]}
						>
							<Select
								options={subscriptionStatusOptions}
								placeholder="Select status"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Customer name"
							name="customer_name"
							rules={[{ required: true, message: 'Please enter customer name' }]}
						>
							<Input placeholder="Enter customer name" disabled={isEdit} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Plan name"
							name="plan_name"
							rules={[{ required: true, message: 'Please select plan name' }]}
						>
							<Select options={planNameOptions} placeholder="Select plan name" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Subscription type"
							name="subscription_type"
							rules={[{ required: true, message: 'Please select subscription type' }]}
						>
							<Select
								options={subscriptionTypeOptions}
								placeholder="Select subscription type"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Start date"
							name="start_date"
							rules={[{ required: true, message: 'Please select start date' }]}
						>
							<CommonFormDate placeholder="Select start date" disabled={isEdit} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Expiry date"
							name="expiry_date"
							rules={[{ required: true, message: 'Please select expiry date' }]}
						>
							<CommonFormDate placeholder="Select expiry date" />
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};

export default SubscriptionFormModal;
