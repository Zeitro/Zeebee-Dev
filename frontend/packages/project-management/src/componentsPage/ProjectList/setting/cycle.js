/**
 * 项目周期换算：截止日期默认 = 开始日期 + 周期。纯计算，不改外部状态。
 */
import { dayjs, formatDate, formatStr } from '@repo/shared-methods/utils';

export const PROJECT_CYCLE_UNITS = [
	{ value: 'day', labelKey: 'projectSetting.cycleDay' },
	{ value: 'week', labelKey: 'projectSetting.cycleWeek' },
	{ value: 'month', labelKey: 'projectSetting.cycleMonth' },
];

/**
 * @param {import('dayjs').Dayjs} start
 * @param {number} amount
 * @param {string} [unit]
 * @returns {import('dayjs').Dayjs}
 */
export function addProjectCycle(start, amount, unit) {
	const n = Math.max(0, Math.floor(Number(amount) || 0));
	if (unit === 'day') return start.add(n || 1, 'day');
	if (unit === 'week') return start.add(n || 1, 'week');
	return start.add(n || 3, 'month');
}

/**
 * 接口日期 → dayjs。优先系统格式 MM/DD/YYYY，兼容 ISO。
 * @param {string|undefined|null} value
 * @returns {import('dayjs').Dayjs|undefined}
 */
function parseProjectDate(value) {
	if (!value) return undefined;
	const asDisplay = dayjs(value, formatStr, true);
	if (asDisplay.isValid()) return asDisplay;
	const parsed = dayjs(value);
	return parsed.isValid() ? parsed : undefined;
}

/**
 * 详情里的周期字段 → 表单初值；接口未返回时不预填日期。
 * @param {object|null|undefined} project
 * @returns {{ cycle_months: number|undefined, cycle_unit: string, start_date: import('dayjs').Dayjs|undefined, end_date: import('dayjs').Dayjs|undefined }}
 */
export function toCycleFormValues(project) {
	const start = parseProjectDate(project?.start_date);
	const amount = project?.cycle_months || undefined;
	const unit = project?.cycle_unit || 'month';
	const end = project?.end_date
		? parseProjectDate(project.end_date)
		: start
			? addProjectCycle(start, amount || 3, unit)
			: undefined;
	return {
		cycle_months: amount,
		cycle_unit: unit,
		start_date: start,
		end_date: end,
	};
}

/**
 * 表单周期 → 更新 body 片段。日期走 `formatDate`（MM/DD/YYYY），与任务日期一致；空日期不传。
 * @param {object} values
 * @returns {object}
 */
export function toCycleUpdateBody(values) {
	const body = {
		cycle_unit: values.cycle_unit || 'month',
	};
	if (values.cycle_months != null && values.cycle_months !== '') {
		body.cycle_months = Number(values.cycle_months) || 0;
	}
	const startDate = formatDate(values.start_date);
	const endDate = formatDate(values.end_date);
	if (startDate) body.start_date = startDate;
	if (endDate) body.end_date = endDate;
	return body;
}
