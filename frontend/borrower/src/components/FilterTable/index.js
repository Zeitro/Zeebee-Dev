import React, { useMemo } from 'react';
import FilterField from '@/components/FilterField';
import styles from './index.module.scss';
import NewTable from '@/components/NewTable';
import { Col, Row, Tag } from 'antd';
const FilterTable = ({
	columns,
	scroll,
	dataSource,
	filterList,
	rowProps={},
	filterName = 'Active filters',
	hasFilterField=true,
	filterWrapClassName,
	...props
}) => {
	const newList = useMemo(() => {
		return filterList?.filter((i, index) => i.type === 'treeSelect' && i.value?.length);
	}, [filterList]);
	return (
		<div className={[styles.filterTable,filterWrapClassName].join(' ')}>
			{hasFilterField &&<FilterField list={filterList} {...rowProps} />}
			{!!newList.length && (
				<div className={styles.filterStatistics}>
					<div className={styles.filterName}>{filterName}</div>
					<Row gutter={12}>
						{newList.map((item, index) => {
							return (
								<Col key={index}>
									<span className="pr-2">{item.label}:</span>
									{item.value?.map((vItem, vIndex) => {
										const labelItem = item?.treeData?.find(
											(i) => i.value === vItem,
										);
										return (
											<Tag
												key={vIndex}
												closeIcon
												color="processing"
												onClose={() => {
													const newValue = item.value.filter(
														(i) => i !== vItem,
													);
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
			)}
			<NewTable dataSource={dataSource} columns={columns} scroll={scroll} {...props}/>
		</div>
	);
};

export default FilterTable;
