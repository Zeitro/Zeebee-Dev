/**
 * 项目操作日志：解析 `/operation/page` 返回的 operation + metadata，生成可读文案。
 * metadata 约定：仅 after=添加、仅 before=删除、两者皆有=修改；
 * 修改时 before 为修改前、after 为修改后。
 */

import { formatDate, formatTime } from '@repo/shared-methods/utils';
import { TASK_DEPENDENCY_RULE } from '@repo/project-management/enum/task';

/** 富文本是否无有效文本内容 */
function isRichTextEmpty(html) {
	if (!html) return true;
	return html.replace(/<[^>]+>/g, '').trim() === '';
}

/** 后端 OperationAction 常量，与接口返回值一致 */
export const OperationAction = {
	PROJECT_CREATED: 'Project Created',
	PROJECT_UPDATED: 'Project Updated',
	PROJECT_DELETED: 'Project Deleted',
	PROJECT_OWNER_TRANSFERRED: 'Project Owner Transferred',
	PROJECT_MEMBER_ADDED: 'Project Member Added',
	PROJECT_MEMBER_ROLE_UPDATED: 'Project Member Role Updated',
	PROJECT_MEMBER_REMOVED: 'Project Member Removed',
	PROJECT_TASK_CREATED: 'Project Task Created',
	PROJECT_TASK_UPDATED: 'Project Task Updated',
	PROJECT_TASK_DELETED: 'Project Task Deleted',
	PROJECT_TASK_TRANSITIONED: 'Project Task Transitioned',
	PROJECT_TASK_ASSIGNEE_ADDED: 'Project Task Assignee Added',
	PROJECT_TASK_ASSIGNEE_REMOVED: 'Project Task Assignee Removed',
	PROJECT_TASK_ATTACHMENT_ADDED: 'Project Task Attachment Added',
	PROJECT_TASK_ATTACHMENT_DELETED: 'Project Task Attachment Deleted',
	PROJECT_TASK_TYPE_CREATED: 'Project Task Type Created',
	PROJECT_TASK_TYPE_UPDATED: 'Project Task Type Updated',
	PROJECT_TASK_TYPE_DELETED: 'Project Task Type Deleted',
	PROJECT_TASK_DEPENDENCY_CREATED: 'Project Task Dependency Created',
	PROJECT_TASK_DEPENDENCY_UPDATED: 'Project Task Dependency Updated',
	PROJECT_TASK_DEPENDENCY_DELETED: 'Project Task Dependency Deleted',
	PROJECT_WORKFLOW_CREATED: 'Project Workflow Created',
	PROJECT_WORKFLOW_UPDATED: 'Project Workflow Updated',
	PROJECT_WORKFLOW_DELETED: 'Project Workflow Deleted',
	PROJECT_WORKFLOW_STAGE_CREATED: 'Project Workflow Stage Created',
	PROJECT_WORKFLOW_STAGE_UPDATED: 'Project Workflow Stage Updated',
	PROJECT_WORKFLOW_STAGE_DELETED: 'Project Workflow Stage Deleted',
	PROJECT_WORKFLOW_TRANSITION_CREATED: 'Project Workflow Transition Created',
	PROJECT_WORKFLOW_TRANSITION_UPDATED: 'Project Workflow Transition Updated',
	PROJECT_WORKFLOW_TRANSITION_DELETED: 'Project Workflow Transition Deleted',
};

