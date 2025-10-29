import React from 'react';
import { Input } from 'antd';
import { formatTINOrSSN } from '@/lib/format';
// SSN和ITIN都是9位数字
let maxLength = 9;

// 格式化不同类型税务ID
const InputTIN = (props) => {
	const { value = '', onChange, type = 'SSN', ...rest } = props;

	const formatValue = (value) => {
		const numbers = value.replace(/\D/g, ''); // 提取纯数字
		const truncated = numbers.substring(0, maxLength);
		return {
			truncated,
			value: formatTINOrSSN(truncated, type),
		};
	};

	const handleChange = (e) => {
		const rawValue = e.target.value;
		const info = formatValue(rawValue);
		onChange?.(info?.truncated);
	};

	// 获取最大长度（包含分隔符）
	const getMaxLength = () => {
		switch (type) {
			case 'SSN':
			case 'ITIN':
				return 11; // XXX-XX-XXXX
			case 'EIN':
				return 10; // XX-XXXXXXX
			default:
				return 11;
		}
	};

	// 获取占位符
	const getPlaceholder = () => {
		switch (type) {
			case 'SSN':
				return 'XXX-XX-XXXX';
			case 'ITIN':
				return '9XX-XX-XXXX';
			case 'EIN':
				return 'XX-XXXXXXX';
			default:
				return 'XXX-XX-XXXX';
		}
	};

	return (
		<Input
			{...rest}
			value={formatValue(value)?.value}
			onChange={handleChange}
			placeholder={getPlaceholder()}
			maxLength={getMaxLength()}
		/>
	);
};

export default InputTIN;
