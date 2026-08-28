/**
 * 流转连线路径与标签中点；纯函数，供 SVG path 渲染使用。
 */

import { NODE_H, NODE_W, getCenter, getEdgePoint } from './flowLayout';

/** 连线中点圆 + 标签留白，用于 SVG 包围盒 */
const EDGE_HANDLE_PAD = 20;

/**
 * @param {{ minX: number, minY: number, maxX: number, maxY: number }} bounds
 * @param {number} x
 * @param {number} y
 */
function addBoundsPoint(bounds, x, y) {
	bounds.minX = Math.min(bounds.minX, x);
	bounds.minY = Math.min(bounds.minY, y);
	bounds.maxX = Math.max(bounds.maxX, x);
	bounds.maxY = Math.max(bounds.maxY, y);
}

/**
 * @param {{ minX: number, minY: number, maxX: number, maxY: number }} bounds
 * @param {number} midX
 * @param {number} midY
 */
function addMidHandleBounds(bounds, midX, midY) {
	addBoundsPoint(bounds, midX - EDGE_HANDLE_PAD, midY - EDGE_HANDLE_PAD - 14);
	addBoundsPoint(bounds, midX + EDGE_HANDLE_PAD, midY + EDGE_HANDLE_PAD);
}

/**
 * @param {object} params
 * @param {object} params.transition
 * @param {object[]} params.transitions
 * @param {object[]} params.states
 * @param {Map<number, { x: number; y: number }>} params.positions
 * @returns {{ pathD: string; midX: number; midY: number; bounds: { minX: number; minY: number; maxX: number; maxY: number } } | null}
 */
export function computeTransitionPath({ transition, transitions, states, positions }) {
	const fromPos = positions.get(transition.from_state_id);
	const toPos = positions.get(transition.to_state_id);
	if (!fromPos || !toPos) return null;

	const from = getEdgePoint(positions, transition.from_state_id, transition.to_state_id);
	const to = getEdgePoint(positions, transition.to_state_id, transition.from_state_id);

	const hasReverse = transitions.some(
		(other) =>
			other.from_state_id === transition.to_state_id &&
			other.to_state_id === transition.from_state_id,
	);
	const fromCenter = getCenter(positions, transition.from_state_id);
	const toCenter = getCenter(positions, transition.to_state_id);
	const physDist = Math.sqrt((toCenter.x - fromCenter.x) ** 2 + (toCenter.y - fromCenter.y) ** 2);
	const isForward = transition.from_state_id < transition.to_state_id;
	const stateDistance =
		states.length > 1
			? Math.abs(
					states.findIndex((s) => s.id === transition.from_state_id) -
						states.findIndex((s) => s.id === transition.to_state_id),
				)
			: 1;

	let pathD;
	let midX;
	let midY;
	const bounds = {
		minX: Infinity,
		minY: Infinity,
		maxX: -Infinity,
		maxY: -Infinity,
	};

	if (stateDistance <= 1 || physDist < 300) {
		if (hasReverse) {
			const firstId = Math.min(transition.from_state_id, transition.to_state_id);
			const lastId = Math.max(transition.from_state_id, transition.to_state_id);
			const fc = getCenter(positions, firstId);
			const lc = getCenter(positions, lastId);
			const fdx = lc.x - fc.x;
			const fdy = lc.y - fc.y;
			const flen = Math.sqrt(fdx * fdx + fdy * fdy) || 1;
			const pnx = -fdy / flen;
			const pny = fdx / flen;
			const offset = 12;
			const sign = isForward ? 1 : -1;
			const fromCx = from.x + pnx * offset * sign;
			const fromCy = from.y + pny * offset * sign;
			const toCx = to.x + pnx * offset * sign;
			const toCy = to.y + pny * offset * sign;
			const curveMag = 8;
			const cpx = (fromCx + toCx) / 2 + pnx * curveMag * sign;
			const cpy = (fromCy + toCy) / 2 + pny * curveMag * sign;
			pathD = `M ${fromCx} ${fromCy} Q ${cpx} ${cpy} ${toCx} ${toCy}`;
			midX = 0.25 * fromCx + 0.5 * cpx + 0.25 * toCx;
			midY = 0.25 * fromCy + 0.5 * cpy + 0.25 * toCy;
			[fromCx, cpx, toCx].forEach((x, i) => {
				const ys = [fromCy, cpy, toCy];
				addBoundsPoint(bounds, x, ys[i]);
			});
		} else {
			pathD = `M ${from.x} ${from.y} L ${to.x} ${to.y}`;
			midX = (from.x + to.x) / 2;
			midY = (from.y + to.y) / 2;
			addBoundsPoint(bounds, from.x, from.y);
			addBoundsPoint(bounds, to.x, to.y);
		}
	} else {
		const routeGap = 25 + stateDistance * 18;
		const cpOffset = Math.abs(to.x - from.x) * 0.15;
		if (isForward) {
			const topY = from.y - routeGap * stateDistance * 0.6;
			pathD = `M ${from.x} ${from.y} C ${from.x + cpOffset} ${topY}, ${to.x - cpOffset} ${topY}, ${to.x} ${to.y}`;
			midX =
				0.125 * from.x +
				0.375 * (from.x + cpOffset) +
				0.375 * (to.x - cpOffset) +
				0.125 * to.x;
			midY = 0.125 * from.y + 0.375 * topY + 0.375 * topY + 0.125 * to.y;
			addBoundsPoint(bounds, from.x, from.y);
			addBoundsPoint(bounds, from.x + cpOffset, topY);
			addBoundsPoint(bounds, to.x - cpOffset, topY);
			addBoundsPoint(bounds, to.x, to.y);
		} else {
			const bottomY = from.y + NODE_H + routeGap * stateDistance * 0.6;
			pathD = `M ${from.x} ${from.y + NODE_H} C ${from.x + cpOffset} ${bottomY}, ${to.x - cpOffset} ${bottomY}, ${to.x} ${to.y + NODE_H}`;
			midX =
				0.125 * from.x +
				0.375 * (from.x + cpOffset) +
				0.375 * (to.x - cpOffset) +
				0.125 * to.x;
			midY =
				0.125 * (from.y + NODE_H) +
				0.375 * bottomY +
				0.375 * bottomY +
				0.125 * (to.y + NODE_H);
			addBoundsPoint(bounds, from.x, from.y + NODE_H);
			addBoundsPoint(bounds, from.x + cpOffset, bottomY);
			addBoundsPoint(bounds, to.x - cpOffset, bottomY);
			addBoundsPoint(bounds, to.x, to.y + NODE_H);
		}
	}

	addMidHandleBounds(bounds, midX, midY);

	return { pathD, midX, midY, bounds };
}

