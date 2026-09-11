/**
 * 行拖拽排序 Table：交互/样式对齐 NewTable；默认前置手柄列。
 * 受控：拖拽结束后通过 onSortEnd 回传新 dataSource，不在内部持久化。
 */
'use client';

import React, { useMemo } from 'react';
import { closestCenter, DndContext, DragOverlay } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { NewTable } from '../NewTable';
import { DragHandle, SortableRow } from './SortableRow';
import { useSortableTable } from './useSortableTable';

export const TABLE_TYPE = {
	DEFAULT: 'default',
	SORTABLE: 'sortable',
};

/**
 * @param {object} props
 * @param {object[]} [props.dataSource]
 * @param {object[]} [props.columns]
 * @param {string|function} [props.rowKey='id']
 * @param {(nextList: object[], meta: { fromIndex: number, toIndex: number }) => void} [props.onSortEnd]
 * @param {boolean} [props.showDragHandle=true] 是否自动插入手柄列；false 时需自行在 columns 里渲染 DragHandle
 * @param {object} [props.components] 与 antd Table `components` 合并，body.row 固定为 SortableRow
 */
export const SortableTable = ({
	dataSource,
	columns,
	pagination,
	rowKey = 'id',
	scroll,
	onSortEnd,
	showDragHandle = true,
	components,
	...props
}) => {
	const {
		sensors,
		sortableIds,
		handleDragStart,
		handleDragEnd,
		handleDragCancel,
		modifiers,
		activeRecord,
	} = useSortableTable({
		dataSource,
		rowKey,
		onSortEnd,
	});

	const mergedColumns = useMemo(() => {
		const next = columns || [];
		if (!showDragHandle) return next;
		return [
			{
				key: '__dragHandle',
				title: '',
				width: 48,
				align: 'center',
				render: () => <DragHandle />,
			},
			...next,
		];
	}, [columns, showDragHandle]);

	const tableComponents = useMemo(
		() => ({
			...components,
			body: {
				...components?.body,
				row: SortableRow,
			},
		}),
		[components],
	);

	return (
		<DndContext
			sensors={sensors}
			modifiers={modifiers}
			collisionDetection={closestCenter}
			onDragStart={handleDragStart}
			onDragEnd={handleDragEnd}
			onDragCancel={handleDragCancel}
		>
			<SortableContext items={sortableIds} strategy={verticalListSortingStrategy}>
				<NewTable
					dataSource={dataSource}
					columns={mergedColumns}
					pagination={pagination}
					rowKey={rowKey}
					scroll={scroll}
					{...props}
					components={tableComponents}
				/>
			</SortableContext>
			<DragOverlay dropAnimation={null}>
				{activeRecord ? (
					<div className="px-3 py-2 bg-white rounded-md shadow-lg border border-base-disabled cursor-move">
						{activeRecord.name ?? activeRecord.title ?? String(activeRecord.id ?? '')}
					</div>
				) : null}
			</DragOverlay>
		</DndContext>
	);
};

export { DragHandle, SortableRow };
