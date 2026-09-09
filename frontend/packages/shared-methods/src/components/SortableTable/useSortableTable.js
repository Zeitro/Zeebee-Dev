/**
 * 行拖拽排序：把当前页 dataSource 按手柄拖拽结果 arrayMove，再交给 onSortEnd。
 * 不改外部 state；受控列表由调用方（或 FilterTable 的 setList）更新。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { arrayMove } from '@dnd-kit/sortable';

/** dnd-kit modifier：只允许纵向位移，避免行跟着指针横向跑 */
function restrictToVerticalAxis({ transform }) {
	return { ...transform, x: 0 };
}

/**
 * @param {object} params
 * @param {object[]} [params.dataSource]
 * @param {string|function} params.rowKey
 * @param {(nextList: object[], meta: { fromIndex: number, toIndex: number }) => void} [params.onSortEnd]
 */
export function useSortableTable({ dataSource, rowKey, onSortEnd }) {
	const sensors = useSensors(
		useSensor(PointerSensor, {
			activationConstraint: { distance: 5 },
		}),
	);

	const getRowKey = useCallback(
		(record, index) => {
			if (typeof rowKey === 'function') return rowKey(record, index);
			return record?.[rowKey];
		},
		[rowKey],
	);

	/** antd 写入的 data-row-key 是字符串，须对齐否则 indexOf 对不上 */
	const sortableIds = useMemo(
		() =>
			(dataSource || []).map((record, index) => {
				const key = getRowKey(record, index);
				return key == null ? key : String(key);
			}),
		[dataSource, getRowKey],
	);

	const [activeId, setActiveId] = useState(null);

	const handleDragStart = useCallback(({ active }) => {
		setActiveId(active.id);
	}, []);

	const handleDragEnd = useCallback(
		({ active, over }) => {
			setActiveId(null);
			if (!over || String(active.id) === String(over.id)) return;
			const fromIndex = sortableIds.indexOf(String(active.id));
			const toIndex = sortableIds.indexOf(String(over.id));
			if (fromIndex < 0 || toIndex < 0) return;
			const nextList = arrayMove(dataSource || [], fromIndex, toIndex);
			onSortEnd?.(nextList, { fromIndex, toIndex });
		},
		[dataSource, onSortEnd, sortableIds],
	);

	const handleDragCancel = useCallback(() => {
		setActiveId(null);
	}, []);

	const activeRecord = useMemo(() => {
		if (activeId == null) return null;
		return (dataSource || []).find(
			(record, index) => String(getRowKey(record, index)) === String(activeId),
		);
	}, [activeId, dataSource, getRowKey]);

	return {
		sensors,
		sortableIds,
		handleDragStart,
		handleDragEnd,
		handleDragCancel,
		modifiers: [restrictToVerticalAxis],
		activeRecord,
	};
}
