import React, {
	useEffect,
	useMemo,
	useRef,
	useState,
	forwardRef,
	useImperativeHandle,
} from 'react';
import FilterField from '@/components/FilterField';
import styles from './index.module.scss';
import NewTable from '@/components/NewTable';
import { Col, Row, Tag } from 'antd';
import { createRequest } from '@/lib/fetch';
import { getPagination } from '@/lib/common';
import { asc, des, ascend, descend } from '@/enum/common';
import { useDisUpdateEffect } from '@/hooks/common';
/**
 * Form筛选 + Table
 * @param {*Boolean} hasFilterField true显示筛选/false隐藏
 * @param {*function} getData table list接口请求函数
 * @param {*string} api table list接口请求的api，没有api时会默认调用getData
 * @param {*string} method 接口请求的类别。与api一同使用
 * @param {*String} filterName 筛选汇总的标题
 */
const FilterTable = (
	{
		columns,
		scroll,
		dataSource,
		dataSourceLoading,
		filterList,
		rowProps = {},
		filterName = 'Active filters',
		hasFilterField = true,
		filterWrapClassName,
		getData = () => {},
		api,
		method = 'POST',
		...props
	},
	ref,
) => {
	const isGet = useRef(true);

	const filterFieldRef = useRef();
	const { request } = createRequest();
	const pagination = useRef(getPagination());
	const sorter = useRef({
		column: 'address',
		sort: des,
	});
	const [tableList, setTableList] = useState([]);
	const [loading, setLoading] = useState(false);
	const newList = useMemo(() => {
		return filterList?.filter((i) => i.type === 'treeSelect' && i.value?.length);
	}, [filterList]);

	const getList = async () => {
		try {
			setLoading(true);
			const filterParams = filterFieldRef?.current?.getApiParams?.();
			const params = {
				page: pagination.current?.current,
				limit: pagination.current?.pageSize,
				filter: {
					...filterParams,
				},
			};
			if (sorter.current?.column) {
				params.sorter = sorter.current;
			}
			let res = {};
			if (api) {
				if (method === 'GET') {
					const params = new URLSearchParams(params);
					api += `?${params}`;
					res = await request(api, { method });
				} else {
					res = await request(api, { method, body: params });
				}
			} else {
				res = await getData(params);
			}
			const { total, list } = res;
			pagination.current.total = total;
			setTableList(list);
			setLoading(false);
		} catch (error) {
			setLoading(false);
		}
	};

	const onChange = (paginationObj, filters, { order, columnKey }) => {
		const orderObj = {
			[descend]: des,
			[ascend]: asc,
		};
		sorter.current.column = columnKey;
		sorter.current.sort = orderObj[order];
		if (paginationObj.current !== pagination.current.current) {
			pagination.current.current = paginationObj.current;
		} else {
			pagination.current.current = 1;
		}
		getList();
	};

	const handleSearch = () => {
		pagination.current.current = 1;
		getList();
	};

	useEffect(() => {
		if (dataSource && Array.isArray(dataSource)) {
			setTableList(dataSource);
		} else if (isGet.current) {
			getList();
			isGet.current = false;
		}
	}, [dataSource]);

	useDisUpdateEffect(() => {
		setLoading(dataSourceLoading);
	}, [dataSourceLoading]);

	useImperativeHandle(ref, () => ({
		handleSearch,
	}));

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
				onChange={onChange}
				{...props}
				pagination={props?.pagination ?? pagination.current}
				className={[styles.filter_table_content, props?.className].join(' ')}
			/>
		</div>
	);
};

export default forwardRef(FilterTable);
