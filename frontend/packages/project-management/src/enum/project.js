/**
 * 项目优先级 / 状态：与后端字段 `priority`（0–3）、`status` 对齐。
 * 文案走 i18n key，勿在 UI 写死中英文。
 */

export const PROJECT_PRIORITY = {
	LOW: 0,
	MEDIUM: 1,
	HIGH: 2,
	URGENT: 3,
};

/**
 * FilterField 对 value=0 会当空，筛选项用字符串。
 * 对象 key 会被转成字符串，故 0 与 '0' 共用同一条。
 */
export const PROJECT_PRIORITY_FILTER_VALUES = ['0', '1', '2', '3'];

export const PROJECT_PRIORITY_I18N_KEY = {
	[PROJECT_PRIORITY.LOW]: 'priority.low',
	[PROJECT_PRIORITY.MEDIUM]: 'priority.medium',
	[PROJECT_PRIORITY.HIGH]: 'priority.high',
	[PROJECT_PRIORITY.URGENT]: 'priority.critical',
};

/** 卡片优先级标签：已注册 token 用 utility，未注册的用 CSS 变量 */
export const PROJECT_PRIORITY_TAG_CLASS = {
	[PROJECT_PRIORITY.LOW]: 'text-base-gray bg-gray-light border-base-border',
	[PROJECT_PRIORITY.MEDIUM]: 'text-blue-medium bg-blue-light-1 border-blue-light',
	[PROJECT_PRIORITY.HIGH]:
		'text-base-warning bg-[var(--orange-light)] border-[var(--orange-medium)]',
	[PROJECT_PRIORITY.URGENT]:
		'text-base-error bg-[var(--error-light)] border-[var(--error-medium)]',
};

export const PROJECT_STATUS = {
	PLANNING: 'planning',
	ACTIVE: 'active',
	ON_HOLD: 'on_hold',
	COMPLETED: 'completed',
	ARCHIVED: 'archived',
};

export const PROJECT_STATUS_I18N_KEY = {
	[PROJECT_STATUS.PLANNING]: 'status.planning',
	[PROJECT_STATUS.ACTIVE]: 'status.active',
	[PROJECT_STATUS.ON_HOLD]: 'status.onHold',
	[PROJECT_STATUS.COMPLETED]: 'status.completed',
	[PROJECT_STATUS.ARCHIVED]: 'status.archived',
};

/**
 * @param {number|string|null|undefined} priority
 * @returns {string|undefined}
 */
export function getPriorityI18nKey(priority) {
	if (priority === null || priority === undefined || priority === '') return undefined;
	return PROJECT_PRIORITY_I18N_KEY[priority];
}

/**
 * @param {number|string|null|undefined} priority
 * @returns {string}
 */
export function getPriorityTagClass(priority) {
	return (
		PROJECT_PRIORITY_TAG_CLASS[priority] || PROJECT_PRIORITY_TAG_CLASS[PROJECT_PRIORITY.MEDIUM]
	);
}

/**
 * @param {string|null|undefined} status
 * @returns {string|undefined}
 */
export function getStatusI18nKey(status) {
	return PROJECT_STATUS_I18N_KEY[status];
}
