/**
 * 看板纯函数：列表解包、query、执行者 id、成员展示；不改外部状态。
 */

import { dayjs } from '@repo/shared-methods/utils';
import { PROJECT_PRIORITY } from '@repo/project-management/enum';

/**
 * 兼容 `{ list }` / `{ data }` / 裸数组。
 * @param {*} res
 * @returns {Array}
 */
export function unwrapList(res) {
	if (Array.isArray(res)) return res;
	if (Array.isArray(res?.list)) return res.list;
	if (Array.isArray(res?.data)) return res.data;
	return [];
}

/**
 * @param {Record<string, string|number|undefined|null>} params
 * @returns {string}
 */
export function buildQuery(params) {
	const search = new URLSearchParams();
	Object.entries(params).forEach(([key, value]) => {
		if (value === undefined || value === null || value === '') return;
		search.set(key, String(value));
	});
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

/**
 * 优先 `assignee_ids`，否则回落 `assignee_id`。
 * @param {object|null|undefined} task
 * @returns {number[]}
 */
export function getAssigneeIds(task) {
	if (Array.isArray(task?.assignee_ids) && task.assignee_ids.length) {
		return task.assignee_ids;
	}
	if (task?.assignee_id) return [task.assignee_id];
	return [];
}

/**
 * @param {object|null|undefined} member
 * @param {number|string} [fallbackId]
 * @returns {string}
 */
export function getMemberLabel(member, fallbackId) {
	return member?.name || member?.email || (fallbackId != null ? String(fallbackId) : '');
}

/**
 * @param {object|null|undefined} member
 * @param {number|string} [fallbackId]
 * @returns {string}
 */
export function getMemberInitial(member, fallbackId) {
	const label = getMemberLabel(member, fallbackId);
	return label ? label.charAt(0).toUpperCase() : '?';
}

/**
 * DatePicker / dayjs 转后端 ISO；空值不传。
 * @param {*} value
 * @returns {string|undefined}
 */
export function toIsoDateTime(value) {
	if (!value) return undefined;
	const parsed = typeof value.format === 'function' ? value : dayjs(value);
	if (!parsed || !parsed.isValid?.()) return undefined;
	return parsed.format('YYYY-MM-DDTHH:mm:ssZ');
}

/**
 * 卡片不展示「中」优先级（与 demo 一致：priority === 1 隐藏）。
 * @param {number|string|null|undefined} priority
 * @returns {boolean}
 */
export function shouldShowCardPriority(priority) {
	if (priority === null || priority === undefined || priority === '') return false;
	return Number(priority) !== PROJECT_PRIORITY.MEDIUM;
}

/**
 * @param {Array} tasks
 * @param {number|string} taskId
 * @returns {number}
 */
export function countSubtasks(tasks, taskId) {
	if (!Array.isArray(tasks) || taskId == null) return 0;
	return tasks.filter((item) => item?.parent_id === taskId).length;
}

/**
 * @param {Array} tasks
 * @param {number|string} parentId
 * @returns {Array}
 */
export function getSubtasks(tasks, parentId) {
	if (!Array.isArray(tasks) || parentId == null) return [];
	return tasks.filter((item) => item?.parent_id === parentId);
}

/**
 * @param {number} bytes
 * @returns {string}
 */
export function formatFileSize(bytes) {
	if (bytes < 1024) return `${bytes} B`;
	if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
	return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

/**
 * @param {string} name
 * @returns {string}
 */
export function getFileExtension(name) {
	return name.split('.').pop()?.toLowerCase() || '';
}

/**
 * data URL 转 Blob，用于附件预览/下载。
 * @param {string} dataUrl
 * @returns {Blob}
 */
export function dataUrlToBlob(dataUrl) {
	const [meta, base64] = dataUrl.split(',');
	const mime = meta.match(/:(.*?);/)?.[1] || 'application/octet-stream';
	const bin = atob(base64);
	const arr = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i += 1) arr[i] = bin.charCodeAt(i);
	return new Blob([arr], { type: mime });
}

const ACTIVITY_FIELD_KEYS = {
	title: 'kanban.activityFieldTitle',
	description: 'kanban.activityFieldDescription',
	priority: 'kanban.activityFieldPriority',
	type: 'kanban.activityFieldType',
	assignee_id: 'kanban.activityFieldAssignee',
	assignee_ids: 'kanban.activityFieldAssignee',
	observer: 'kanban.activityFieldObserver',
	due_date: 'kanban.activityFieldDueDate',
	start_date: 'kanban.activityFieldStartDate',
	estimated_hours: 'kanban.activityFieldEstimatedHours',
	workflow_state: 'kanban.activityFieldStatus',
	status: 'kanban.activityFieldStatus',
};

/**
 * 将任务活动记录格式化为可读文案；`t` 来自 next-intl。
 * @param {object} act
 * @param {(key: string) => string} t
 * @param {(id: number) => string} [resolveUserName]
 * @returns {string}
 */
export function formatActivityDescription(act, t, resolveUserName) {
	const priorityLabels = [
		t('priority.low'),
		t('priority.medium'),
		t('priority.high'),
		t('priority.urgent'),
	];
	const formatValue = (field, val) => {
		if (!val) return val;
		if (field === 'priority') {
			const idx = Number(val);
			return Number.isNaN(idx) ? val : priorityLabels[idx];
		}
		if (field === 'due_date' || field === 'start_date') {
			return val.length >= 16 ? val.substring(0, 16) : val;
		}
		if (
			(field === 'assignee_id' || field === 'assignee_ids' || field === 'observer') &&
			resolveUserName
		) {
			return val
				.split(',')
				.map((segment) => {
					const id = Number(segment.trim());
					return !Number.isNaN(id) && id > 0 ? resolveUserName(id) : segment;
				})
				.join(', ');
		}
		return val;
	};
	const fieldLabel = ACTIVITY_FIELD_KEYS[act.field]
		? t(ACTIVITY_FIELD_KEYS[act.field])
		: act.field;
	if (act.action === 'created') return t('kanban.activityCreated');
	if (act.action === 'transitioned') {
		const oldVal = formatValue(act.field, act.old_value);
		const newVal = formatValue(act.field, act.new_value);
		return t('kanban.activityTransitioned', {
			field: fieldLabel,
			old: oldVal || t('kanban.activityNone'),
			new: newVal,
		});
	}
	if (act.action === 'updated') {
		const oldVal = formatValue(act.field, act.old_value);
		const newVal = formatValue(act.field, act.new_value);
		return t('kanban.activityUpdated', {
			field: fieldLabel,
			old: oldVal || t('kanban.activityEmpty'),
			new: newVal || t('kanban.activityEmpty'),
		});
	}
	if (act.action === 'commented') {
		return t('kanban.activityCommented', { content: act.new_value });
	}
	return act.new_value || act.action;
}
