/**
 * 项目任务表格：自拉任务、列头排序、前端分页。
 * 打开任务走 URL `?task=`。排序由 columns.sorter 处理已加载列表。
 */
'use client';

import { useMemo, useState } from 'react';
import { Button, Input, Spin, Tag } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { NewTable, NoData } from '@repo/shared-methods';
import { displayValue, formatDate, formatTime, getPagination } from '@repo/shared-methods/utils';
import {
	getInitialTablePageSize,
	setStoredTablePageSize,
} from '@repo/shared-methods/lib/tablePagination';
import { getPriorityI18nKey, getPriorityTagClass } from '@repo/project-management/enum/project';
import { toTaskTypeColorHex } from '@repo/project-management/enum/taskType';
import { getWorkflowStateColor } from '@repo/project-management/enum/workflow';
import { useProjectDetailQuery } from '../hooks/useProjectDetailQuery';
import { useProjectTaskCatalog } from '../hooks/useProjectTaskCatalog';
import {
	buildProjectTaskTree,
	countChildTasks,
	filterProjectTasks,
	taskColumnSorter,
} from '../lib/projectTableRows';
import { spanDays } from '../lib/taskSchedule';
import { ProjectViewFrame } from './ProjectViewFrame';
import styles from './ProjectTableView.module.scss';

/**
 * @param {object} task
 * @returns {string}
 */
function assigneeText(task) {
	return (task.assignees || [])
		.map((item) => item.organization_user_full_name || item.full_name || item.name)
		.filter(Boolean)
		.join(', ');
}

/**
 * @param {{ projectId?: number|string, workflow?: object|null }} props
 */
