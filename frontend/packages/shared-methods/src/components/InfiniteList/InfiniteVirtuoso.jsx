'use client';
import React from 'react';
import { Virtuoso } from 'react-virtuoso';
import styles from './index.module.scss';
import useInfiniteVirtuoso from './useInfiniteVirtuoso';

/**
 * 虚拟无限滚动列表（无拖拽）
 * - DOM：Virtuoso 必须限定高度（style.height + 父级 min-h-0），只渲染可视区 + overscan，与 data 条数无关
 * - 数据：getData 可一次返回全量 { list, total }；enableLoadMore=false 时不分页追加
 * - itemHeight：固定等高（Virtuoso fixedItemHeight）
 * - minItemHeight：不等高时的默认最小高度
 */
function getVirtuosoItemConfig({ itemHeight, minItemHeight }) {
	if (itemHeight != null) {
		return {
			itemStyle: {
				height: itemHeight,
				minHeight: itemHeight,
				boxSizing: 'border-box',
			},
			virtuosoProps: { fixedItemHeight: itemHeight },
			scrollMetric: { itemHeight },
		};
	}
	return {
		itemStyle: minItemHeight != null ? { minHeight: minItemHeight } : undefined,
		virtuosoProps: minItemHeight != null ? { defaultItemHeight: minItemHeight } : {},
		scrollMetric: {},
	};
}

export function InfiniteVirtuoso({
	itemHeight,
	minItemHeight,
	enableLoadMore = true,
	pageSize = 10,
	overscan = 2,
	renderItem,
	className,
	boxClassName,
	getData,
	params,
	setItems,
	items,
	itemKey,
	multiple = false,
	style,
	...props
}) {
	const { listBoxRef, newList, handleFooter, handleEndReached } = useInfiniteVirtuoso({
		getData,
		params,
		pageSize,
		enableLoadMore,
		multiple,
		items,
		setItems,
		itemKey,
	});

	const { itemStyle, virtuosoProps, scrollMetric } = getVirtuosoItemConfig({
		itemHeight,
		minItemHeight,
	});

	return (
		<div
			className={['h-full min-h-0 overflow-hidden', boxClassName].filter(Boolean).join(' ')}
			ref={listBoxRef}
		>
			<Virtuoso
				data={newList}
				className={[styles.virtuosoList, className].filter(Boolean).join(' ')}
				style={{ height: '100%', ...style }}
				itemContent={(index, item) =>
					renderItem(item, index, itemStyle ? { style: itemStyle } : {})
				}
				endReached={(endIndex) => handleEndReached(endIndex, scrollMetric)}
				overscan={overscan}
				components={{
					Footer: () => handleFooter(),
				}}
				{...virtuosoProps}
				{...props}
			/>
		</div>
	);
}

export default InfiniteVirtuoso;
