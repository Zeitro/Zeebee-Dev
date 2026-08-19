/**
 * 我的任务列表：标题/优先级/状态筛选，FilterTable 分页。
 * 任务范围：创建者或执行者为当前用户；接口就绪后切 GET /tasks。
 */
'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Empty, Tag, Button, message } from 'antd';
import { ReloadOutlined, SearchOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest, displayValue, formatDate, dayjs } from '@repo/shared-methods/utils';
import { TASK_API } from '@repo/project-management/api/task';
import {
	PROJECT_PRIORITY_FILTER_VALUES,
	PROJECT_PRIORITY_I18N_KEY,
	getPriorityI18nKey,
} from '@repo/project-management/enum/project';
import { WORKFLOW_STATE_TYPE } from '@repo/project-management/enum/workflow';
import { getMockMyTaskPage, getMockProjectName, MOCK_MEMBER_MAP } from '../mockData';

const PRIORITY_TAG_COLOR = {
	0: 'green',
	1: 'default',
	2: 'blue',
	3: 'orange',
};

const STATUS_FILTER_OPTIONS = [
	{ value: WORKFLOW_STATE_TYPE.INITIAL, labelKey: 'myTask.statusInitial' },
	{ value: WORKFLOW_STATE_TYPE.IN_PROGRESS, labelKey: 'myTask.statusInProgress' },
	{ value: WORKFLOW_STATE_TYPE.COMPLETED, labelKey: 'myTask.statusCompleted' },
	{ value: WORKFLOW_STATE_TYPE.PAUSED, labelKey: 'myTask.statusPaused' },
];

/**
 * @param {number|string|undefined} userId
 * @param {number|string|undefined} memberId
 * @param {Record<number, string>} memberMap
 * @param {string|undefined} userName
 * @returns {string}
 */
function resolveMemberName(userId, memberId, memberMap, userName) {
	if (memberId == null || memberId === '') return '';
	if (String(memberId) === String(userId)) {
		return userName || memberMap[memberId] || '';
	}
	return memberMap[memberId] || '';
}

