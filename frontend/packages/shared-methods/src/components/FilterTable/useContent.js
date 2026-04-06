import React, { useRef, useState, useCallback } from 'react';
import { createRequest } from '../../lib/fetch';
import { getPagination } from '../../lib/common';
import { asc, des, ascend, descend } from '@repo/shared-methods/enum';
const defaultSortOrder = {
	// column: 'address',
	sort: des,
};
const useContent = ({ getData = () => {}, api, method = 'POST' }) => {
	const scrollRef = useRef();
	const isGet = useRef(true);
	const filterFieldRef = useRef();
	const { request, cancel } = createRequest();
	// 因为pagination、sorter改变后就要调用接口，useState是异步，导致getList访问的是旧值，所以用useRef。
	const pagination = useRef(getPagination());
	const sorter = useRef(defaultSortOrder);
	// 因为有的页面需要展示排序和total，useRef的pagination、sorter会报错，页面上渲染的数据必须是useState的数据，所以此处做备份展示
	const [paramsData, setParamsData] = useState({
		...defaultSortOrder,
		...(getPagination() || {}),
	});
	const [list, setList] = useState([]);
	const [loading, setLoading] = useState(false);

	const getList = useCallback(async () => {
		try {
			setLoading(true);
			const filterParams = filterFieldRef?.current?.getApiParams?.();
			const params = {
				page: pagination?.current?.current,
				limit: pagination?.current?.pageSize,
				filter: {
					...filterParams,
				},
			};
			if (sorter?.current?.column) {
				params.sorter = sorter?.current;
			}
			let res = {};
			let newApi = api;
			if (newApi) {
				cancel(newApi);
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
			setParamsData((v) => ({ ...v, total }));
			setList(list);
			setLoading(false);
			return list;
		} catch {
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [api, getData, method]);

	const handleSearch = async () => {
		pagination.current.current = 1;
		setParamsData((v) => ({ ...v, current: 1 }));
		return await getList();
	};

	const handleChangePagination = (page) => {
		pagination.current.current = page;
		setParamsData((v) => ({ ...v, current: page }));
		getList();
		scrollRef.current.scrollTo({
			top: 0,
		});
	};

	const onChange = (paginationObj, filters, { order, columnKey }) => {
		const orderObj = {
			[descend]: des,
			[ascend]: asc,
		};
		pagination.current.current = paginationObj.current;
		sorter.current.sort = orderObj[order];
		sorter.current.column = columnKey;
		setParamsData((v) => ({ ...v, ...sorter.current, ...pagination.current }));
		getList();
	};
	return {
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
	};
};
export default useContent;
