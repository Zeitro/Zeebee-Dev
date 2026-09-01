/**
 * 任务阶段 / 优先级：对齐下拉接口 value（如 Pending、InProgress、High、Critical）。
 * 请求仍传接口原值；展示文案走 i18n，便于中英切换。
 */

import { PROJECT_PRIORITY } from './project';

export const TASK_PHASE = {
	PENDING: 'Pending',
	IN_PROGRESS: 'InProgress',
	FINISHED: 'Finished',
	COMPLETED: 'Completed',
	PAUSED: 'Paused',
};

export const TASK_PHASE_I18N_KEY = {
	[TASK_PHASE.PENDING]: 'myTask.phasePending',
	[TASK_PHASE.IN_PROGRESS]: 'myTask.phaseInProgress',
	'In Progress': 'myTask.phaseInProgress',
	[TASK_PHASE.FINISHED]: 'myTask.phaseFinished',
	[TASK_PHASE.COMPLETED]: 'myTask.phaseCompleted',
	[TASK_PHASE.PAUSED]: 'myTask.phasePaused',
};

export const TASK_PRIORITY = {
	LOW: 'Low',
	MEDIUM: 'Medium',
	HIGH: 'High',
	CRITICAL: 'Critical',
};

export const TASK_PRIORITY_I18N_KEY = {
	[TASK_PRIORITY.LOW]: 'priority.low',
	[TASK_PRIORITY.MEDIUM]: 'priority.medium',
	[TASK_PRIORITY.HIGH]: 'priority.high',
	[TASK_PRIORITY.CRITICAL]: 'priority.critical',
};

export const TASK_PRIORITY_TAG_COLOR = {
	[TASK_PRIORITY.LOW]: 'green',
	[TASK_PRIORITY.MEDIUM]: 'default',
	[TASK_PRIORITY.HIGH]: 'blue',
	[TASK_PRIORITY.CRITICAL]: 'orange',
	[TASK_PRIORITY.URGENT]: 'orange',
};

/** 看板数字优先级（0–3）→ 接口 / 下拉原值（Low、Medium…） */
export const PROJECT_TO_TASK_PRIORITY = {
	[PROJECT_PRIORITY.LOW]: TASK_PRIORITY.LOW,
	[PROJECT_PRIORITY.MEDIUM]: TASK_PRIORITY.MEDIUM,
	[PROJECT_PRIORITY.HIGH]: TASK_PRIORITY.HIGH,
	[PROJECT_PRIORITY.URGENT]: TASK_PRIORITY.CRITICAL,
};

/** 接口 / 下拉原值 → 看板数字优先级（供 PriorityIcon 等） */
export const TASK_TO_PROJECT_PRIORITY = {
	[TASK_PRIORITY.LOW]: PROJECT_PRIORITY.LOW,
	[TASK_PRIORITY.MEDIUM]: PROJECT_PRIORITY.MEDIUM,
	[TASK_PRIORITY.HIGH]: PROJECT_PRIORITY.HIGH,
	[TASK_PRIORITY.CRITICAL]: PROJECT_PRIORITY.URGENT,
};

/**
 * 表单 / 保存用接口 priority 字符串；看板内部可能是 0–3 数字。
 * @param {string|number|null|undefined} priority
 * @returns {string}
 */
export function toTaskPriorityValue(priority) {
	if (typeof priority === 'string' && TASK_PRIORITY_I18N_KEY[priority]) return priority;
	if (typeof priority === 'number' && !Number.isNaN(priority)) {
		return PROJECT_TO_TASK_PRIORITY[priority] ?? TASK_PRIORITY.MEDIUM;
	}
	return TASK_PRIORITY.MEDIUM;
}

/**
 * @param {string|number|null|undefined} priority
 * @returns {number}
 */
export function toProjectPriorityLevel(priority) {
	if (typeof priority === 'number' && !Number.isNaN(priority)) return priority;
	if (typeof priority === 'string' && priority) {
		return TASK_TO_PROJECT_PRIORITY[priority] ?? PROJECT_PRIORITY.MEDIUM;
	}
	return PROJECT_PRIORITY.MEDIUM;
}

/**
 * @param {string|null|undefined} phase
 * @returns {string|undefined}
 */
export function getPhaseI18nKey(phase) {
	if (phase === null || phase === undefined || phase === '') return undefined;
	return TASK_PHASE_I18N_KEY[phase];
}

/**
 * @param {string|null|undefined} priority
 * @returns {string|undefined}
 */
export function getTaskPriorityI18nKey(priority) {
	if (priority === null || priority === undefined || priority === '') return undefined;
	return TASK_PRIORITY_I18N_KEY[priority];
}

/**
 * 接口 options 的 value 不变，label 换成当前语言。
 * @param {Array<{ value?: string, label?: string }>} options
 * @param {(key: string) => string} t
 * @param {(value: string) => string|undefined} getI18nKey
 * @returns {Array<{ value?: string, label?: string }>}
 */
export function localizeSelectOptions(options, t, getI18nKey) {
	return (options || []).map((item) => {
		const key = getI18nKey(item.value);
		return {
			...item,
			label: key ? t(key) : item.label || item.value,
		};
	});
}

/** 任务依赖方向：详情 `dependencies[].dependency_type` */
export const TASK_DEPENDENCY_DIRECTION = {
	PRE_TASK: 'PreTask',
	POST_TASK: 'PostTask',
};

/** 任务依赖规则：详情 `dependencies[].dependency_rule`（FS / FF / SS / SF） */
export const TASK_DEPENDENCY_RULE = {
	FINISH_TO_START: 'FS',
	FINISH_TO_FINISH: 'FF',
	START_TO_START: 'SS',
	START_TO_FINISH: 'SF',
};
