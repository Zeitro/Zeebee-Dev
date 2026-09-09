/**
 * 我的任务列表：POST `/project_management/task/my`。
 * 筛选只传 keyword / phase（treeSelect 多选数组）/ priority（字符串，如 High）。
 */
'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Empty, Tag, Button, message } from 'antd';
import { ReloadOutlined, SearchOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import {
	createRequest,
	displayValue,
	formatDate,
	transformArrayStructure,
} from '@repo/shared-methods/utils';
import { TASK_MY_API, createTaskApi } from '@repo/project-management/api/task';
import {
	TASK_PRIORITY_TAG_COLOR,
	getPhaseI18nKey,
	getTaskPriorityI18nKey,
	localizeSelectOptions,
} from '@repo/project-management/enum/task';
import { toTaskTypeColorHex } from '@repo/project-management/enum/taskType';
import { useOptionsStore } from '@repo/project-management/store/options';
import { useAuth } from '@repo/shared-methods/hooks';

/**
 * @param {Array<{ value?: string, label?: string }>} options
 * @param {string} value
 * @returns {string|undefined}
 */
function getOptionLabel(options, value) {
	return options?.find((item) => item.value === value)?.label;
}

/**
 * 只提交 keyword / phase / priority；空值不传。
 * @param {object} raw
 * @returns {{ keyword?: string, phase?: string[], priority?: string }}
 */
function pickMyTaskFilter(raw = {}) {
	const filter = {};
	const keyword = typeof raw.keyword === 'string' ? raw.keyword.trim() : '';
	if (keyword) filter.keyword = keyword;

	const phase = Array.isArray(raw.phase) ? raw.phase.filter(Boolean) : [];
	if (phase.length) filter.phase = phase;

	if (typeof raw.priority === 'string' && raw.priority) {
		filter.priority = raw.priority;
	}
	return filter;
}

export function useMyTaskList(PERMISSIONS) {
	const t = useTranslations();
	const router = useRouter();
	const { hasPermission } = useAuth();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const tableFilterRef = useRef();

	const [keyword, setKeyword] = useState('');
	const [priority, setPriority] = useState();
	const [phase, setPhase] = useState([]);
	const [messageApi, messageContextHolder] = message.useMessage();
	const priorityOptions = useOptionsStore((state) => state.priorityOptions);
	const stageOptions = useOptionsStore((state) => state.stageOptions);

	const localizedPriorityOptions = useMemo(
		() => localizeSelectOptions(priorityOptions, t, getTaskPriorityI18nKey),
		[priorityOptions, t],
	);

	const localizedStageOptions = useMemo(
		() => localizeSelectOptions(stageOptions, t, getPhaseI18nKey),
		[stageOptions, t],
	);

	const phaseTreeData = useMemo(
		() =>
			transformArrayStructure({
				sourceArray: localizedStageOptions,
				outputKey: 'title',
			}).map((item) => ({ ...item, key: item.value })),
		[localizedStageOptions],
	);

	const handleRefresh = useCallback(() => {
		cancelMatch(TASK_MY_API);
		tableFilterRef.current?.handleSearch?.();
	}, [cancelMatch]);

	const getData = useCallback(
		async (params) => {
			try {
				const res = await taskApi.page({
					page: params?.page || 1,
					limit: params?.limit || 20,
					filter: pickMyTaskFilter(params?.filter),
				});
				const list = Array.isArray(res?.list) ? res.list : [];
				return { list, total: res?.total ?? 0 };
			} catch {
				messageApi.error(t('myTask.loadFailed'));
				return { list: [], total: 0 };
			}
		},
		[messageApi, t, taskApi],
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
				render: (value, record) =>
					hasPermission(PERMISSIONS?.PROJECT_LIST?.DETAILS) ? (
						<Button
							type="link"
							className="font-medium !p-0 h-auto"
							onClick={() => goTaskDetail(record)}
						>
							{displayValue(value)}
						</Button>
					) : (
						displayValue(value)
					),
			},
			{
				title: t('myTask.statusColumn'),
				dataIndex: 'phase',
				key: 'phase',
				width: 120,
				render: (value) => {
					const key = getPhaseI18nKey(value);
					return (
						<Tag>
							{key
								? t(key)
								: displayValue(
										getOptionLabel(localizedStageOptions, value) || value,
									)}
						</Tag>
					);
				},
			},
			{
				title: t('myTask.priorityColumn'),
				dataIndex: 'priority',
				key: 'priority',
				width: 90,
				render: (value) => {
					const key = getTaskPriorityI18nKey(value);
					return (
						<Tag color={TASK_PRIORITY_TAG_COLOR[value] || 'default'}>
							{key
								? t(key)
								: displayValue(
										getOptionLabel(localizedPriorityOptions, value) || value,
									)}
						</Tag>
					);
				},
			},
			{
				title: t('myTask.typeColumn'),
				key: 'task_type',
				width: 120,
				render: (_, record) => {
					const name = record.task_type?.name;
					const hex = toTaskTypeColorHex(record.task_type?.color);
					if (!name) return displayValue();
					return hex ? <Tag color={hex}>{name}</Tag> : displayValue(name);
				},
			},
			{
				title: t('myTask.creatorColumn'),
				dataIndex: 'creator_name',
				key: 'creator_name',
				width: 120,
				render: (value) => displayValue(value),
			},
			{
				title: t('myTask.startAtColumn'),
				dataIndex: 'start_at',
				key: 'start_at',
				width: 120,
				render: (value) => displayValue(value, formatDate),
			},
			{
				title: t('myTask.dueDateColumn'),
				dataIndex: 'end_at',
				key: 'end_at',
				width: 120,
				render: (value, record) =>
					displayValue(value, (v) => (
						<span className={record.overdue ? 'text-base-error' : undefined}>
							{formatDate(v)}
						</span>
					)),
			},
		],
		[goTaskDetail, hasPermission, localizedPriorityOptions, localizedStageOptions, t],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: t('myTask.searchPlaceholder'),
				paramsKey: 'keyword',
				value: keyword,
				onChange: ({ target: { value } }) => setKeyword(value),
				style: { width: 240 },
			},
			{
				filterType: 'treeSelect',
				label: t('myTask.statusColumn'),
				paramsKey: 'phase',
				placeholder: t('myTask.statusColumn'),
				value: phase,
				onChange: (value) => setPhase(value || []),
				treeData: phaseTreeData,
				style: { width: 160 },
				popupMatchSelectWidth: 172,
			},
			{
				filterType: 'select',
				placeholder: t('myTask.priorityColumn'),
				paramsKey: 'priority',
				value: priority,
				onChange: (value) => setPriority(value),
				options: localizedPriorityOptions,
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
					cancelMatch(TASK_MY_API);
					setKeyword('');
					setPriority(undefined);
					setPhase([]);
					tableFilterRef.current?.handleSearch?.({
						filter: {
							keyword: '',
							phase: undefined,
							priority: undefined,
						},
					});
				},
			},
		],
		[
			cancelMatch,
			handleRefresh,
			keyword,
			localizedPriorityOptions,
			phase,
			phaseTreeData,
			priority,
			t,
		],
	);

	return {
		tableFilterRef,
		filterList,
		columns,
		getData,
		messageContextHolder,
		// emptyText: <Empty description={t('myTask.empty')} />,
	};
}
