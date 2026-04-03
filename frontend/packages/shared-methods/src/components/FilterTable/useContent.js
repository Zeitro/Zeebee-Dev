import { useRef, useState, useCallback } from 'react';
import { createRequest } from '../../lib/fetch';
import { getPagination } from '../../lib/common';
import { asc, des, ascend, descend } from '@repo/shared-methods/enum';
const defaultSortOrder = {
	// column: 'address',
	sort: des,
};
// FilterTable与FilterTableCard共用的方法
const useContent = ({
	getData = () => {},
	api,
	method = 'POST',
	hasPaginationRequest,
	settingPagination = {},
}) => {
	const scrollRef = useRef();
	const isGetRef = useRef(true); // 是否请求接口
	const filterFieldRef = useRef();
	const { request, cancel } = createRequest();
	// 因为pagination、sorter改变后就要调用接口，useState是异步，导致getList访问的是旧值，所以用useRef。
	const pagination = useRef(getPagination(settingPagination));
	const sorter = useRef(defaultSortOrder);
	// 因为有的页面需要展示排序和total，useRef的pagination、sorter会报错，页面上渲染的数据必须是useState的数据，所以此处做备份展示
	const [paramsData, setParamsData] = useState({
		...defaultSortOrder,
		...(getPagination() || {}),
	});
	const [list, setList] = useState([]);
	const [loading, setLoading] = useState(false);

	const getList = useCallback(
		async (modifyParams = {}) => {
			try {
				setLoading(true);
				const filterParams = filterFieldRef?.current?.getApiParams?.();
				const params = {
					page: pagination?.current?.current,
					limit: pagination?.current?.pageSize,
					filter: {
						...filterParams,
					},
					sorter: sorter?.current,
					...modifyParams, // 传入需要直接修改的参数
				};
				if (!sorter?.current?.column) {
					delete params.sorter;
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
		},
		// eslint-disable-next-line react-hooks/exhaustive-deps
		[api, getData, method],
	);

	const handleSearch = async (modifyParams) => {
		pagination.current.current = 1;
		setParamsData((v) => ({ ...v, current: 1 }));
		return await getList(modifyParams);
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
		// 不请求分页接口
		if (!hasPaginationRequest) {
			setParamsData((v) => {
				if (paginationObj.current !== v.current) {
					v.current = paginationObj.current;
				} else {
					v.current = 1;
				}
				return { ...v };
			});
			return;
		}
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
		isGetRef,
		setLoading,
		paramsData,
		scrollRef,
		handleChangePagination,
	};
};
export default useContent;
