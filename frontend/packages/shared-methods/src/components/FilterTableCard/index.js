'use client';
import React, { useEffect, forwardRef, useImperativeHandle } from 'react';
import { useDisUpdateEffect } from '@repo/shared-methods/hooks';
import { FilterTag } from '../FilterField/filterTag.js';
import { FilterField } from '../FilterField/index';
import { NewTable } from '../NewTable';
import { Cards } from './cards.js';
import styles from './index.module.scss';
import useContent from '../FilterTable/useContent';

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
			rowProps = {}, // 过滤搜索栏的Row配置
			filterName = '',
			hasFilterField = true, // 有过滤搜索栏
			hasInitRequest = true, //true初始化请求接口，当有dataSource时默认‘不请求’
			filterWrapClassName,
			getData = () => {}, //请求接口的方法（无api时使用）
			hasPaginationRequest = true, // true需要请求分页数据，false不请求分页接口
			api,
			method = 'POST',
			renderItem = () => {}, // card模式下使用
			itemClassName,
			typeKey = 'table', // table/card
			cardTitle,
			settingPagination, // 初始化修改分页请求参数
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
			scrollRef,
			handleChangePagination,
		} = useContent({
			api,
			method,
			getData,
			hasPaginationRequest,
			settingPagination,
			hasInitRequest,
		});

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
			} else if (isGetRef.current) {
				// isGetRef代表需要请求接口
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
