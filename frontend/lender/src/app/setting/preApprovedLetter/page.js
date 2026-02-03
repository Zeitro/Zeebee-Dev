'use client';
import { Button, Form, message, Spin } from 'antd';
import { useEffect, useState } from 'react';
import _ from 'lodash';
import { FieldTips, AssemblyFormItem } from '@repo/shared-methods';
import {
	createRequest,
	formatAmount,
	formatPhone,
	injectHeadResources,
	getBodyInnerHTML,
} from '@repo/shared-methods/utils';
const PreApprovedLetter = () => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [valueHtml, setValueHtml] = useState('');
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [data, setData] = useState({});
	const [loading, setLoading] = useState(false);

	const formatValue = (format, value) => {
		switch (format) {
			case 'Currency':
				return formatAmount(value);
			case 'Phone':
				return formatPhone(value);
			default:
				break;
		}
		return value;
	};

	// 更新HTML中的变量
	const updateVariablesInHtml = (html, value, idKey, format, isHave = true) => {
		// 创建一个临时DOM元素来处理HTML
		const tempDiv = document.createElement('div');
		tempDiv.innerHTML = html;

		// 查找所有带有 data-variable 属性的元素
		const variableElements = tempDiv.querySelectorAll('[data-variable]');

		variableElements.forEach((element) => {
			const variableName = element.getAttribute('data-variable');

			// 检查是否在表单数据中有对应的值
			if (idKey === variableName) {
				// 更新元素的文本内容
				element.textContent = formatValue(format, value);
			}
		});

		// 更新预览内容
		if (isHave) {
			setValueHtml(tempDiv.innerHTML);
		}
		return tempDiv.innerHTML;
	};

	const handlePreview = async (result) => {
		const { content } = await request('/organization_template/preview', {
			method: 'POST',
			body: {
				...result,
			},
		});
		injectHeadResources(content);
		let bodyHtml = getBodyInnerHTML(content);
		let newHtml = bodyHtml;

		const exampleVariablesObj = JSON.parse(result?.example_variables);

		result?.variable_group_schema?.map((item) => {
			item.variables.map((i) => {
				const vKey = i.key.replace('.', '');
				if (i.source !== 'Template') {
					const exampleValue = exampleVariablesObj[vKey] || `{${vKey}}`;
					let v = i.source === 'Fixed' ? i.fixed_value : exampleValue;
					newHtml = updateVariablesInHtml(newHtml, v, i.key, i.format, false);
				}
			});
		});

		setValueHtml(newHtml);
	};

	const getDetails = async () => {
		try {
			setLoading(true);
			const params = new URLSearchParams({ purpose: 'preApprovedLetter' });
			const res = await request(`/organization_template/detail_by_purpose?${params}`);
			setData(res);
			form.setFieldsValue(res);
			await handlePreview(res);
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};

	const handleReset = async () => {
		try {
			setLoading(true);
			const res = await request(`/organization_template/${data?.id}/reset`, {
				method: 'POST',
			});
			setData(res);
			form.setFieldsValue(res);
			await handlePreview(res);
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			await request(`/organization_template/${data?.id}`, {
				method: 'PUT',
				body: {
					...values,
				},
			});
			message.destroy();
			message.open({
				type: 'success',
				content: 'Update successful！',
			});
			setConfirmLoading(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	// 输入框修改更改html
	const handleFieldsChange = (changedValues, allValues) => {
		const formNameArr = _.cloneDeep(changedValues?.[0]?.name);
		let itemNameArr = formNameArr.slice(0, -1);
		if (formNameArr[formNameArr.length - 1] === 'fixed_value') {
			formNameArr[formNameArr.length - 1] = 'key';
		}
		const idKey = _.get(data, formNameArr);
		const item = _.get(data, itemNameArr);
		const changeItem = allValues.find((i) => _.isEqual(i.name, changedValues?.[0]?.name));
		updateVariablesInHtml(valueHtml, changeItem?.value, idKey, item.format);
	};

	useEffect(() => {
		const fetchDetails = async () => {
			await getDetails();
		};
		fetchDetails();
	}, []);

	return (
		<div className="flex-1 flex overflow-hidden bg-[#f5f5f5] mt-[2px]">
			<div className="border-r border-base-border w-[450px] flex-shrink-0 h-full bg-background flex flex-col relative px-8 pt-6">
				<h2 className="border-b-2 border-base-border text-xl font-medium pb-4 mb-4 flex items-center justify-between flex-shrink-0 text-black-dark">
					<span>
						Configuration
						<FieldTips
							title={
								<>
									For example data, please set your own pre-approved letter, which
									will be used for the pre-approved loan submitted by the
									borrower. Zeebee will automatically send this attachment to the
									borrower.
									<br />
									<br />
									Note: The "Default" option is calculated by the system, while
									"Customize" is set by yourself.
								</>
							}
							isShowDescription={true}
							rootClassName="!max-w-[300px]"
							className="ml-1"
						/>
					</span>
					<FieldTips title={'Reset data and get the latest configuration template.'}>
						<Button size="small" onClick={handleReset}>
							Reset
						</Button>
					</FieldTips>
				</h2>
				<Form
					form={form}
					name="preApproved-form"
					layout="vertical"
					autoComplete="off"
					className={[
						'flex-1 overflow-hidden overflow-y-auto [&_.ant-form-item]:mb-4 [&_.ant-form-item-explain-error]:leading-4',
					].join(' ')}
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
					onFieldsChange={handleFieldsChange}
				>
					<Spin spinning={loading}>
						<Form.List name="variable_group_schema">
							{(fields) =>
								fields.map(({ key, name, ...restField }, index) => {
									const item = data?.variable_group_schema?.[index];
									return (
										<div key={key}>
											<div className="border-b border-base-border font-medium py-1 mb-4 text-base/[22px] text-black-dark">
												{item?.name}
											</div>
											<Form.List
												name={[name, 'variables']}
												{...restField}
												rules={[{ required: true }]}
											>
												{(fields2) =>
													fields2.map(
														(
															{
																key: variableItemKey,
																name: variableItemName,
																...variableItemRestField
															},
															i,
														) => {
															const variableItem =
																item?.variables?.[i];
															return (
																<AssemblyFormItem
																	key={variableItemKey}
																	item={variableItem}
																	formListRestField={
																		variableItemRestField
																	}
																	formListName={[
																		variableItemName,
																	]}
																	allFormListName={[
																		'variable_group_schema',
																		name,
																		'variables',
																		variableItemName,
																	]}
																	example_variables={
																		data?.example_variables
																	}
																	changeHtml={(
																		key,
																		value,
																		format,
																	) => {
																		updateVariablesInHtml(
																			valueHtml,
																			value,
																			key,
																			format,
																		);
																	}}
																/>
															);
														},
													)
												}
											</Form.List>
										</div>
									);
								})
							}
						</Form.List>
					</Spin>
				</Form>
				<div className="w-full flex-shrink-0 bg-background pt-2 pb-6">
					<Button
						type="primary"
						className="w-full"
						loading={confirmLoading}
						onClick={handleSave}
						size="large"
					>
						Save
					</Button>
				</div>
			</div>

			<div
				className="flex-1 overflow-hidden overflow-y-auto py-6 flex justify-center items-center [&>.page]:max-w-[800px]"
				dangerouslySetInnerHTML={{ __html: valueHtml }}
			></div>
		</div>
	);
};
export default PreApprovedLetter;
