'use client';

import React, { useCallback, useMemo, useState } from 'react';
import { Row, Col, Form, Select, message } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { layoutCol } from '@repo/shared-methods/enum';
import { createRequest, formatAmount } from '@repo/shared-methods/utils';

/**
 * 将单条 policy 转为可搜索的字符串 label。
 *
 * @param {{ type?: string, actual_amount?: number, original_amount?: number }} policy
 * @returns {string}
 */
const buildPolicySearchLabel = (policy = {}) => {
	const hasAmount =
		policy?.actual_amount !== undefined &&
		policy?.actual_amount !== null &&
		policy?.original_amount !== undefined &&
		policy?.original_amount !== null;
	const amountLabel = hasAmount
		? ` (${formatAmount(policy.actual_amount)} / ${formatAmount(policy.original_amount)})`
		: '';
	return `${policy.type || ''}${amountLabel}`.trim();
};

/**
 * 将 plan policies 转为下拉选项；label 为字符串以支持 Select 搜索。
 *
 * @param {Array<{ id: number, type?: string, actual_amount?: number, original_amount?: number }>} policies
 * @returns {Array<{ label: string, value: number, policy: object }>}
 */
const buildPolicyOptions = (policies = []) =>
	policies.map((policy) => ({
		label: buildPolicySearchLabel(policy),
		value: policy.id,
		policy,
	}));

/**
 * 新增组织订阅弹窗；字段与 POST `/organization_subscription` body 一致。
 *
 * @param {object} props
 * @param {boolean} props.open 弹窗显隐
 * @param {(open: boolean) => void} props.setOpen 控制弹窗显隐
 * @param {object} [props.rowData] 至少含 type（Borrower / Lender / Admin）
 * @param {() => void} [props.successCallback] 保存成功后刷新列表
 */
const SubscriptionFormModal = ({ open, setOpen, rowData = {}, successCallback }) => {
	const systemType = rowData?.type || 'Borrower';
	const [form] = Form.useForm();
	const pricingPlanId = Form.useWatch('pricing_plan_id', form);
	const { request } = useMemo(() => createRequest(), []);
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [organizationLoading, setOrganizationLoading] = useState(false);
	const [pricingPlanLoading, setPricingPlanLoading] = useState(false);
	const [organizationOptions, setOrganizationOptions] = useState([]);
	const [pricingPlanList, setPricingPlanList] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();

	const pricingPlanOptions = useMemo(
		() =>
			pricingPlanList.map((item) => ({
				label: item.name,
				value: item.id,
			})),
		[pricingPlanList],
	);

	const organizationSelectOptions = useMemo(
		() =>
			organizationOptions.map((item) => ({
				label: item.name,
				value: item.id,
			})),
		[organizationOptions],
	);
	const policyOptions = useMemo(() => {
		const pricingPlan = pricingPlanList.find((item) => item.id === pricingPlanId);
		return buildPolicyOptions(pricingPlan?.policies || []);
	}, [pricingPlanList, pricingPlanId]);

	const loadOrganizationOptions = useCallback(async () => {
		try {
			setOrganizationLoading(true);
			const params = new URLSearchParams({ type: systemType });
			const { list } = await request(`/organization/select_options?${params}`);
			setOrganizationOptions(list || []);
		} catch {
			setOrganizationOptions([]);
		} finally {
			setOrganizationLoading(false);
		}
	}, [request, systemType]);

	const loadPricingPlanOptions = useCallback(async () => {
		try {
			setPricingPlanLoading(true);
			const params = new URLSearchParams({ type: systemType });
			const { list } = await request(`/pricing_plan/all?${params}`);
			setPricingPlanList(list || []);
		} catch {
			setPricingPlanList([]);
		} finally {
			setPricingPlanLoading(false);
		}
	}, [request, systemType]);

	useDeferredSyncEffect({
		deferred: () => {
			loadOrganizationOptions();
			loadPricingPlanOptions();
		},
		deps: [loadOrganizationOptions, loadPricingPlanOptions],
		delay: 0,
	});

	const handleClose = () => {
		form.resetFields();
		setOpen(false);
	};

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			await request('/organization_subscription', {
				method: 'POST',
				body: {
					organization_id: values.organization_id,
					pricing_plan_id: values.pricing_plan_id,
					pricing_plan_policy_id: values.pricing_plan_policy_id,
					type: systemType,
				},
			});
			messageApi.success('Subscription added successfully.');
			successCallback?.();
			handleClose();
		} catch {
			// 表单校验失败或请求异常
		} finally {
			setConfirmLoading(false);
		}
	};

	const filterSelectOption = (input, option) => {
		const label = option?.label ?? option?.data?.label ?? '';
		return String(label).toLowerCase().includes(input.toLowerCase());
	};

	const renderPolicyOption = (option) => {
		const policy = option?.data?.policy;
		if (!policy) return option.label;
		return (
			<div className="flex flex-wrap items-center gap-x-1">
				<span>{policy.type}</span>
				{policy.actual_amount !== undefined && policy.actual_amount !== null ? (
					<span className="text-xs text-gray-500">
						(actual: {formatAmount(policy.actual_amount)} / original:{' '}
						{formatAmount(policy.original_amount)})
					</span>
				) : (
					<span className="text-xs text-gray-500">
						{formatAmount(policy.original_amount)}
					</span>
				)}
			</div>
		);
	};

	return (
		<>
			{contextHolder}
			<NewModal
				open={open}
				width={720}
				title="Add Subscription"
				okText="Save"
				confirmLoading={confirmLoading}
				onOk={handleSave}
				onCancel={handleClose}
			>
				<Form
					form={form}
					name="subscription-form"
					layout="vertical"
					autoComplete="off"
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
				>
					<Row gutter={16}>
						<Col {...layoutCol}>
							<Form.Item
								label="Organization"
								name="organization_id"
								rules={[{ required: true, message: 'Please select organization' }]}
							>
								<Select
									showSearch
									allowClear
									placeholder="Select organization"
									options={organizationSelectOptions}
									loading={organizationLoading}
									filterOption={filterSelectOption}
									optionFilterProp="label"
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Pricing plan"
								name="pricing_plan_id"
								rules={[{ required: true, message: 'Please select pricing plan' }]}
							>
								<Select
									showSearch
									allowClear
									placeholder="Select pricing plan"
									options={pricingPlanOptions}
									loading={pricingPlanLoading}
									filterOption={filterSelectOption}
									optionFilterProp="label"
									onChange={() => {
										form.setFieldValue('pricing_plan_policy_id', undefined);
									}}
								/>
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Pricing plan policy"
								name="pricing_plan_policy_id"
								rules={[
									{
										required: true,
										message: 'Please select pricing plan policy',
									},
								]}
							>
								<Select
									showSearch
									allowClear
									placeholder="Select pricing plan policy"
									options={policyOptions}
									disabled={!pricingPlanId}
									filterOption={filterSelectOption}
									optionFilterProp="label"
									optionRender={renderPolicyOption}
								/>
							</Form.Item>
						</Col>
					</Row>
				</Form>
			</NewModal>
		</>
	);
};

export default SubscriptionFormModal;
