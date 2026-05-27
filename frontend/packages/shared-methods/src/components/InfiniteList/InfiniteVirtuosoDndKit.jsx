'use client';
import React, { useState } from 'react';
import { Virtuoso } from 'react-virtuoso';
import { DndContext, closestCenter, DragOverlay } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import useContent from './useContent';
import useInfiniteVirtuoso from './useInfiniteVirtuoso';
import styles from './index.module.scss';

/**
 * 可拖拽的虚拟列表
 * 跨列拖拽必传：items、setItems、itemKey、multiple=true （单列不需要传）
 * items：[ { value: 'a', label:'A', list:[], total:0 } ]
 * itemKey: 'a'
 * 单列拖拽排序
 */
export function InfiniteVirtuosoDndKit({
	itemHeight,
	minItemHeight,
	renderItem,
	enableDrag,
	className,
	boxClassName,
	getData,
	params,
	setItems,
	items,
	itemKey,
	multiple = false,
	isSort = false,
	onDragEnd,
	...props
}) {
	const [activeCard, setActiveCard] = useState(null);
	const { listBoxRef, list, setList, newList, handleFooter, handleEndReached } =
		useInfiniteVirtuoso({
			getData,
			params,
			multiple,
			items,
			setItems,
			itemKey,
		});

	const { handleDragEnd, sensors, EmptyColumnDroppable, SortableItem } = useContent({
		onDragEnd,
		isSort,
		list,
		setList,
	});

	const isFixedHeight = itemHeight != null;
	const virtuosoHeightProps = isFixedHeight
		? { fixedItemHeight: itemHeight }
		: minItemHeight != null
			? { defaultItemHeight: minItemHeight }
			: {};

	const virtuosoSort = (
		<SortableContext items={newList?.map((i) => i.id)} strategy={verticalListSortingStrategy}>
			{/* 只有跨多列时，无数据时需要EmptyColumnDroppable占位符 */}
			{(multiple ? newList.length > 0 : true) ? (
				<Virtuoso
					data={newList}
					className={[styles.virtuosoList, className].filter(Boolean).join(' ')}
					style={{ height: '100%', ...(props?.style || {}) }}
					{...virtuosoHeightProps}
					itemContent={(index, item, style) => (
						<SortableItem
							id={item.id}
							item={item}
							index={index}
							renderItem={renderItem}
							style={style}
							itemHeight={itemHeight}
							minItemHeight={minItemHeight}
							isFixedHeight={isFixedHeight}
						/>
					)}
					endReached={(endIndex) =>
						handleEndReached(
							endIndex,
							isFixedHeight ? { itemHeight } : {},
						)
					}
					overscan={2}
					components={{
						Footer: () => handleFooter(),
					}}
					{...props}
				/>
			) : (
				<EmptyColumnDroppable colValue={itemKey}>
					{handleFooter({ loadingStyle: { height: '90%' } })}
				</EmptyColumnDroppable>
			)}
		</SortableContext>
	);

	return (
		<div
			className={['h-full min-h-0 overflow-hidden', boxClassName].filter(Boolean).join(' ')}
			ref={listBoxRef}
		>
			{multiple ? (
				virtuosoSort
			) : (
				<>
					<DndContext
						sensors={sensors}
						collisionDetection={closestCenter}
						onDragEnd={(e) => {
							if (!enableDrag) return;
							setActiveCard(null);
							handleDragEnd(e);
						}}
						onDragStart={({ active }) => {
							if (!enableDrag) return;
							setActiveCard(active?.data?.current);
						}}
					>
						{virtuosoSort}
					</DndContext>
					<DragOverlay>{activeCard ? renderItem(activeCard, 0, {}) : null}</DragOverlay>
				</>
			)}
		</div>
	);
}
export default InfiniteVirtuosoDndKit;
