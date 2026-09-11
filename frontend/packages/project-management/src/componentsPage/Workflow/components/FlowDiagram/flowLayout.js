/** 流程图节点尺寸、端口与自动布局；纯函数，不改外部状态。 */

export const NODE_W = 140;
export const NODE_H = 50;

/** @type {Record<string, { dx: number; dy: number }>} */
export const portOffsets = {
	'top-left': { dx: 0.2, dy: 0 },
	'top-center': { dx: 0.5, dy: 0 },
	'top-right': { dx: 0.8, dy: 0 },
	'bottom-left': { dx: 0.2, dy: 1 },
	'bottom-center': { dx: 0.5, dy: 1 },
	'bottom-right': { dx: 0.8, dy: 1 },
};

/** @type {string[]} */
export const portDirs = [
	'top-left',
	'top-center',
	'top-right',
	'bottom-left',
	'bottom-center',
	'bottom-right',
];

/**
 * @param {object[]} states
 * @returns {Map<number, { x: number; y: number }>}
 */
export function computeLayout(states) {
	const positions = new Map();
	const sorted = [...states].sort((a, b) => a.sort_order - b.sort_order);
	const gapX = 100;
	const gapY = 80;
	const startX = 80;
	const startY = 160;

	const initial = sorted.find((s) => s.is_initial);
	const rest = sorted.filter((s) => !s.is_initial);
	const ordered = initial ? [initial, ...rest] : rest;

	const cols = Math.min(ordered.length, 5);
	ordered.forEach((state, i) => {
		const row = Math.floor(i / cols);
		const col = i % cols;
		positions.set(state.id, {
			x: startX + col * (NODE_W + gapX),
			y: startY + row * (NODE_H + gapY),
		});
	});

	return positions;
}

/**
 * @param {object} transition
 * @param {number} index
 * @returns {string}
 */
export function getTransitionKey(transition, index) {
	if (transition.id !== undefined && transition.id !== null) {
		return `id:${transition.id}`;
	}

	return [
		'edge',
		transition.workflow_id,
		transition.from_state_id,
		transition.to_state_id,
		transition.sort_order,
		transition.name,
		index,
	].join(':');
}

/**
 * @param {Map<number, { x: number; y: number }>} positions
 * @param {number} stateId
 * @returns {{ x: number; y: number }}
 */
export function getCenter(positions, stateId) {
	const pos = positions.get(stateId);
	if (!pos) return { x: 0, y: 0 };
	return { x: pos.x + NODE_W / 2, y: pos.y + NODE_H / 2 };
}

/**
 * @param {Map<number, { x: number; y: number }>} positions
 * @param {number} fromId
 * @param {number} toId
 * @returns {{ x: number; y: number }}
 */
export function getEdgePoint(positions, fromId, toId) {
	const from = getCenter(positions, fromId);
	const to = getCenter(positions, toId);
	const dx = to.x - from.x;
	const dy = to.y - from.y;
	const pos = positions.get(fromId);
	if (!pos) return from;

	const hw = NODE_W / 2;
	const hh = NODE_H / 2;
	let ex;
	let ey;
	if (Math.abs(dx) * NODE_H > Math.abs(dy) * NODE_W) {
		ex = pos.x + hw + (dx > 0 ? hw : -hw);
		ey = from.y + (dx !== 0 ? (hw * dy) / Math.abs(dx) : 0);
		ey = Math.max(pos.y, Math.min(pos.y + NODE_H, ey));
	} else {
		ey = pos.y + hh + (dy > 0 ? hh : -hh);
		ex = from.x + (dy !== 0 ? (hh * dx) / Math.abs(dy) : 0);
		ex = Math.max(pos.x, Math.min(pos.x + NODE_W, ex));
	}
	return { x: ex, y: ey };
}
