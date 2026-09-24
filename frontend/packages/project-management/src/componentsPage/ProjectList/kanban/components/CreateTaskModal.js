/**
 * 新建任务/子任务：表单校验 + POST `/project_management/task`；
 * 优先级走 options store，任务类型 / 执行者 / 项目工作流在弹窗内请求。
 * 周期 / 开始 / 截止与项目设置同一套字段，仅用于推导 start_at / end_at，不向任务接口传 cycle。
 * 工作流默认选项目默认项；提交带 `workflow_id`。
 */
'use client';

import dynamic from 'next/dynamic';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { Form, Input, InputNumber, Select } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest, dayjs, displayValue } from '@repo/shared-methods/utils';
import { getProjectDetailApi } from '@repo/project-management/api/project';
import {
	PROJECT_MEMBER_CANDIDATE_PURPOSE_PROJECT_MEMBER,
	createProjectMemberApi,
} from '@repo/project-management/api/projectMember';
import { createTaskApi } from '@repo/project-management/api/task';
import { createTaskTypeApi } from '@repo/project-management/api/taskType';
import {
	getTaskPriorityI18nKey,
	localizeSelectOptions,
	toProjectPriorityLevel,
} from '@repo/project-management/enum/task';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { useOptionsStore } from '@repo/project-management/store/options';
import { addProjectCycle } from '../../setting/cycle';
import { ProjectCycleFields } from '../../setting/ProjectCycleFields';
import { normalizeRichTextValue, toApiDate, unwrapList } from './helpers';
import { PriorityIcon } from './TaskDetail/PriorityIcon';

const WangEditor = dynamic(() => import('@repo/shared-methods/components/WangEditor'), {
	ssr: false,
});

/**
 * 创建任务时的周期初值：今天起，默认 1 周，截止日期 = 开始 + 周期。
 * @returns {{ cycle_months: number, cycle_unit: string, start_date: import('dayjs').Dayjs, end_date: import('dayjs').Dayjs }}
 */
function createTaskCycleDefaults() {
	const start = dayjs();
	return {
		cycle_months: 1,
		cycle_unit: 'week',
		start_date: start,
		end_date: addProjectCycle(start, 1, 'week'),
	};
}

/**
 * 详情 `workflows` 优先，否则用单个 `workflow`。
 * @param {object|null|undefined} project
 * @returns {object[]}
 */
function listProjectWorkflows(project) {
	if (Array.isArray(project?.workflows) && project.workflows.length) return project.workflows;
	return project?.workflow ? [project.workflow] : [];
}

/**
 * 默认工作流：default_workflow_id → is_default → 第一项。
 * @param {object|null|undefined} project
 * @param {object[]} workflows
 * @returns {number|string|undefined}
 */
function pickDefaultWorkflowId(project, workflows) {
	if (!workflows.length) return undefined;
	const preferred = project?.default_workflow_id;
	if (preferred != null && workflows.some((row) => String(row.id) === String(preferred))) {
		return preferred;
	}
	return workflows.find((row) => row.is_default)?.id ?? workflows[0]?.id;
}

/**
 * @param {{
 *   open: boolean,
 *   projectId?: number|string,
 *   parentTask?: object|null,
 *   onCancel: () => void,
 *   onSuccess?: () => void | Promise<void>,
 *   messageApi?: import('antd/es/message/interface').MessageInstance,
 * }} props
 */
