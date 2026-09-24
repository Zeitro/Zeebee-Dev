/**
 * gantt-task-react 左侧列表：样式对齐项目表格（NewTable / antd Table）。
 */
'use client';

import { Button } from 'antd';
import { displayValue, formatDate } from '@repo/shared-methods/utils';

const COL = {
	name: 200,
	duration: 72,
	date: 108,
};

/**
 * @param {{ headerHeight: number, fontFamily: string, fontSize: string, labels: { name: string, duration: string, start: string, end: string } }} props
 */
export function ProjectGanttTaskListHeader({ headerHeight, fontFamily, fontSize, labels }) {
	const cell = (width, text) => (
		<div
			className="box-border flex shrink-0 items-center border-r border-base-border bg-gray-light px-2 text-sm font-medium text-foreground"
			style={{ width, minWidth: width, height: headerHeight - 2 }}
		>
			{text}
		</div>
	);

	return (
		<div
			className="box-border flex border border-b-0 border-r-0 border-base-border bg-gray-light"
			style={{ fontFamily, fontSize }}
		>
			{cell(COL.name, labels.name)}
			{cell(COL.duration, labels.duration)}
			{cell(COL.date, labels.start)}
			{cell(COL.date, labels.end)}
		</div>
	);
}

/**
 * @param {{
 *   rowHeight: number,
 *   fontFamily: string,
 *   fontSize: string,
 *   tasks: Array,
 *   selectedTaskId: string,
 *   setSelectedTask: (id: string) => void,
 *   onExpanderClick: (task: object) => void,
 *   onTaskNameClick?: (task: object) => void,
 *   formatDuration: (days: number) => string,
 * }} props
 */
export function ProjectGanttTaskListTable({
	rowHeight,
	fontFamily,
	fontSize,
	tasks,
	selectedTaskId,
	setSelectedTask,
	onExpanderClick,
	onTaskNameClick,
	formatDuration,
}) {
	const cell = (width, children, className = '') => (
		<div
			className={`box-border flex min-w-0 shrink-0 items-center overflow-hidden border-r border-base-border px-2 text-sm text-foreground ${className}`}
			style={{ width, minWidth: width, maxWidth: width, height: rowHeight }}
		>
			{children}
		</div>
	);

	return (
		<div className="border-l border-base-border bg-white" style={{ fontFamily, fontSize }}>
			{tasks.map((task) => {
				const expander =
					task.hideChildren === false ? '▼' : task.hideChildren === true ? '▶' : '';
				const selected = task.id === selectedTaskId;
				const title = task.taskTitle || task.name;
				return (
					<div
						key={`${task.id}-row`}
						className={`flex cursor-pointer border-b border-base-border ${
							selected ? 'bg-base-active' : 'bg-white hover:bg-gray-light'
						}`}
						style={{ height: rowHeight }}
						onClick={() => setSelectedTask(task.id)}
					>
						{cell(
							COL.name,
							<div className="flex min-w-0 flex-1 items-center gap-1 overflow-hidden">
								<span
									className={`inline-flex w-3 shrink-0 select-none justify-center text-[10px] text-foreground ${
										expander ? 'cursor-pointer' : ''
									}`}
									onClick={(event) => {
										if (!expander) return;
										event.stopPropagation();
										onExpanderClick(task);
									}}
								>
									{expander || '\u00a0'}
								</span>
								<Button
									type="link"
									className="!inline-flex h-auto min-w-0 flex-1 overflow-hidden p-0 text-left font-medium !text-foreground hover:!text-foreground"
									title={title}
									onClick={(event) => {
										event.stopPropagation();
										onTaskNameClick?.(task);
									}}
								>
									<span className="block min-w-0 truncate !text-foreground">
										{title}
									</span>
								</Button>
							</div>,
						)}
						{cell(
							COL.duration,
							<span className="min-w-0 truncate">
								{formatDuration(task.durationDays || 0)}
							</span>,
						)}
						{cell(
							COL.date,
							<span className="min-w-0 truncate">
								{displayValue(task.displayStart || task.start, formatDate)}
							</span>,
						)}
						{cell(
							COL.date,
							<span className="min-w-0 truncate">
								{displayValue(task.displayEnd || task.end, formatDate)}
							</span>,
						)}
					</div>
				);
			})}
		</div>
	);
}
