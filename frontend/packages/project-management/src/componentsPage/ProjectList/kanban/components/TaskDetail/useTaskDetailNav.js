/**
 * 详情内任务跳转栈：点子任务/依赖入栈，并记录来路 via（subtask | dependency）。
 * 宿主改 taskId（看板点卡）时清空；用 render 期对齐 taskId，避免 effect 内同步 setState。
 */
'use client';

import { useCallback, useState } from 'react';

/** @typedef {'subtask' | 'dependency'} TaskNavVia */

/**
 * @param {{
 *   taskId?: number,
 *   currentTask?: object|null,
 *   onTaskIdChange?: (taskId?: number) => void,
 * }} params
 */
export function useTaskDetailNav({ taskId, currentTask, onTaskIdChange }) {
	const [nav, setNav] = useState({
		/** @type {Array<{ id: number, title?: string, via: TaskNavVia }>} */
		stack: [],
		/** 内部跳转期望到达的 taskId；与外部切换区分 */
		expectedId: /** @type {number|null} */ (null),
		trackedTaskId: taskId,
	});

	// taskId 变化时对齐栈：内部跳转保留来路，外部切换 / 关闭则清空
	if (taskId !== nav.trackedTaskId) {
		if (!taskId) {
			setNav({ stack: [], expectedId: null, trackedTaskId: taskId });
		} else if (Number(nav.expectedId) === Number(taskId)) {
			setNav({
				stack: nav.stack,
				expectedId: null,
				trackedTaskId: taskId,
			});
		} else {
			setNav({ stack: [], expectedId: null, trackedTaskId: taskId });
		}
	}

	const stack = nav.stack;
	const previousTask = stack.length > 0 ? stack[stack.length - 1] : null;

	/**
	 * @param {object} nextTask
	 * @param {TaskNavVia} via 从当前任务点进下一任务的路径
	 */
	const navigateTo = useCallback(
		(nextTask, via = 'subtask') => {
			const nextId = Number(nextTask?.id);
			if (!nextId || nextId === Number(taskId)) return;
			setNav((prev) => ({
				...prev,
				stack: taskId
					? [...prev.stack, { id: Number(taskId), title: currentTask?.title, via }]
					: prev.stack,
				expectedId: nextId,
			}));
			onTaskIdChange?.(nextId);
		},
		[currentTask?.title, onTaskIdChange, taskId],
	);

	const goBack = useCallback(() => {
		let nextId;
		setNav((prev) => {
			if (prev.stack.length > 0) {
				const prevEntry = prev.stack[prev.stack.length - 1];
				nextId = prevEntry.id;
				return {
					...prev,
					stack: prev.stack.slice(0, -1),
					expectedId: prevEntry.id,
				};
			}
			const parentId = currentTask?.parent_id;
			if (parentId) {
				nextId = Number(parentId);
				return { ...prev, expectedId: nextId };
			}
			return prev;
		});
		if (nextId != null) onTaskIdChange?.(nextId);
	}, [currentTask?.parent_id, onTaskIdChange]);

	return {
		previousTask,
		/** 返回按钮旁标注：栈上来路，或无栈但有父任务时视为子任务 */
		backVia: previousTask?.via || (currentTask?.parent_id ? 'subtask' : null),
		navigateTo,
		goBack,
	};
}
