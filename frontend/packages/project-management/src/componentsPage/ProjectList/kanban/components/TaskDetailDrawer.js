/**
 * 任务详情抽屉：标题失焦保存；状态走 transition；执行者即时落库；依赖/附件/日志。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { CalendarOutlined, PlusOutlined, UserOutlined } from '@ant-design/icons';
import { Button, DatePicker, Divider, Drawer, Form, Input, InputNumber, Select } from 'antd';
import { useTranslations } from 'next-intl';
import { dayjs } from '@repo/shared-methods/utils';
import {
	PROJECT_PRIORITY_FILTER_VALUES,
	PROJECT_PRIORITY_I18N_KEY,
} from '../../../../enum/project';
import { getWorkflowStateColor } from '../../../../enum/workflow';
import { TASK_DETAIL_DRAWER_WIDTH } from './constants';
import { AssigneeChips } from './AssigneeChips';
import { getAssigneeIds } from './helpers';
import { PriorityIcon } from './PriorityIcon';
import { SubtaskList } from './SubtaskList';
import { TaskAttachments } from './TaskAttachments';
import { TaskDependencies } from './TaskDependencies';
import { TaskDrawerHeader } from './TaskDrawerHeader';

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
	const [form] = Form.useForm();
	const [dependencyFocusMode, setDependencyFocusMode] = useState(null);
	const parentTask = tasks.find((item) => item.id === task?.parent_id);
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

	useEffect(() => {
		if (!task || !open) return;
		setDependencyFocusMode(null);
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
	}, [form, open, task]);

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
			styles={{ body: { padding: 0 } }}
		>
			{task ? (
				<div className="flex flex-col h-full">
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
					<div className="flex-1 overflow-y-auto px-6 py-4">
						<Form
							form={form}
							layout="vertical"
							onValuesChange={(changed) => {
								if ('assignee_ids' in changed) {
									onPersistAssignees(changed.assignee_ids || []);
								}
							}}
						>
							<div className="flex items-center gap-0.5 mb-3 flex-wrap">
								<Form.Item name="workflow_state_id" noStyle>
									<Select
										variant="borderless"
										size="small"
										className="min-w-[100px]"
									>
										{(workflow?.states || []).map((state) => (
											<Select.Option key={state.id} value={state.id}>
												<span
													className="inline-block w-2 h-2 rounded-full mr-1.5 align-middle"
													style={{
														background: getWorkflowStateColor(
															state.color,
														),
													}}
												/>
												{state.name}
											</Select.Option>
										))}
									</Select>
								</Form.Item>
								<Divider type="vertical" className="h-4" />
								<Form.Item name="type" noStyle>
									<Select
										variant="borderless"
										size="small"
										className="min-w-20"
										allowClear
									>
										{taskTypes.map((item) => (
											<Select.Option key={item.id} value={item.name}>
												{item.name}
											</Select.Option>
										))}
									</Select>
								</Form.Item>
								<Divider type="vertical" className="h-4" />
								<Form.Item name="priority" noStyle>
									<Select
										variant="borderless"
										size="small"
										className="min-w-[90px]"
										options={priorityOptions}
									/>
								</Form.Item>
								<Divider type="vertical" className="h-4" />
								<Form.Item name="start_date" noStyle>
									<DatePicker
										variant="borderless"
										size="small"
										placeholder={t('kanban.startDate')}
										prefix={<CalendarOutlined className="text-base-disabled" />}
										suffixIcon={null}
										className="w-[150px]"
									/>
								</Form.Item>
								<Form.Item name="due_date" noStyle>
									<DatePicker
										variant="borderless"
										size="small"
										placeholder={t('kanban.dueDate')}
										prefix={<CalendarOutlined className="text-base-disabled" />}
										suffixIcon={null}
										className="w-[150px]"
									/>
								</Form.Item>
								<Divider type="vertical" className="h-4" />
								<Form.Item name="estimated_hours" noStyle>
									<InputNumber
										variant="borderless"
										size="small"
										placeholder={t('kanban.estimatedHours')}
										min={0}
										step={0.5}
										suffix={t('kanban.hoursSuffix')}
										className="w-[90px]"
									/>
								</Form.Item>
							</div>
							<div className="flex items-center gap-2 mb-3">
								<span className="text-[13px] text-base-gray shrink-0">
									<UserOutlined className="mr-1.5 text-base-label" />
									{t('kanban.assignee')}
								</span>
								<Form.Item name="assignee_ids" noStyle>
									<AssigneeChips options={memberOptions} fallbacks={members} />
								</Form.Item>
							</div>
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
							<Button
								className="mb-4"
								icon={<PlusOutlined />}
								onClick={() => onCreateSubtask(task)}
							>
								{t('kanban.addSubtask')}
							</Button>
							<SubtaskList
								parentId={task.id}
								tasks={tasks}
								workflow={workflow}
								onTaskClick={onTaskClick}
							/>
							<Form.Item name="description" noStyle>
								<Input.TextArea
									rows={6}
									placeholder={t('kanban.descriptionPlaceholder')}
								/>
							</Form.Item>
							<TaskAttachments
								key={task.id}
								taskId={task.id}
								projectId={project?.id}
								members={members}
							/>
							<div className="mt-4">
								<Button type="primary" onClick={handleSave}>
									{t('common.save')}
								</Button>
							</div>
						</Form>
					</div>
				</div>
			) : null}
		</Drawer>
	);
}