/** operation 字符串 → i18n key（operationLog.operations.*） */
const OPERATION_I18N_KEYS = {
	[OperationAction.PROJECT_CREATED]: 'operationLog.operations.projectCreated',
	[OperationAction.PROJECT_UPDATED]: 'operationLog.operations.projectUpdated',
	[OperationAction.PROJECT_DELETED]: 'operationLog.operations.projectDeleted',
	[OperationAction.PROJECT_OWNER_TRANSFERRED]: 'operationLog.operations.projectOwnerTransferred',
	[OperationAction.PROJECT_MEMBER_ADDED]: 'operationLog.operations.projectMemberAdded',
	[OperationAction.PROJECT_MEMBER_ROLE_UPDATED]:
		'operationLog.operations.projectMemberRoleUpdated',
	[OperationAction.PROJECT_MEMBER_REMOVED]: 'operationLog.operations.projectMemberRemoved',
	[OperationAction.PROJECT_TASK_CREATED]: 'operationLog.operations.projectTaskCreated',
	[OperationAction.PROJECT_TASK_UPDATED]: 'operationLog.operations.projectTaskUpdated',
	[OperationAction.PROJECT_TASK_DELETED]: 'operationLog.operations.projectTaskDeleted',
	[OperationAction.PROJECT_TASK_TRANSITIONED]: 'operationLog.operations.projectTaskTransitioned',
	[OperationAction.PROJECT_TASK_ASSIGNEE_ADDED]:
		'operationLog.operations.projectTaskAssigneeAdded',
	[OperationAction.PROJECT_TASK_ASSIGNEE_REMOVED]:
		'operationLog.operations.projectTaskAssigneeRemoved',
	[OperationAction.PROJECT_TASK_ATTACHMENT_ADDED]:
		'operationLog.operations.projectTaskAttachmentAdded',
	[OperationAction.PROJECT_TASK_ATTACHMENT_DELETED]:
		'operationLog.operations.projectTaskAttachmentDeleted',
	[OperationAction.PROJECT_TASK_TYPE_CREATED]: 'operationLog.operations.projectTaskTypeCreated',
	[OperationAction.PROJECT_TASK_TYPE_UPDATED]: 'operationLog.operations.projectTaskTypeUpdated',
	[OperationAction.PROJECT_TASK_TYPE_DELETED]: 'operationLog.operations.projectTaskTypeDeleted',
	[OperationAction.PROJECT_TASK_DEPENDENCY_CREATED]:
		'operationLog.operations.projectTaskDependencyCreated',
	[OperationAction.PROJECT_TASK_DEPENDENCY_UPDATED]:
		'operationLog.operations.projectTaskDependencyUpdated',
	[OperationAction.PROJECT_TASK_DEPENDENCY_DELETED]:
		'operationLog.operations.projectTaskDependencyDeleted',
	[OperationAction.PROJECT_WORKFLOW_CREATED]: 'operationLog.operations.projectWorkflowCreated',
	[OperationAction.PROJECT_WORKFLOW_UPDATED]: 'operationLog.operations.projectWorkflowUpdated',
	[OperationAction.PROJECT_WORKFLOW_DELETED]: 'operationLog.operations.projectWorkflowDeleted',
	[OperationAction.PROJECT_WORKFLOW_STAGE_CREATED]:
		'operationLog.operations.projectWorkflowStageCreated',
	[OperationAction.PROJECT_WORKFLOW_STAGE_UPDATED]:
		'operationLog.operations.projectWorkflowStageUpdated',
	[OperationAction.PROJECT_WORKFLOW_STAGE_DELETED]:
		'operationLog.operations.projectWorkflowStageDeleted',
	[OperationAction.PROJECT_WORKFLOW_TRANSITION_CREATED]:
		'operationLog.operations.projectWorkflowTransitionCreated',
	[OperationAction.PROJECT_WORKFLOW_TRANSITION_UPDATED]:
		'operationLog.operations.projectWorkflowTransitionUpdated',
	[OperationAction.PROJECT_WORKFLOW_TRANSITION_DELETED]:
		'operationLog.operations.projectWorkflowTransitionDeleted',
};

/** metadata 字段名 → i18n key（operationLog.fields.*） */
const FIELD_I18N_KEYS = {
	title: 'operationLog.fields.title',
	description: 'operationLog.fields.description',
	priority: 'operationLog.fields.priority',
	task_type_name: 'operationLog.fields.taskType',
	type: 'operationLog.fields.type',
	work_hours: 'operationLog.fields.workHours',
	estimated_hours: 'operationLog.fields.estimatedHours',
	start_at: 'operationLog.fields.startAt',
	end_at: 'operationLog.fields.endAt',
	start_date: 'operationLog.fields.startDate',
	due_date: 'operationLog.fields.dueDate',
	stage_name: 'operationLog.fields.stage',
	status: 'operationLog.fields.status',
	workflow_state: 'operationLog.fields.status',
	name: 'operationLog.fields.name',
	role: 'operationLog.fields.role',
	member_name: 'operationLog.fields.member',
	assignee_name: 'operationLog.fields.assignee',
	file_name: 'operationLog.fields.fileName',
	attachment_name: 'operationLog.fields.attachmentName',
	dependency_type: 'operationLog.fields.dependencyType',
	dependency_rule: 'operationLog.fields.dependencyRule',
	color: 'operationLog.fields.color',
	sort_order: 'operationLog.fields.sortOrder',
};

