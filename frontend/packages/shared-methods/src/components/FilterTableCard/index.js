'use client';
import React, {
	useEffect,
	useRef,
	useState,
	forwardRef,
	useImperativeHandle,
	useCallback,
} from 'react';
import { createRequest } from '../../lib/fetch';
import { getPagination } from '../../lib/common';
import { asc, des, ascend, descend } from '@repo/shared-methods/enum';
import { useDisUpdateEffect } from '@repo/shared-methods/hooks';
import { FilterTag } from '../FilterField/filterTag.js';
import { FilterField } from '../FilterField/index';
import { NewTable } from '../NewTable';
import { Cards } from './cards.js';
import styles from './index.module.scss';
import useContent from '../FilterTable/useContent';
const defaultSortOrder = {
	// column: 'address',
	sort: des,
};
/**
 * Form筛选 + Table/Card（两个可切换）
 * 与FilterTable传一样的参数
 */
export const FilterTableCard = forwardRef(
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
			hasPaginationRequest = true, // 是否请求分页数据
			api,
			method = 'POST',
			renderItem = () => {},
			itemClassName,
			typeKey = 'table',
			cardTitle,
			...props
		},
		ref,
	) => {
		const {
			filterFieldRef,
			list,
			setList,
			loading,
			onChange,
			pagination,
			handleSearch,
			getList,
			isGet,
			setLoading,
			paramsData,
			scrollRef,
			handleChangePagination,
		} = useContent({
			api,
			method,
			getData,
			hasPaginationRequest,
		});
		// const scrollRef = useRef();
		// const isGet = useRef(true);
		// const filterFieldRef = useRef();
		// const { request, cancel } = createRequest();
		// // 因为pagination、sorter改变后就要调用接口，useState是异步，导致getList访问的是旧值，所以用useRef。
		// const pagination = useRef(getPagination());
		// const sorter = useRef(defaultSortOrder);
		// // 因为有的页面需要展示排序和total，useRef的pagination、sorter会报错，页面上渲染的数据必须是useState的数据，所以此处做备份展示
		// const [paramsData, setParamsData] = useState({
		// 	...defaultSortOrder,
		// 	...(getPagination() || {}),
		// });
		// const [list, setList] = useState([]);
		// const [loading, setLoading] = useState(false);

		// const getList = useCallback(async () => {
		// 	try {
		// 		setLoading(true);
		// 		const filterParams = filterFieldRef?.current?.getApiParams?.();
		// 		const params = {
		// 			page: pagination?.current?.current,
		// 			limit: pagination?.current?.pageSize,
		// 			filter: {
		// 				...filterParams,
		// 			},
		// 		};
		// 		if (sorter?.current?.column) {
		// 			params.sorter = sorter?.current;
		// 		}
		// 		let res = {};
		// 		let newApi = api;
		// 		if (newApi) {
		// 			cancel(newApi);
		// 			if (method === 'GET') {
		// 				const params = new URLSearchParams(params);
		// 				newApi += `?${params}`;
		// 				res = await request(newApi, { method });
		// 			} else {
		// 				res = await request(newApi, { method, body: params });
		// 			}
		// 		} else {
		// 			res = await getData(params);
		// 		}
		// 		const { total, list } = res;
		// 		pagination.current.total = total;
		// 		setParamsData((v) => ({ ...v, total }));
		// 		setList(list);
		// 		setLoading(false);
		// 		return list;
		// 	} catch {
		// 		setLoading(false);
		// 	}
		// 	// eslint-disable-next-line react-hooks/exhaustive-deps
		// }, [api, getData, method]);

		// const handleSearch = async () => {
		// 	pagination.current.current = 1;
		// 	setParamsData((v) => ({ ...v, current: 1 }));
		// 	return await getList();
		// };

		// const handleChangePagination = (page) => {
		// 	pagination.current.current = page;
		// 	setParamsData((v) => ({ ...v, current: page }));
		// 	getList();
		// 	scrollRef.current.scrollTo({
		// 		top: 0,
		// 	});
		// };

		// const onChange = (paginationObj, filters, { order, columnKey }) => {
		// 	const orderObj = {
		// 		[descend]: des,
		// 		[ascend]: asc,
		// 	};
		// 	pagination.current.current = paginationObj.current;
		// 	sorter.current.sort = orderObj[order];
		// 	sorter.current.column = columnKey;
		// 	setParamsData((v) => ({ ...v, ...sorter.current, ...pagination.current }));
		// 	getList();
		// };
		const getMoreTypesContent = (type) => {
			switch (type) {
				case 'table':
					return (
						<NewTable
							dataSource={list}
							loading={loading}
							columns={columns}
							scroll={scroll}
							onChange={onChange}
							{...props}
							pagination={props?.pagination ?? paramsData}
							className={[
								'flex-1 overflow-hidden',
								styles.filter_table_content,
								props?.className,
							].join(' ')}
						/>
					);
				case 'card':
					return (
						<Cards
							pagination={paramsData}
							loading={loading}
							scrollRef={scrollRef}
							list={list}
							handleChangePagination={handleChangePagination}
							itemClassName={itemClassName}
							renderItem={renderItem}
							cardTitle={cardTitle}
						/>
					);
				default:
					break;
			}
		};

		useEffect(() => {
			if (dataSource && Array.isArray(dataSource)) {
				setList(dataSource);
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
			<div className={['w-full h-full flex flex-col gap-4', filterWrapClassName].join(' ')}>
				{hasFilterField && (
					<FilterField ref={filterFieldRef} list={filterList} {...rowProps} />
				)}
				<FilterTag filterList={filterList} filterName={filterName} />
				{getMoreTypesContent(typeKey)}
			</div>
		);
	},
);
FilterTableCard.displayName = 'FilterTableCard';
