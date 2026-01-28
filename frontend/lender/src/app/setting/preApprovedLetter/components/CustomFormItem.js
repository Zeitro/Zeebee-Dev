import { InputNumAmount, InputPhone } from '@repo/shared-methods';
import { emailRegex, phoneRegex, validateRule } from '@repo/shared-methods/utils';
import { DatePicker, Form, Input, Select } from 'antd';
import _ from 'lodash';

export const CustomFormItem = ({
	item,
	isHaveFormItem = true,
	formListName = [],
	formListRestField = {},
	allFormListName = [],
}) => {
	const getFormItem = (children, options) => {
		return (
			<Form.Item
				{...formListRestField}
				label={item.label}
				name={[...formListName, 'fixed_value']}
				rules={[{ required: true }]}
				{...options}
			>
				{children}
			</Form.Item>
		);
	};
	const getCurrentElement = (options = {}) => {
		let defaultProps = {
			placeholder: item.label,
			disabled: !item.is_editable,
		};
		let ele = <Input {...defaultProps} />;
		if (item.format === 'Input') {
			// ele = <Input {...defaultProps} />;
		}
		if (item.format === 'TextArea') {
			ele = <Input.TextArea {...defaultProps} />;
		}
		if (item.format === 'Date') {
			ele = <DatePicker {...defaultProps} />;
		}
		if (item.format === 'Phone') {
			ele = <InputPhone {...defaultProps} />;
			return isHaveFormItem
				? getFormItem(ele, {
						rules: [
							{ required: true },
							{
								validator: (_, value) =>
									validateRule(value, phoneRegex, 'Invalid phone number'),
							},
						],
						...options,
					})
				: ele;
		}
		if (item.format === 'Email') {
			return isHaveFormItem
				? getFormItem(ele, {
						rules: [
							{ required: true },
							{
								validator: (_, value) =>
									validateRule(value, emailRegex, 'Invalid email format'),
							},
						],
						...options,
					})
				: ele;
		}
		if (item.format === 'Currency') {
			ele = <InputNumAmount {...defaultProps} />;
			// 因为在模板的情况下,后端框架只支持字符串,所以此处需将金额相关转成字符串(仅限于模板)
			return isHaveFormItem
				? getFormItem(ele, {
						getValueFromEvent: (e) => e?.toString?.(),
						...options,
					})
				: ele;
		}
		return isHaveFormItem ? getFormItem(ele, options) : ele;
	};

	const getComponents = (item) => {
		if (item.source === 'Template') {
			return getCurrentElement();
		}
		return (
			<>
				<Form.Item
					{...formListRestField}
					label={item.label}
					name={[...formListName, 'source']}
					rules={[{ required: true }]}
				>
					<Select
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
					/>
				</Form.Item>
				<Form.Item
					shouldUpdate={(prevValues, curValues) => {
						const prevSource = _.get(prevValues, [...allFormListName, 'source']);
						const curSource = _.get(curValues, [...allFormListName, 'source']);
						return prevSource !== curSource;
					}}
					noStyle
				>
					{({ getFieldValue }) => {
						const currentSource = getFieldValue([...allFormListName, 'source']);
						return (
							currentSource === 'Fixed' &&
							getCurrentElement({ noStyle: false, label: null })
						);
					}}
				</Form.Item>
			</>
		);
	};
	return getComponents(item);
};
