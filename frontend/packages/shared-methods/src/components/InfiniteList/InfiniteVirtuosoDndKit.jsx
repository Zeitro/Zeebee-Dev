'use client';
import React, { useCallback, useMemo, useRef, useState } from 'react';
import { Virtuoso } from 'react-virtuoso';
import { DndContext, closestCenter } from '@dnd-kit/core';
import { SortableContext, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { getPagination } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import useContent from './useContent';
import styles from './index.module.scss';

function SortableItem({ id, item, index, renderItem, style, itemHeight }) {
	const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
		id,
	});
	return renderItem(item, index, {
		ref: setNodeRef,
		...attributes,
		...listeners,
		style: {
			...style,
			transform: CSS.Transform.toString(transform),
			transition,
			minHeight: itemHeight,
			zIndex: isDragging ? 9999 : 'auto',
			userSelect: 'none',
		},
	});
}
// 虚拟列表 + 可拖拽
export function InfiniteVirtuosoDndKit({
	itemHeight = 160,
	renderItem,
	enableDrag,
	className,
	boxClassName,
	getData,
	params,
	setItems,
	items, //多个时为对象，(setItems/items/itemKey/multiple只有一个虚拟列表时，这些不传)
	itemKey,
	multiple = false,
	setAllTotalData,
	onDragEnd,
}) {
	const [list, setList] = useState([]);
	const [loading, setLoading] = useState(false);
	const [hasMore, setHasMore] = useState(true);
	const [total, setTotal] = useState(0);
	const pagination = useRef(getPagination());
	const newList = useMemo(() => {
		return multiple ? items[itemKey] || [] : list;
	}, [items, multiple, itemKey, list]);

	// multiple为false时使用
	const { handleDragMove, handleDragEnd, sensors } = useContent({ items, setItems, onDragEnd });

	const getList = useCallback(
		async (isInit) => {
			try {
				if (multiple && !itemKey) return;
				setLoading(true);
				const { total, list } = await getData({
					page: pagination?.current?.current,
					limit: pagination?.current?.pageSize,
					...(params || {}),
					itemKey,
				});
				if (multiple) {
					setItems((v) => ({
						...v,
						[itemKey]: isInit ? list : [...v[itemKey], ...list],
					}));
				} else {
					setList((v) => (isInit ? list : [...v, ...list]));
				}
				setAllTotalData((v) => ({ ...v, [itemKey]: total }));
				setTotal(total);
				setLoading(false);
			} catch (error) {
				console.log(error);
				setLoading(false);
			}
		},
		[multiple, itemKey, params],
	);

	const loadMore = () => {
		const arr = multiple ? items[itemKey] : list;
		if (arr.length < total && arr.length > 0) {
			pagination.current.current += 1;
			getList();
		} else {
			setHasMore(false);
		}
	};

	const handleSearch = useCallback(() => {
		pagination.current.current = 1;
		getList();
	}, [getList]);

	useDeferredSyncEffect({
		deferred: () => {
			handleSearch();
		},
		deps: [handleSearch],
		isUseMicrotask: true,
	});

	const virtuosoSort = (
		<SortableContext items={newList?.map((i) => i.id)} strategy={verticalListSortingStrategy}>
			<Virtuoso
				data={newList}
				className={[styles.virtuosoList, className].join(' ')}
				itemContent={(index, item, style) => (
					<SortableItem
						id={item.id}
						item={item}
						index={index}
						renderItem={renderItem}
						style={style}
						itemHeight={itemHeight}
						disabled={!enableDrag}
					/>
				)}
				endReached={() => {
					if (hasMore && !loading) loadMore();
				}}
				overscan={2}
				components={{
					Footer: () =>
						loading ? (
							<div style={{ textAlign: 'center', padding: 8 }}>Loading...</div>
						) : null,
				}}
			/>
		</SortableContext>
	);

	return (
		<div className={['h-full', boxClassName].join(' ')}>
			{multiple ? (
				virtuosoSort
			) : (
				<DndContext
					sensors={sensors}
					collisionDetection={closestCenter}
					onDragEnd={enableDrag ? handleDragEnd : undefined}
					onDragMove={handleDragMove}
				>
					{virtuosoSort}
				</DndContext>
			)}
		</div>
	);
}
export default InfiniteVirtuosoDndKit;
