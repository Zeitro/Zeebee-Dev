import { Col, Input, Row, Select, Space, TreeSelect } from 'antd';
const { Search } = Input;
import React, { useRef } from 'react';
import styles from './index.module.scss';
import { debounce } from 'lodash';
const { SHOW_ALL } = TreeSelect;
/**
 * Form过滤数据组件，配合table使用
 * @returns
 */
const FilterField = ({ list, rowProps = {} }) => {
	const isFirst = useRef(true);
	// selectProps用于select和input组合的select-input类型
	const getItemDom = ({
		type,
		options,
		value,
		label,
		treeData,
		selectValue,
		selectProps,
		...eventProps
	}) => {
		let filterItem = '';
		switch (type) {
			case 'input-search':
				filterItem = (
					<Search
						size="large"
						placeholder="input search text"
						value={value}
						onSearch={eventProps?.onSearch}
						style={{ width: 200 }}
						onChange={debounce(eventProps?.onChange)}
						{...(eventProps || {})}
					/>
				);
				break;
			case 'select-input':
				filterItem = (
					<Space direction="vertical" size="large">
						<Space.Compact>
							<Select
								size="large"
								style={{ width: 120 }}
								placeholder="select"
								options={options}
								value={selectValue}
								onChange={selectProps?.onChange}
								{...(selectProps || {})}
							/>
							<Input
								size="large"
								style={{ width: 256 }}
								placeholder="Search by field"
								value={value}
								onChange={debounce(eventProps?.onChange)}
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
						style={{ width: 148 }}
						value={value}
						onChange={debounce(eventProps?.onChange)}
						{...(eventProps || {})}
					/>
				);
				break;
			case 'select':
				filterItem = (
					<Select
						size="large"
						placeholder="select"
						style={{ width: 148 }}
						options={options}
						value={value}
						onChange={eventProps?.onChange}
						{...(eventProps || {})}
					/>
				);
				break;
			case 'treeSelect':
				filterItem = (
					<TreeSelect
						size="large"
						placeholder="select"
						style={{ width: 148 }}
						maxTagCount={0}
						maxTagPlaceholder={() => `${label}(${value?.length})`}
						treeCheckable={true}
						value={value}
						showSearch={false}
						allowClear={false}
						showCheckedStrategy={SHOW_ALL}
						onChange={eventProps?.onChange}
						treeData={treeData}
						{...(eventProps || {})}
						classNames={{
							root: [
								styles.treeSelectRoot,
								value?.length && styles.treeSelectRoot_border,
								eventProps?.classNames?.root,
							].join(' '),
							popup: {
								root: [styles.treeSelectPopup,eventProps?.classNames?.popup?.root].join(' '),
							},
						}}
					/>
				);
				break;
			default:
				break;
		}
		return filterItem;
	};
	return (
		<Row
			gutter={8}
			{...rowProps}
			className={`${styles.filterField} w-full ${rowProps.className || ''}`}
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
};
export default FilterField;
