'use client';
import { InputNumber } from 'antd';
import React from 'react';

/**
 * 金额输入框
 * InputNumber所需值
 * */
const InputNumAmount = ({
	value,
	onChange,
	placeholder,
	style = {},
	min = 0,
	max = 1000000000000,
	...props
}) => {
	return (
		<InputNumber
			style={{ width: '100%', ...style }}
			prefix={<span className="text-base pt-[2px]">$</span>}
			formatter={(value) => {
				// 千分位逗号隔开，保留两位小数
				const [int, dec] = value.split('.');
				const str = dec?.slice?.(0, 2);
				return `${int.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}` + (dec ? `.${str}` : '');
			}}
			parser={(value) => parseFloat(value.replace(/[^\d.]/g, '')) || 0}
			min={min}
			max={max}
			value={value}
			onChange={onChange}
			controls={false}
			placeholder={placeholder}
			{...props}
		/>
	);
};
export default InputNumAmount;
