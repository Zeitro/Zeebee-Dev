import { Form, Select, Space } from 'antd';
import React from 'react';
import { InputTIN } from '.';
import { tinOptions } from '@/enum/loan';
import { validateRule, tinOrSSNRegex } from '@/lib/format';
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

export const SelectFormTINType = ({ isIndividual }) => {
	const form = Form.useFormInstance();
	const tin_type = Form.useWatch('tin_type', form);
	return (
		<Space direction="vertical" size="large" className="w-full">
			<Space.Compact className="w-full">
				<Form.Item label="TIN" name="tin_type" rules={[{ required: true }]}>
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
					/>
				</Form.Item>
				<Form.Item
					label={null}
					name="tin"
					getValueFromEvent={(e) => {
						// 因ssn与ITIN的格式相同，需特殊处理
						if (e.length === 9 && isIndividual) {
							const tin_type = e[0] === '9' ? 'ITIN' : 'SSN';
							form.setFieldsValue({ tin_type });
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
				>
					<InputTIN
						allowClear
						size="large"
						style={{ width: '100%', marginTop: 23 }}
						placeholder="Search by field"
						type={tin_type}
					/>
				</Form.Item>
			</Space.Compact>
		</Space>
	);
};