export default function CreateTaskModal({
	open,
	projectId,
	parentTask,
	onCancel,
	onSuccess,
	messageApi,
}) {
	const t = useTranslations();
	const locale = useLocale();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [optionsLoading, setOptionsLoading] = useState(true);
	const [typeOptions, setTypeOptions] = useState([]);
	const [assigneeOptions, setAssigneeOptions] = useState([]);
	const [workflowOptions, setWorkflowOptions] = useState([]);
	const [cycleDefaults] = useState(createTaskCycleDefaults);
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);

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

	const defaultTaskTypeId = typeOptions[0]?.value;

	const fetchTaskTypeOptions = useCallback(async () => {
		const data = await taskTypeApi.list();
		return unwrapList(data)
			.sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
			.map((item) => ({
				value: item.id,
				label: displayValue(item.name),
			}));
	}, [taskTypeApi]);

	const fetchAssigneeOptions = useCallback(async () => {
		if (!projectId) return [];
		const res = await projectMemberApi.candidate({
			project_id: Number(projectId),
			keyword: '',
			purpose: PROJECT_MEMBER_CANDIDATE_PURPOSE_PROJECT_MEMBER,
		});
		return unwrapList(res)
			.map(mapProjectMember)
			.filter(Boolean)
			.map((item) => ({
				value: item.organization_user_id ?? item.user_id,
				label: displayValue(item.full_name || item.name || item.email),
			}));
	}, [projectId, projectMemberApi]);

	/** @returns {{ options: Array<{ value: number|string, label: string }>, defaultId: number|string|undefined }} */
	const fetchWorkflowOptions = useCallback(async () => {
		if (!projectId) return { options: [], defaultId: undefined };
		const project = await request(getProjectDetailApi(projectId), {
			method: 'GET',
			isHideTips: true,
		});
		const workflows = listProjectWorkflows(project);
		const options = workflows.map((row) => ({
			value: row.id,
			label: row.is_default
				? `${displayValue(row.name)} (${t('projectSetting.workflowDefault')})`
				: displayValue(row.name),
		}));
		return { options, defaultId: pickDefaultWorkflowId(project, workflows) };
	}, [projectId, request, t]);

	useEffect(() => {
		if (!open) return undefined;

		let cancelled = false;
		(async () => {
			setOptionsLoading(true);
			try {
				const [, nextTypes, nextAssignees, nextWorkflows] = await Promise.all([
					getOptions('project_management_priority', 'priorityOptions').catch(() => {}),
					fetchTaskTypeOptions().catch(() => []),
					fetchAssigneeOptions().catch(() => []),
					fetchWorkflowOptions().catch(() => ({ options: [], defaultId: undefined })),
				]);
				if (cancelled) return;

				setTypeOptions(nextTypes);
				setAssigneeOptions(nextAssignees);
				setWorkflowOptions(nextWorkflows.options);

				const { priorityOptions } = useOptionsStore.getState();
				form.setFieldsValue({
					...cycleDefaults,
					priority: priorityOptions[0]?.value,
					task_type_id: nextTypes[0]?.value,
					workflow_id: nextWorkflows.defaultId,
				});
			} finally {
				if (!cancelled) setOptionsLoading(false);
			}
		})();

		return () => {
			cancelled = true;
		};
	}, [
		cycleDefaults,
		fetchAssigneeOptions,
		fetchTaskTypeOptions,
		fetchWorkflowOptions,
		form,
		getOptions,
		open,
	]);

	const handleClose = useCallback(() => {
		form.resetFields();
		setTypeOptions([]);
		setAssigneeOptions([]);
		setWorkflowOptions([]);
		onCancel();
	}, [form, onCancel]);

	const handleSave = async () => {
		if (!projectId) return;

		try {
			const values = await form.validateFields();
			setConfirmLoading(true);

			const assigneeIds = Array.isArray(values.assignee_ids)
				? values.assignee_ids.map((id) => Number(id)).filter(Boolean)
				: [];
			const taskTypeId = values.task_type_id ?? defaultTaskTypeId;
			const workflowId = values.workflow_id;

			await taskApi.create({
				project_id: Number(projectId),
				parent_task_id: parentTask?.id,
				title: values.title,
				description: values.description || '',
				priority: values.priority,
				task_type_id: taskTypeId ? Number(taskTypeId) : undefined,
				workflow_id: workflowId ? Number(workflowId) : undefined,
				start_at: toApiDate(values.start_date),
				end_at: toApiDate(values.end_date),
				work_hours: values.estimated_hours ? Number(values.estimated_hours) : undefined,
				assignee_organization_user_ids: assigneeIds,
			});

			messageApi?.success(t('kanban.createSuccess'));
			await onSuccess?.();
			form.resetFields();
			setTypeOptions([]);
			setAssigneeOptions([]);
			setWorkflowOptions([]);
			onCancel();
		} catch (error) {
			if (error?.errorFields) return;
			messageApi?.error(error?.message || t('kanban.createFailed'));
		} finally {
			setConfirmLoading(false);
		}
	};

	if (!open) return null;

	return (
		<NewModal
			width="70%"
			title={parentTask ? t('kanban.newSubtask') : t('kanban.newTask')}
			okText={t('common.create')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			destroyOnHidden
			onOk={handleSave}
			onCancel={handleClose}
		>
			<div className="max-h-[70vh] overflow-y-auto">
				{parentTask ? (
					<p className="mb-3 px-3 py-2 rounded-md bg-base-text2 text-[13px] text-base-gray">
						{t('kanban.parentTask', { title: displayValue(parentTask.title) })}
					</p>
				) : null}
				<Form
					form={form}
					layout="vertical"
					autoComplete="off"
					size="large"
					disabled={confirmLoading || optionsLoading}
					initialValues={cycleDefaults}
				>
					<Form.Item name="title" label={t('kanban.title')} rules={[{ required: true }]}>
						<Input placeholder={t('kanban.titlePlaceholder')} />
					</Form.Item>
					<Form.Item
						name="description"
						label={t('kanban.description')}
						getValueFromEvent={normalizeRichTextValue}
						className="h-[422px]"
					>
						{optionsLoading ? (
							<div className="h-[400px] bg-base-text2 rounded-md border border-base-text2"></div>
						) : (
							<WangEditor
								locale={locale}
								wrapperProps={{ className: 'h-[400px]' }}
								editorProps={{
									defaultConfig: {
										placeholder: t('kanban.descriptionPlaceholder'),
									},
									style: { overflowY: 'hidden' },
								}}
							/>
						)}
					</Form.Item>
					<div className="grid grid-cols-2 gap-4">
						<Form.Item
							name="task_type_id"
							label={t('kanban.type')}
							rules={[{ required: true }]}
						>
							<Select
								options={typeOptions}
								allowClear
								loading={optionsLoading}
								placeholder={t('kanban.type')}
							/>
						</Form.Item>
						<Form.Item
							name="priority"
							label={t('projectList.priority')}
							rules={[{ required: true }]}
						>
							<Select options={priorityOptions} loading={optionsLoading} />
						</Form.Item>
					</div>

					<Form.Item
						name="workflow_id"
						label={t('kanban.workflow')}
						rules={[{ required: true, message: t('kanban.workflowRequired') }]}
					>
						<Select
							options={workflowOptions}
							loading={optionsLoading}
							showSearch
							optionFilterProp="label"
							placeholder={t('kanban.workflowPlaceholder')}
							notFoundContent={optionsLoading ? undefined : t('kanban.emptyWorkflow')}
						/>
					</Form.Item>

					<div className="grid grid-cols-2 gap-4">
						<Form.Item name="estimated_hours" label={t('kanban.estimatedHours')}>
							<InputNumber className="!w-full" min={0} step={0.5} />
						</Form.Item>
						<Form.Item name="assignee_ids" label={t('kanban.assignee')}>
							<Select
								mode="multiple"
								allowClear
								showSearch
								optionFilterProp="label"
								loading={optionsLoading}
								placeholder={t('kanban.assigneePlaceholder')}
								options={assigneeOptions}
							/>
						</Form.Item>
					</div>
					<ProjectCycleFields />
				</Form>
			</div>
		</NewModal>
	);
}
