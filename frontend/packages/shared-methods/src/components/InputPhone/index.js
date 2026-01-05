'use client';
import React, { useState, useEffect } from 'react';
import { Input } from 'antd';
import { formatPhone } from '../../lib';

export const InputPhone = (props) => {
	const { value = '', onChange, ...rest } = props;
	const [inputValue, setInputValue] = useState('');

	// 从格式化字符串中提取纯数字
	const extractNumbers = (formatted) => {
		return formatted.replace(/\D/g, '').substring(0, 10);
	};

	// 初始化和更新输入值
	useEffect(() => {
		setInputValue(formatPhone(value));
	}, [value]);

	const handleChange = (e) => {
		const rawValue = e.target.value;
		const numbersOnly = extractNumbers(rawValue);
		const formattedValue = formatPhone(numbersOnly);

		setInputValue(formattedValue);
		onChange?.(numbersOnly);
	};

	return (
		<Input
			{...rest}
			value={inputValue}
			onChange={handleChange}
			placeholder="(XXX) XXX-XXXX"
			maxLength={14}
		/>
	);
};
