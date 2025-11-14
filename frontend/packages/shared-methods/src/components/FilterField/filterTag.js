import React, { useMemo } from 'react';
import styles from './index.module.scss';
import { Col, Row, Tag } from 'antd';
/**
 * FilterField组件下拉多选筛选出来的字段展示组件
 * @param {*Array} filterList 筛选字段
 * @param {*string} filterName 筛选标题
 */
export const FilterTag = ({ filterList, filterName = '' }) => {
	const newList = useMemo(() => {
		return filterList?.filter((i) => i.filterType === 'treeSelect' && i.value?.length);
	}, [filterList]);

	return newList.length > 0 ? (
		<div className={styles.filterStatistics}>
			{filterName && <div className={styles.filterName}>{filterName}</div>}
			<Row gutter={12}>
				{newList.map((item, index) => {
					return (
						<Col key={index}>
							<span className="pr-2">{item.label}:</span>
							{item.value?.map((vItem, vIndex) => {
								const labelItem = item?.treeData?.find((i) => i.value === vItem);
								return (
									<Tag
										key={vIndex}
										closeIcon
										color="processing"
										onClose={() => {
											const newValue = item.value.filter((i) => i !== vItem);
											item?.onChange(newValue);
										}}
									>
										{labelItem?.title}
									</Tag>
								);
							})}
						</Col>
					);
				})}
			</Row>
		</div>
	) : (
		''
	);
};

FilterTag.displayName = 'FilterTag';
