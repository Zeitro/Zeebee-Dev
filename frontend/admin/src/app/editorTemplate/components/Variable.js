'use client';
import { NewModal } from '@repo/shared-methods';
import { Button, Form, Input, Select, Space } from 'antd';
import { useState } from 'react';
import { MinusCircleOutlined, PlusCircleOutlined } from '@ant-design/icons';

// 增加变量
export const VariableModal = ({ setOpen, onChange, value }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			onChange(values?.variable_group_schema || []);
			setConfirmLoading(false);
			setOpen(false);
		} catch (errpr) {
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
		type: 'String',
	};
	const initValue = {
		name: '',
		variables: [initVariable],
	};
	return (
		<NewModal
			width={900}
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
												) => (
													<Space
														key={key2}
														style={
															{
																// display: 'flex',
															}
														}
														// align="baseline"
													>
														<Form.Item
															{...restField2}
															label="ID"
															name={[name2, 'key']}
															rules={[{ required: true }]}
														>
															<Input placeholder="ID" />
														</Form.Item>
														<Form.Item
															{...restField2}
															label="Label"
															name={[name2, 'label']}
															rules={[{ required: true }]}
														>
															<Input placeholder="Label" />
														</Form.Item>
														<Form.Item
															{...restField2}
															label="Value"
															name={[name2, 'fixed_value']}
															rules={[{ required: true }]}
														>
															<Input placeholder="Value" />
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
																options={[
																	{
																		label: 'true',
																		value: true,
																	},
																	{
																		label: 'false',
																		value: false,
																	},
																]}
																style={{ width: 90 }}
															/>
														</Form.Item>
														{/* <Form.Item
															{...restField2}
															name={[name2, 'type']}
															rules={[{ required: true }]}
														>
															<Select
																placeholder="Type"
																allowClear
																options={[
																	{
																		label: 'String',
																		value: 'String',
																	},
																	{
																		label: 'Object',
																		value: 'Object',
																	},
																	{
																		label: 'Array',
																		value: 'Array',
																	},
																]}
																style={{ width: 150 }}
															/>
														</Form.Item> */}
														{fields2?.length === index2 + 1 && (
															<PlusCircleOutlined
																onClick={() =>
																	actions?.add(initVariable)
																}
															/>
														)}
														{fields2?.length > 0 && (
															<MinusCircleOutlined
																onClick={() =>
																	actions?.remove(name2)
																}
															/>
														)}
													</Space>
												),
											);
										}}
									</Form.List>
								</div>
							))}
						</>
					)}
				</Form.List>
			</Form>
		</NewModal>
	);
};

export const VariableBtn = ({ value, onChange }) => {
	const [open, setOpen] = useState(false);
	return (
		<>
			<Button
				icon={<PlusCircleOutlined />}
				type="primary"
				ghost
				onClick={() => setOpen(true)}
			>
				Add Variable ({value?.length || 0})
			</Button>
			{open && <VariableModal value={value} onChange={onChange} setOpen={setOpen} />}
		</>
	);
};
