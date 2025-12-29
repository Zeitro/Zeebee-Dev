import dayjs from 'dayjs';

export const formatStr = 'MM/DD/YYYY';
export const formatTimeStr = 'MM/DD/YYYY HH:mm:ss';

// 格式化时间
export function formatTime(date, format = formatTimeStr) {
	if (!date || date === 'Invalid Date') return '--';
	return dayjs(date).format(format);
}
// 格式化日期
export function formatDate(date, format = formatStr) {
	if (!date || date === 'Invalid Date') return '--';
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