/** 字符串优先级 → i18n key */
const PRIORITY_VALUE_KEYS = {
	Low: 'priority.low',
	Medium: 'priority.medium',
	High: 'priority.high',
	Urgent: 'priority.critical',
	Critical: 'priority.critical',
};

/** 依赖规则 FS / FF / SS / SF → kanban 已有文案 */
const DEPENDENCY_RULE_I18N_KEYS = {
	[TASK_DEPENDENCY_RULE.FINISH_TO_START]: 'kanban.dependencyFinishToStart',
	[TASK_DEPENDENCY_RULE.FINISH_TO_FINISH]: 'kanban.dependencyFinishToFinish',
	[TASK_DEPENDENCY_RULE.START_TO_START]: 'kanban.dependencyStartToStart',
	[TASK_DEPENDENCY_RULE.START_TO_FINISH]: 'kanban.dependencyStartToFinish',
};

const TRANSITION_OPERATIONS = new Set([OperationAction.PROJECT_TASK_TRANSITIONED]);

/**
 * @param {string} key
 * @returns {boolean}
 */
function isIdField(key) {
	if (!key || key === 'id') return true;
	return /_id(s)?$/i.test(key);
}

/**
 * @param {object|null|undefined} obj
 * @returns {boolean}
 */
function hasPayload(obj) {
	return obj != null && typeof obj === 'object' && Object.keys(obj).length > 0;
}

/**
 * @param {object|null|undefined} metadata
 * @returns {'add'|'delete'|'modify'|'none'}
 */
function getMetadataChangeType(metadata) {
	const before = metadata?.before;
	const after = metadata?.after;
	const hasBefore = hasPayload(before);
	const hasAfter = hasPayload(after);
	if (hasBefore && hasAfter) return 'modify';
	if (hasAfter) return 'add';
	if (hasBefore) return 'delete';
	return 'none';
}

/**
 * 修改：metadata.before=旧值，metadata.after=新值。
 * @param {object|null|undefined} metadata
 * @returns {Array<{ field: string, oldValue: *, newValue: * }>}
 */
function getModifiedFields(metadata) {
	const oldState = metadata?.before || {};
	const newState = metadata?.after || {};
	const keys = new Set([...Object.keys(oldState), ...Object.keys(newState)]);
	const changes = [];

	keys.forEach((key) => {
		if (isIdField(key)) return;
		const oldValue = oldState[key];
		const newValue = newState[key];
		if (oldValue === newValue) return;
		if (oldValue === undefined && newValue === undefined) return;
		changes.push({ field: key, oldValue, newValue });
	});

	return changes;
}

/**
 * @param {object|null|undefined} metadata
 * @returns {Array<{ field: string, value: * }>}
 */
function getAddedFields(metadata) {
	const state = metadata?.after || {};
	return Object.entries(state)
		.filter(([key]) => !isIdField(key))
		.map(([field, value]) => ({ field, value }));
}

/**
 * @param {object|null|undefined} metadata
 * @returns {Array<{ field: string, value: * }>}
 */
function getDeletedFields(metadata) {
	const state = metadata?.before || {};
	return Object.entries(state)
		.filter(([key]) => !isIdField(key))
		.map(([field, value]) => ({ field, value }));
}

/**
 * @param {string} operation
 * @returns {string}
 */
function getOperationI18nKey(operation) {
	return OPERATION_I18N_KEYS[operation] || 'operationLog.operations.unknown';
}

/**
 * @param {string} field
 * @returns {string}
 */
function getFieldI18nKey(field) {
	return FIELD_I18N_KEYS[field] || 'operationLog.fields.unknown';
}

/**
 * @param {string} field
 * @param {*} value
 * @param {(key: string, params?: object) => string} t
 * @param {{ resolveUserName?: (id: number) => string }} [options]
 * @returns {string}
 */
