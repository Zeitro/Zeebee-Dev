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
import { asc, des, ascend, descend } from '@/enum/common';
import { useDisUpdateEffect } from '@repo/shared-methods/hooks';
import { FilterTag } from '../FilterField/filterTag.js';
import { FilterField } from '../FilterField/index';
import { NewTable } from '../NewTable';
import { Cards } from './cards.js';
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
		const scrollRef = useRef();
		const isGet = useRef(true);
		const filterFieldRef = useRef();
		const { request, cancel } = createRequest();
		const [pagination, setPagination] = useState(getPagination());
		const [sorter, setSorter] = useState({
			// column: 'address',
			sort: des,
		});
		const [list, setList] = useState([]);
		const [loading, setLoading] = useState(false);

		const getList = useCallback(async () => {
			try {
				setLoading(true);
				const filterParams = filterFieldRef?.current?.getApiParams?.();
				const params = {
					page: pagination?.current,
					limit: pagination?.pageSize,
					filter: {
						...filterParams,
					},
				};
				if (sorter?.column) {
					params.sorter = sorter;
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
				setPagination((v) => ({ ...v, total }));
				setList(list);
				setLoading(false);
			} catch {
				setLoading(false);
			}
			// eslint-disable-next-line react-hooks/exhaustive-deps
		}, [api, getData, method, pagination?.current, pagination?.pageSize, sorter]);

		const handleSearch = () => {
			setPagination((v) => ({ ...v, current: 1 }));
			getList();
		};

		const handleChangePagination = (page) => {
			setPagination((v) => {
				const next = { ...v, current: page };
				return next;
			});
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
			setSorter({
				column: columnKey,
				sort: orderObj[order],
			});

			setPagination((v) => ({
				...v,
				current: paginationObj.current,
			}));
			getList();
		};
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
							pagination={props?.pagination ?? pagination}
							className={['flex-1 overflow-hidden', props?.className].join(' ')}
						/>
					);
				case 'card':
					return (
						<Cards
							pagination={pagination}
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