export function useMyTaskList() {
	const t = useTranslations();
	const router = useRouter();
	const { cancelMatch } = useMemo(() => createRequest(), []);
	const userId = useUserStore((state) => state.userInfo?.id);
	const userName = useUserStore((state) => state.userInfo?.full_name);
	const tableFilterRef = useRef();

	const [title, setTitle] = useState('');
	const [priority, setPriority] = useState();
	const [statusType, setStatusType] = useState();
	const [messageApi, messageContextHolder] = message.useMessage();

	const priorityOptions = useMemo(
		() =>
			PROJECT_PRIORITY_FILTER_VALUES.map((value) => ({
				value,
				label: t(PROJECT_PRIORITY_I18N_KEY[value]),
			})),
		[t],
	);

	const statusOptions = useMemo(
		() =>
			STATUS_FILTER_OPTIONS.map(({ value, labelKey }) => ({
				value,
				label: t(labelKey),
			})),
		[t],
	);

	const handleRefresh = useCallback(() => {
		cancelMatch(TASK_API);
		tableFilterRef.current?.handleSearch?.();
	}, [cancelMatch]);

	const getData = useCallback(
		async (params) => {
			const filter = { ...(params?.filter || {}) };
			if (filter.priority !== undefined && filter.priority !== '') {
				filter.priority = Number(filter.priority);
			} else {
				delete filter.priority;
			}
			if (!filter.statusType) {
				delete filter.statusType;
			}

			try {
				// API:
				// const res = await request(
				// 	`${TASK_API}?${new URLSearchParams({
				// 		page: String(params?.page || 1),
				// 		page_size: String(params?.limit || 15),
				// 		user_id: String(userId || ''),
				// 		...(filter.title ? { title: filter.title } : {}),
				// 		...(filter.priority != null ? { priority: String(filter.priority) } : {}),
				// 		...(filter.statusType ? { status_type: filter.statusType } : {}),
				// 	})}`,
				// 	{ method: 'GET', isHideTips: true },
				// );
				// const list = res?.list ?? res?.data ?? [];
				// return { list, total: res?.total ?? list.length };

				return getMockMyTaskPage({ ...params, filter }, userId, userName);
			} catch {
				messageApi.error(t('myTask.loadFailed'));
				return { list: [], total: 0 };
			}
		},
		[messageApi, t, userId, userName],
	);

	const goTaskDetail = useCallback(
		(record) => {
			if (record?.project_id == null || record?.id == null) return;
			router.push(`/project/projectList/${record.project_id}?task=${record.id}`);
		},
		[router],
	);

	const columns = useMemo(
		() => [
			{
				title: t('myTask.titleColumn'),
				dataIndex: 'title',
				key: 'title',
				width: 260,
				ellipsis: true,
				render: (value, record) => (
					<Button
						type="link"
						className="font-medium !p-0 h-auto"
						onClick={() => goTaskDetail(record)}
					>
						{displayValue(value)}
					</Button>
				),
			},
			{
				title: t('myTask.projectColumn'),
				key: 'project',
				width: 140,
				render: (_, record) =>
					displayValue(record.project_name || getMockProjectName(record.project_id)),
			},
			{
				title: t('myTask.statusColumn'),
				key: 'status',
				width: 100,
				render: (_, record) => <Tag>{displayValue(record.status)}</Tag>,
			},
			{
				title: t('myTask.priorityColumn'),
				dataIndex: 'priority',
				key: 'priority',
				width: 80,
				render: (value) => {
					const key = getPriorityI18nKey(value);
					if (!key) return displayValue(value);
					return <Tag color={PRIORITY_TAG_COLOR[value] || 'default'}>{t(key)}</Tag>;
				},
			},
			{
				title: t('myTask.typeColumn'),
				key: 'type',
				width: 80,
				render: (_, record) => displayValue(record.type_label || record.type),
			},
			{
				title: t('myTask.creatorColumn'),
				key: 'creator',
				width: 100,
				render: (_, record) =>
					displayValue(
						record.creator_name ||
							resolveMemberName(userId, record.creator_id, MOCK_MEMBER_MAP, userName),
					),
			},
			{
				title: t('myTask.assigneeColumn'),
				key: 'assignee',
				width: 120,
				render: (_, record) => {
					if (record.assignee_id == null) return displayValue();
					const name =
						record.assignee_name ||
						resolveMemberName(userId, record.assignee_id, MOCK_MEMBER_MAP, userName);
					if (!name) return displayValue();
					const isMe = String(record.assignee_id) === String(userId);
					return (
						<span className={isMe ? 'font-semibold text-base-primary' : undefined}>
							{name}
							{isMe ? t('myTask.meSuffix') : ''}
						</span>
					);
				},
			},
			{
				title: t('myTask.dueDateColumn'),
				dataIndex: 'due_date',
				key: 'due_date',
				width: 120,
				render: (value) => {
					if (!value) return displayValue();
					const isOverdue = dayjs(value).isBefore(dayjs(), 'day');
					return (
						<span className={isOverdue ? 'text-base-error' : undefined}>
							{formatDate(value, 'YYYY/M/D')}
						</span>
					);
				},
			},
			{
				title: t('myTask.roleColumn'),
				key: 'role',
				width: 120,
				render: (_, record) => {
					const tags = [];
					if (String(record.creator_id) === String(userId)) {
						tags.push({
							key: 'creator',
							label: t('myTask.roleCreator'),
							color: 'purple',
						});
					}
					if (String(record.assignee_id) === String(userId)) {
						tags.push({
							key: 'assignee',
							label: t('myTask.roleAssignee'),
							color: 'cyan',
						});
					}
					if (!tags.length) return displayValue();
					return tags.map((item) => (
						<Tag key={item.key} color={item.color} className="mr-1">
							{item.label}
						</Tag>
					));
				},
			},
		],
		[goTaskDetail, t, userId, userName],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: t('myTask.searchPlaceholder'),
				paramsKey: 'title',
				value: title,
				onChange: ({ target: { value } }) => setTitle(value),
				style: { width: 240 },
			},
			{
				filterType: 'select',
				placeholder: t('myTask.priorityColumn'),
				paramsKey: 'priority',
				value: priority,
				onChange: (value) => setPriority(value),
				options: priorityOptions,
				style: { width: 120 },
			},
			{
				filterType: 'select',
				placeholder: t('myTask.statusColumn'),
				paramsKey: 'statusType',
				value: statusType,
				onChange: (value) => setStatusType(value),
				options: statusOptions,
				style: { width: 120 },
			},
			{
				filterType: 'button',
				text: t('common.search'),
				icon: <SearchOutlined />,
				onClick: handleRefresh,
			},
			{
				filterType: 'button',
				type: 'default',
				text: t('common.reset'),
				icon: <ReloadOutlined />,
				onClick: () => {
					cancelMatch(TASK_API);
					setTitle('');
					setPriority(undefined);
					setStatusType(undefined);
					tableFilterRef.current?.handleSearch?.({
						filter: {
							title: '',
							priority: undefined,
							statusType: undefined,
						},
					});
				},
			},
		],
		[
			cancelMatch,
			handleRefresh,
			priority,
			priorityOptions,
			statusOptions,
			statusType,
			t,
			title,
		],
	);

	return {
		tableFilterRef,
		filterList,
		columns,
		getData,
		messageContextHolder,
		emptyText: <Empty description={t('myTask.empty')} />,
	};
}
