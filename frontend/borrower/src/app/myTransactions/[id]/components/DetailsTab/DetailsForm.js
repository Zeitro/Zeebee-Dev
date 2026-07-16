'use client';

/**
 * Details Tab 编辑态表单：Property Info / Transaction Info / Title/Escrow。
 * Property Info 字段均可编辑；物业数值字段初始值来自 rowData.borrower_property。
 */

import React, { useMemo } from 'react';
import { Checkbox, Form, Input, InputNumber, Select } from 'antd';
import { CommonFormDate, InputNumAmount, InputPhone } from '@repo/shared-methods';
import { useGetOptions } from '@repo/shared-methods/hooks';
import {
	changeUnitByPropertyType,
	emailRegex,
	livingAreaProps,
	phoneRegex,
	validateRule,
} from '@repo/shared-methods/utils';

/**
 * @param {object} props
 * @param {object|null} props.rowData 交易详情（物业字段初始值取自 rowData.borrower_property）
 * @param {React.ReactNode} [props.propertyInfoExtra] Property Info 标题右侧操作（Save），在 fieldset 外以免保存时被禁用
 * @param {boolean} [props.saving] 保存中禁用字段区
 */
const DetailsForm = ({ rowData, propertyInfoExtra, saving = false }) => {
	const form = Form.useFormInstance();
	const { propertyTypeOptions, businessStrategyOptions } = useGetOptions([
		{ key: 'property_type', valueName: 'propertyTypeOptions' },
		{ key: 'business_strategy', valueName: 'businessStrategyOptions' },
	]);

	const propertyType = Form.useWatch(['property', 'property_type'], form);
	const unit = Form.useWatch(['property', 'unit'], form);
	const unitInfo = useMemo(
		() =>
			changeUnitByPropertyType(
				propertyType ?? rowData?.property?.property_type,
				unit ?? rowData?.property?.unit,
			),
		[propertyType, rowData?.property?.property_type, unit, rowData?.property?.unit],
	);

	const sectionClassName = (isLast) =>
		['flex flex-col gap-4', isLast ? '' : 'border-b border-base-border pb-8']
			.filter(Boolean)
			.join(' ');

	return (
		<div className="flex w-full flex-col gap-8">
			{/* Property Info：全部可编辑 */}
			<section className={sectionClassName(false)}>
				<div className="flex items-center justify-between gap-4">
					<h2 className="m-0 text-xl font-medium leading-7 text-foreground">
						Property Info
					</h2>
					{propertyInfoExtra}
				</div>
				<fieldset
					disabled={saving}
					className="m-0 min-w-0 border-0 p-0 disabled:opacity-60"
				>
					<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						<Form.Item label="Property type" name={['property', 'property_type']}>
							<Select
								size="large"
								options={propertyTypeOptions}
								allowClear
								placeholder="Select property type"
								onChange={(value) => {
									const info = changeUnitByPropertyType(value);
									form.setFieldValue(['property', 'unit'], info.unit);
								}}
							/>
						</Form.Item>
						<Form.Item label="Business strategy" name="business_strategy">
							<Select
								size="large"
								options={businessStrategyOptions}
								allowClear
								mode="multiple"
								maxCount={5}
								style={{ width: '100%' }}
								placeholder="Select business strategy"
							/>
						</Form.Item>
						<Form.Item label="Title holder" name="title_holder">
							<Input size="large" allowClear placeholder="Enter title holder" />
						</Form.Item>
						<Form.Item label="Living area" name={['property', 'living_area_sqft']}>
							<InputNumber
								size="large"
								style={{ width: '100%' }}
								{...livingAreaProps}
								placeholder="Enter living area"
								suffix="sqft"
							/>
						</Form.Item>
						<Form.Item label="Units" name={['property', 'unit']}>
							<InputNumber
								size="large"
								style={{ width: '100%' }}
								min={unitInfo?.minUnit}
								max={unitInfo?.maxUnit}
								controls={false}
								placeholder="Enter units"
								disabled={unitInfo?.disabled}
							/>
						</Form.Item>
						<Form.Item label="Lot area" name={['property', 'lot_area_sqft']}>
							<InputNumber
								size="large"
								style={{ width: '100%' }}
								{...livingAreaProps}
								placeholder="Enter lot area"
								suffix="sqft"
							/>
						</Form.Item>
						<Form.Item label="Beds" name={['property', 'bedrooms']}>
							<InputNumber
								size="large"
								style={{ width: '100%' }}
								min={0}
								max={100}
								controls={false}
								placeholder="Enter beds"
							/>
						</Form.Item>
						<Form.Item label="Baths" name={['property', 'bathrooms']}>
							<InputNumber
								size="large"
								style={{ width: '100%' }}
								min={0}
								max={100}
								controls={false}
								placeholder="Enter baths"
							/>
						</Form.Item>
					</div>
				</fieldset>
			</section>

			{/* Transaction Info：全部可编辑 */}
			<section className={sectionClassName(false)}>
				<h2 className="m-0 text-xl font-medium leading-7 text-foreground">
					Transaction Info
				</h2>
				<fieldset
					disabled={saving}
					className="m-0 min-w-0 border-0 p-0 disabled:opacity-60"
				>
					<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2">
						<Form.Item label="Purchase price" name="purchase_amount">
							<InputNumAmount placeholder="$" size="large" />
						</Form.Item>
						<Form.Item label="EMD" name="emd_amount">
							<InputNumAmount placeholder="$" size="large" />
						</Form.Item>
					</div>
					<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
						<Form.Item label="Contract date" name="contract_date">
							<CommonFormDate size="large" />
						</Form.Item>
						<Form.Item label="EMD due date" name="emd_due_date">
							<CommonFormDate size="large" />
						</Form.Item>
						<Form.Item label="Contingency removal date" name="contingency_removal_date">
							<CommonFormDate size="large" />
						</Form.Item>
						<Form.Item label="COE date" name="coe_date">
							<CommonFormDate size="large" />
						</Form.Item>
					</div>
					<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2">
						<Form.Item label="PM" name="pm">
							<Input size="large" allowClear placeholder="Enter PM" />
						</Form.Item>
						<Form.Item label="TC" name="tc">
							<Input size="large" allowClear placeholder="Enter TC" />
						</Form.Item>
					</div>
					<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2">
						<Form.Item label="Who pays what" name="cost_allocation" className="mb-0">
							<Input.TextArea allowClear placeholder="Enter who pays what" rows={5} />
						</Form.Item>
						<Form.Item label="TC notes" name="note" className="mb-0">
							<Input.TextArea allowClear placeholder="Enter TC notes" rows={5} />
						</Form.Item>
					</div>
				</fieldset>
			</section>

			{/* Title/Escrow：escrows 为数组，Form.List 目前只展示第一项，增删按钮暂隐 */}
			<section className={sectionClassName(true)}>
				<h2 className="m-0 text-xl font-medium leading-7 text-foreground">Title/Escrow</h2>
				<fieldset
					disabled={saving}
					className="m-0 min-w-0 border-0 p-0 disabled:opacity-60"
				>
					<Form.List name="escrows">
						{(fields) =>
							fields.map((field) => {
								const id = form.getFieldValue(['escrows', field.name, 'id']);
								return (
									<div key={field.key} className="flex flex-col gap-4">
										<div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
											{id && (
												<Form.Item
													name={[field.name, 'id']}
													noStyle
												></Form.Item>
											)}
											<Form.Item
												label="Title company"
												name={[field.name, 'escrow_company']}
											>
												<Input
													size="large"
													allowClear
													placeholder="Enter title company"
												/>
											</Form.Item>
											<Form.Item
												label="Escrow officer"
												name={[field.name, 'officer_name']}
											>
												<Input
													size="large"
													allowClear
													placeholder="Enter escrow officer"
												/>
											</Form.Item>
											<Form.Item
												label="Escrow officer phone"
												name={[field.name, 'officer_phone']}
												rules={[
													{
														validator: (_, value) =>
															validateRule(
																value,
																phoneRegex,
																'Invalid phone number',
															),
													},
												]}
											>
												<InputPhone size="large" />
											</Form.Item>
											<Form.Item
												label="Escrow officer email"
												name={[field.name, 'officer_email']}
												rules={[
													{
														validator: (_, value) =>
															validateRule(
																value,
																emailRegex,
																'Invalid email format',
															),
													},
												]}
											>
												<Input
													size="large"
													allowClear
													placeholder="Enter email"
												/>
											</Form.Item>
											<Form.Item label="Escrow #" name={[field.name, 'no']}>
												<Input
													size="large"
													allowClear
													placeholder="Enter escrow #"
												/>
											</Form.Item>
										</div>
										<div className="grid w-full grid-cols-2 gap-4 sm:grid-cols-3 xl:grid-cols-5">
											<Form.Item
												name={[
													field.name,
													'checklist',
													'has_wire_instruction',
												]}
												valuePropName="checked"
												className="mb-0"
											>
												<Checkbox>Wire instruction</Checkbox>
											</Form.Item>
											<Form.Item
												name={[field.name, 'checklist', 'has_emd_receipt']}
												valuePropName="checked"
												className="mb-0"
											>
												<Checkbox>EMD receipt</Checkbox>
											</Form.Item>
											<Form.Item
												name={[
													field.name,
													'checklist',
													'has_prelim_report',
												]}
												valuePropName="checked"
												className="mb-0"
											>
												<Checkbox>Prelim report</Checkbox>
											</Form.Item>
											<Form.Item
												name={[field.name, 'checklist', 'has_nhd']}
												valuePropName="checked"
												className="mb-0"
											>
												<Checkbox>NHD</Checkbox>
											</Form.Item>
											<Form.Item
												name={[field.name, 'checklist', 'has_colored_map']}
												valuePropName="checked"
												className="mb-0"
											>
												<Checkbox>Colored map</Checkbox>
											</Form.Item>
										</div>
									</div>
								);
							})
						}
					</Form.List>
				</fieldset>
			</section>
		</div>
	);
};

export default DetailsForm;
