'use client';

import React, { useCallback, useMemo, useState } from 'react';
import {
	Row,
	Col,
	Form,
	Input,
	Select,
	InputNumber,
	Radio,
	message,
	Spin,
	Drawer,
	Button,
	Divider,
} from 'antd';
import { createRequest } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { layoutCol, systemTypeOptions, titleCol, yesNoOptions } from '@repo/shared-methods/enum';
import { planTypeOptions, PRICING_PLAN_TYPE } from '../../constants';
import PricingStrategyList from './PricingStrategyList.js';

const { TextArea } = Input;

/**
 * Add / Edit pricing plan 弹窗表单。
 * 字段与 POST/PUT `/pricing_plan` body 一致；编辑时 GET `/pricing_plan/{id}` 初始化表单。
 *
 * @param {object} props
 * @param {(open: boolean) => void} props.setOpen 控制弹窗显隐
 * @param {object} [props.rowData] 有 id 时为编辑；新增时至少含 type（Borrower / Lender / Admin）
 * @param {() => void} [props.successCallback] 保存成功后刷新列表
 */
const PlanFormDrawer = ({ setOpen, rowData = {}, successCallback }) => {
	const isEdit = Boolean(rowData?.id);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [existingPolicies, setExistingPolicies] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const handleClose = () => setOpen(false);

	/** 编辑时拉取详情，字段与 PUT body 一致 */
	const getDetails = useCallback(async () => {
		if (!rowData?.id) return;
		try {
			setLoading(true);
			const res = await request(`/pricing_plan/${rowData.id}`);
			form.setFieldsValue(res);
			setExistingPolicies(res?.policies || []);
		} catch {
			// 请求失败时保留空表单，由用户关闭弹窗
		} finally {
			setLoading(false);
		}
	}, [form, request, rowData]);

	useDeferredSyncEffect({
		predicate: () => isEdit,
		deferred: getDetails,
		deps: [isEdit, getDetails],
		delay: 0,
	});

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			const restValues = { ...values };
			delete restValues.policies;
			const body = {
				...restValues,
				type: rowData.type,
			};

			let api = '/pricing_plan';
			let method = 'POST';
			if (isEdit) {
				api = `/pricing_plan/${rowData.id}`;
				method = 'PUT';
			}
			await request(api, { method, body });

			messageApi.success(`Plan ${isEdit ? 'updated' : 'added'} successfully.`);
			successCallback?.();
		} catch {
			// 表单校验失败或请求异常
		} finally {
			setConfirmLoading(false);
		}
	};

	return (
		<Drawer
			width={760}
			title={`${isEdit ? 'Edit' : 'Add'} Plan`}
			open
			onClose={handleClose}
			destroyOnHidden
		>
			{contextHolder}
			<Spin spinning={loading}>
				<Form
					form={form}
					name="pricing-plan-form"
					layout="vertical"
					autoComplete="off"
					size="large"
					disabled={confirmLoading || loading}
					scrollToFirstError
					initialValues={
						isEdit
							? undefined
							: {
									plan_type: PRICING_PLAN_TYPE.STANDARD,
									is_default: false,
									is_public: true,
									trial_days: 0,
									sort_order: 1,
									type: rowData?.type,
								}
					}
				>
					<div className="text-base-gray text-base font-medium mb-4">
						Basic information
					</div>
					<Row gutter={16}>
						<Col {...layoutCol}>
							<Form.Item
								label="System type"
								name={['type']}
								rules={[{ required: true }]}
							>
								<Select
									options={systemTypeOptions}
									placeholder="Select system type"
									disabled={isEdit}
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Plan type"
								name="plan_type"
								rules={[{ required: true, message: 'Please select plan type' }]}
							>
								<Select options={planTypeOptions} placeholder="Select plan type" />
							</Form.Item>
						</Col>
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
							<Form.Item label="Plan subtitle" name="subtitle">
								<Input placeholder="Enter plan subtitle" />
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item label="Trial days" name="trial_days">
								<InputNumber
									className="!w-full"
									min={0}
									precision={0}
									placeholder="Enter trial days"
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Sort order"
								name="sort_order"
								rules={[{ required: true, message: 'Please enter sort order' }]}
							>
								<InputNumber
									className="!w-full"
									min={0}
									precision={0}
									placeholder="Enter sort order"
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Default plan"
								name="is_default"
								rules={[{ required: true, message: 'Please select default plan' }]}
								extra="Used for the default plan displayed after user registration. If Yes is selected here, the rest of the plan's plan will automatically be changed to No。"
							>
								<Radio.Group
									options={yesNoOptions}
									block
									className="[&_.ant-radio-wrapper]:!justify-start"
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Public plan"
								name="is_public"
								rules={[{ required: true, message: 'Please select public plan' }]}
								extra="Whether permission details are shown to users, such as comparing all permissions across different plans. Generally, only custom models are not shown to users"
							>
								<Radio.Group
									options={yesNoOptions}
									block
									className="[&_.ant-radio-wrapper]:!justify-start"
								/>
							</Form.Item>
						</Col>
						<Col {...titleCol}>
							<Form.Item label="Description" name="description">
								<TextArea rows={4} placeholder="Enter description" />
							</Form.Item>
						</Col>
						<Button type="primary" loading={confirmLoading} onClick={handleSave}>
							Save basic information
						</Button>
						<Divider />
					</Row>
				</Form>
				<PricingStrategyList
					request={request}
					pricingPlanId={rowData?.id}
					existingPolicies={existingPolicies}
					isEdit={isEdit}
					messageApi={messageApi}
					successCallback={successCallback}
				/>
			</Spin>
		</Drawer>
	);
};

export default PlanFormDrawer;
