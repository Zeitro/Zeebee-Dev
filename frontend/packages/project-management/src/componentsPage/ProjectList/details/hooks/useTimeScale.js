/**
 * 时间轴窗口：周 / 月 / 季切换与前后翻页。只改本组件 state，不写 URL。
 */
'use client';

import { useMemo, useState } from 'react';
import { dayjs, formatDate } from '@repo/shared-methods/utils';
import { buildHoursUnits, formatScheduleRange } from '../lib/taskSchedule';

/**
 * @param {'week'|'month'|'quarter'} [initialScale]
 */
export function useTimeScale(initialScale = 'week') {
	const [scale, setScale] = useState(initialScale);
	const [anchor, setAnchor] = useState(() => dayjs());

	const units = useMemo(() => buildHoursUnits(anchor, scale), [anchor, scale]);

	const rangeLabel = useMemo(() => {
		if (!units.length) return '';
		if (scale === 'quarter') return formatScheduleRange(units[0], units[units.length - 1].endOf('month'));
		return formatScheduleRange(units[0], units[units.length - 1]);
	}, [scale, units]);

	const shift = (direction) => {
		setAnchor((current) => {
			if (scale === 'week') return current.add(direction * 7, 'day');
			if (scale === 'month') return current.add(direction, 'month');
			return current.add(direction, 'quarter');
		});
	};

	return {
		scale,
		setScale,
		units,
		rangeLabel,
		shift,
		goToday: () => setAnchor(dayjs()),
		todayLabel: formatDate(dayjs()),
	};
}
