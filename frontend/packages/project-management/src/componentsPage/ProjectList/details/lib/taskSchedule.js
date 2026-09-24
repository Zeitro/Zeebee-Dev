/**
 * 时间线 / 工时共用的日期跨度与坐标。只读计算，不改任务数据。
 * 接口日期优先按 MM/DD/YYYY 解析，失败再兼容 ISO。
 */
import { dayjs, formatDate, formatStr } from '@repo/shared-methods/utils';

/**
 * @param {*} value
 * @returns {import('dayjs').Dayjs|null}
 */
export function parseTaskDate(value) {
	if (!value) return null;
	const strict = dayjs(value, formatStr, true);
	if (strict.isValid()) return strict.startOf('day');
	const parsed = dayjs(value);
	return parsed.isValid() ? parsed.startOf('day') : null;
}

/**
 * 无起止日时用创建日顶一个格子，undated 供样式降权。
 * @param {object} task
 * @returns {{ start: import('dayjs').Dayjs, end: import('dayjs').Dayjs, undated: boolean }}
 */
export function getTaskSpan(task) {
	const start =
		parseTaskDate(task?.start_date) ||
		parseTaskDate(task?.created_at) ||
		dayjs().startOf('day');
	const hours = Number(task?.estimated_hours) || 0;
	let end = parseTaskDate(task?.due_date);
	if (!end) end = start.add(Math.max(0, Math.ceil(hours / 8) - 1), 'day');
	if (end.isBefore(start))
		return { start: end, end: start, undated: !task?.start_date && !task?.due_date };
	return { start, end, undated: !task?.start_date && !task?.due_date };
}

/**
 * @param {object} task
 * @returns {number}
 */
export function spanDays(task) {
	const { start, end } = getTaskSpan(task);
	return end.diff(start, 'day') + 1;
}

/**
 * @param {Array} tasks
 * @returns {{ min: import('dayjs').Dayjs, max: import('dayjs').Dayjs }}
 */
export function getScheduleRange(tasks) {
	const today = dayjs().startOf('day');
	if (!tasks?.length) {
		return { min: today.subtract(3, 'day'), max: today.add(14, 'day') };
	}
	const spans = tasks.map(getTaskSpan);
	let min = spans[0].start;
	let max = spans[0].end;
	spans.forEach((span) => {
		if (span.start.isBefore(min)) min = span.start;
		if (span.end.isAfter(max)) max = span.end;
	});
	if (today.isBefore(min)) min = today;
	if (today.isAfter(max)) max = today;
	return { min: min.subtract(2, 'day'), max: max.add(7, 'day') };
}

/**
 * @param {'day'|'week'|'month'} scale
 * @param {{ min: import('dayjs').Dayjs, max: import('dayjs').Dayjs }} range
 * @returns {import('dayjs').Dayjs[]}
 */
export function buildTimelineUnits(scale, range) {
	if (scale === 'month') {
		const start = range.min.startOf('month');
		const count = Math.max(1, range.max.startOf('month').diff(start, 'month') + 1);
		return Array.from({ length: count }, (_, index) => start.add(index, 'month'));
	}
	if (scale === 'week') {
		const start = range.min.startOf('day').subtract((range.min.day() + 6) % 7, 'day');
		const count = Math.max(1, Math.ceil(range.max.diff(start, 'day') / 7));
		return Array.from({ length: count }, (_, index) => start.add(index, 'week'));
	}
	const count = Math.max(1, range.max.diff(range.min, 'day') + 1);
	return Array.from({ length: count }, (_, index) => range.min.add(index, 'day'));
}

/**
 * @param {import('dayjs').Dayjs} anchor
 * @param {'week'|'month'|'quarter'} scale
 * @returns {import('dayjs').Dayjs[]}
 */
export function buildHoursUnits(anchor, scale) {
	if (scale === 'quarter') {
		const start = anchor.startOf('quarter');
		return Array.from({ length: 3 }, (_, index) => start.add(index, 'month'));
	}
	if (scale === 'month') {
		const start = anchor.startOf('month');
		return Array.from({ length: start.daysInMonth() }, (_, index) => start.add(index, 'day'));
	}
	const start = anchor.startOf('day').subtract((anchor.day() + 6) % 7, 'day');
	return Array.from({ length: 14 }, (_, index) => start.add(index, 'day'));
}

/**
 * @param {import('dayjs').Dayjs} date
 * @param {import('dayjs').Dayjs[]} units
 * @param {'day'|'week'|'month'} scale
 * @returns {number}
 */
