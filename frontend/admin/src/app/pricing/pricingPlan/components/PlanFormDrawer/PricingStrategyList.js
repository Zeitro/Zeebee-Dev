'use client';

import React, { useEffect, useState } from 'react';
import { Button, Col, Form, Row, Select, Space } from 'antd';
import { InputNumAmount, NoData } from '@repo/shared-methods';

const strategyTypeOptions = [
	{ label: 'Monthly', value: 'Monthly' },
	{ label: 'Quarterly', value: 'Quarterly' },
	{ label: 'Yearly', value: 'Yearly' },
];

/** Save 按钮可用条件：type/actual_amount/original_amount 都有值（0 视为有效值）。 */
const canSavePolicyItem = (policy = {}) => {
	const hasType = Boolean(policy?.type);
	const hasActual = policy?.actual_amount !== undefined && policy?.actual_amount !== null;
	const hasOriginal = policy?.original_amount !== undefined && policy?.original_amount !== null;
	return hasType && hasActual && hasOriginal;
};

/** 统一金额比较：兼容 number/string，保留数值语义。 */
const normalizeAmount = (value) => {
	if (value === undefined || value === null || value === '') return null;
	const num = Number(value);
	return Number.isNaN(num) ? null : num;
};

/**
 * 编辑态判断当前行是否有变更：
 * - 有 id：按 id 与 existingPolicies 对比
 * - 无 id：视为新增行，只要有值即视为变更
 */
const isPolicyChanged = (policy = {}, existingPolicies = []) => {
	if (!policy?.id) return true;
	const original = existingPolicies.find((item) => item.id === policy.id);
	if (!original) return true;
	const isTypeChanged = (policy?.type || '') !== (original?.type || '');
	const isActualChanged =
		normalizeAmount(policy?.actual_amount) !== normalizeAmount(original?.actual_amount);
	const isOriginalChanged =
		normalizeAmount(policy?.original_amount) !== normalizeAmount(original?.original_amount);
	return isTypeChanged || isActualChanged || isOriginalChanged;
};

/**
 * 定价策略列表：使用 Form.List 维护 policies，每条策略支持单独保存。
 * 保存规则：
 * - 详情 `policies` 为空时：POST `/pricing_plan_policy`
 * - 详情 `policies` 非空时：优先 PUT `/pricing_plan_policy/{id}`，新增项走 POST
 *
 * @param {object} props
 * @param {import('antd').FormInstance} props.form
 * @param {(url: string, options?: object) => Promise<object>} props.request
 * @param {number|string} [props.pricingPlanId]
 * @param {object[]} props.existingPolicies
 * @param {boolean} [props.disabled]
 * @param {boolean} [props.isEdit]
 * @param {import('antd').MessageInstance} props.messageApi
 * @param {() => void} [props.successCallback]
 */
