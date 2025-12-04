'use client';
import { DndContext, closestCenter } from '@dnd-kit/core';
import InfiniteVirtuosoDndKit from './InfiniteVirtuosoDndKit';
import _ from 'lodash';
import { useState } from 'react';
import useContent from './useContent';

// 多列虚拟滚动+拖拽
export function MultiVirtuosoDndKit({
	items, // {Quoting: [...], ...}
	setItems, // (newColumns) => void
	itemHeight = 160,
	height = 500,
	width = 400,
	enableDrag = false,
	renderHeader,
	renderFooter,
	renderItem,
	getData,
	className,
	boxClassName,
	parentClassName,
	columnClassName,
	onDragEnd,
}) {
	const [allTotalData, setAllTotalData] = useState({});
	const { handleDragMove, handleDragEnd, sensors } = useContent({ items, setItems, onDragEnd });

	return (
		<DndContext
			sensors={sensors}
			collisionDetection={closestCenter}
			onDragEnd={enableDrag ? handleDragEnd : undefined}
			onDragMove={handleDragMove}
		>
			<div className={['flex h-full', parentClassName].join(' ')}>
				{Object.keys(items).map((col, index) => (
					<div
						className={['h-full flex flex-col', columnClassName].join(' ')}
						key={index}
					>
						{renderHeader?.(col, allTotalData[col], items[col])}

						<InfiniteVirtuosoDndKit
							items={items}
							setItems={setItems}
							multiple={true}
							itemKey={col}
							renderItem={renderItem}
							enableDrag={enableDrag} // false 时只虚拟滚动
							height={height}
							width={width}
							itemHeight={itemHeight}
							boxClassName={['flex-1', boxClassName].join(' ')}
							className={className}
							getData={getData}
							setAllTotalData={setAllTotalData}
						/>

						{renderFooter?.(col, allTotalData[col], items[col])}
					</div>
				))}
			</div>
		</DndContext>
	);
}
