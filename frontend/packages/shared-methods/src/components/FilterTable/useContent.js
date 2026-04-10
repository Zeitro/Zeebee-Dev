import {useCallback, useRef, useState} from 'react';
import {createRequest} from '../../lib/fetch';
import {getPagination} from '../../lib/common';
import {asc, ascend, des, descend} from '@repo/shared-methods/enum';

const defaultSortOrder = {
	// column: 'address',
	sort: des,
};
// FilterTable与FilterTableCard共用的方法
const useContent = ({
						getData = () => {
	                    },
						api,
						method = 'POST',
						hasPaginationRequest,
						settingPagination = {},
						hasInitRequest,
					}) => {
	const scrollRef = useRef();
	const isGet = useRef(hasInitRequest);
	const filterFieldRef = useRef();
	const { request, cancel } = createRequest();

	// 保留 ref 以便 getList 读取最新值（避免异步 state 导致的旧值问题），
	// 但提供一套封装方法同步更新 ref 与用于渲染的 state（paramsData）
	// const pagination = useRef(getPagination());
	// const sorter = useRef(defaultSortOrder);
	const [paramsData, setParamsData] = useState({
		...defaultSortOrder,
		...(getPagination(settingPagination) || {}),
	});
	const [list, setList] = useState([]);
	const [loading, setLoading] = useState(false);

	const getList = useCallback(async () => {
		try {
			const {page, limit, sorter} = modifyParams;
			setLoading(true);
			const filterParams = filterFieldRef?.current?.getApiParams?.();

			const params = {
				// page: pagination?.current?.current,
				// limit: pagination?.current?.pageSize,
				page: page || paramsData?.current,
				limit: limit || paramsData?.pageSize,
				filter: {
					...filterParams,
				},
				sorter: sorter || {
					column: paramsData?.column,
					sort: paramsData?.sort,
				},
				...modifyParams, // 传入需要直接修改的参数
			};
			if (!params?.sorter?.column) {
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
			// 使用封装的设置函数，保持 ref 与 state 同步
			// pagination.current.total = total;
			setParamsData((v) => ({ ...v, total }));
			setList(list);
			setLoading(false);
			return list;
		} catch {
			setLoading(false);
		}
		},
		[api, getData, method, paramsData],
	);

	const handleSearch = async (modifyParams = {}) => {
		// 重置到第一页并请求
		const {page, limit, sorter} = modifyParams;
		// pagination.current.current = 1;
		let info = {current: page || 1};
		if (limit) {
			info.pageSize = limit;
		}
		if (sorter) {
			info.column = sorter?.column;
			info.sort = sorter?.sort;
		}
		setParamsData((v) => ({...v, ...info}));
		return await getList({...modifyParams, page: page || 1});
	};

	const handleChangePagination = (page) => {
		// pagination.current.current = page;
		setParamsData((v) => ({ ...v, current: page }));
		getList({page});
		if (scrollRef.current && typeof scrollRef.current.scrollTo === 'function') {
			scrollRef.current.scrollTo({
				top: 0,
			});
		}
	};

	const onChange = (paginationObj, filters, { order, columnKey }) => {
		// 不请求分页接口，只更新 paramsData（用于显示）
		if (!hasPaginationRequest) {
			setParamsData((v) => {
				const next = {...v};
				if (paginationObj.current !== v.current) {
					next.current = paginationObj.current;
				} else {
					next.current = 1;
				}
				return next;
			});
			return;
		}

		const orderObj = {
			[descend]: des,
			[ascend]: asc,
		};

		// 同步更新 ref 与 paramsData
		// pagination.current.current = paginationObj.current;
		let sorter = {
			column: columnKey,
			sort: orderObj[order],
		};
		setParamsData((v) => ({...v, ...sorter, current: paginationObj.current}));
		getList({sorter, page: paginationObj.current});
	};
	return {
		filterFieldRef,
		list,
		setList,
		loading,
		onChange,
		// 保持向后兼容，返回原来的 pagination ref
		// pagination,
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
