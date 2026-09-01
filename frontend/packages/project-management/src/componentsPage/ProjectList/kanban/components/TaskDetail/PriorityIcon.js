/**
 * Linear 风格优先级图标：低 1 灰条、中 2 蓝条、高 3 黄条、紧急红叹号。
 */

import { PROJECT_PRIORITY } from '@repo/project-management/enum/project';

/**
 * @param {{ level?: number, size?: number }} props
 */
export function PriorityIcon({ level = PROJECT_PRIORITY.MEDIUM, size = 13 }) {
	const priority = Number(level);

	if (priority === PROJECT_PRIORITY.URGENT) {
		return (
			<svg width={size} height={size} viewBox="0 0 16 16" className="align-[-2px]">
				<rect x="1" y="1" width="14" height="14" rx="3.5" fill="var(--base-error)" />
				<rect x="7.2" y="3.5" width="1.6" height="6" rx="0.8" fill="#fff" />
				<circle cx="8" cy="11.8" r="1" fill="#fff" />
			</svg>
		);
	}

	const barStyles = {
		[PROJECT_PRIORITY.LOW]: { filled: 1, activeColor: 'var(--base-disabled)' },
		[PROJECT_PRIORITY.MEDIUM]: { filled: 2, activeColor: 'var(--blue-medium)' },
		[PROJECT_PRIORITY.HIGH]: { filled: 3, activeColor: 'var(--base-warning)' },
	};
	const { filled, activeColor } = barStyles[priority] ?? barStyles[PROJECT_PRIORITY.MEDIUM];

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
