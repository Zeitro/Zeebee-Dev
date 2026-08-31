/**
 * 新建任务/子任务：表单校验 + POST `/project_management/task`；
 * 优先级走 options store，任务类型 / 执行者在弹窗内请求。
 */
'use client';

import dynamic from 'next/dynamic';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { DatePicker, Form, Input, InputNumber, Select } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest, dayjs, displayValue, formatDate } from '@repo/shared-methods/utils';
import {
	PROJECT_MEMBER_CANDIDATE_PURPOSE_PROJECT_MEMBER,
	createProjectMemberApi,
} from '@repo/project-management/api/projectMember';
import { createTaskApi } from '@repo/project-management/api/task';
import { createTaskTypeApi } from '@repo/project-management/api/taskType';
import { getTaskPriorityI18nKey, localizeSelectOptions } from '@repo/project-management/enum/task';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { useOptionsStore } from '@repo/project-management/store/options';
import { normalizeRichTextValue, unwrapList } from './helpers';

const WangEditor = dynamic(() => import('@repo/shared-methods/components/WangEditor'), {
	ssr: false,
});

/**
 * DatePicker / dayjs → 接口日期 MM/DD/YYYY；空值不传。
 * @param {*} value
 * @returns {string|undefined}
 */
function toApiDate(value) {
	if (!value) return undefined;
	const parsed = typeof value.format === 'function' ? value : dayjs(value);
	if (!parsed?.isValid?.()) return undefined;
	return formatDate(parsed) || undefined;
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
	const [optionsLoading, setOptionsLoading] = useState(false);
	const [typeOptions, setTypeOptions] = useState([]);
	const [assigneeOptions, setAssigneeOptions] = useState([]);
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);

	const priorityOptionsRaw = useOptionsStore((state) => state.priorityOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);

	const priorityOptions = useMemo(
		() => localizeSelectOptions(priorityOptionsRaw, t, getTaskPriorityI18nKey),
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

	useEffect(() => {
		if (!open) return undefined;

		let cancelled = false;
		(async () => {
			setOptionsLoading(true);
			try {
				const [, nextTypes, nextAssignees] = await Promise.all([
					getOptions('project_management_priority', 'priorityOptions').catch(() => {}),
					fetchTaskTypeOptions().catch(() => []),
					fetchAssigneeOptions().catch(() => []),
				]);
				if (cancelled) return;

				setTypeOptions(nextTypes);
				setAssigneeOptions(nextAssignees);

				const { priorityOptions } = useOptionsStore.getState();
				form.setFieldsValue({
					priority: priorityOptions[0]?.value,
					task_type_id: nextTypes[0]?.value,
				});
			} finally {
				if (!cancelled) setOptionsLoading(false);
			}
		})();

		return () => {
			cancelled = true;
		};
	}, [fetchAssigneeOptions, fetchTaskTypeOptions, form, getOptions, open]);

	const handleClose = useCallback(() => {
		form.resetFields();
		setTypeOptions([]);
		setAssigneeOptions([]);
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

			await taskApi.create({
				project_id: Number(projectId),
				parent_task_id: parentTask?.id,
				title: values.title,
				description: values.description || '',
				priority: values.priority,
				task_type_id: taskTypeId ? Number(taskTypeId) : undefined,
				start_at: toApiDate(values.start_date),
				end_at: toApiDate(values.due_date),
				work_hours: values.estimated_hours ? Number(values.estimated_hours) : undefined,
				assignee_organization_user_ids: assigneeIds,
			});

			messageApi?.success(t('kanban.createSuccess'));
			await onSuccess?.();
			form.resetFields();
			setTypeOptions([]);
			setAssigneeOptions([]);
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
				>
					<Form.Item name="title" label={t('kanban.title')} rules={[{ required: true }]}>
						<Input placeholder={t('kanban.titlePlaceholder')} />
					</Form.Item>
					<Form.Item
						name="description"
						label={t('kanban.description')}
						getValueFromEvent={normalizeRichTextValue}
					>
						<WangEditor
							locale={locale}
							wrapperProps={{ className: 'h-[400px]' }}
							editorProps={{
								defaultConfig: {
									placeholder: t('kanban.descriptionPlaceholder'),
								},
								style: { height: '100%', overflowY: 'hidden' },
							}}
						/>
					</Form.Item>
					<div className="grid grid-cols-2 gap-4">
						<Form.Item name="task_type_id" label={t('kanban.type')}>
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

					<div className="grid grid-cols-2 gap-4">
						<Form.Item name="start_date" label={t('kanban.startDate')}>
							<DatePicker className="w-full" />
						</Form.Item>
						<Form.Item name="due_date" label={t('kanban.dueDate')}>
							<DatePicker className="w-full" />
						</Form.Item>
					</div>
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
				</Form>
			</div>
		</NewModal>
	);
}
