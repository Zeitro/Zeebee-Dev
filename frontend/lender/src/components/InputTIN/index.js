// src/components/TaxIdInput/index.js
import React, { useState, useEffect } from 'react';
import { Input } from 'antd';
import { formatTINOrSSN } from '@/lib/format';
// SSN和ITIN都是9位数字
let maxLength = 9;

// 格式化不同类型税务ID
const InputTIN = (props) => {
	const { value = '', onChange, type = 'SSN', ...rest } = props;
	const [inputValue, setInputValue] = useState('');

	// 提取纯数字
	const extractNumbers = (formatted) => {
		return formatted.replace(/\D/g, '');
	};

	// 初始化和更新输入值
	useEffect(() => {
		const numbers = extractNumbers(value);
		const truncated = numbers.substring(0, maxLength);
		setInputValue(formatTINOrSSN(truncated, type));
	}, [value, type]);

	const handleChange = (e) => {
		const rawValue = e.target.value;
		const numbersOnly = extractNumbers(rawValue);
		const truncated = numbersOnly.substring(0, maxLength);
		const formattedValue = formatTINOrSSN(truncated, type);

		setInputValue(formattedValue);
		onChange?.(truncated);
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
			value={inputValue}
			onChange={handleChange}
			placeholder={getPlaceholder()}
			maxLength={getMaxLength()}
		/>
	);
};

export default InputTIN;
