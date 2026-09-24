/**
 * 工时甘特：仅执行者一行，汇总其任务工时与日期跨度；不展开任务子行。
 */

import { dayjs } from '@repo/shared-methods/utils';
import { getTaskSpan } from './taskSchedule';

const BAR = 'var(--base-primary)';
const BAR_ACTIVE = 'var(--blue-dark)';

/**
 * @param {Array<{ id: number|string, name: string, tasks: Array, hours: number }>} assigneeRows
 * @param {{ formatHours: (hours: number) => string, unassignedLabel: string }} options
 * @returns {Array}
 */
export function mapHoursGanttTasks(assigneeRows, { formatHours, unassignedLabel }) {
	if (!assigneeRows?.length) return [];

	const today = dayjs().startOf('day');

	return assigneeRows.map((row, index) => {
		const spans = row.tasks.map((task) => getTaskSpan(task));
		let start = today;
		let end = today;
		if (spans.length) {
			start = spans.reduce(
				(min, span) => (span.start.isBefore(min) ? span.start : min),
				spans[0].start,
			);
			end = spans.reduce(
				(max, span) => (span.end.isAfter(max) ? span.end : max),
				spans[0].end,
			);
		}

		const hours = Number(row.hours) || 0;
		const title = row.id === 'unassigned' ? unassignedLabel : row.name || unassignedLabel;

		return {
			id: `assignee-${row.id}`,
			name: formatHours(hours),
			taskTitle: title,
			type: 'task',
			start: start.toDate(),
			end: end.add(1, 'day').startOf('day').toDate(),
			progress: 0,
			isDisabled: true,
			displayOrder: index + 1,
			hours,
			displayStart: start.toDate(),
			displayEnd: end.toDate(),
			styles: {
				backgroundColor: BAR,
				backgroundSelectedColor: BAR_ACTIVE,
				progressColor: BAR,
				progressSelectedColor: BAR_ACTIVE,
			},
		};
	});
}
