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

/** SVG / 色块选择器用 hex，与 `WORKFLOW_STATE_COLOR_VAR` 键一致 */
export const WORKFLOW_STATE_COLOR_HEX = {
	gray: '#8c8c8c',
	blue: '#1890ff',
	purple: '#722ed1',
	teal: '#13c2c2',
	green: '#52c41a',
	orange: '#fa8c16',
	red: '#f5222d',
};

export const WORKFLOW_STATE_COLOR_OPTIONS = Object.keys(WORKFLOW_STATE_COLOR_HEX);

/** 与 `WORKFLOW_STATE_COLOR_HEX` 键一致，对应 `workflow.editor.color*` i18n */
export const WORKFLOW_STATE_COLOR_LABEL_KEYS = {
	gray: 'workflow.editor.colorGray',
	blue: 'workflow.editor.colorBlue',
	purple: 'workflow.editor.colorPurple',
	teal: 'workflow.editor.colorTeal',
	green: 'workflow.editor.colorGreen',
	orange: 'workflow.editor.colorOrange',
	red: 'workflow.editor.colorRed',
};

export const WORKFLOW_STATE_TYPE_OPTIONS = [
	{ value: WORKFLOW_STATE_TYPE.INITIAL, labelKey: 'workflow.editor.stateTypeInitial' },
	{ value: WORKFLOW_STATE_TYPE.IN_PROGRESS, labelKey: 'workflow.editor.stateTypeInProgress' },
	{ value: WORKFLOW_STATE_TYPE.COMPLETED, labelKey: 'workflow.editor.stateTypeCompleted' },
	{ value: WORKFLOW_STATE_TYPE.PAUSED, labelKey: 'workflow.editor.stateTypePaused' },
];

export const WORKFLOW_TRANSITION_TRIGGER_OPTIONS = [
	{ value: 'manual', labelKey: 'workflow.editor.triggerManual' },
	{ value: 'automatic', labelKey: 'workflow.editor.triggerAutomatic' },
	{ value: 'scheduled', labelKey: 'workflow.editor.triggerScheduled' },
];

export const WORKFLOW_TRANSITION_MEMBER_OPTIONS = [
	{ value: 'assignee', labelKey: 'workflow.editor.memberAssignee' },
	{ value: 'creator', labelKey: 'workflow.editor.memberCreator' },
	{ value: 'owner', labelKey: 'workflow.editor.memberOwner' },
	{ value: 'admin', labelKey: 'workflow.editor.memberAdmin' },
];

/**
 * @param {string|null|undefined} color
 * @returns {string}
 */
export function getWorkflowStateColor(color) {
	return WORKFLOW_STATE_COLOR_VAR[color] || WORKFLOW_STATE_COLOR_VAR.gray;
}

/**
 * @param {string|null|undefined} color
 * @returns {string}
 */
export function getWorkflowStateColorHex(color) {
	return WORKFLOW_STATE_COLOR_HEX[color] || WORKFLOW_STATE_COLOR_HEX.gray;
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
