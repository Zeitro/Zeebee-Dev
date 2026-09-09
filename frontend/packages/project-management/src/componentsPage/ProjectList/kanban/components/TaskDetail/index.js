/**
 * 任务详情抽屉：自持数据/导航/日志/建子任务。
 * 宿主最小接口：taskId + projectId + onClose；可选 onTaskIdChange / onMutated / seedTask。
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
import CreateTaskModal from '../CreateTaskModal';
import { TASK_DETAIL_DRAWER_WIDTH } from '../constants';
import {
	getAssigneeIds,
	getMemberLabel,
	isTaskDetailFormDirty,
	normalizeRichTextValue,
} from '../helpers';
import { isTerminalState } from '../../../../../enum/workflow';
import { PriorityIcon } from './PriorityIcon';
import { SubtaskList } from './SubtaskList';
import { TaskActivityDrawer } from './TaskActivityDrawer';
import { TaskAttachments } from './TaskAttachments';
import { TaskDependencies } from './TaskDependencies';
import { TaskDetailPanel, TaskDetailPanelAction } from './TaskDetailPanel';
import { TaskDetailProperties } from './TaskDetailProperties';
import { TaskDetailSection } from './TaskDetailSection';
import { TaskDescription } from './TaskDescription';
import { TaskDrawerHeader } from './TaskDrawerHeader';
import { useTaskDetail } from './useTaskDetail';
import { useTaskDetailNav } from './useTaskDetailNav';

export { TaskActivityDrawer } from './TaskActivityDrawer';

/**
 * @param {{
 *   taskId?: number,
 *   projectId?: number|string,
 *   seedTask?: object|null,
 *   onClose: () => void,
 *   onTaskIdChange?: (taskId?: number) => void,
 *   onMutated?: () => Promise<void>|void,
 * }} props
 */
export function TaskDetailDrawer({
	taskId,
	projectId,
	seedTask = null,
	onClose,
	onTaskIdChange,
	onMutated,
}) {
	const t = useTranslations();
	const locale = useLocale();
	const [form] = Form.useForm();
	const [dependencyFocusMode, setDependencyFocusMode] = useState(null);
	const [detailTask, setDetailTask] = useState(null);
	const [detailLoading, setDetailLoading] = useState(false);
	const [logOpen, setLogOpen] = useState(false);
	const [createOpen, setCreateOpen] = useState(false);
	const detailRequestRef = useRef(0);
	const open = Boolean(taskId);

	const {
		messageApi,
		messageContextHolder,
		project,
		workflow,
		members,
		canManage,
		taskTypes,
		persistAssignees,
		saveTask,
		saveTitle,
		deleteTask,
		deleting,
		fetchTaskDetail,
		fetchTaskTypes,
		fetchTaskCandidates,
		updateTaskDetail,
	} = useTaskDetail({
		projectId,
		enabled: open,
		onMutated,
	});

	const displayTask = detailTask ?? (open ? seedTask || { id: taskId } : null);
	const { previousTask, backVia, navigateTo, goBack } = useTaskDetailNav({
		taskId,
		currentTask: displayTask,
		onTaskIdChange,
	});

	const subtasks = useMemo(() => {
		if (!displayTask?.id) return [];
		return Array.isArray(displayTask.children) ? displayTask.children : [];
	}, [displayTask]);

	const subtaskCompletedCount = useMemo(
		() =>
			subtasks.filter((item) => {
				const state = workflow?.states?.find((row) => row.id === item.workflow_state_id);
				return isTerminalState(state);
			}).length,
		[subtasks, workflow?.states],
	);

	const relatedTasks = useMemo(() => {
		const list = [...subtasks];
		if (displayTask?.id) list.push(displayTask);
		if (previousTask?.id) list.push(previousTask);
		return list;
	}, [displayTask, previousTask, subtasks]);

	const backTask =
		previousTask ||
		(displayTask?.parent_id
			? { id: displayTask.parent_id, title: displayTask.parent_title }
			: null);

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
		predicate: () => Boolean(open && taskId && projectId),
		notPassedCallback: () => {
			setDetailTask(null);
			setLogOpen(false);
			setCreateOpen(false);
		},
		immediate: () => {
			const initial = seedTask?.id === taskId ? seedTask : { id: taskId };
			setDetailTask(initial);
			form.setFieldsValue({
				title: initial.title,
				description: initial.description,
				priority: toTaskPriorityValue(initial.priority),
				type: initial.type,
				workflow_state_id: initial.workflow_state_id,
				assignee_ids: getAssigneeIds(initial),
				start_date: initial.start_date ? formatDate(initial.start_date) : undefined,
				due_date: initial.due_date ? formatDate(initial.due_date) : undefined,
				estimated_hours: initial.estimated_hours,
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

			Promise.all([priorityPromise, fetchTaskDetail(taskId, projectId), fetchTaskTypes()])
				.then(([, detail]) => {
					if (requestId !== detailRequestRef.current || !detail) return;
					const merged = { ...(seedTask?.id === taskId ? seedTask : null), ...detail };
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
		deps: [fetchTaskDetail, fetchTaskTypes, form, getOptions, open, projectId, taskId],
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

	const handleDelete = () => deleteTask(displayTask, onClose);
	const handleTitleBlur = (title) => saveTitle(displayTask, title);
	const handleClose = () => {
		setLogOpen(false);
		setCreateOpen(false);
		onClose();
	};

	return (
		<>
			{messageContextHolder}
			<Drawer
				title={null}
				width={TASK_DETAIL_DRAWER_WIDTH}
				open={open}
				onClose={handleClose}
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
								backVia={backVia}
								assigneeIds={liveAssigneeIds}
								members={members}
								onBack={goBack}
								onTitleBlur={handleTitleBlur}
								onDelete={handleDelete}
								deleteLoading={deleting}
								onOpenLog={() => setLogOpen(true)}
								onAddDependency={setDependencyFocusMode}
								canManage={canManage}
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
										projectId={projectId}
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
										projectId={projectId}
										tasks={relatedTasks}
										onTaskClick={(task) => navigateTo(task, 'dependency')}
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
												onClick={() => setCreateOpen(true)}
											/>
										}
									>
										<SubtaskList
											parentId={displayTask.id}
											tasks={subtasks.map((item) => ({
												...item,
												parent_id: item.parent_id ?? displayTask.id,
											}))}
											workflow={workflow}
											onTaskClick={(task) => navigateTo(task, 'subtask')}
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

			<TaskActivityDrawer
				open={logOpen}
				projectId={projectId}
				taskId={taskId}
				members={members}
				onClose={() => setLogOpen(false)}
			/>
			<CreateTaskModal
				open={createOpen}
				projectId={projectId}
				parentTask={displayTask}
				messageApi={messageApi}
				onCancel={() => setCreateOpen(false)}
				onSuccess={async () => {
					setCreateOpen(false);
					const detail = await fetchTaskDetail(taskId, projectId);
					if (detail) applyDetailTask({ ...displayTask, ...detail });
					await onMutated?.();
				}}
			/>
		</>
	);
}
