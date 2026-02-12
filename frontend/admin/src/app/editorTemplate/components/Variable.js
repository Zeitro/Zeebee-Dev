'use client';
import { NewModal } from '@repo/shared-methods';
import { Button, Form, Input, Select, Space, Spin } from 'antd';
import { useEffect, useMemo, useState } from 'react';
import { MinusCircleOutlined, PlusCircleOutlined } from '@ant-design/icons';
import _ from 'lodash';
const { TextArea } = Input;
// 增加变量
export const VariableModal = ({ setOpen, onChange, value }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(true);
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			onChange(values?.variable_group_schema || []);
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};
	const initVariable = {
		key: '',
		label: '',
		fixed_value: '',
		source: 'Dynamic',
		format: 'Input',
		is_editable: true,
		is_required: false,
		type: 'String',
	};
	const initValue = {
		name: '',
		variables: [initVariable],
	};
	const options = [
		{
			label: 'true',
			value: true,
		},
		{
			label: 'false',
			value: false,
		},
	];

	useEffect(() => {
		// 因为数据太多,导致弹框点击了打开很卡顿，所以加个定时器增加loading效果缓解
		let loadingTime = setTimeout(() => {
			setLoading(false);
		}, 100);
		return () => {
			clearTimeout(loadingTime);
		};
	}, []);
	return (
		<NewModal
			width={'80vw'}
			title="Add variable"
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				name="addVariableForm"
				size="large"
				layout="vertical"
				autoComplete="off"
				form={form}
				initialValues={{
					variable_group_schema: value?.length > 0 ? value : [initValue],
				}}
				className="h-[70vh] overflow-hidden overflow-y-auto"
			>
				{loading ? (
					<div className="flex justify-center items-center h-full">
						<Spin spinning={loading} />
					</div>
				) : (
					<Form.List name="variable_group_schema">
						{(fields, { add, remove }) => (
							<>
								{fields.map(({ key, name, ...restField }, index) => (
									<div
										key={key}
										className="border border-base-border rounded-md pt-4 px-2 bg-gray-light mb-4"
									>
										<div className="flex items-start justify-between">
											<Form.Item
												{...restField}
												label="Title"
												name={[name, 'name']}
												rules={[{ required: true, message: 'Enter name' }]}
											>
												<Input placeholder="name" />
											</Form.Item>
											<Space>
												{fields?.length === index + 1 && (
													<Button
														type="link"
														onClick={() => {
															add(initValue);
														}}
														size="small"
													>
														Add
													</Button>
												)}
												{fields?.length > 0 && (
													<Button
														type="link"
														danger
														onClick={() => {
															remove(name);
														}}
														size="small"
													>
														Delete
													</Button>
												)}
											</Space>
										</div>
										<Form.List
											name={[name, 'variables']}
											{...restField}
											rules={[{ required: true }]}
										>
											{(fields2, actions) => {
												return fields2.map(
													(
														{ name: name2, key: key2, restField2 },
														index2,
													) => {
														const fieldName = [
															'variable_group_schema',
															name,
															'variables',
															name2,
														];
														return (
															<div
																key={key2}
																className="flex items-center gap-2"
															>
																<Form.Item
																	{...restField2}
																	label="ID"
																	name={[name2, 'key']}
																	rules={[{ required: true }]}
																	className="flex-1"
																>
																	<Input placeholder="ID" />
																</Form.Item>
																<Form.Item
																	{...restField2}
																	label="Label"
																	name={[name2, 'label']}
																	rules={[{ required: true }]}
																	className="flex-1"
																>
																	<Input placeholder="Label" />
																</Form.Item>
																<Form.Item
																	shouldUpdate={(
																		prevValues,
																		curValues,
																	) => {
																		const prevV = _.get(
																			prevValues,
																			[
																				...fieldName,
																				'format',
																			],
																		);
																		const curV = _.get(
																			curValues,
																			[
																				...fieldName,
																				'format',
																			],
																		);

																		return prevV !== curV;
																	}}
																	noStyle
																>
																	{({ getFieldValue }) => {
																		const currentFormat =
																			getFieldValue([
																				...fieldName,
																				'format',
																			]);
																		return (
																			<Form.Item
																				{...restField2}
																				label="Value"
																				name={[
																					name2,
																					'fixed_value',
																				]}
																				className="flex-1"
																			>
																				{currentFormat ===
																				'TextArea' ? (
																					<TextArea
																						placeholder="Value"
																						autoSize={{
																							minRows: 1,
																							maxRows: 1,
																						}}
																					/>
																				) : (
																					<Input placeholder="Value" />
																				)}
																			</Form.Item>
																		);
																	}}
																</Form.Item>
																<Form.Item
																	{...restField2}
																	label="Source"
																	name={[name2, 'source']}
																	rules={[{ required: true }]}
																>
																	<Select
																		placeholder="Source"
																		allowClear
																		options={[
																			{
																				label: 'Fixed',
																				value: 'Fixed',
																			},
																			{
																				label: 'Dynamic',
																				value: 'Dynamic',
																			},
																			{
																				label: 'Template',
																				value: 'Template',
																			},
																		]}
																		style={{ width: 120 }}
																	/>
																</Form.Item>
																<Form.Item
																	{...restField2}
																	label="Format"
																	name={[name2, 'format']}
																	rules={[{ required: true }]}
																>
																	<Select
																		placeholder="Format"
																		allowClear
																		options={[
																			{
																				label: 'Phone',
																				value: 'Phone',
																			},
																			{
																				label: 'Money',
																				value: 'Currency',
																			},
																			{
																				label: 'Email',
																				value: 'Email',
																			},
																			{
																				label: 'TextArea',
																				value: 'TextArea',
																			},
																			{
																				label: 'Input',
																				value: 'Input',
																			},
																		]}
																		style={{ width: 120 }}
																	/>
																</Form.Item>
																<Form.Item
																	{...restField2}
																	label="Is editable"
																	name={[name2, 'is_editable']}
																	rules={[{ required: true }]}
																>
																	<Select
																		placeholder="Is editable"
																		allowClear
																		options={options}
																		style={{ width: 90 }}
																	/>
																</Form.Item>
																<Form.Item
																	{...restField2}
																	label="Is required"
																	name={[name2, 'is_required']}
																	rules={[{ required: true }]}
																>
																	<Select
																		placeholder="Is required"
																		allowClear
																		options={options}
																		style={{ width: 90 }}
																	/>
																</Form.Item>
																<div className="w-[40px] flex gap-2">
																	{fields2?.length ===
																		index2 + 1 && (
																		<PlusCircleOutlined
																			onClick={() =>
																				actions?.add(
																					initVariable,
																				)
																			}
																		/>
																	)}
																	{fields2?.length > 0 && (
																		<MinusCircleOutlined
																			onClick={() =>
																				actions?.remove(
																					name2,
																				)
																			}
																		/>
																	)}
																</div>
															</div>
														);
													},
												);
											}}
										</Form.List>
									</div>
								))}
							</>
						)}
					</Form.List>
				)}
			</Form>
		</NewModal>
	);
};

export const VariableBtn = ({ value, onChange }) => {
	const [open, setOpen] = useState(false);
	const variableLength = useMemo(() => {
		const list = [];
		value?.map?.((i) => {
			list.push(...i.variables);
		});
		return list.length;
	}, [value]);
	return (
		<>
			<Button
				icon={<PlusCircleOutlined />}
				type="primary"
				ghost
				onClick={() => setOpen(true)}
			>
				Add Variable ({variableLength || 0})
			</Button>
			{open && <VariableModal value={value} onChange={onChange} setOpen={setOpen} />}
		</>
	);
};
