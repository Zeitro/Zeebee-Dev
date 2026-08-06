/**
 * 将 operations API 单条记录（operation + metadata.after）映射为 ActivityLogItem 展示行。
 * 新增 operation 类型时，在 ACTIVITY_LOG_OPERATION_FIELDS 追加字段配置即可。
 */

/** 与后端 operations 列表项 operation 字段一致 */
export const ACTIVITY_LOG_OPERATION = {
	STATUS_CHANGED: 'Status changed',
	// Google Drive sync: the detail shape is special (counts/folder link/per-file list),
	// so it does not use the generic metadata.after label/value and is rendered
	// dedicatedly by ActivityLogItem.
	SYNCED_TO_DRIVE: 'Synced to Google Drive',
};

/**
 * @typedef {object} ActivityLogFieldConfig
 * @property {string} label 左侧展示标签
 * @property {string} key metadata.after 中的字段名
 * @property {'amount'} [format] 传给 ActivityLogItem 的值格式化类型
 */

/** @type {Record<string, ActivityLogFieldConfig[]>} */
const ACTIVITY_LOG_OPERATION_FIELDS = {
	[ACTIVITY_LOG_OPERATION.STATUS_CHANGED]: [{ label: 'Status', key: 'status' }],
	// 示例：后续 Document created 等在此追加
	// [ACTIVITY_LOG_OPERATION.DOCUMENT_CREATED]: [{ label: 'Document', key: 'name' }],
};

/**
 * 按 operation 从 metadata.after 取出需展示的字段行。
 * @param {object|null|undefined} item operations list 单条（含 operation、metadata）
 * @returns {{ label: string, value: unknown, format?: string }[]}
 */
export function getActivityLogDisplayFields(item) {
	const configs = ACTIVITY_LOG_OPERATION_FIELDS[item?.operation];
	if (!configs?.length) {
		return [];
	}

	const after = item?.metadata?.after;
	if (!after || typeof after !== 'object') {
		return [];
	}

	return configs.map(({ label, key, format }) => ({
		label,
		value: after[key],
		...(format ? { format } : {}),
	}));
}
