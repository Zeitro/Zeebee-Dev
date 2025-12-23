'use client';
import { PointerSensor, useSensor, useSensors, useDroppable } from '@dnd-kit/core';
import { arrayMove } from '@dnd-kit/sortable';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
const useContent = ({ items, setItems, onDragEnd, isSort, list, setList }) => {
	const sensors = useSensors(
		useSensor(PointerSensor, {
			activationConstraint: {
				distance: 5, // 移动5像素后才开始拖拽
			},
		}),
	);

	// 列空占位符
	const EmptyColumnDroppable = ({ colValue, children }) => {
		const { setNodeRef } = useDroppable({ id: colValue });
		return (
			<div
				ref={setNodeRef}
				style={{
					height: '100%',
					cursor: 'pointer',
				}}
			>
				{children}
			</div>
		);
	};

	// 列表项
	const SortableItem = ({ id, item, index, renderItem, style, itemHeight }) => {
		const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
			useSortable({
				id,
				data: item,
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
				opacity: isDragging ? 0.7 : 1,
				boxShadow: isDragging ? '0 4px 8px rgba(0,0,0,0.2)' : 'none',
			},
		});
	};

	function findContainer(id) {
		return items.find((item) => item?.list.some((item) => item?.id === id))?.value;
	}

	// 跨列拖拽，支持空列拖拽
	function handleCrossColumnDragEnd(e) {
		const { active, over } = e;
		if (!over) return;
		const activeId = active.id;
		const overId = over.id;
		if (activeId === overId) return;

		const fromColValue = findContainer(activeId);
		// 目标是卡片还是空列
		const isOverColumn = items.some((col) => col.value === overId);
		let toColValue, overIndex;
		if (isOverColumn) {
			// 拖到空列
			toColValue = overId;
			overIndex = 0;
		} else {
			toColValue = findContainer(overId);
			const toCol = items.find((col) => col.value === toColValue);
			overIndex = toCol ? toCol.list.findIndex((item) => item.id === overId) : -1;
		}

		if (!fromColValue || !toColValue) return;

		const fromCol = items.find((col) => col.value === fromColValue);
		const toCol = items.find((col) => col.value === toColValue);
		const activeIndex = fromCol.list.findIndex((item) => item.id === activeId);

		if (fromColValue === toColValue && overIndex !== -1) {
			// 列内排序
			if (!isSort) return;
			const newList = arrayMove(fromCol.list, activeIndex, overIndex);
			const newItems = items.map((col) =>
				col.value === fromColValue ? { ...col, list: newList } : col,
			);
			setItems(newItems);
			onDragEnd({ index: activeIndex, key: toColValue, type: 'columnSorting' });
		} else if (toCol && activeIndex !== -1) {
			// 跨列移动（含空列）
			const movingItem = fromCol.list[activeIndex];
			const newFromList = fromCol.list.filter((item) => item.id !== activeId);
			const newToList = [
				...toCol.list.slice(0, overIndex),
				movingItem,
				...toCol.list.slice(overIndex),
			];
			const newItems = items.map((col) => {
				if (col.value === fromColValue) return { ...col, list: newFromList };
				if (col.value === toColValue) return { ...col, list: newToList };
				return col;
			});
			setItems(newItems);
			onDragEnd({ movingItem, key: toColValue, type: 'CrossColumnDragSorting' });
		}
	}

	// 只有单列，列内排序
	const handleDragEnd = (e) => {
		if (!isSort) return;
		const { active, over } = e;
		const activeId = active.id;
		const overId = over.id;
		if (activeId === overId) return;
		const activeIndex = list.findIndex((item) => item.id === activeId);
		const overIndex = list.findIndex((item) => item.id === overId);
		const newList = arrayMove(list, activeIndex, overIndex);
		setList(newList);
		onDragEnd({ index: activeIndex, type: 'columnSorting' });
	};

	return { handleDragEnd, sensors, handleCrossColumnDragEnd, SortableItem, EmptyColumnDroppable };
};
export default useContent;
