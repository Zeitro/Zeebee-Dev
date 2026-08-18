/**
 * Linear 风格优先级图标：1 灰条、2 黄条、3 红叹号；0 与 1 同形。
 */

import { PROJECT_PRIORITY } from '@repo/project-management/enum/project';

/**
 * @param {{ level?: number, size?: number }} props
 */
export function PriorityIcon({ level = PROJECT_PRIORITY.MEDIUM, size = 13 }) {
	if (Number(level) === PROJECT_PRIORITY.URGENT) {
		return (
			<svg width={size} height={size} viewBox="0 0 16 16" className="align-[-2px]">
				<rect x="1" y="1" width="14" height="14" rx="3.5" fill="var(--base-error)" />
				<rect x="7.2" y="3.5" width="1.6" height="6" rx="0.8" fill="#fff" />
				<circle cx="8" cy="11.8" r="1" fill="#fff" />
			</svg>
		);
	}
	const isHigh = Number(level) === PROJECT_PRIORITY.HIGH;
	const activeColor = isHigh ? 'var(--base-warning)' : 'var(--base-disabled)';
	const filled = isHigh ? 3 : 2;
	const bars = [
		{ x: 2, h: 11 },
		{ x: 6.5, h: 8 },
		{ x: 11, h: 5 },
	];
	return (
		<svg width={size} height={size} viewBox="0 0 16 16" className="align-[-2px]">
			{bars.map((bar, index) => (
				<rect
					key={bar.x}
					x={bar.x}
					y={14 - bar.h}
					width="3"
					height={bar.h}
					rx="1"
					fill={index < filled ? activeColor : 'var(--base-border)'}
				/>
			))}
		</svg>
	);
}