function formatFieldValue(field, value, t, options = {}) {
	const { resolveUserName } = options;
	if (value === null || value === undefined || value === '') {
		return t('operationLog.empty');
	}

	if (field === 'priority') {
		const key = PRIORITY_VALUE_KEYS[value];
		if (key) return t(key);
		const idx = Number(value);
		if (!Number.isNaN(idx)) {
			const labels = [
				t('priority.low'),
				t('priority.medium'),
				t('priority.high'),
				t('priority.critical'),
			];
			return labels[idx] ?? String(value);
		}
	}

	if (field === 'description' && typeof value === 'string') {
		return isRichTextEmpty(value) ? t('operationLog.empty') : t('operationLog.hasContent');
	}

	if (field === 'dependency_rule') {
		const ruleKey = DEPENDENCY_RULE_I18N_KEYS[String(value).toUpperCase()];
		if (ruleKey) return t(ruleKey);
	}

	if (
		(field === 'assignee_id' ||
			field === 'assignee_ids' ||
			field === 'observer' ||
			field === 'user_id') &&
		resolveUserName
	) {
		return String(value)
			.split(',')
			.map((segment) => {
				const id = Number(segment.trim());
				return !Number.isNaN(id) && id > 0 ? resolveUserName(id) : segment;
			})
			.join(', ');
	}

	if (field === 'start_date' || field === 'due_date') {
		return formatDate(value) || String(value);
	}

	if (
		field === 'start_at' ||
		field === 'end_at' ||
		(typeof value === 'string' && /^\d{4}-\d{2}-\d{2}/.test(value))
	) {
		return formatTime(value) || String(value);
	}

	if (field === 'work_hours' || field === 'estimated_hours') {
		return `${value}${t('kanban.hoursSuffix')}`;
	}

	return String(value);
}

/**
 * @param {string} field
 * @param {(key: string, params?: object) => string} t
 * @returns {string}
 */
function getFieldLabel(field, t) {
	return t(getFieldI18nKey(field), { field });
}

/**
 * @param {{ field: string, oldValue: *, newValue: * }} change
 * @param {(key: string, params?: object) => string} t
 * @param {{ resolveUserName?: (id: number) => string, isTransition?: boolean }} [options]
 * @returns {string}
 */
function formatFieldChange(change, t, options = {}) {
	const { field, oldValue, newValue } = change;
	const fieldLabel = getFieldLabel(field, t);
	const oldText = formatFieldValue(field, oldValue, t, options);
	const newText = formatFieldValue(field, newValue, t, options);
	const templateKey = options.isTransition
		? 'operationLog.changeTransitioned'
		: 'operationLog.changeUpdated';

	return t(templateKey, {
		field: fieldLabel,
		old: oldText ?? t('operationLog.none'),
		new: newText ?? t('operationLog.empty'),
	});
}

/**
 * @param {{ field: string, value: * }} item
 * @param {(key: string, params?: object) => string} t
 * @param {{ resolveUserName?: (id: number) => string }} [options]
 * @returns {string}
 */
function formatAddedField(item, t, options = {}) {
	return t('operationLog.changeAdded', {
		field: getFieldLabel(item.field, t),
		value: formatFieldValue(item.field, item.value, t, options),
	});
}

/**
 * @param {{ field: string, value: * }} item
 * @param {(key: string, params?: object) => string} t
 * @param {{ resolveUserName?: (id: number) => string }} [options]
 * @returns {string}
 */
function formatDeletedField(item, t, options = {}) {
	return t('operationLog.changeDeleted', {
		field: getFieldLabel(item.field, t),
		value: formatFieldValue(item.field, item.value, t, options),
	});
}

/**
 * @param {object} act 单条操作记录
 * @param {(key: string, params?: object) => string} t
 * @param {{ resolveUserName?: (id: number) => string }} [options]
 * @returns {string}
 */
function formatOperationDescription(act, t, options = {}) {
	const operation = act?.operation || '';
	const metadata = act?.metadata || {};
	const operationText = t(getOperationI18nKey(operation), { operation });
	const changeType = getMetadataChangeType(metadata);
	const isTransition = TRANSITION_OPERATIONS.has(operation);

	let details = [];

	if (changeType === 'modify') {
		details = getModifiedFields(metadata).map((change) =>
			formatFieldChange(change, t, { ...options, isTransition }),
		);
	} else if (changeType === 'add') {
		details = getAddedFields(metadata).map((item) => formatAddedField(item, t, options));
	} else if (changeType === 'delete') {
		details = getDeletedFields(metadata).map((item) => formatDeletedField(item, t, options));
	}

	if (!details.length) return operationText;
	return `${operationText} (${details.join(t('operationLog.detailSeparator'))})`;
}

/** 日志解析与文案封装，便于集中维护 operation / metadata 相关逻辑 */
export const OperationLog = {
	OperationAction,
	isIdField,
	getMetadataChangeType,
	getModifiedFields,
	getAddedFields,
	getDeletedFields,
	getOperationI18nKey,
	getFieldI18nKey,
	formatFieldValue,
	formatFieldChange,
	formatOperationDescription,
};

export { formatOperationDescription };
