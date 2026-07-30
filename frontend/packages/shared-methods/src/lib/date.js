import dayjs from 'dayjs';
import advancedFormat from 'dayjs/plugin/advancedFormat';
import customParseFormat from 'dayjs/plugin/customParseFormat';
import localeData from 'dayjs/plugin/localeData';
import quarterOfYear from 'dayjs/plugin/quarterOfYear';
import weekOfYear from 'dayjs/plugin/weekOfYear';
import weekYear from 'dayjs/plugin/weekYear';
import weekday from 'dayjs/plugin/weekday';

dayjs.extend(customParseFormat);
dayjs.extend(weekday);
dayjs.extend(localeData);
dayjs.extend(weekOfYear);
dayjs.extend(weekYear);
dayjs.extend(advancedFormat);
dayjs.extend(quarterOfYear);

/** 与 antd DatePicker / RangePicker 兼容的 dayjs 实例 */
export { dayjs };

export const formatStr = 'MM/DD/YYYY';
export const formatTimeStr = 'MM/DD/YYYY HH:mm:ss';

// 格式化时间
export function formatTime(date, format = formatTimeStr) {
	if (!date || date === 'Invalid Date') return '';
	return dayjs(date).format(format);
}
// 格式化日期
export function formatDate(date, format = formatStr) {
	if (!date || date === 'Invalid Date') return '';
	return dayjs(date).format(format);
}

// 计算两日期相差几天（绝对值）
export function diffDays(date1, date2) {
	const d1 = dayjs(date1);
	const d2 = dayjs(date2);
	return Math.abs(d1.diff(d2, 'day'));
}

//  判断日期1是否在日期2之前
export function isBefore(date1, date2) {
	return dayjs(date1).isBefore(dayjs(date2));
}

// 获取某日期加减指定天数后的日期
export function addDays(date, days) {
	return dayjs(date).add(days, 'day').toDate();
}

/**
 * 计算时间区间，返回 { startDate, endDate }
 * @param {number|string} days - 天数
 * @param {Date|string|number} [date=new Date()] - 基准日期，默认为当天
 * @param {string|boolean} [direction='past'] - 'past' | 'future'，或布尔值(true 表示过去)
 * @returns {{startDate: Date, endDate: Date}}
 */
export function getDateRange(days = 7, date = new Date(), direction = 'past') {
	// 特殊值处理：'All' 或 空字符串 返回空对象
	if (days === 'All' || !days) return {};

	const n = Math.abs(Number(days)) || 0;
	const count = Math.max(1, n);
	const base = dayjs(date);

	// 支持布尔型参数：true 表示过去，false 表示未来
	const isPast =
		typeof direction === 'boolean'
			? Boolean(direction)
			: String(direction).toLowerCase() !== 'future';

	if (isPast) {
		const startDate = base
			.subtract(count - 1, 'day')
			.startOf('day')
			.toDate();
		const endDate = base.endOf('day').toDate();
		return {startDate: formatDate(startDate), endDate: formatDate(endDate)};
	}

	const startDate = base.startOf('day').toDate();
	const endDate = base
		.add(count - 1, 'day')
		.endOf('day')
		.toDate();
	return {start_date: formatDate(startDate), end_date: formatDate(endDate)};
}

/**
 * 秒数转换为 Date 对象（更精确）
 * @param {number} seconds - 秒数
 * @param {number} defaultDays - 默认天数（秒数无效时使用）
 * @returns {Date} 过期时间 Wed Mar 18 2026 11:18:28 GMT+0800 (GMT+08:00)
 */
export const secondsToDate = (seconds, defaultDays = 7) => {
	const date = new Date();
	if (!seconds || typeof seconds !== 'number' || seconds <= 0) {
		date.setDate(date.getDate() + defaultDays);
		return date;
	}
	date.setSeconds(date.getSeconds() + seconds);
	return date;
};
