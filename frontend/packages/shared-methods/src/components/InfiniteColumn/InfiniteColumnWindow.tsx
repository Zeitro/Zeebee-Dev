import React, { useRef } from 'react';
import { FixedSizeList as List, ListOnItemsRenderedProps } from 'react-window';
import { Droppable, Draggable } from '@hello-pangea/dnd';
import type { DragInfo } from './dnd-types';

export interface InfiniteColumnWindowProps<T> {
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

export function InfiniteColumnWindow<T>({
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
}: InfiniteColumnWindowProps<T>) {
	const loadingRef = useRef(false);

	const handleItemsRendered = (props: ListOnItemsRenderedProps) => {
		const { visibleStopIndex } = props;
		if (hasMore && !isLoading && !loadingRef.current && visibleStopIndex >= items.length - 3) {
			loadingRef.current = true;
			loadMore();
			setTimeout(() => {
				loadingRef.current = false;
			}, 500);
		}
	};

	if (!draggable) {
		return (
			<List
				height={height}
				itemCount={items.length}
				itemSize={itemHeight}
				width={width}
				onItemsRendered={handleItemsRendered}
			>
				{({ index, style }: any) => (
					<div style={style} key={index}>
						{renderItem(items[index], index)}
					</div>
				)}
			</List>
		);
	}

	// 拖拽模式下，外层 Droppable，item 用 Draggable 包裹
	return (
		<Droppable
			droppableId={columnId}
			mode="virtual"
			renderClone={({ index, style, ...provided }: any) => (
				<div
					ref={provided.innerRef}
					{...provided.draggableProps}
					{...provided.dragHandleProps}
					style={{ ...style, ...provided.draggableProps.style }}
				>
					{renderItem(items[index], index)}
				</div>
			)}
		>
			{(provided) => (
				<div ref={provided.innerRef} {...provided.droppableProps} style={{ width, height }}>
					<List
						height={height}
						itemCount={items.length}
						itemSize={itemHeight}
						width={width}
						onItemsRendered={handleItemsRendered}
						outerRef={provided.innerRef}
					>
						{({ index, style }: any) => (
							<Draggable
								draggableId={getDraggableId(items[index], index)}
								index={index}
								key={getDraggableId(items[index], index)}
							>
								{(dragProvided) => (
									<div
										ref={dragProvided.innerRef}
										{...dragProvided.draggableProps}
										{...dragProvided.dragHandleProps}
										style={{ ...style, ...dragProvided.draggableProps.style }}
									>
										{renderItem(items[index], index)}
									</div>
								)}
							</Draggable>
						)}
					</List>
					{provided.placeholder}
				</div>
			)}
		</Droppable>
	);
}
