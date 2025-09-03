import { useOptionsStore } from '@/store/options';
import { Col, Form, Radio, Select } from 'antd';
import { useCallback, useEffect, useRef } from 'react';
/**
 * 请求options下拉数据公共接口
 * list = [{
 *      key: loan_purpose,
 *      valueName: loanPurposeOptions
 *      type: 'lender' //默认可不要type，会请求getOptions方法
 * }]
 */
const MAX_COUNT = 5;
export const useGetOptions = (list = []) => {
	const optionsInfo = useOptionsStore();
	const isInitGet = useRef(true);
	const { getOptions, getLender, getBorrowerOptions, ...state } = optionsInfo;
	const getInitOptions = useCallback(async () => {
		await Promise.all(
			list.map(async (item) => {
				try {
					if (item.type === 'lender') {
						return await getLender();
					} else if (item.type === 'borrower') {
						return await getBorrowerOptions();
					} else if (!optionsInfo?.[item.valueName]?.length && item.key) {
						return await optionsInfo?.getOptions(item.key, item.valueName);
					}
					return Promise.resolve();
				} catch (error) {
					return Promise.resolve();
				}
			}),
		);
	}, [getOptions, list, state]);

	useEffect(() => {
		if (list?.length && isInitGet.current) {
			getInitOptions();
			isInitGet.current = false;
		}
	}, [list, getInitOptions]);

	const getSelectFormDom = (
		{
			label,
			name,
			rules = [],
			requiredMessage = 'Please enter value!',
			formProps = {},
			valueName, // 对应存储在store里面的值
			mode = null,
			placeholder = 'Select',
			componentType = 'select',
			...selectProps
		},
		layoutCol = {},
	) => {
		let selectMaxCount = mode ? { maxCount: MAX_COUNT } : {};
		const getComponentDom = (componentType) => {
			switch (componentType) {
				case 'select':
					return (
						<Select
							allowClear
							showSearch
							mode={mode}
							style={{ width: '100%' }}
							placeholder={placeholder}
							options={optionsInfo?.[valueName] || []}
							{...selectMaxCount}
							{...selectProps}
						/>
					);
				case 'radio':
					return (
						<Radio.Group
							options={optionsInfo?.[valueName] || []}
							{...selectProps}
							style={{ width: '100%' }}
						/>
					);
				default:
					break;
			}
		};
		return (
			<Col xs={24} sm={24} {...layoutCol}>
				<Form.Item
					label={label}
					name={name}
					rules={[
						{
							required: true,
							message: requiredMessage,
						},
						...rules,
					]}
					{...formProps}
				>
					{getComponentDom(componentType)}
				</Form.Item>
			</Col>
		);
	};
	return { ...optionsInfo, getSelectFormDom };
};
