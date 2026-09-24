/**
 * 在甘特「今天」列内画精确刻度线 + 表头「今天」标签。
 * 周/月视图下列很宽，仅靠整列底色不够定位。
 * 插件 today rect 的 fill 须与 Gantt todayColor（transparent）一致才能命中定位。
 */

import { useEffect } from 'react';

const TODAY_FILL = 'transparent';
const MARK = 'data-today-tick';
const TICK_COLOR = 'var(--base-error)';

/**
 * 今天在当前列内的水平比例（0–1）。
 * @param {'day'|'week'|'month'|'year'} scale
 * @returns {number}
 */
function todayFractionInColumn(scale) {
	const now = new Date();
	if (scale === 'week') {
		const dayFromMon = (now.getDay() + 6) % 7;
		return (dayFromMon + now.getHours() / 24) / 7;
	}
	if (scale === 'month') {
		const days = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate();
		return (now.getDate() - 1 + now.getHours() / 24) / days;
	}
	if (scale === 'year') {
		const start = new Date(now.getFullYear(), 0, 0);
		const dayOfYear = Math.floor((now - start) / 86400000);
		const daysInYear =
			now.getFullYear() % 4 === 0 && (now.getFullYear() % 100 !== 0 || now.getFullYear() % 400 === 0)
				? 366
				: 365;
		return dayOfYear / daysInYear;
	}
	return 0.5;
}

/**
 * @param {SVGSVGElement} svg
 * @param {number} x
 * @param {number} y1
 * @param {number} y2
 */
function appendLine(svg, x, y1, y2) {
	const line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
	line.setAttribute(MARK, 'line');
	line.setAttribute('x1', String(x));
	line.setAttribute('x2', String(x));
	line.setAttribute('y1', String(y1));
	line.setAttribute('y2', String(y2));
	line.setAttribute('stroke', TICK_COLOR);
	line.setAttribute('stroke-width', '1.5');
	line.setAttribute('stroke-dasharray', '4 3');
	line.setAttribute('pointer-events', 'none');
	svg.appendChild(line);
}

/**
 * @param {SVGSVGElement} svg
 * @param {number} x
 * @param {number} headerHeight
 * @param {string} label
 */
function appendHeaderBadge(svg, x, headerHeight, label) {
	const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
	g.setAttribute(MARK, 'badge');
	g.setAttribute('pointer-events', 'none');

	const chipW = Math.max(36, label.length * 12 + 12);
	const chipH = 18;
	const chipX = x - chipW / 2;
	const chipY = 4;

	const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
	rect.setAttribute('x', String(chipX));
	rect.setAttribute('y', String(chipY));
	rect.setAttribute('width', String(chipW));
	rect.setAttribute('height', String(chipH));
	rect.setAttribute('rx', '4');
	rect.setAttribute('fill', TICK_COLOR);

	const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
	text.setAttribute('x', String(x));
	text.setAttribute('y', String(chipY + chipH / 2 + 1));
	text.setAttribute('text-anchor', 'middle');
	text.setAttribute('dominant-baseline', 'middle');
	text.setAttribute('fill', '#fff');
	text.setAttribute('font-size', '11');
	text.setAttribute('font-weight', '500');
	text.textContent = label;

	const tick = document.createElementNS('http://www.w3.org/2000/svg', 'line');
	tick.setAttribute('x1', String(x));
	tick.setAttribute('x2', String(x));
	tick.setAttribute('y1', String(chipY + chipH));
	tick.setAttribute('y2', String(headerHeight));
	tick.setAttribute('stroke', TICK_COLOR);
	tick.setAttribute('stroke-width', '1.5');
	tick.setAttribute('stroke-dasharray', '4 3');

	g.appendChild(rect);
	g.appendChild(text);
	g.appendChild(tick);
	svg.appendChild(g);
}

/**
 * @param {HTMLElement|null} container
 * @param {'day'|'week'|'month'|'year'} scale
 * @param {string} label
 */
export function useTodayTick(container, scale, label) {
	useEffect(() => {
		if (!container || !label) return undefined;

		let raf = 0;

		const clearMarks = () => {
			container.querySelectorAll(`[${MARK}]`).forEach((node) => node.remove());
		};

		const paint = () => {
			clearMarks();
			const todayRect = [...container.querySelectorAll('rect')]
				.filter((node) => node.getAttribute('fill') === TODAY_FILL)
				.sort(
					(a, b) =>
						(Number(b.getAttribute('height')) || 0) -
						(Number(a.getAttribute('height')) || 0),
				)[0];
			if (!todayRect) return;

			const colX = Number(todayRect.getAttribute('x')) || 0;
			const colW = Number(todayRect.getAttribute('width')) || 0;
			const bodyH = Number(todayRect.getAttribute('height')) || 0;
			const frac = Math.min(0.96, Math.max(0.04, todayFractionInColumn(scale)));
			const x = colX + colW * frac;

			const bodySvg = todayRect.ownerSVGElement;
			if (bodySvg) appendLine(bodySvg, x, 0, bodyH);

			const vert = todayRect.closest('div')?.parentElement;
			const calendarSvg = vert?.querySelector(':scope > svg');
			if (calendarSvg instanceof SVGSVGElement) {
				const headerH = Number(calendarSvg.getAttribute('height')) || 50;
				appendHeaderBadge(calendarSvg, x, headerH, label);
			}
		};

		const schedule = () => {
			cancelAnimationFrame(raf);
			raf = requestAnimationFrame(paint);
		};

		schedule();
		const observer = new MutationObserver((records) => {
			const selfOnly = records.every((record) =>
				[...record.addedNodes, ...record.removedNodes].every(
					(node) =>
						node.nodeType === 1 &&
						/** @type {Element} */ (node).hasAttribute?.(MARK),
				),
			);
			if (selfOnly) return;
			schedule();
		});
		observer.observe(container, { childList: true, subtree: true, attributes: true });

		const scrollParents = [...container.querySelectorAll('div')];
		scrollParents.forEach((el) => el.addEventListener('scroll', schedule, { passive: true }));

		return () => {
			cancelAnimationFrame(raf);
			observer.disconnect();
			scrollParents.forEach((el) => el.removeEventListener('scroll', schedule));
			clearMarks();
		};
	}, [container, label, scale]);
}