const PricingStrategyList = ({
	request,
	pricingPlanId,
	existingPolicies,
	disabled = false,
	isEdit = false,
	messageApi,
	successCallback,
}) => {
	const [form] = Form.useForm();
	const [savingMap, setSavingMap] = useState({});
	const [removingMap, setRemovingMap] = useState({});
	const policyValues = Form.useWatch('policies', form) || [];

	/** 仅校验并保存当前行，避免一次性校验全部策略项。 */
	const handleSaveItem = async (index) => {
		try {
			if (!pricingPlanId) {
				messageApi.warning('Please save plan first, then save pricing strategy.');
				return;
			}
			setSavingMap((prev) => ({ ...prev, [index]: true }));
			await form.validateFields([
				['policies', index, 'type'],
				['policies', index, 'actual_amount'],
				['policies', index, 'original_amount'],
			]);
			const policies = form.getFieldValue('policies') || [];
			const current = policies[index] || {};
			const shouldCreate = !existingPolicies.length || !current.id;

			if (shouldCreate) {
				const res = await request('/pricing_plan_policy', {
					method: 'POST',
					body: {
						actual_amount: current.actual_amount,
						original_amount: current.original_amount,
						pricing_plan_id: Number(pricingPlanId),
						type: current.type,
					},
				});
				const createdId = res?.id;
				if (createdId) {
					form.setFieldValue(['policies', index, 'id'], createdId);
				}
			} else {
				await request(`/pricing_plan_policy/${current.id}`, {
					method: 'PUT',
					body: {
						actual_amount: current.actual_amount,
						original_amount: current.original_amount,
						type: current.type,
					},
				});
			}

			messageApi.success('Pricing strategy saved.');
			successCallback?.();
		} catch {
			// 当前行校验失败或请求异常
		} finally {
			setSavingMap((prev) => ({ ...prev, [index]: false }));
		}
	};

	/**
	 * 删除策略项：
	 * - 详情返回过的项（有 id）先调后端删除
	 * - 未落库项（无 id）仅移除本地行
	 */
	const handleRemoveItem = async (index, fieldName, remove) => {
		try {
			setRemovingMap((prev) => ({ ...prev, [index]: true }));
			const policies = form.getFieldValue('policies') || [];
			const current = policies[index] || {};
			if (current?.id) {
				await request(`/pricing_plan_policy/${current.id}`, { method: 'DELETE' });
				messageApi.success('Pricing strategy deleted.');
			}
			remove(fieldName);
			successCallback?.();
		} catch {
			// 删除失败时保留当前行
		} finally {
			setRemovingMap((prev) => ({ ...prev, [index]: false }));
		}
	};

	useEffect(() => {
		form.setFieldsValue({
			policies: existingPolicies?.length
				? existingPolicies
				: isEdit
					? []
					: [
							{
								type: 'Monthly',
								actual_amount: undefined,
								original_amount: undefined,
							},
						],
		});
	}, [existingPolicies, form, isEdit]);

	return (
		<Form
			form={form}
			name="pricing-strategy-list"
			layout="vertical"
			autoComplete="off"
			size="large"
			disabled={disabled}
			scrollToFirstError
		>
			<Row>
				<Col span={24}>
					<Form.List name="policies">
						{(fields, { add, remove }) => (
							<div className="w-full">
								<div className="flex items-center justify-between mb-3">
									<div className="text-base-gray text-base font-medium">
										Pricing strategy
									</div>
									<Button
										type="dashed"
										onClick={() => add()}
										disabled={disabled}
										size="middle"
									>
										Add strategy
									</Button>
								</div>
								<Space direction="vertical" className="w-full" size={12}>
									{isEdit && !existingPolicies?.length && fields.length === 0 && (
										<NoData height={160} text="No pricing strategy" />
									)}
									{fields.map((field, index) => {
										const currentPolicy = policyValues?.[field.name] || {};
										const canSave =
											canSavePolicyItem(currentPolicy) &&
											isPolicyChanged(currentPolicy, existingPolicies);
										return (
											<div
												key={field.key}
												className="border border-base-border rounded-lg p-4"
											>
												<Row gutter={12}>
													<Col span={8}>
														<Form.Item
															label="Type"
															name={[field.name, 'type']}
															rules={[
																{
																	required: true,
																	message: 'Please select type',
																},
															]}
														>
															<Select
																options={strategyTypeOptions}
																placeholder="Select type"
															/>
														</Form.Item>
													</Col>
													<Col span={8}>
														<Form.Item
															label="Actual amount"
															name={[field.name, 'actual_amount']}
															rules={[
																{
																	required: true,
																	message:
																		'Please enter actual amount',
																},
															]}
														>
															<InputNumAmount
																min={0}
																precision={2}
																placeholder="Enter actual amount"
															/>
														</Form.Item>
													</Col>
													<Col span={8}>
														<Form.Item
															label="Original amount"
															name={[field.name, 'original_amount']}
															rules={[
																{
																	required: true,
																	message:
																		'Please enter original amount',
																},
															]}
														>
															<InputNumAmount
																min={0}
																precision={2}
																placeholder="Enter original amount"
															/>
														</Form.Item>
													</Col>
												</Row>
												<div className="flex justify-end gap-2">
													<Button
														onClick={() =>
															handleRemoveItem(
																index,
																field.name,
																remove,
															)
														}
														loading={Boolean(removingMap[index])}
														disabled={
															disabled ||
															Boolean(savingMap[index]) ||
															Boolean(removingMap[index])
														}
														size="middle"
													>
														Remove
													</Button>
													<Button
														type="primary"
														onClick={() => handleSaveItem(index)}
														loading={Boolean(savingMap[index])}
														disabled={
															disabled ||
															Boolean(savingMap[index]) ||
															Boolean(removingMap[index]) ||
															!canSave
														}
														size="middle"
													>
														Save
													</Button>
												</div>
											</div>
										);
									})}
								</Space>
							</div>
						)}
					</Form.List>
				</Col>
			</Row>
		</Form>
	);
};

export default PricingStrategyList;
