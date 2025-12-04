'use client';
import { PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
// import { arrayMove } from '@dnd-kit/sortable';
import _ from 'lodash';
const useContent = ({ items, setItems, onDragEnd }) => {
	const sensors = useSensors(useSensor(PointerSensor));

	function findContainer(id, items) {
		if (id in items) {
			return id;
		}
		return Object.keys(items).find((key) => items[key].find((item) => item?.id === id));
	}

	const handleDragEnd = async (e) => {
		const { active, over } = e;
		const overId = over?.id;
		const activeId = active?.id;
		if (!overId) return;
		const activeContainer = findContainer(activeId) || '';
		const overContainer = findContainer(overId) || '';
		if (activeContainer === overContainer) return;
		const activeRow = items[activeContainer]?.find((item) => item.id === activeId);
		onDragEnd?.(activeRow, activeContainer);
	};

	function handleDragMove(e) {
		const { active, over } = e;
		const overId = over?.id;
		const activeId = active?.id;
		const activeContainer = findContainer(activeId) || '';
		const overContainer = findContainer(overId) || '';

		if (!overContainer || !activeContainer) {
			return;
		}
		const activeArr = _.cloneDeep(items[activeContainer]);
		const overArr = _.cloneDeep(items[overContainer]);
		// 跨列移动:将activeContainer里删除拖拽元素，在overContainer中添加拖拽元素
		if (activeContainer !== overContainer) {
			const overIndex = items[overContainer].findIndex((item) => item.id === overId);
			const activeItem = items[activeContainer]?.find((item) => item.id === activeId);
			const newIndex = overIndex >= 0 ? overIndex : items[overContainer].length + 1;
			if (!activeItem) return;
			const data = {
				...items,
				[activeContainer]: activeArr.filter((item) => {
					return item.id !== activeId;
				}),
				[overContainer]: [
					...overArr.slice(0, newIndex),
					activeItem,
					...overArr.slice(newIndex, overArr.length),
				],
			};
			setItems(data);
			return;
		}
		// 列内排序
	}
	return { handleDragMove, handleDragEnd, sensors, findContainer };
};
export default useContent;
