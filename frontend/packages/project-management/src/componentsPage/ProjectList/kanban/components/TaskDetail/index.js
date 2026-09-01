/**
 * 任务详情抽屉（本目录入口）：标题失焦保存；执行者/依赖即时走 PUT 全量更新；附件/日志。
 */
'use client';

import { useMemo, useRef, useState } from 'react';
import {
	CalendarOutlined,
	CloseOutlined,
	ExclamationCircleFilled,
	PlusOutlined,
	UserOutlined,
} from '@ant-design/icons';
import { Button, Drawer, Form, Spin, Tooltip } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { displayValue, formatDate, formatTime } from '@repo/shared-methods/utils';
import {
	getTaskPriorityI18nKey,
	localizeSelectOptions,
	toProjectPriorityLevel,
	toTaskPriorityValue,
} from '@repo/project-management/enum/task';
import { useOptionsStore } from '@repo/project-management/store/options';
import { useKanbanTaskDetail } from '../../hooks/useKanbanTaskDetail';
import { TASK_DETAIL_DRAWER_WIDTH } from '../constants';
import {
	getAssigneeIds,
	getMemberLabel,
	getSubtasks,
	isTaskDetailFormDirty,
	normalizeRichTextValue,
} from '../helpers';
import { isTerminalState } from '../../../../../enum/workflow';
import { PriorityIcon } from './PriorityIcon';
import { SubtaskList } from './SubtaskList';
import { TaskAttachments } from './TaskAttachments';
import { TaskDependencies } from './TaskDependencies';
import { TaskDetailPanel, TaskDetailPanelAction } from './TaskDetailPanel';
import { TaskDetailProperties } from './TaskDetailProperties';
import { TaskDetailSection } from './TaskDetailSection';
import { TaskDescription } from './TaskDescription';
import { TaskDrawerHeader } from './TaskDrawerHeader';

export { TaskActivityDrawer } from './TaskActivityDrawer';

/**
 * @param {{
 *   project?: object|null,
 *   workflow?: object|null,
 *   tasks?: Array,
 *   memberOptions?: Array,
 *   members?: Array,
 *   onClose: () => void,
 *   previousTask?: object|null,
 *   onTaskClick: (task: object) => void,
 *   onBack?: () => void,
 *   onCreateSubtask: (parent: object) => void,
 *   onOpenLog?: () => void,
 * }} props
 */
