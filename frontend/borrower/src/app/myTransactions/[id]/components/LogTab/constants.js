import dayjs from 'dayjs';
import { formatStr } from '@repo/shared-methods/utils';

/** @typedef {'all' | 'info_changes' | 'documents_uploaded' | 'imported_google' | 'synced_google' | 'status_change'} ActivityFilterKey */

export const ACTIVITY_LOG_ITEM_HEIGHT = 148;

export const ACTIVITY_FILTER_OPTIONS = [
	{ key: 'all', label: 'All activities' },
	{ key: 'info_changes', label: 'Info changes' },
	{ key: 'documents_uploaded', label: 'Documents uploaded' },
	{ key: 'imported_google', label: 'Imported from Google Drive' },
	{ key: 'synced_google', label: 'Synced to Google Drive' },
	{ key: 'status_change', label: 'Status change' },
];

/** Figma Activity summary 计数（接口就绪后改由 API 返回） */
export const MOCK_ACTIVITY_SUMMARY = {
	all: 23,
	info_changes: 3,
	documents_uploaded: 4,
	imported_google: 5,
	synced_google: 6,
	status_change: 7,
};

const MOCK_ACTIVITY_LOG_SOURCE = [
	{
		id: 1,
		title: 'Transaction created',
		category: 'info_changes',
		operator_name: 'Sarah Chen',
		operator_initials: 'SC',
		created_at: '2026-05-01T12:30:00',
		fields: [
			{ label: 'Transaction type', value: 'Purchase' },
			{ label: 'Property address', value: '1482 Elmwood Drive, San Mateo, CA, 94403' },
		],
	},
	{
		id: 2,
		title: 'Basic information updated',
		category: 'info_changes',
		operator_name: 'Sarah Chen',
		operator_initials: 'SC',
		created_at: '2026-05-01T12:30:00',
		fields: [
			{ label: 'COE', value: 1845000, format: 'amount' },
			{ label: 'PM', value: 'Nancy' },
		],
	},
	{
		id: 3,
		title: 'Document uploaded',
		category: 'documents_uploaded',
		operator_name: 'Sarah Chen',
		operator_initials: 'SC',
		created_at: '2026-05-01T12:30:00',
		fields: [{ label: 'Payment Receipt', value: 'Filename.pdf' }],
	},
	{
		id: 4,
		title: 'Imported from Google Drive',
		category: 'imported_google',
		operator_name: 'Sarah Chen',
		operator_initials: 'SC',
		created_at: '2026-05-01T12:30:00',
		fields: [{ label: 'Payment Receipt', value: 'Filename.pdf' }],
	},
	{
		id: 5,
		title: 'Sync to Google Drive',
		category: 'synced_google',
		operator_name: 'Sarah Chen',
		operator_initials: 'SC',
		created_at: '2026-05-01T12:30:00',
		fields: [{ label: 'Payment Receipt', value: 'Filename.pdf' }],
	},
	{
		id: 6,
		title: 'Status change',
		category: 'status_change',
		operator_name: 'Sarah Chen',
		operator_initials: 'SC',
		created_at: '2026-05-01T12:30:00',
		fields: [{ label: 'Status', value: 'Closed' }],
	},
];

const filterLogs = (logs, { activityType, dateRange }) => {
	let result = logs;
	if (activityType && activityType !== 'all') {
		result = result.filter((item) => item.category === activityType);
	}
	if (dateRange?.length === 2) {
		const start = dayjs(dateRange[0], formatStr);
		const end = dayjs(dateRange[1], formatStr);
		result = result.filter((item) => {
			const logDate = dayjs(item.created_at);
			return (
				start.isValid() &&
				end.isValid() &&
				logDate.isValid() &&
				!logDate.isBefore(start, 'day') &&
				!logDate.isAfter(end, 'day')
			);
		});
	}
	return result;
};

/**
 * 活动日志分页（mock；后续替换为 transaction activity API）
 */
export async function fetchActivityLogs({ page = 1, limit = 10, activityType, dateRange } = {}) {
	const filtered = filterLogs(MOCK_ACTIVITY_LOG_SOURCE, { activityType, dateRange });
	const start = (page - 1) * limit;
	const list = filtered.slice(start, start + limit);
	return {
		list,
		total: filtered.length,
	};
}
