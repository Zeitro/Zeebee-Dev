/**
 * 任务详情抽屉（本目录入口）：标题失焦保存；状态走 transition；执行者即时落库；依赖/附件/日志。
 */
'use client';

import { useMemo, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Drawer, Form } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { dayjs } from '@repo/shared-methods/utils';
import {
	PROJECT_PRIORITY_FILTER_VALUES,
	PROJECT_PRIORITY_I18N_KEY,
} from '../../../../../enum/project';
import { TASK_DETAIL_DRAWER_WIDTH } from '../constants';
import { getAssigneeIds, getSubtasks, normalizeRichTextValue } from '../helpers';
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
 *   open: boolean,
 *   task?: object|null,
 *   project?: object|null,
 *   workflow?: object|null,
 *   tasks?: Array,
 *   taskTypes?: Array,
 *   memberOptions?: Array,
 *   members?: Array,
 *   onClose: () => void,
 *   onTaskClick: (task: object) => void,
 *   onCreateSubtask: (parent: object) => void,
 *   onTitleBlur: (title: string) => void,
 *   onPersistAssignees: (ids: number[]) => void,
 *   onSave: (values: object) => Promise<void>,
 *   onDelete: () => void,
 *   onOpenLog?: () => void,
 *   fetchTaskDependencies: (taskId: number) => Promise<Array>,
 *   fetchTaskDependents: (taskId: number) => Promise<Array>,
 *   onAddDependency: (taskId: number, dependsOnTaskId: number) => Promise<void>,
 *   onAddDependent: (taskId: number, dependentTaskId: number) => Promise<void>,
 *   onUpdateDependency: (taskId: number, dependencyId: number, type: string) => Promise<void>,
 *   onUpdateDependent: (taskId: number, dependencyId: number, type: string) => Promise<void>,
 *   onDeleteDependency: (taskId: number, dependencyId: number) => Promise<void>,
 *   onDeleteDependent: (taskId: number, dependencyId: number) => Promise<void>,
 * }} props
 */
export function TaskDetailDrawer({
	open,
	task,
	project,
	workflow,
	tasks = [],
	taskTypes = [],
	memberOptions = [],
	members = [],
	onClose,
	onTaskClick,
	onCreateSubtask,
	onTitleBlur,
	onPersistAssignees,
	onSave,
	onDelete,
	onOpenLog,
	fetchTaskDependencies,
	fetchTaskDependents,
	onAddDependency,
	onAddDependent,
	onUpdateDependency,
	onUpdateDependent,
	onDeleteDependency,
	onDeleteDependent,
}) {
	const t = useTranslations();
	const locale = useLocale();
	const [form] = Form.useForm();
	const [dependencyFocusMode, setDependencyFocusMode] = useState(null);
	const parentTask = tasks.find((item) => item.id === task?.parent_id);
	const { subtasks, subtaskCompletedCount } = useMemo(() => {
		const list = task?.id ? getSubtasks(tasks, task.id) : [];
		const completed = list.filter((item) => {
			const state = workflow?.states?.find((row) => row.id === item.workflow_state_id);
			return isTerminalState(state);
		}).length;
		return { subtasks: list, subtaskCompletedCount: completed };
	}, [task, tasks, workflow?.states]);
	const watchedAssigneeIds = Form.useWatch('assignee_ids', form);
	const liveAssigneeIds = watchedAssigneeIds ?? getAssigneeIds(task);

	const priorityOptions = useMemo(
		() =>
			PROJECT_PRIORITY_FILTER_VALUES.map((value) => ({
				value: Number(value),
				label: (
					<span className="inline-flex items-center gap-1.5">
						<PriorityIcon level={Number(value)} />
						{t(PROJECT_PRIORITY_I18N_KEY[value])}
					</span>
				),
			})),
		[t],
	);

	useDeferredSyncEffect({
		predicate: () => Boolean(task && open),
		immediate: () => {
			form.setFieldsValue({
				title: task.title,
				description: task.description,
				priority: task.priority,
				type: task.type,
				workflow_state_id: task.workflow_state_id,
				assignee_ids: getAssigneeIds(task),
				start_date: task.start_date ? dayjs(task.start_date) : undefined,
				due_date: task.due_date ? dayjs(task.due_date) : undefined,
				estimated_hours: task.estimated_hours,
			});
		},
		deferred: () => {
			setDependencyFocusMode(null);
		},
		deps: [form, open, task],
		isUseMicrotask: true,
	});

	const handleSave = async () => {
		const values = form.getFieldsValue();
		await onSave(values);
	};

	return (
		<Drawer
			title={null}
			width={TASK_DETAIL_DRAWER_WIDTH}
			open={open}
			onClose={onClose}
			styles={{ body: { padding: 0, background: 'var(--base-text2)' } }}
		>
			{task ? (
				<div className="flex flex-col h-full bg-base-text2">
					<TaskDrawerHeader
						task={task}
						project={project}
						parentTask={parentTask}
						assigneeIds={liveAssigneeIds}
						members={members}
						onTaskClick={onTaskClick}
						onTitleBlur={onTitleBlur}
						onDelete={onDelete}
						onOpenLog={onOpenLog}
						onAddDependency={setDependencyFocusMode}
					/>

					<div className="flex-1 overflow-y-auto px-6 py-5 pb-24">
						<Form
							form={form}
							layout="vertical"
							onValuesChange={(changed) => {
								if ('assignee_ids' in changed) {
									onPersistAssignees(changed.assignee_ids || []);
								}
							}}
						>
							<TaskDetailProperties
								t={t}
								workflow={workflow}
								taskTypes={taskTypes}
								priorityOptions={priorityOptions}
								memberOptions={memberOptions}
								members={members}
							/>

							<TaskDetailSection title={t('kanban.description')}>
								<Form.Item
									name="description"
									noStyle
									getValueFromEvent={normalizeRichTextValue}
								>
									<TaskDescription
										key={`${task.id}-${open}`}
										locale={locale}
										placeholder={t('kanban.descriptionPlaceholder')}
									/>
								</Form.Item>
							</TaskDetailSection>

							<TaskDependencies
								key={`${task.id}-${dependencyFocusMode || 'idle'}`}
								task={task}
								tasks={tasks}
								workflow={workflow}
								members={members}
								onTaskClick={onTaskClick}
								fetchDependencies={fetchTaskDependencies}
								fetchDependents={fetchTaskDependents}
								onAddDependency={onAddDependency}
								onAddDependent={onAddDependent}
								onUpdateDependency={onUpdateDependency}
								onUpdateDependent={onUpdateDependent}
								onDeleteDependency={onDeleteDependency}
								onDeleteDependent={onDeleteDependent}
								initialAddingMode={dependencyFocusMode}
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
										onClick={() => onCreateSubtask(task)}
									/>
								}
							>
								<SubtaskList
									parentId={task.id}
									tasks={tasks}
									workflow={workflow}
									onTaskClick={onTaskClick}
								/>
							</TaskDetailPanel>

							<TaskAttachments
								key={`${project?.id ?? 'default'}-${task.id}`}
								taskId={task.id}
								projectId={project?.id}
								members={members}
							/>
						</Form>
					</div>

					<div className="sticky bottom-0 shrink-0 px-6 py-3 border-t border-base-border bg-white/90 backdrop-blur-sm">
						<div className="flex justify-end">
							<Button type="primary" onClick={handleSave}>
								{t('common.save')}
							</Button>
						</div>
					</div>
				</div>
			) : null}
		</Drawer>
	);
}
