import { useOptionsStore } from '@/store/options';
import { Col, Form, Select } from 'antd';
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
		if (list.length && isInitGet.current) {
			getInitOptions();
			isInitGet.current = false;
		}
	}, [JSON.stringify(list)]);

	// useEffect(() => {
	// 	if (!loanPurposeOptions?.length && loanPurposeName) {
	// 		getOptions('loan_purpose', loanPurposeName);
	// 	}
	// }, [loanPurposeOptions, loanPurposeName]);

	// useEffect(() => {
	// 	if (!loanTypeOptions?.length && loanTypeName) {
	// 		getOptions('loan_type', loanTypeName);
	// 	}
	// }, [loanTypeOptions, loanTypeName]);

	// useEffect(() => {
	// 	if (!loanPendingStatus?.length && pendingStatusName) {
	// 		getOptions('loan_pending_status', 'loanPendingStatus');
	// 	}
	// }, [loanPendingStatus, pendingStatusName]);

	// useEffect(() => {
	// 	!loanStatusOptions?.length && getOptions('loan_stage', 'loanStatusOptions');
	// }, [loanStatusOptions]);
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
			...selectProps
		},
		layoutCol = {},
	) => {
		let selectMaxCount = mode ? { maxCount: MAX_COUNT } : {};
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
				</Form.Item>
			</Col>
		);
	};
	return { ...optionsInfo, getSelectFormDom };
};
