// src/components/TaxIdInput/index.js
import React, { useState, useEffect } from 'react';
import { Input } from 'antd';
// SSN和ITIN都是9位数字
let maxLength = 9;
// 格式化不同类型税务ID
const InputTIN = (props) => {
	const { value = '', onChange, type = 'SSN', ...rest } = props;
	const [inputValue, setInputValue] = useState('');

	// 格式化不同类型税务ID
	const formatTaxId = (taxId) => {
		// 移除所有非数字字符并限制长度
		const cleaned = taxId.replace(/\D/g, '');

		switch (type) {
			case 'SSN':
				// SSN格式: XXX-XX-XXXX (9位数字)
				if (cleaned.length <= 3) return cleaned;
				if (cleaned.length <= 5) return `${cleaned.slice(0, 3)}-${cleaned.slice(3)}`;
				return `${cleaned.slice(0, 3)}-${cleaned.slice(3, 5)}-${cleaned.slice(5, 9)}`;

			case 'ITIN':
				// ITIN格式: XXX-XX-XXXX (9位数字，与SSN相同格式)
				if (cleaned.length <= 3) return cleaned;
				if (cleaned.length <= 5) return `${cleaned.slice(0, 3)}-${cleaned.slice(3)}`;
				return `${cleaned.slice(0, 3)}-${cleaned.slice(3, 5)}-${cleaned.slice(5, 9)}`;

			case 'EIN':
				// EIN格式: XX-XXXXXXX (9位数字)
				if (cleaned.length <= 2) return cleaned;
				return `${cleaned.slice(0, 2)}-${cleaned.slice(2, 9)}`;

			default:
				return cleaned;
		}
	};

	// 提取纯数字
	const extractNumbers = (formatted) => {
		return formatted.replace(/\D/g, '');
	};

	// 初始化和更新输入值
	useEffect(() => {
		const numbers = extractNumbers(value);
		const truncated = numbers.substring(0, maxLength);
		setInputValue(formatTaxId(truncated));
	}, [value, type]);

	const handleChange = (e) => {
		const rawValue = e.target.value;
		const numbersOnly = extractNumbers(rawValue);
		const truncated = numbersOnly.substring(0, maxLength);
		const formattedValue = formatTaxId(truncated);

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
