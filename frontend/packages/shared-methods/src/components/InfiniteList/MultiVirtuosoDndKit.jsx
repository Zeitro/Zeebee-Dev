'use client';
import { DndContext, closestCenter, DragOverlay } from '@dnd-kit/core';
import InfiniteVirtuosoDndKit from './InfiniteVirtuosoDndKit';
import _ from 'lodash';
import { useState } from 'react';
import useContent from './useContent';

// 多列虚拟滚动+拖拽
export function MultiVirtuosoDndKit({
	items, // [ { value: 'a', label:'A', list:[], total:0 } ]
	setItems, 
	itemHeight = 160,
	enableDrag = false,
	renderHeader,
	renderFooter,
	renderItem,
	getData,
	className,
	boxClassName,
	parentClassName,
	columnClassName,
	isSort,
	onDragEnd,
}) {
	const { handleCrossColumnDragEnd, sensors } = useContent({ items, setItems, onDragEnd });
	const [activeCard, setActiveCard] = useState(null);
	return (
		<DndContext
			sensors={sensors}
			collisionDetection={closestCenter}
			onDragEnd={(e) => {
				if (!enableDrag) return;
				setActiveCard(null);
				handleCrossColumnDragEnd(e);
			}}
			onDragStart={({ active }) => {
				if (!enableDrag) return;
				setActiveCard(active?.data?.current);
			}}
		>
			<div className={['flex h-full', parentClassName].join(' ')}>
				{items.map((col, index) => (
					<div
						className={['h-full flex flex-col', columnClassName].join(' ')}
						key={index}
					>
						{renderHeader?.(col.value, col?.total, col)}

						<InfiniteVirtuosoDndKit
							items={items}
							setItems={setItems}
							multiple={true}
							itemKey={col.value}
							renderItem={renderItem}
							enableDrag={enableDrag} // false 时只虚拟滚动
							itemHeight={itemHeight}
							boxClassName={['flex-1', boxClassName].join(' ')}
							className={className}
							getData={getData}
							isSort={isSort}
							key={col.value}
						/>

						{renderFooter?.(col.value, col?.total, col)}
					</div>
				))}
				<DragOverlay>{activeCard ? renderItem(activeCard, 0, {}) : null}</DragOverlay>
			</div>
		</DndContext>
	);
}
