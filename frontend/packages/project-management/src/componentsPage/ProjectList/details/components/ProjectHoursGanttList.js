/**
 * 工时甘特左侧列表：仅执行者 / 工时两列。
 */
'use client';

const COL = {
	name: 200,
	hours: 88,
};

/**
 * @param {{ headerHeight: number, fontFamily: string, fontSize: string, labels: { assignee: string, hours: string } }} props
 */
export function ProjectHoursGanttListHeader({ headerHeight, fontFamily, fontSize, labels }) {
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
			{cell(COL.name, labels.assignee)}
			{cell(COL.hours, labels.hours)}
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
 *   formatHours: (hours: number) => string,
 * }} props
 */
export function ProjectHoursGanttListTable({
	rowHeight,
	fontFamily,
	fontSize,
	tasks,
	selectedTaskId,
	setSelectedTask,
	formatHours,
}) {
	const cell = (width, children) => (
		<div
			className="box-border flex min-w-0 shrink-0 items-center overflow-hidden border-r border-base-border px-2 text-sm text-foreground"
			style={{ width, minWidth: width, maxWidth: width, height: rowHeight }}
		>
			{children}
		</div>
	);

	return (
		<div className="border-l border-base-border bg-white" style={{ fontFamily, fontSize }}>
			{tasks.map((task) => {
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
							<span className="min-w-0 truncate font-medium" title={title}>
								{title}
							</span>,
						)}
						{cell(
							COL.hours,
							<span className="min-w-0 truncate">
								{formatHours(task.hours || 0)}
							</span>,
						)}
					</div>
				);
			})}
		</div>
	);
}
