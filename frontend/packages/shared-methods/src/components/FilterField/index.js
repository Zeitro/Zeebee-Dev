'use client';
import {Button, Col, Input, Row, Select, Space, TreeSelect} from 'antd';
import React, {forwardRef, useImperativeHandle} from 'react';
import {SearchOutlined} from '@ant-design/icons';
import styles from './index.module.scss';
import {debounce} from 'lodash';

const {Search} = Input;
const { SHOW_ALL } = TreeSelect;
/**
 * Form过滤数据组件，配合table使用
 * type 区分哪种组件（input-search、select-input、input、select、treeSelect）
 * paramsKey 请求接口时的key
 * value 为paramsKey对应的值
 * selectParamsKey （只有type为select-input有）请求接口时的key
 * selectValue 为selectParamsKey对应的值
 * options 当有select组件时的下拉列表数据
 * treeData 当type为treeSelect时，传入的tree下拉数据
 * label 现只用在type为treeSelect时，Placeholder显示
 */
export const FilterField = forwardRef(({ list, className, ...rowProps }, ref) => {
	// selectProps用于select和input组合的select-input类型
	const getItemDom = ({
		filterType,
		options,
		value = null,
		label,
		treeData,
		selectValue,
		selectProps,
		text = 'Search',
		children,
		selectParamsKey,
		paramsKey,
		...eventProps
	}) => {
		let filterItem = '';
		switch (filterType) {
			case 'input-search':
				filterItem = (
					<Search
						allowClear
						size="large"
						placeholder="input search text"
						value={value}
						onSearch={eventProps?.onSearch}
						style={{ width: 240 }}
						onChange={debounce((e) => eventProps?.onChange(e))}
						{...(eventProps || {})}
					/>
				);
				break;
			case 'select-input':
				filterItem = (
					<Space direction="vertical" size="large">
						<Space.Compact>
							<Select
								allowClear
								size="large"
								style={{ width: 120 }}
								placeholder="select"
								options={options}
								value={selectValue}
								onChange={selectProps?.onChange}
								{...(selectProps || {})}
							/>
							<Input
								allowClear
								size="large"
								style={{ width: 240 }}
								placeholder="Search by field"
								value={value}
								onChange={debounce((e) => eventProps?.onChange(e))}
								{...(eventProps || {})}
							/>
						</Space.Compact>
					</Space>
				);
				break;
			case 'input':
				filterItem = (
					<Input
						size="large"
						placeholder="input text"
						style={{ width: 240 }}
						value={value}
						onChange={debounce((e) => eventProps?.onChange(e))}
						allowClear
						{...(eventProps || {})}
					/>
				);
				break;
			case 'select':
				filterItem = (
					<Select
						size="large"
						allowClear
						style={{ width: 180 }}
						options={options}
						value={value}
						onChange={(e) => eventProps?.onChange(e)}
						{...(eventProps || {})}
					/>
				);
				break;
			case 'treeSelect':
				filterItem = (
					<TreeSelect
						size="large"
						placeholder="select"
						style={{ width: 180 }}
						maxTagCount={0}
						maxTagPlaceholder={() => `${label}(${value?.length})`}
						treeCheckable={true}
						value={value}
						showSearch={true}
						allowClear={true}
						showCheckedStrategy={SHOW_ALL}
						onChange={(e) => eventProps?.onChange(e)}
						treeData={treeData}
						{...(eventProps || {})}
						classNames={{
							root: [
								styles.treeSelectRoot,
								value?.length && styles.treeSelectRoot_border,
								eventProps?.classNames?.root,
							].join(' '),
							popup: {
								root: [
									styles.treeSelectPopup,
									eventProps?.classNames?.popup?.root,
								].join(' '),
							},
						}}
					/>
				);
				break;
			case 'button':
				filterItem = children || (
					<Button
						type="primary"
						size="large"
						icon={<SearchOutlined />}
						onClick={debounce((e) => eventProps?.onClick(e), 300)}
						{...(eventProps || {})}
					>
						{text}
					</Button>
				);
				break;
			case 'custom':
				filterItem = children;
				break;
			default:
				break;
		}
		return filterItem;
	};

	// 调用此方法，将filterList转成请求接口参数的形式，必传paramsKey（selectParamsKey在select-input时必传）
	const getApiParams = () => {
		const params = {};
		list.map((item) => {
			if (item?.selectParamsKey && item?.selectValue) {
				params[item.selectParamsKey] = item.selectValue;
			}
			if (item?.paramsKey && item?.value) {
				params[item.paramsKey] = item.value;
			}
		});
		return params;
	};

	useImperativeHandle(ref, () => ({
		getApiParams,
	}));
	return (
		<Row
			gutter={[8, 8]}
			{...rowProps}
			className={`${styles.filterField} w-full ${className || ''}`}
		>
			{list?.map(({ colProps, ...item }, index) => {
				return (
					<Col {...(colProps || {})} key={index}>
						{getItemDom(item, index)}
					</Col>
				);
			})}
		</Row>
	);
});

FilterField.displayName = 'FilterField';
