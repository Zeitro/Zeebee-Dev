/**
 * 工作流列颜色 / 类型：与后端 `states.color`、`states.type` 对齐。
 * 色值走已有 CSS token；teal 无独立 token，回落 success。
 */

export const WORKFLOW_STATE_TYPE = {
	INITIAL: 'initial',
	IN_PROGRESS: 'in_progress',
	COMPLETED: 'completed',
	PAUSED: 'paused',
	ANY: 'any',
};

export const WORKFLOW_STATE_COLOR_VAR = {
	gray: 'var(--gray-medium)',
	blue: 'var(--blue-medium)',
	purple: 'var(--purple-dark)',
	teal: 'var(--base-success)',
	green: 'var(--green-medium)',
	orange: 'var(--orange-dark)',
	red: 'var(--base-error)',
};

/**
 * @param {string|null|undefined} color
 * @returns {string}
 */
export function getWorkflowStateColor(color) {
	return WORKFLOW_STATE_COLOR_VAR[color] || WORKFLOW_STATE_COLOR_VAR.gray;
}

/**
 * 初始列才展示「加任务」；兼容 type 与 is_initial。
 * @param {object|null|undefined} state
 * @returns {boolean}
 */
export function isInitialState(state) {
	return state?.type === WORKFLOW_STATE_TYPE.INITIAL || state?.is_initial === true;
}

/**
 * 终态任务标题划线。
 * @param {object|null|undefined} state
 * @returns {boolean}
 */
export function isTerminalState(state) {
	return state?.is_terminal === true || state?.type === WORKFLOW_STATE_TYPE.COMPLETED;
}