export function unitIndexOf(date, units, scale) {
	if (!units.length) return -1;
	const unit = scale === 'month' || scale === 'quarter' ? 'month' : 'day';
	return units.findIndex((item) => item.isSame(date, unit));
}

/**
 * 条形相对时间轴的 left / width（px）。不可见返回 null。
 * @param {{ start: import('dayjs').Dayjs, end: import('dayjs').Dayjs }} span
 * @param {import('dayjs').Dayjs[]} units
 * @param {number} colWidth
 * @param {string} scale
 * @returns {{ left: number, width: number }|null}
 */
export function getBarGeometry(span, units, colWidth, scale) {
	if (!units.length) return null;
	const origin = units[0];
	const last = units[units.length - 1];
	const monthScale = scale === 'month' || scale === 'quarter';
	const visibleEnd = monthScale
		? last.endOf('month')
		: scale === 'week'
			? last.add(6, 'day')
			: last;
	if (span.end.isBefore(origin, 'day') || span.start.isAfter(visibleEnd, 'day')) return null;

	const offset = (date) => {
		if (monthScale) {
			const months = date.startOf('month').diff(origin.startOf('month'), 'month');
			return months + (date.date() - 1) / Math.max(1, date.daysInMonth());
		}
		const days = date.startOf('day').diff(origin.startOf('day'), 'day');
		return scale === 'week' ? days / 7 : days;
	};

	const endPad = monthScale
		? 1 / Math.max(1, span.end.daysInMonth())
		: scale === 'week'
			? 1 / 7
			: 1;
	const start = Math.max(0, offset(span.start));
	const end = Math.min(units.length, offset(span.end) + endPad);
	return {
		left: start * colWidth + 4,
		width: Math.max(colWidth * 0.35, (end - start) * colWidth - 8),
	};
}

/**
 * @param {Array} tasks
 * @returns {Map<number, number>}
 */
export function buildDepthMap(tasks) {
	const byId = new Map(tasks.map((task) => [task.id, task]));
	const depths = new Map();
	tasks.forEach((task) => {
		let depth = 0;
		let cursor = task;
		const seen = new Set();
		while (cursor?.parent_id && byId.has(cursor.parent_id) && !seen.has(cursor.id)) {
			seen.add(cursor.id);
			depth += 1;
			cursor = byId.get(cursor.parent_id);
		}
		depths.set(task.id, depth);
	});
	return depths;
}

/**
 * 按执行者拆行；一人多任务合并。未分配单独成行。
 * @param {Array} tasks
 * @returns {Array<{ id: number|string, name: string, tasks: Array, hours: number }>}
 */
export function groupTasksByAssignee(tasks) {
	const rows = new Map();
	const unassigned = [];

	tasks.forEach((task) => {
		const people = Array.isArray(task.assignees)
			? task.assignees
					.map((item) => ({
						id: item.organization_user_id ?? item.user_id,
						name: item.organization_user_full_name || item.full_name || item.name || '',
					}))
					.filter((item) => item.id != null)
			: [];
		if (!people.length) {
			unassigned.push(task);
			return;
		}
		people.forEach((person) => {
			const current = rows.get(person.id) || { id: person.id, name: person.name, tasks: [] };
			if (!current.name && person.name) current.name = person.name;
			current.tasks.push(task);
			rows.set(person.id, current);
		});
	});

	const result = [...rows.values()].map((row) => ({
		...row,
		hours: row.tasks.reduce((sum, task) => sum + (Number(task.estimated_hours) || 0), 0),
	}));
	if (unassigned.length) {
		result.push({
			id: 'unassigned',
			name: '',
			tasks: unassigned,
			hours: unassigned.reduce((sum, task) => sum + (Number(task.estimated_hours) || 0), 0),
		});
	}
	return result;
}

/**
 * 重叠任务分泳道，避免工时条叠在一起。
 * @param {Array} tasks
 * @returns {Array<{ task: object, span: object, lane: number }>}
 */
export function assignTaskLanes(tasks) {
	const placed = [];
	return tasks.map((task) => {
		const span = getTaskSpan(task);
		let lane = 0;
		while (
			placed.some(
				(item) =>
					item.lane === lane &&
					!(span.end.isBefore(item.start, 'day') || span.start.isAfter(item.end, 'day')),
			)
		) {
			lane += 1;
		}
		placed.push({ lane, start: span.start, end: span.end });
		return { task, span, lane };
	});
}

/**
 * @param {import('dayjs').Dayjs} start
 * @param {import('dayjs').Dayjs} end
 * @returns {string}
 */
export function formatScheduleRange(start, end) {
	return `${formatDate(start)} – ${formatDate(end)}`;
}
