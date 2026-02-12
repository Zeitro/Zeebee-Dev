'use client';
import { InputNumAmount, InputPhone } from '@repo/shared-methods';
import { emailRegex, phoneRegex, validateRule } from '@repo/shared-methods/utils';
import { DatePicker, Form, Input, Radio } from 'antd';
import _ from 'lodash';
import { useEffect } from 'react';

/**
 *  获取输入框/选择框/日期选择框/金额输入框/手机输入框/email输入/动态选择(可自定义/系统指定值)
 * @param {*} item
 * @param {*} isHaveFormItem 是否有Form.Item包裹
 * @param {*} formListName 表单名称(当用于Form.List时，需要传入)
 * @param {*} formListRestField 当用于Form.List时，需要传入
 * @param {*} allFormListName 当用于多层嵌套Form.List时,动态选择组件获取source值时需要传入完整的表单名称,此时formListName传入的不是完整的表单名
 * @returns
 */
export const AssemblyFormItem = ({
	item = {},
	isHaveFormItem = true,
	formListName = [],
	formListRestField = {},
	allFormListName = [],
	example_variables = '{}', //示例数据用于source==='Fixed'/'Dynamic'
	changeHtml,
	isAllEditable, // 全部可编辑
}) => {
	const form = Form.useFormInstance();
	const exampleVariablesObj = JSON.parse(example_variables);

	useEffect(() => {
		const vKey = item.key.replace('.', '');
		const exampleValue = exampleVariablesObj[vKey] || `{${vKey}}`;

		if (item.source === 'Dynamic') {
			form.setFieldValue(
				[...(allFormListName?.length ? allFormListName : formListName), 'dynamicValue'],
				exampleValue,
			);
			// changeHtml(item.key, exampleValue, item.format);
		} else {
			form.setFieldValue(
				[...(allFormListName?.length ? allFormListName : formListName), 'fixed_value'],
				item.fixed_value,
			);
			// changeHtml(item.key, item.fixed_value, item.format);
		}
	}, [item]);
	const getFormItem = (children, options = {}) => {
		return (
			<Form.Item
				{...formListRestField}
				label={item.label}
				name={[...formListName, 'fixed_value']}
				rules={[{ required: item.is_required }]}
				{...options}
			>
				{children}
			</Form.Item>
		);
	};
	const getCurrentElement = (options = {}, formItemOptions = {}, isHave = isHaveFormItem) => {
		let defaultProps = {
			placeholder: item.label,
			disabled: isAllEditable ? false : !item.is_editable,
			...options,
		};
		let ele = <Input allowClear {...defaultProps} />;
		if (item.format === 'Input') {
			// ele = <Input {...defaultProps} />;
		}
		if (item.format === 'TextArea') {
			ele = (
				<Input.TextArea
					autoSize={{ minRows: 5, maxRows: 5 }}
					allowClear
					{...defaultProps}
				/>
			);
		}
		if (item.format === 'Date') {
			ele = <DatePicker allowClear {...defaultProps} />;
		}
		if (item.format === 'Phone') {
			ele = <InputPhone {...defaultProps} />;
			return isHave
				? getFormItem(ele, {
						rules: [
							{ required: item.is_required },
							{
								validator: (_, value) =>
									validateRule(value, phoneRegex, 'Invalid phone number'),
							},
						],
						...formItemOptions,
					})
				: ele;
		}
		if (item.format === 'Email') {
			return isHave
				? getFormItem(ele, {
						rules: [
							{ required: item.is_required },
							{
								validator: (_, value) =>
									validateRule(value, emailRegex, 'Invalid email format'),
							},
						],
						...formItemOptions,
					})
				: ele;
		}
		if (item.format === 'Currency') {
			ele = <InputNumAmount {...defaultProps} />;
			// 因为在模板的情况下,后端框架只支持字符串,所以此处需将金额相关转成字符串(仅限于模板)
			return isHave
				? getFormItem(ele, {
						getValueFromEvent: (e) => e?.toString?.(),
						...formItemOptions,
					})
				: ele;
		}
		return isHave ? getFormItem(ele, formItemOptions) : ele;
	};

	const getComponents = (item = {}) => {
		if (item?.source === 'Template') {
			return getCurrentElement();
		}
		const names = allFormListName?.length ? allFormListName : formListName;
		const vKey = item?.key?.replace?.('.', '');
		// const exampleValue = exampleVariablesObj[vKey] || `{${vKey}}`;
		return (
			<>
				<Form.Item
					{...formListRestField}
					label={item.label}
					name={[...formListName, 'source']}
					rules={[{ required: item.is_required }]}
					size="small"
					style={{ marginBottom: 0 }}
				>
					<Radio.Group
						options={[
							{
								value: 'Dynamic',
								label: 'Default',
							},
							{
								value: 'Fixed',
								label: 'Customize',
							},
						]}
						disabled={isAllEditable ? false : !item.is_editable}
						size="small"
						onChange={({ target: { value } }) => {
							const fixedValue = form.getFieldValue([...names, 'fixed_value']);
							const dynamicValue = form.getFieldValue([...names, 'dynamicValue']);
							const exampleValue =
								dynamicValue || exampleVariablesObj[vKey] || `{${vKey}}`;
							let v = value === 'Fixed' ? fixedValue : exampleValue;
							changeHtml(item.key, v, item.format);
						}}
					/>
					{/* <Select
						placeholder={item.label}
						options={[
							{
								value: 'Fixed',
								label: 'Custom input',
							},
							{
								value: 'Dynamic',
								label: 'Default',
							},
						]}
					/> */}
				</Form.Item>
				<Form.Item
					shouldUpdate={(prevValues, curValues) => {
						const prevSource = _.get(prevValues, [...names, 'source']);
						const curSource = _.get(curValues, [...names, 'source']);
						return prevSource !== curSource;
					}}
					noStyle
				>
					{({ getFieldValue }) => {
						const currentSource = getFieldValue([...names, 'source']);

						return currentSource === 'Fixed'
							? getCurrentElement(
									{},
									{
										label: null,
									},
								)
							: getCurrentElement(
									{
										// value: exampleValue,
										// onChange: (v) => {},

										disabled: isAllEditable ? false : true,
									},
									{
										className: '!mb-4',
										name: [...formListName, 'dynamicValue'],
										label: null,
									},
									// false,
								);
					}}
				</Form.Item>
			</>
		);
	};
	return getComponents(item);
};