const DIAGRAM_PAD_X = 80;
const DIAGRAM_PAD_Y = 48;
const DIAGRAM_MIN_HEIGHT = 500;

/**
 * 节点 + 连线路径包围盒，用于 SVG / 画布尺寸，避免长距弧线被裁切。
 * @param {Map<number, { x: number; y: number }>} positions
 * @param {object[]} states
 * @param {object[]} transitions
 * @returns {{ width: number, height: number, offsetY: number }}
 */
export function computeDiagramExtents(positions, states, transitions) {
	let minX = 0;
	let minY = 0;
	let maxX = NODE_W + 100;
	let maxY = NODE_H + 200;

	positions.forEach((pos) => {
		minX = Math.min(minX, pos.x);
		minY = Math.min(minY, pos.y);
		maxX = Math.max(maxX, pos.x + NODE_W);
		maxY = Math.max(maxY, pos.y + NODE_H);
	});

	transitions.forEach((transition) => {
		const path = computeTransitionPath({
			transition,
			transitions,
			states,
			positions,
		});
		if (!path?.bounds) return;
		minX = Math.min(minX, path.bounds.minX);
		minY = Math.min(minY, path.bounds.minY);
		maxX = Math.max(maxX, path.bounds.maxX);
		maxY = Math.max(maxY, path.bounds.maxY);
	});

	const offsetY = minY < 0 ? -minY + DIAGRAM_PAD_Y : 0;

	return {
		width: maxX + DIAGRAM_PAD_X,
		height: Math.max(DIAGRAM_MIN_HEIGHT, maxY + offsetY + DIAGRAM_PAD_Y),
		offsetY,
	};
}
