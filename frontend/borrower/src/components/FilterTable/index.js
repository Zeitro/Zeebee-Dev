import React, { useEffect, useMemo, useRef, useState } from 'react';
import FilterField from '@/components/FilterField';
import styles from './index.module.scss';
import NewTable from '@/components/NewTable';
import { Col, Row, Tag } from 'antd';
import { createRequest } from '@/lib/fetch';
import { getPagination } from '@/lib/common';
const FilterTable = ({
	columns,
	scroll,
	dataSource,
	filterList,
	rowProps = {},
	filterName = 'Active filters',
	hasFilterField = true,
	filterWrapClassName,
	getData = () => {},
	...props
}) => {
	const filterFieldRef = useRef();
	const { request } = createRequest();
	const pagination = useRef(getPagination());
	const [sorter, setsorter] = useState({
		column: 'address',
		sort: 'asc',
	});
	const [tableList, setTableList] = useState([]);
	const [loading, setLoading] = useState(false);
	const newList = useMemo(() => {
		return filterList?.filter((i, index) => i.type === 'treeSelect' && i.value?.length);
	}, [filterList]);

	const getList = async () => {
		try {
			setLoading(true);
			const filterParams = filterFieldRef?.current?.getApiParams?.();
			const { total, list } = await getData({
				page: pagination.current?.current,
				limit: pagination.current?.pageSize,
				sorter,
				filter: {
					...filterParams,
				},
			});
			pagination.current.total = total;
			setTableList(list);
			setLoading(false);
		} catch (error) {
			setLoading(false);
		}
	};

	const onChange = (pagination, filters, sorter, extra) => {
  console.log('params', pagination, filters, sorter, extra);
};

	useEffect(() => {
		if (dataSource) {
			setTableList(dataSource);
		} else {
			getList();
		}
	}, [dataSource]);
	return (
		<div className={[styles.filterTable, filterWrapClassName].join(' ')}>
			{hasFilterField && <FilterField ref={filterFieldRef} list={filterList} {...rowProps} />}
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
			<NewTable
				dataSource={tableList}
				loading={loading}
				columns={columns}
				scroll={scroll}
				pagination={pagination}
				onChange={onChange}
				{...props}
			/>
		</div>
	);
};

export default FilterTable;
