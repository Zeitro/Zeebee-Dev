import { Select, Space } from 'antd';
import React from 'react';
import InputTIN from '.';
import { tinOptions } from '@/enum/loan';
const SelectTINType = ({ value, onChange, isIndividual }) => {
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
					value={value?.type}
					onChange={(e) => onChange({ type: e, tin: value?.tin })}
				/>
				<InputTIN
					allowClear
					size="large"
					style={{ width: '70%' }}
					placeholder="Search by field"
					value={value?.tin}
					onChange={(e) => {
						let type = value?.type;
						// 因ssn与ITIN的格式相同，需特殊处理
						if (e.length === 9 && isIndividual) {
							type = e[0] === '9' ? 'ITIN' : 'SSN';
						}

						onChange({ type, tin: e });
					}}
					type={value?.type}
				/>
			</Space.Compact>
		</Space>
	);
};
export default SelectTINType;
