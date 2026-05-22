'use client';
import React from 'react';
import { Virtuoso } from 'react-virtuoso';
import styles from './index.module.scss';
import useInfiniteVirtuoso from './useInfiniteVirtuoso';

/**
 * 虚拟无限滚动列表（无拖拽）
 * 多列模式：items、setItems、itemKey、multiple=true
 * items：[ { value: 'a', label:'A', list:[], total:0 } ]
 */
export function InfiniteVirtuoso({
	itemHeight = 160,
	renderItem,
	className,
	boxClassName,
	getData,
	params,
	setItems,
	items,
	itemKey,
	multiple = false,
	...props
}) {
	const { listBoxRef, newList, handleFooter, handleEndReached } =
		useInfiniteVirtuoso({
			getData,
			params,
			multiple,
			items,
			setItems,
			itemKey,
		});

	return (
		<div className={['h-full', boxClassName].filter(Boolean).join(' ')} ref={listBoxRef}>
			<Virtuoso
				data={newList}
				className={[styles.virtuosoList, className].filter(Boolean).join(' ')}
				itemContent={(index, item) =>
					renderItem(item, index, {
						style: { minHeight: itemHeight },
					})
				}
				endReached={(endIndex) =>
					handleEndReached(endIndex, { itemHeight })
				}
				overscan={2}
				components={{
					Footer: () => handleFooter(),
				}}
				{...props}
			/>
		</div>
	);
}

export default InfiniteVirtuoso;
