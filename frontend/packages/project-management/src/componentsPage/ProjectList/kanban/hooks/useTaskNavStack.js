/**
 * 任务详情来路栈：看板点开清空；抽屉内点子任务/依赖入栈，供「返回上级」。
 */
'use client';

import { useCallback, useState } from 'react';

/**
 * @param {{
 *   tasks?: Array,
 *   activeTask?: object|null,
 *   openTask: (task: object) => void,
 *   setTaskQuery: (taskId?: number) => void,
 * }} params
 * @returns {{
 *   previousTask: object|null,
 *   openFromBoard: (task: object) => void,
 *   openFromDrawer: (nextTask: object) => void,
 *   backFromDrawer: () => void,
 *   reset: () => void,
 * }}
 */
export function useTaskNavStack({ tasks = [], activeTask, openTask, setTaskQuery }) {
	const [stack, setStack] = useState([]);

	const previousTask =
		stack.length > 0 ? tasks.find((item) => item.id === stack[stack.length - 1]) || null : null;

	const openFromBoard = useCallback(
		(task) => {
			if (!task?.id) return;
			setStack([]);
			openTask(task);
			setTaskQuery(task.id);
		},
		[openTask, setTaskQuery],
	);

	const openFromDrawer = useCallback(
		(nextTask) => {
			if (!nextTask?.id || nextTask.id === activeTask?.id) return;
			if (activeTask?.id) {
				setStack((prev) => [...prev, activeTask.id]);
			}
			openTask(nextTask);
			setTaskQuery(nextTask.id);
		},
		[activeTask, openTask, setTaskQuery],
	);

	const backFromDrawer = useCallback(() => {
		if (stack.length > 0) {
			const prevId = stack[stack.length - 1];
			const prevTask = tasks.find((item) => item.id === prevId);
			setStack((prev) => prev.slice(0, -1));
			if (prevTask) {
				openTask(prevTask);
				setTaskQuery(prevTask.id);
			}
			return;
		}
		const parentTask = tasks.find((item) => item.id === activeTask?.parent_id);
		if (parentTask) {
			openTask(parentTask);
			setTaskQuery(parentTask.id);
		}
	}, [activeTask, openTask, setTaskQuery, stack, tasks]);

	const reset = useCallback(() => {
		setStack([]);
	}, []);

	return {
		previousTask,
		openFromBoard,
		openFromDrawer,
		backFromDrawer,
		reset,
	};
}
