'use client';
import { Button, Form, message, Spin } from 'antd';
import { useEffect, useState } from 'react';
import _ from 'lodash';
import { FieldTips, AssemblyFormItem } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
const PreApprovedLetter = () => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [valueHtml, setValueHtml] = useState('');
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [data, setData] = useState({});
	const [loading, setLoading] = useState(false);

	// 工具函数：插入head资源
	const injectHeadResources = (html) => {
		// 清理旧的
		document.querySelectorAll('[data-template-preview]').forEach((el) => el.remove());
		try {
			const doc = new window.DOMParser().parseFromString(html, 'text/html');
			const headNodes = doc.head ? Array.from(doc.head.children) : [];
			headNodes.forEach((node) => {
				if (['LINK', 'STYLE'].includes(node.tagName) || node.tagName.startsWith('FONT')) {
					const clone = node.cloneNode(true);
					clone.setAttribute('data-template-preview', '1');
					document.head.appendChild(clone);
				}
			});
		} catch {
			// dd
		}
	};

	// 工具函数：取body内容
	const getBodyInnerHTML = (html) => {
		try {
			const doc = new window.DOMParser().parseFromString(html, 'text/html');
			return doc.body.innerHTML;
		} catch {
			return html;
		}
	};

	const handlePreview = async (result) => {
		const { content } = await request('/organization_template/preview', {
			method: 'POST',
			body: {
				...result,
			},
		});
		injectHeadResources(content);
		const newHtml = getBodyInnerHTML(content);
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

	// 更新HTML中的变量
	const updateVariablesInHtml = (html, value, idKey) => {
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
				element.textContent = value;
			}
		});

		// 更新预览内容
		setValueHtml(tempDiv.innerHTML);
	};

	const handleFieldsChange = (changedValues, allValues) => {
		const formNameArr = changedValues?.[0]?.name;
		if (formNameArr[formNameArr.length - 1] === 'fixed_value') {
			formNameArr[formNameArr.length - 1] = 'key';
		}
		const idKey = _.get(data, formNameArr);
		const changeItem = allValues.find((i) => _.isEqual(i.name, changedValues?.[0]?.name));
		updateVariablesInHtml(valueHtml, changeItem?.value, idKey);
	};

	useEffect(() => {
		const fetchDetails = async () => {
			await getDetails();
		};
		fetchDetails();
	}, []);

	return (
		<div className="flex-1 flex overflow-hidden bg-[#f5f5f5] mt-[2px]">
			<div className="border-r border-base-border w-[300px] flex-shrink-0 h-full bg-background flex flex-col relative">
				<h2 className="border-b border-base-border text-xl font-medium py-2 px-4 flex items-center justify-between flex-shrink-0">
					<span>Configuration</span>
					<FieldTips title={'Reset data and get the latest configuration.'}>
						<Button type="primary" size="small" onClick={handleReset} ghost>
							Reset
						</Button>
					</FieldTips>
				</h2>
				<Form
					form={form}
					name="preApproved-form"
					layout="vertical"
					autoComplete="off"
					className={['flex-1 overflow-hidden overflow-y-auto !px-4'].join(' ')}
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
											<div className="border-b border-base-border font-medium pt-3 mb-2">
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
				<div className="w-full flex-shrink-0 bg-background py-2 px-4">
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
