import { Form, Select, Space } from 'antd';
import React from 'react';
import { InputTIN } from '.';
import { tinOptions } from '@repo/shared-methods/enum';
import { validateRule, tinOrSSNRegex } from '../../lib/format';
export const SelectTINType = ({ value, onChange, isIndividual }) => {
	return (
		<Space direction="vertical" size="large" className="w-full">
			<Space.Compact className="w-full">
				<Select
					allowClear
					size="large"
					style={{ width: '30%' }}
					placeholder="select"
					options={tinOptions
						.map((i) => ({
							...i,
							disabled: isIndividual ? i.value === 'EIN' : i.value !== 'EIN',
						}))
						.sort((a, b) => a.disabled - b.disabled)}
					value={value?.tin_type}
					onChange={(e) => onChange({ tin_type: e, tin: value?.tin })}
				/>
				<InputTIN
					allowClear
					size="large"
					style={{ width: '70%' }}
					placeholder="Search by field"
					value={value?.tin}
					onChange={(e) => {
						let tin_type = value?.tin_type;
						// 因ssn与ITIN的格式相同，需特殊处理
						if (e.length === 9 && isIndividual) {
							tin_type = e[0] === '9' ? 'ITIN' : 'SSN';
						}

						onChange({ tin_type, tin: e });
					}}
					type={value?.tin_type}
				/>
			</Space.Compact>
		</Space>
	);
};

export const SelectFormTINType = ({
	isIndividual,
	disabled,
	typeName = ['tin_type'],
	tinName = ['tin'],
}) => {
	const form = Form.useFormInstance();
	const tin_type = Form.useWatch(typeName, form);

	return (
		<Space direction="vertical" size="large" className="w-full">
			<Space.Compact className="w-full">
				<Form.Item label="TIN" name={typeName} rules={[{ required: true }]}>
					<Select
						allowClear
						size="large"
						style={{ width: 90 }}
						placeholder="select"
						options={tinOptions
							.map((i) => ({
								...i,
								disabled: isIndividual ? i.value === 'EIN' : i.value !== 'EIN',
							}))
							.sort((a, b) => a.disabled - b.disabled)}
						disabled={disabled}
					/>
				</Form.Item>
				<Form.Item
					label={null}
					name={tinName}
					getValueFromEvent={(e) => {
						// 因ssn与ITIN的格式相同，需特殊处理
						if (e.length === 9 && isIndividual) {
							const tin_type = e[0] === '9' ? 'ITIN' : 'SSN';
							form.setFieldsValue({ [typeName]: tin_type });
						}
						return e;
					}}
					rules={[
						{
							required: true,
						},
						{
							validator: (_, value) =>
								validateRule(value, tinOrSSNRegex, 'Please enter valid TIN!'),
						},
					]}
					className="flex-1"
				>
					<InputTIN
						allowClear
						size="large"
						style={{ width: '100%', marginTop: 23 }}
						placeholder="Search by field"
						type={tin_type}
						disabled={disabled}
					/>
				</Form.Item>
			</Space.Compact>
		</Space>
	);
};
