/**
 * FilterTable / FilterTableCard 共用：按 tableType 渲染 NewTable 或 SortableTable。
 * sortable 时先 setList 再调 onSortEnd，保证当前页顺序立刻更新。
 */
'use client';

import React from 'react';
import { NewTable } from '../NewTable';
import { SortableTable, TABLE_TYPE } from '../SortableTable';

/**
 * @param {object} props
 * @param {'default'|'sortable'} [props.tableType='default']
 * @param {(nextList: object[], meta: { fromIndex: number, toIndex: number }) => void} [props.onSortEnd]
 * @param {(nextList: object[]) => void} [props.setList]
 * @param {boolean} [props.showDragHandle]
 */
export function ListTable({
	tableType = TABLE_TYPE.DEFAULT,
	onSortEnd,
	setList,
	showDragHandle,
	...props
}) {
	if (tableType === TABLE_TYPE.SORTABLE) {
		const handleSortEnd = (nextList, meta) => {
			setList?.(nextList);
			onSortEnd?.(nextList, meta);
		};
		return (
			<SortableTable {...props} showDragHandle={showDragHandle} onSortEnd={handleSortEnd} />
		);
	}
	return <NewTable {...props} />;
}
