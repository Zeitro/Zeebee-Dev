'use client';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Divider, Spin } from 'antd';
import { getPagination } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect, useElementSize } from '@repo/shared-methods/hooks';
import { NoData } from '../NoData';

/**
 * 虚拟无限滚动列表：分页拉取、多列/单列数据、底部加载与空态（与 DnD 版共用）
 */
export function useInfiniteVirtuoso({
	getData,
	params,
	multiple = false,
	items,
	setItems,
	itemKey,
}) {
	const [rect, listBoxRef] = useElementSize();
	const pagination = useRef(getPagination());
	const [list, setList] = useState([]);
	const [loading, setLoading] = useState(true);
	const [hasMore, setHasMore] = useState(true);
	const [total, setTotal] = useState(0);
	const [isScroll, setIsScroll] = useState(false);

	// 避免父组件内联 getData / params 导致 getList → handleSearch 变化，触发 useDeferredSyncEffect 无限请求
	const getDataRef = useRef(getData);
	const paramsRef = useRef(params);
	useEffect(() => {
		getDataRef.current = getData;
	}, [getData]);
	useEffect(() => {
		paramsRef.current = params;
	}, [params]);

	const newList = useMemo(() => {
		return multiple ? items?.find((i) => i.value === itemKey)?.list || [] : list;
	}, [items, multiple, itemKey, list]);

	const getList = useCallback(
		async (isInit) => {
			try {
				if (multiple && !itemKey) return;
				setLoading(true);
				const { total: nextTotal, list: nextList } = await getDataRef.current({
					page: pagination?.current?.current,
					limit: pagination?.current?.pageSize,
					...(paramsRef.current || {}),
					itemKey,
				});
				if (multiple) {
					setItems((v) => {
						return v?.map((i) => {
							if (i?.value === itemKey) {
								i.list = isInit ? nextList : [...(i?.list || []), ...nextList];
								i.total = nextTotal;
							}
							return i;
						});
					});
				} else {
					setList((v) => (isInit ? nextList : [...v, ...nextList]));
				}
				setTotal(nextTotal);
				setLoading(false);
			} catch (error) {
				console.log(error);
				setLoading(false);
			}
		},
		[multiple, itemKey, setItems],
	);

	const loadMore = useCallback(() => {
		const arr = multiple ? items?.find((i) => i.value === itemKey)?.list || [] : list;

		if (arr.length < total && arr.length > 0) {
			pagination.current.current += 1;
			getList();
		} else {
			setHasMore(false);
		}
	}, [getList, items, itemKey, list, multiple, total]);

	const handleSearch = useCallback(() => {
		pagination.current.current = 1;
		setHasMore(true);
		getList(true);
	}, [getList]);

	const handleFooter = useCallback(
		({ loadingClassName, loadingStyle = {} } = {}) => {
			if (loading) {
				return (
					<div
						className={['flex items-center justify-center', loadingClassName]
							.filter(Boolean)
							.join(' ')}
						style={{ height: 48, ...loadingStyle }}
					>
						<Spin size="small" />
					</div>
				);
			}
			if (newList?.length > 0 && newList.length === total && isScroll) {
				return (
					<Divider plain>
						<span className="text-base-label">No more data</span>
					</Divider>
				);
			}
			if (!newList?.length) {
				return <NoData height="90%" />;
			}
		},
		[loading, newList, total, isScroll],
	);

	const handleEndReached = useCallback(
		(endIndex, { itemHeight } = {}) => {
			if (hasMore && !loading) loadMore();
			const rectHeight = rect?.height;
			if (itemHeight && rectHeight) {
				const maxNum = Math.floor(rectHeight / itemHeight);
				if (endIndex + 1 > maxNum) {
					setIsScroll(true);
				}
			}
		},
		[hasMore, loading, loadMore, rect],
	);

	useDeferredSyncEffect({
		deferred: () => {
			handleSearch();
		},
		deps: [handleSearch],
		isUseMicrotask: true,
	});

	return {
		rect,
		listBoxRef,
		list,
		setList,
		newList,
		loading,
		hasMore,
		total,
		isScroll,
		setIsScroll,
		loadMore,
		handleSearch,
		handleFooter,
		handleEndReached,
	};
}

export default useInfiniteVirtuoso;
