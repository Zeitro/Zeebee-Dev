'use client';
import React, { createContext, useContext, useMemo } from 'react';

/**
 * FilterTable / FilterTableCard 的 UI 注入上下文。
 * platform-core 不内置 NewTable / NewPagination；由应用在根布局提供。
 */

const FilterTableUiContext = createContext({
	TableComponent: null,
	PaginationComponent: null,
});

/**
 * @param {object} props
 * @param {React.ComponentType} [props.TableComponent] 应用侧表格（如 NewTable）
 * @param {React.ComponentType} [props.PaginationComponent] 应用侧分页（如 NewPagination）
 * @param {React.ReactNode} props.children
 */
export function FilterTableProvider({ TableComponent, PaginationComponent, children }) {
	const value = useMemo(
		() => ({ TableComponent, PaginationComponent }),
		[TableComponent, PaginationComponent],
	);
	return <FilterTableUiContext.Provider value={value}>{children}</FilterTableUiContext.Provider>;
}

/**
 * 解析最终表格组件：props.TableComponent 优先，否则用 Context。
 * @param {React.ComponentType|undefined|null} TableComponent
 * @returns {React.ComponentType}
 */
export function useResolvedTableComponent(TableComponent) {
	const { TableComponent: contextTable } = useContext(FilterTableUiContext);
	const Table = TableComponent ?? contextTable;
	if (!Table) {
		throw new Error(
			'FilterTable: pass TableComponent or wrap with <FilterTableProvider TableComponent={NewTable} />',
		);
	}
	return Table;
}

/**
 * 解析最终分页组件：props.PaginationComponent 优先，否则用 Context。
 * @param {React.ComponentType|undefined|null} PaginationComponent
 * @returns {React.ComponentType}
 */
export function useResolvedPaginationComponent(PaginationComponent) {
	const { PaginationComponent: contextPagination } = useContext(FilterTableUiContext);
	const Pagination = PaginationComponent ?? contextPagination;
	if (!Pagination) {
		throw new Error(
			'FilterTableCard: pass PaginationComponent or wrap with <FilterTableProvider PaginationComponent={NewPagination} />',
		);
	}
	return Pagination;
}
