'use client';

import React, { useEffect, useState } from 'react';
import { Row, Col, Form, Input, Select, InputNumber, message } from 'antd';
import { InputNumAmount, NewModal } from '@repo/shared-methods';
import { layoutCol, titleCol } from '@repo/shared-methods/enum';
import { planStatusOptions, PRICING_PLAN_STATUS, upsertMockPricingPlan } from '../../constants';

const { TextArea } = Input;

/**
 * Add / Edit pricing plan 弹窗表单。
 * Edit 模式下 plan name 只读；monthly_price / yearly_price 为 0 时表示无需支付。
 *
 * @param {object} props
 * @param {(open: boolean) => void} props.setOpen 控制弹窗显隐
 * @param {object} [props.rowData] 有 id 时为编辑
 * @param {() => void} [props.successCallback] 保存成功后刷新列表
 */
const PlanFormModal = ({ setOpen, rowData = {}, successCallback }) => {
	const isEdit = Boolean(rowData?.id);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);

	useEffect(() => {
		if (!isEdit) return;
		form.setFieldsValue({
			name: rowData.name,
			status: rowData.status,
			subtitle: rowData.subtitle,
			trial_days: rowData.trial_days,
			monthly_price: rowData.monthly_price ?? undefined,
			yearly_price: rowData.yearly_price ?? undefined,
			description: rowData.description,
		});
	}, [form, isEdit, rowData]);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			const body = {
				...values,
				monthly_price: values.monthly_price ?? null,
				yearly_price: values.yearly_price ?? null,
			};

			// 接入后端后启用真实请求
			// let api = '/pricing_plan';
			// let method = 'POST';
			// if (isEdit) {
			// 	api = `/pricing_plan/${rowData.id}`;
			// 	method = 'PUT';
			// }
			// await request(api, { method, body });
			upsertMockPricingPlan(body, rowData?.id);

			message.success(`Plan ${isEdit ? 'updated' : 'added'} successfully.`);
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
			title={`${isEdit ? 'Edit' : 'Add'} Plan`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				name="pricing-plan-form"
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
				initialValues={{
					status: PRICING_PLAN_STATUS.ACTIVE,
				}}
			>
				<Row gutter={16}>
					<Col {...layoutCol}>
						<Form.Item
							label="Plan name"
							name="name"
							rules={[{ required: true, message: 'Please enter plan name' }]}
						>
							<Input placeholder="Enter plan name" disabled={isEdit} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Status"
							name="status"
							rules={[{ required: true, message: 'Please select status' }]}
						>
							<Select options={planStatusOptions} placeholder="Select status" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item label="Plan subtitle" name="subtitle">
							<Input placeholder="Enter plan subtitle" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Trial days"
							name="trial_days"
							rules={[{ required: true, message: 'Please enter trial days' }]}
						>
							<InputNumber
								className="!w-full"
								min={0}
								precision={0}
								placeholder="Enter trial days"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item label="Monthly price" name="monthly_price">
							<InputNumAmount placeholder="Enter monthly price" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item label="Yearly price" name="yearly_price">
							<InputNumAmount placeholder="Enter yearly price" />
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item label="Description" name="description">
							<TextArea rows={4} placeholder="Enter description" />
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};

export default PlanFormModal;
