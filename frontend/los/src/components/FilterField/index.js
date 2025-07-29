import { Col, Input, Row, Select, Space, TreeSelect } from 'antd';
const { Search } = Input;
import React from 'react';
import styles from './index.module.scss';
/**
 * Form过滤数据组件，配合table使用
 * @returns
 */
const FilterField = ({ list, rowProps = {} }) => {
	const getItemDom = ({ type, options, value, label,selectValue,selectProps, ...eventProps }, index) => {
		let filterItem = '';
		switch (type) {
			case 'input-search':
				filterItem = (
					<Search
						size="large"
						placeholder="input search text"
                        value={value}
						// onSearch={onSearch}
						style={{ width: 200 }}
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
								{...(selectProps || {})}
							/>
							<Input
								size="large"
								style={{ width: 256 }}
								placeholder="Search by field"
                                value={value}
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
						treeCheckable={true}
                        value={value}
						tagRender={() => `${label}(${value?.length})`}
						{...(eventProps || {})}
						classNames={{
							popup: {
								root: styles.treeSelectPopup,
							},
                            ...(eventProps?.classNames || {})
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
