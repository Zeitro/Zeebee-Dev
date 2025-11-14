import React, {
	useEffect,
	useRef,
	useState,
	forwardRef,
	useImperativeHandle,
	useCallback,
} from 'react';
import styles from './index.module.scss';
import { createRequest } from '@/lib/fetch';
import { getPagination } from '@/lib/common';
import { asc, des, ascend, descend } from '@/enum/common';
import { useDisUpdateEffect } from '@/hooks/common';
import { FilterTag } from '../FilterField/filterTag.js';
import { FilterField } from '../FilterField/index';
import { NewTable } from '../NewTable';
/**
 * Form筛选 + Table
 * @param {*Boolean} hasFilterField true显示筛选/false隐藏
 * @param {*function} getData table list接口请求函数
 * @param {*string} api table list接口请求的api，没有api时会默认调用getData
 * @param {*string} method 接口请求的类别。与api一同使用
 * @param {*String} filterName 筛选汇总的标题
 */
export const FilterTable = forwardRef(
	(
		{
			columns,
			scroll,
			dataSource,
			dataSourceLoading,
			filterList,
			rowProps = {},
			filterName = '',
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
			// column: 'address',
			sort: des,
		});
		const [tableList, setTableList] = useState([]);
		const [loading, setLoading] = useState(false);

		const getList = useCallback(async () => {
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
				let newApi = api;
				if (newApi) {
					if (method === 'GET') {
						const params = new URLSearchParams(params);
						newApi += `?${params}`;
						res = await request(newApi, { method });
					} else {
						res = await request(newApi, { method, body: params });
					}
				} else {
					res = await getData(params);
				}
				const { total, list } = res;
				pagination.current.total = total;
				setTableList(list);
				setLoading(false);
			} catch {
				setLoading(false);
			}
			// eslint-disable-next-line react-hooks/exhaustive-deps
		}, [api, getData, method]);

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
		}, [dataSource, getList]);

		useDisUpdateEffect(() => {
			setLoading(dataSourceLoading);
		}, [dataSourceLoading]);

		useImperativeHandle(ref, () => ({
			handleSearch,
		}));

		return (
			<div className={[styles.filterTable, filterWrapClassName].join(' ')}>
				{hasFilterField && (
					<>
						<FilterField ref={filterFieldRef} list={filterList} {...rowProps} />
						<FilterTag filterList={filterList} filterName={filterName} />
					</>
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
	},
);
FilterTable.displayName = 'FilterTable';