export function ProjectTableView({ projectId, workflow }) {
	const t = useTranslations();
	const { setTaskQuery } = useProjectDetailQuery();
	const { tasks, loading, truncated } = useProjectTaskCatalog(projectId);
	const [keyword, setKeyword] = useState('');
	const [page, setPage] = useState(1);
	const [pageSize, setPageSize] = useState(() => getInitialTablePageSize());

	const tree = useMemo(
		() => buildProjectTaskTree(filterProjectTasks(tasks, keyword)),
		[keyword, tasks],
	);
	const childCounts = useMemo(() => countChildTasks(tasks), [tasks]);
	const pageCount = Math.max(1, Math.ceil(tree.length / pageSize) || 1);
	const current = Math.min(page, pageCount);
	const stageById = useMemo(
		() => new Map((workflow?.states || []).map((state) => [state.id, state])),
		[workflow],
	);

	const columns = useMemo(
		() => [
			{
				title: t('myTask.titleColumn'),
				dataIndex: 'title',
				key: 'title',
				width: 260,
				fixed: 'left',
				ellipsis: true,
				render: (value, record) => (
					<Button
						type="link"
						className="h-auto max-w-full p-0 font-medium"
						onClick={() => setTaskQuery(record.id)}
					>
						<span className="truncate">{displayValue(value)}</span>
						<span className="ml-2 text-xs font-normal text-base-label">
							{t('projectDetail.spanDays', { count: spanDays(record) })}
						</span>
					</Button>
				),
			},
			{
				title: t('myTask.statusColumn'),
				key: 'stage',
				width: 130,
				render: (_, record) => {
					const state =
						stageById.get(record.workflow_state_id) || stageById.get(record.stage_id);
					const name = state?.name || record.stage_name;
					if (!name) return displayValue();
					return (
						<span className="inline-flex items-center gap-1.5">
							<span
								className="inline-block h-2 w-2 rounded-full"
								style={{ background: getWorkflowStateColor(state?.color) }}
							/>
							{name}
						</span>
					);
				},
			},
			{
				title: t('myTask.priorityColumn'),
				dataIndex: 'priority',
				key: 'priority',
				width: 100,
				sorter: taskColumnSorter('priority'),
				render: (value) => {
					const key = getPriorityI18nKey(value);
					if (!key) return displayValue();
					return <Tag className={getPriorityTagClass(value)}>{t(key)}</Tag>;
				},
			},
			{
				title: t('myTask.typeColumn'),
				key: 'type',
				width: 120,
				ellipsis: true,
				render: (_, record) => {
					const name = record.task_type?.name || record.type;
					const hex = toTaskTypeColorHex(record.task_type?.color);
					if (!name) return displayValue();
					return hex ? <Tag color={hex}>{name}</Tag> : displayValue(name);
				},
			},
			{
				title: t('kanban.assignee'),
				key: 'assignee',
				width: 160,
				ellipsis: true,
				render: (_, record) => displayValue(assigneeText(record)),
			},
			{
				title: t('myTask.startAtColumn'),
				dataIndex: 'start_date',
				width: 120,
				sorter: taskColumnSorter('start_date'),
				render: (value) => displayValue(value, formatDate),
			},
			{
				title: t('myTask.dueDateColumn'),
				dataIndex: 'due_date',
				width: 120,
				sorter: taskColumnSorter('due_date'),
				render: (value, record) =>
					displayValue(value, (date) => (
						<span className={record.overdue ? 'text-base-error' : undefined}>
							{formatDate(date)}
						</span>
					)),
			},
			{
				title: t('projectDetail.hoursColumn'),
				dataIndex: 'estimated_hours',
				width: 90,
				render: (value) =>
					displayValue(value, (hours) => t('projectDetail.hoursValue', { hours })),
			},
			{
				title: t('kanban.subtasks'),
				key: 'subtasks',
				width: 110,
				render: (_, record) =>
					t('projectDetail.subtaskCount', { count: childCounts.get(record.id) || 0 }),
			},
			{
				title: t('myTask.creatorColumn'),
				dataIndex: 'creator_name',
				width: 120,
				ellipsis: true,
				render: (value) => displayValue(value),
			},
			{
				title: t('projectDetail.createdAtColumn'),
				dataIndex: 'created_at',
				width: 170,
				sorter: taskColumnSorter('created_at'),
				defaultSortOrder: 'descend',
				render: (value) => displayValue(value, formatTime),
			},
			{
				title: t('projectDetail.updatedAtColumn'),
				dataIndex: 'updated_at',
				width: 170,
				sorter: taskColumnSorter('updated_at'),
				render: (value) => displayValue(value, formatTime),
			},
		],
		[childCounts, setTaskQuery, stageById, t],
	);

	const resetPage = () => setPage(1);

	const toolbar = (
		<>
			<Input
				allowClear
				prefix={<SearchOutlined className="text-base-label" />}
				placeholder={t('projectDetail.tableSearch')}
				value={keyword}
				onChange={(event) => {
					setKeyword(event.target.value);
					resetPage();
				}}
				className="max-w-xs"
			/>
			<span className="shrink-0 text-xs text-base-label">
				{truncated
					? t('projectDetail.catalogTruncated')
					: t('projectDetail.taskCount', { count: tree.length })}
			</span>
		</>
	);

	return (
		<ProjectViewFrame toolbar={toolbar}>
			{loading ? (
				<div className="flex flex-1 items-center justify-center">
					<Spin />
				</div>
			) : (
				<div className={styles.tableFill}>
					<NewTable
						rowKey="id"
						columns={columns}
						dataSource={tree}
						scroll={{ x: 1680 }}
						showSorterTooltip={false}
						locale={{ emptyText: <NoData text={t('projectDetail.tableEmpty')} /> }}
						pagination={getPagination({
							current,
							pageSize,
							total: tree.length,
							// 表格 size=middle/small 时 antd 会把分页强制成 small，统一按钮样式会变小且边框错位
							size: 'default',
							onChange: (nextPage, nextSize) => {
								const size = Number(nextSize);
								if (size !== pageSize) {
									setStoredTablePageSize(size);
									setPageSize(size);
									setPage(1);
									return;
								}
								setPage(nextPage);
							},
						})}
					/>
				</div>
			)}
		</ProjectViewFrame>
	);
}
