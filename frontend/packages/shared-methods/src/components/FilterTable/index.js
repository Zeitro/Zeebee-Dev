'use client';
import React, { useEffect, forwardRef, useImperativeHandle } from 'react';
import styles from './index.module.scss';
import { useDisUpdateEffect } from '@repo/shared-methods/hooks';
import { FilterTag } from '../FilterField/filterTag.js';
import { FilterField } from '../FilterField/index';
import { NewTable } from '../NewTable';
import useContent from './useContent';
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
			rowProps = {}, // 过滤搜索栏的Row配置
			filterName = '',
			hasFilterField = true, // true有过滤搜索栏
			hasInitRequest = true, //true初始化请求接口，当有dataSource时默认‘不请求’
			filterWrapClassName,
			getData = () => {}, //请求接口的方法（无api时使用）
			hasPaginationRequest = true, // true需要请求分页接口，false不请求
			api,
			method = 'POST',
			settingPagination, // 初始化修改分页请求参数
			pageSizeStorageKey, // localStorage 键，默认全局 zb_table_page_size
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
			handleSearch,
			getList,
			isGetRef,
			setLoading,
			paramsData,
		} = useContent({
			api,
			method,
			getData,
			hasPaginationRequest,
			settingPagination,
			hasInitRequest,
			pageSizeStorageKey,
		});

		useEffect(() => {
			if (dataSource && Array.isArray(dataSource)) {
				setList(dataSource);
			} else if (isGetRef.current) {
				getList();
				isGetRef.current = false;
			}
			// eslint-disable-next-line react-hooks/exhaustive-deps
		}, [dataSource, getList]);

		useDisUpdateEffect(() => {
			setLoading(dataSourceLoading);
		}, [dataSourceLoading]);

		useImperativeHandle(ref, () => ({
			handleSearch,
			getList,
			getApiParams: () => filterFieldRef?.current?.getApiParams(),
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
					dataSource={list}
					loading={loading}
					columns={columns}
					scroll={scroll}
					onChange={onChange}
					{...props}
					pagination={props?.pagination ?? paramsData}
					className={[styles.filter_table_content, props?.className].join(' ')}
				/>
			</div>
		);
	},
);
FilterTable.displayName = 'FilterTable';
