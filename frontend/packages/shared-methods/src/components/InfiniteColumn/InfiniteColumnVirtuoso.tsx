import React from 'react';
import { Virtuoso } from 'react-virtuoso';
import { Droppable, Draggable } from '@hello-pangea/dnd';
import type { DragInfo } from './dnd-types';

export interface InfiniteColumnVirtuosoProps<T> {
	items: T[];
	hasMore: boolean;
	loadMore: () => void;
	isLoading: boolean;
	itemHeight?: number;
	height?: number;
	width?: number | string;
	renderItem: (item: T, index: number) => React.ReactNode;
	draggable?: boolean;
	columnId?: string;
	onDragEnd?: (info: DragInfo) => void;
	getDraggableId?: (item: T, index: number) => string;
}

export function InfiniteColumnVirtuoso<T>({
	items,
	hasMore,
	loadMore,
	isLoading,
	itemHeight = 160,
	height = 500,
	width = 400,
	renderItem,
	draggable = false,
	columnId = '',
	onDragEnd,
	getDraggableId = (item: T, index: number) => String(index),
}: InfiniteColumnVirtuosoProps<T>) {
	if (!draggable) {
		return (
			<div style={{ width, height }}>
				<Virtuoso
					style={{ width: '100%', height: '100%' }}
					data={items}
					itemContent={(index: any, item: any) => renderItem(item, index)}
					endReached={() => {
						if (hasMore && !isLoading) loadMore();
					}}
					overscan={2}
					components={{
						Footer: () =>
							isLoading ? (
								<div style={{ textAlign: 'center', padding: 8 }}>加载中...</div>
							) : null,
					}}
				/>
			</div>
		);
	}

	// 拖拽模式下，外层 Droppable，item 用 Draggable 包裹
	return (
		<Droppable
			droppableId={columnId}
			mode="virtual"
			renderClone={({ index, ...provided }: any) => (
				<div
					ref={provided.innerRef}
					{...provided.draggableProps}
					{...provided.dragHandleProps}
					style={provided.draggableProps.style}
				>
					{renderItem(items[index], index)}
				</div>
			)}
		>
			{(provided) => (
				<div ref={provided.innerRef} {...provided.droppableProps} style={{ width, height }}>
					<Virtuoso
						style={{ width: '100%', height: '100%' }}
						data={items}
						itemContent={(index: any, item: any) => (
							<Draggable
								draggableId={getDraggableId(item, index)}
								index={index}
								key={getDraggableId(item, index)}
							>
								{(dragProvided) => (
									<div
										ref={dragProvided.innerRef}
										{...dragProvided.draggableProps}
										{...dragProvided.dragHandleProps}
										style={dragProvided.draggableProps.style}
									>
										{renderItem(item, index)}
									</div>
								)}
							</Draggable>
						)}
						endReached={() => {
							if (hasMore && !isLoading) loadMore();
						}}
						overscan={2}
						components={{
							Footer: () =>
								isLoading ? (
									<div style={{ textAlign: 'center', padding: 8 }}>加载中...</div>
								) : null,
						}}
					/>
					{provided.placeholder}
				</div>
			)}
		</Droppable>
	);
}