export function TaskDetailDrawer({
	project,
	workflow,
	tasks = [],
	members = [],
	onClose,
	previousTask = null,
	onTaskClick,
	onBack,
	onCreateSubtask,
	onOpenLog,
	messageApi,
}) {
	const t = useTranslations();
	const locale = useLocale();
	const [form] = Form.useForm();
	const [dependencyFocusMode, setDependencyFocusMode] = useState(null);
	const [detailTask, setDetailTask] = useState(null);
	const [detailLoading, setDetailLoading] = useState(false);
	const detailRequestRef = useRef(0);
	const {
		activeTask: task,
		taskTypes,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
		fetchTaskDetail,
		fetchTaskTypes,
		fetchTaskCandidates,
		updateTaskDetail,
	} = useKanbanTaskDetail({ messageApi });
	const open = Boolean(task);
	const displayTask = detailTask ?? task;
	const subtaskTasks = useMemo(() => {
		const children = displayTask?.children;
		if (!Array.isArray(children) || !children.length) return tasks;
		const existingIds = new Set(tasks.map((item) => item.id));
		const extra = children.filter((item) => !existingIds.has(item.id));
		return extra.length ? [...tasks, ...extra] : tasks;
	}, [displayTask?.children, tasks]);
	const parentTask = tasks.find((item) => item.id === displayTask?.parent_id);
	const backTask = previousTask || parentTask;
	const { subtasks, subtaskCompletedCount } = useMemo(() => {
		const list = displayTask?.id
			? displayTask.children?.length
				? displayTask.children
				: getSubtasks(subtaskTasks, displayTask.id)
			: [];
		const completed = list.filter((item) => {
			const state = workflow?.states?.find((row) => row.id === item.workflow_state_id);
			return isTerminalState(state);
		}).length;
		return { subtasks: list, subtaskCompletedCount: completed };
	}, [displayTask, subtaskTasks, workflow?.states]);
	const watchedAssigneeIds = Form.useWatch('assignee_ids', form);
	const liveAssigneeIds = watchedAssigneeIds ?? getAssigneeIds(displayTask);
	const watchedFormValues = Form.useWatch([], form);
	const hasUnsavedChanges = useMemo(() => {
		if (!displayTask || detailLoading) return false;
		return isTaskDetailFormDirty(displayTask, watchedFormValues);
	}, [detailLoading, displayTask, watchedFormValues]);
	const creatorLabel = useMemo(() => {
		if (!displayTask) return '';
		if (displayTask.creator_name) return displayTask.creator_name;
		const member = members.find(
			(item) => item.id === displayTask.creator_id || item.user_id === displayTask.creator_id,
		);
		return getMemberLabel(member, displayTask.creator_id);
	}, [displayTask, members]);

	const priorityOptionsRaw = useOptionsStore((state) => state.priorityOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);
	const priorityOptions = useMemo(
		() =>
			localizeSelectOptions(priorityOptionsRaw, t, getTaskPriorityI18nKey).map((item) => ({
				...item,
				label: (
					<span className="inline-flex items-center gap-1.5">
						<PriorityIcon level={toProjectPriorityLevel(item.value)} />
						{item.label}
					</span>
				),
			})),
		[priorityOptionsRaw, t],
	);

	useDeferredSyncEffect({
		predicate: () => Boolean(task && open),
		notPassedCallback: () => {
			setDetailTask(null);
		},
		immediate: () => {
			setDetailTask(task);
			form.setFieldsValue({
				title: task.title,
				description: task.description,
				priority: toTaskPriorityValue(task.priority),
				type: task.type,
				workflow_state_id: task.workflow_state_id,
				assignee_ids: getAssigneeIds(task),
				start_date: task.start_date ? formatDate(task.start_date) : undefined,
				due_date: task.due_date ? formatDate(task.due_date) : undefined,
				estimated_hours: task.estimated_hours,
			});
		},
		deferred: () => {
			setDependencyFocusMode(null);
			const requestId = detailRequestRef.current + 1;
			detailRequestRef.current = requestId;
			setDetailLoading(true);

			const priorityPromise = getOptions(
				'project_management_priority',
				'priorityOptions',
			).catch(() => {});

			const detailPromise =
				task?.id && project?.id
					? fetchTaskDetail(task.id, project.id)
					: Promise.resolve(null);

			Promise.all([priorityPromise, detailPromise, fetchTaskTypes()])
				.then(([, detail]) => {
					if (requestId !== detailRequestRef.current || !detail) return;
					const merged = { ...task, ...detail };
					setDetailTask(merged);
					form.setFieldsValue({
						title: merged.title,
						description: merged.description,
						priority: toTaskPriorityValue(merged.priority),
						type: merged.type,
						workflow_state_id: merged.workflow_state_id,
						assignee_ids: getAssigneeIds(merged),
						start_date: merged.start_date ? formatDate(merged.start_date) : undefined,
						due_date: merged.due_date ? formatDate(merged.due_date) : undefined,
						estimated_hours: merged.work_hours ?? merged.estimated_hours,
					});
				})
				.finally(() => {
					if (requestId === detailRequestRef.current) {
						setDetailLoading(false);
					}
				});
		},
		destroyCallback: () => {
			detailRequestRef.current += 1;
		},
		deps: [fetchTaskDetail, fetchTaskTypes, form, getOptions, open, project?.id, task?.id],
		isUseMicrotask: true,
	});

	const applyDetailTask = (updated) => {
		if (!updated) return;
		setDetailTask(updated);
		form.setFieldsValue({
			title: updated.title,
			description: updated.description,
			priority: toTaskPriorityValue(updated.priority),
			type: updated.type,
			workflow_state_id: updated.workflow_state_id,
			assignee_ids: getAssigneeIds(updated),
			start_date: updated.start_date ? formatDate(updated.start_date) : undefined,
			due_date: updated.due_date ? formatDate(updated.due_date) : undefined,
			estimated_hours: updated.work_hours ?? updated.estimated_hours,
		});
	};

	const handleSave = async () => {
		const values = form.getFieldsValue();
		const updated = await saveTask(displayTask, values);
		applyDetailTask(updated);
	};

	const handleDelete = () => {
		deleteTask(onClose);
	};

	return (
		<Drawer
			title={null}
			width={TASK_DETAIL_DRAWER_WIDTH}
			open={open}
			onClose={onClose}
			closeIcon={
				<Tooltip
					title={hasUnsavedChanges ? t('kanban.unsavedChangesCloseHint') : undefined}
				>
					<CloseOutlined />
				</Tooltip>
			}
			styles={{ body: { padding: 0, background: 'var(--base-text2)' } }}
			destroyOnHidden
			footer={
				<div className="flex items-center justify-between gap-4">
					<div className="flex items-center gap-4 min-w-0 text-xs text-base-label">
						<span className="inline-flex items-center gap-1.5 shrink-0">
							<UserOutlined className="text-base-disabled text-xs" />
							{t('kanban.taskCreator', { name: displayValue(creatorLabel) })}
						</span>
						<span className="inline-flex items-center gap-1.5 shrink-0">
							<CalendarOutlined className="text-base-disabled text-xs" />
							{t('kanban.taskCreatedAt', {
								time: displayValue(displayTask?.created_at, formatTime),
							})}
						</span>
					</div>
					<Button
						type="primary"
						onClick={handleSave}
						icon={
							hasUnsavedChanges ? (
								<Tooltip title={t('kanban.unsavedChangesHint')}>
									<ExclamationCircleFilled
										style={{ color: 'var(--base-warning)' }}
									/>
								</Tooltip>
							) : null
						}
					>
						{t('common.save')}
					</Button>
				</div>
			}
			keyboard={false}
			maskClosable={false}
		>
			{displayTask ? (
				<Spin spinning={detailLoading} wrapperClassName="flex flex-col h-full">
					<div className="flex flex-col h-full bg-base-text2">
						<TaskDrawerHeader
							task={displayTask}
							project={project}
							backTask={backTask}
							assigneeIds={liveAssigneeIds}
							members={members}
							onBack={onBack}
							onTitleBlur={saveTitle}
							onDelete={handleDelete}
							onOpenLog={onOpenLog}
							onAddDependency={setDependencyFocusMode}
						/>

						<div className="flex-1 overflow-y-auto px-6 py-5 pb-24">
							<Form
								form={form}
								layout="vertical"
								onValuesChange={(changed) => {
									if ('assignee_ids' in changed) {
										persistAssignees(
											displayTask,
											changed.assignee_ids || [],
										).then(applyDetailTask);
									}
								}}
							>
								<TaskDetailProperties
									t={t}
									workflow={workflow}
									taskTypes={taskTypes}
									priorityOptions={priorityOptions}
									projectId={project?.id}
								/>

								<TaskDetailSection title={t('kanban.description')}>
									<Form.Item
										name="description"
										noStyle
										getValueFromEvent={normalizeRichTextValue}
									>
										<TaskDescription
											key={`${displayTask.id}-${open}`}
											locale={locale}
											placeholder={t('kanban.descriptionPlaceholder')}
										/>
									</Form.Item>
								</TaskDetailSection>

								<TaskDependencies
									key={`${displayTask.id}-${dependencyFocusMode || 'idle'}`}
									task={displayTask}
									projectId={project?.id}
									tasks={tasks}
									onTaskClick={onTaskClick}
									fetchCandidates={fetchTaskCandidates}
									onSaveDependencies={async (dependencies) => {
										const updated = await updateTaskDetail(displayTask, {
											dependencies,
										});
										applyDetailTask(updated);
										return updated;
									}}
									initialAddingMode={dependencyFocusMode}
									messageApi={messageApi}
								/>

								<TaskDetailPanel
									title={t('kanban.subtasks')}
									titleSuffix={
										subtasks.length > 0
											? `${subtaskCompletedCount}/${subtasks.length}`
											: undefined
									}
									isEmpty={subtasks.length === 0}
									emptyText={t('kanban.subtasksEmpty')}
									action={
										<TaskDetailPanelAction
											title={t('kanban.addSubtask')}
											icon={<PlusOutlined />}
											onClick={() => onCreateSubtask(displayTask)}
										/>
									}
								>
									<SubtaskList
										parentId={displayTask.id}
										tasks={subtaskTasks}
										workflow={workflow}
										onTaskClick={onTaskClick}
									/>
								</TaskDetailPanel>

								<TaskAttachments
									task={displayTask}
									saveTask={saveTask}
									updateTaskDetail={updateTaskDetail}
									messageApi={messageApi}
									onTaskUpdated={applyDetailTask}
								/>
							</Form>
						</div>
					</div>
				</Spin>
			) : null}
		</Drawer>
	);
}
