/**
 * 工作流列颜色 / 类型：与后端 `states.color`、`states.type` 对齐。
 * 预设色走 CSS token；teal 无独立 token，回落 success。
 * ColorPicker 用 hex；接口兼容命名色（blue）与自定义 hex。
 * 编辑器里 `stages.phase` 文案与任务列表不同：Pending 起始、InProgress 进行、Finished 终止。
 */

import { TASK_PHASE } from './task';

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

export const WORKFLOW_STATE_DEFAULT_COLOR = WORKFLOW_STATE_COLOR_HEX.blue;

export const WORKFLOW_STATE_COLOR_PRESETS = Object.values(WORKFLOW_STATE_COLOR_HEX);

const HEX_TO_NAME = Object.fromEntries(
	Object.entries(WORKFLOW_STATE_COLOR_HEX).map(([name, hex]) => [hex.toLowerCase(), name]),
);

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

/** 工作流编辑器 `stages.phase` 下拉文案，与 myTask 阶段文案分开 */
export const WORKFLOW_PHASE_I18N_KEY = {
	[TASK_PHASE.PENDING]: 'workflow.editor.phasePending',
	[TASK_PHASE.IN_PROGRESS]: 'workflow.editor.phaseInProgress',
	'In Progress': 'workflow.editor.phaseInProgress',
	[TASK_PHASE.FINISHED]: 'workflow.editor.phaseFinished',
	[TASK_PHASE.COMPLETED]: 'workflow.editor.phaseFinished',
	[TASK_PHASE.PAUSED]: 'workflow.editor.phasePaused',
};

/**
 * @param {string|null|undefined} phase
 * @returns {string|undefined}
 */
export function getWorkflowPhaseI18nKey(phase) {
	if (phase === null || phase === undefined || phase === '') return undefined;
	return WORKFLOW_PHASE_I18N_KEY[phase];
}

export const WORKFLOW_TRANSITION_TRIGGER_OPTIONS = [
	{ value: 'manual', labelKey: 'workflow.editor.triggerManual' },
	{ value: 'automatic', labelKey: 'workflow.editor.triggerAutomatic' },
	{ value: 'scheduled', labelKey: 'workflow.editor.triggerScheduled' },
];

export const WORKFLOW_TRANSITION_MEMBER_OPTIONS = [
	{ value: 'ProjectOwner', labelKey: 'workflow.editor.roleProjectOwner' },
	{ value: 'ProjectAdmin', labelKey: 'workflow.editor.roleProjectAdmin' },
	{ value: 'TaskExecutor', labelKey: 'workflow.editor.roleTaskExecutor' },
	{ value: 'TaskCreator', labelKey: 'workflow.editor.roleTaskCreator' },
];

/** `project_management_member_role` 的 value → i18n；对不上则用接口 label */
export const WORKFLOW_MEMBER_ROLE_I18N_KEY = {
	ProjectOwner: 'workflow.editor.roleProjectOwner',
	ProjectAdmin: 'workflow.editor.roleProjectAdmin',
	TaskExecutor: 'workflow.editor.roleTaskExecutor',
	TaskCreator: 'workflow.editor.roleTaskCreator',
};

/**
 * @param {string|null|undefined} role
 * @returns {string|undefined}
 */
export function getWorkflowMemberRoleI18nKey(role) {
	if (role === null || role === undefined || role === '') return undefined;
	return WORKFLOW_MEMBER_ROLE_I18N_KEY[role];
}

/**
 * 规范为 hex，供 ColorPicker / 色块；命名色走预设表。
 * @param {string|{toHexString?: () => string}|null|undefined} value
 * @returns {string}
 */
export function toWorkflowStateColorHex(value) {
	if (!value) return '';
	if (typeof value !== 'string') {
		return value.toHexString?.() || '';
	}
	const trimmed = value.trim();
	if (trimmed.startsWith('#')) return trimmed;
	const name = trimmed === 'grey' ? 'gray' : trimmed;
	return WORKFLOW_STATE_COLOR_HEX[name] || trimmed;
}

/**
 * 提交给创建 / 更新接口的 color：预设色用英文名（如 blue），自定义色保留 hex。
 * @param {string|{toHexString?: () => string}|null|undefined} value
 * @returns {string}
 */
export function toWorkflowStateColorPayload(value) {
	const hex = toWorkflowStateColorHex(value);
	if (!hex) return '';
	return HEX_TO_NAME[hex.toLowerCase()] || hex;
}

/**
 * 看板等走 CSS token；自定义 hex 原样返回。
 * @param {string|{toHexString?: () => string}|null|undefined} color
 * @returns {string}
 */
export function getWorkflowStateColor(color) {
	const hex = toWorkflowStateColorHex(color);
	if (!hex) return WORKFLOW_STATE_COLOR_VAR.gray;
	const name = HEX_TO_NAME[hex.toLowerCase()];
	if (name) return WORKFLOW_STATE_COLOR_VAR[name];
	return hex.startsWith('#') ? hex : WORKFLOW_STATE_COLOR_VAR.gray;
}

/**
 * SVG / Tag / 色块用 hex；命名色与自定义 hex 均可。
 * @param {string|{toHexString?: () => string}|null|undefined} color
 * @returns {string}
 */
export function getWorkflowStateColorHex(color) {
	const hex = toWorkflowStateColorHex(color);
	return hex.startsWith('#') ? hex : WORKFLOW_STATE_COLOR_HEX.gray;
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
