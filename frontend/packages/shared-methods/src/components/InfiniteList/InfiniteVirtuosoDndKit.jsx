'use client';
import React, { useCallback, useMemo, useRef, useState } from 'react';
import { Virtuoso } from 'react-virtuoso';
import { DndContext, closestCenter, DragOverlay } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { getPagination } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect, useElementSize } from '@repo/shared-methods/hooks';
import useContent from './useContent';
import styles from './index.module.scss';
import { Divider, Spin } from 'antd';
import { NoData } from '../NoData';

/**
 * 可拖拽的虚拟列表
 * 跨列拖拽必传：items、setItems、itemKey、multiple=true （单列不需要传）
 * items：[ { value: 'a', label:'A', list:[], total:0 } ]
 * itemKey: 'a'
 * 单列拖拽排序
 */
export function InfiniteVirtuosoDndKit({
	itemHeight = 160,
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
	const [rect, listBoxRef] = useElementSize();
	const pagination = useRef(getPagination());
	const [activeCard, setActiveCard] = useState(null);
	const [list, setList] = useState([]);
	const [loading, setLoading] = useState(true);
	const [hasMore, setHasMore] = useState(true);
	const [total, setTotal] = useState(0);
	const [isScroll, setIsScroll] = useState(false);

	const newList = useMemo(() => {
		return multiple ? items?.find((i) => i.value === itemKey)?.list || [] : list;
	}, [items, multiple, itemKey, list]);

	//  multiple：false时调用list、 setList
	const { handleDragEnd, sensors, EmptyColumnDroppable, SortableItem } = useContent({
		onDragEnd,
		isSort,
		list,
		setList,
	});

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
					setItems((v) => {
						return v?.map((i) => {
							if (i?.value === itemKey) {
								i.list = isInit ? list : [...(i?.list || []), ...list];
								i.total = total;
							}
							return i;
						});
					});
				} else {
					setList((v) => (isInit ? list : [...v, ...list]));
				}
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
		const arr = multiple ? items?.find((i) => i.value === itemKey)?.list || [] : list;
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

	const handleFooter = ({ loadingClassName, loadingStyle = {} }) => {
		if (loading) {
			return (
				<div
					className={['flex items-center justify-center', loadingClassName].join(' ')}
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
	};

	useDeferredSyncEffect({
		deferred: () => {
			handleSearch();
		},
		deps: [handleSearch],
		isUseMicrotask: true,
	});

	const virtuosoSort = (
		<SortableContext items={newList?.map((i) => i.id)} strategy={verticalListSortingStrategy}>
			{/* 只有跨多列时，无数据时需要EmptyColumnDroppable占位符 */}
			{(multiple ? newList.length > 0 : true) ? (
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
					endReached={(e) => {
						if (hasMore && !loading) loadMore(e);
						// 有滚动条时才显示‘no more data’，不够满屏不显示
						if (itemHeight && rect?.height) {
							const maxNum = Math.floor(rect.height / itemHeight);
							if (e + 1 > maxNum) {
								setIsScroll(true);
							}
						}
					}}
					overscan={2}
					components={{
						Footer: handleFooter,
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
		<div className={['h-full', boxClassName].join(' ')} ref={listBoxRef}>
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
