/**
 * 创建项目弹窗：模板横滑选择与表单在同一 NewModal 内，模板默认排在项目名称前。
 * `initialWorkflowId` 有值且在列表中时预选该模板，否则仍选默认项。
 * POST `/project_management/project`，模板来自 GET `/project_management/workflow/all`；
 * 周期 / 开始 / 截止与项目设置同一套字段，成功后回调 handleSearch 刷新列表。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { Form, Input, Select } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest, dayjs } from '@repo/shared-methods/utils';
import { PROJECT_API } from '../../../api/project';
import { WORKFLOW_ALL_API } from '../../../api/workflow';
import { getTaskPriorityI18nKey, localizeSelectOptions, TASK_PRIORITY } from '../../../enum/task';
import { useOptionsStore } from '../../../store/options';
import { addProjectCycle, toCycleUpdateBody } from '../setting/cycle';
import { ProjectCycleFields } from '../setting/ProjectCycleFields';
import { WorkflowTemplatePicker } from './WorkflowTemplatePicker';

/**
 * 默认选 is_default，否则第一项。
 * @param {object[]} list
 * @returns {number|undefined}
 */
function pickDefaultWorkflowId(list) {
	return list.find((workflow) => workflow.is_default)?.id ?? list[0]?.id;
}

/**
 * 创建时的周期初值：今天起，默认 3 个月，截止日期 = 开始 + 周期。
 * @returns {{ cycle_months: number, cycle_unit: string, start_date: import('dayjs').Dayjs, end_date: import('dayjs').Dayjs }}
 */
function createCycleDefaults() {
	const start = dayjs();
	return {
		cycle_months: 3,
		cycle_unit: 'month',
		start_date: start,
		end_date: addProjectCycle(start, 3, 'month'),
	};
}

/**
 * 列表里与 initialWorkflowId 同 id 的项优先，否则默认模板。保留接口返回的 id 类型，供选择器 === 比较。
 * @param {object[]} list
 * @param {number|string|null|undefined} initialWorkflowId
 * @returns {number|string|undefined}
 */
function resolveWorkflowId(list, initialWorkflowId) {
	if (initialWorkflowId != null) {
		const matched = list.find((workflow) => String(workflow.id) === String(initialWorkflowId));
		if (matched) return matched.id;
	}
	return pickDefaultWorkflowId(list);
}

/**
 * @param {{
 *   messageApi: { success: (msg: string) => void },
 *   setOpen: (open: boolean) => void,
 *   handleSearch?: () => void,
 *   initialWorkflowId?: number|string|null
 * }} props
 */
export default function CreateProjectModal({
	messageApi,
	setOpen,
	handleSearch,
	initialWorkflowId,
}) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [workflows, setWorkflows] = useState([]);
	const [workflowLoading, setWorkflowLoading] = useState(true);
	const [priorityLoading, setPriorityLoading] = useState(true);
	const [cycleDefaults] = useState(createCycleDefaults);

	const priorityOptionsRaw = useOptionsStore((state) => state.priorityOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);

	const priorityOptions = useMemo(
		() => localizeSelectOptions(priorityOptionsRaw, t, getTaskPriorityI18nKey),
		[priorityOptionsRaw, t],
	);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const [data] = await Promise.all([
					request(WORKFLOW_ALL_API, {
						method: 'GET',
						isHideTips: true,
					}),
					getOptions('project_management_priority', 'priorityOptions').catch(() => {}),
				]);
				if (cancelled) return;
				const list = Array.isArray(data?.list) ? data.list : [];
				setWorkflows(list);

				const { priorityOptions: loadedPriority } = useOptionsStore.getState();
				const defaultPriority =
					loadedPriority.find((item) => item.value === TASK_PRIORITY.MEDIUM)?.value ??
					loadedPriority[0]?.value;
				const defaultId = resolveWorkflowId(list, initialWorkflowId);
				const fields = {};
				if (defaultId != null) fields.workflow_template_id = defaultId;
				if (defaultPriority != null) fields.priority = defaultPriority;
				if (Object.keys(fields).length) form.setFieldsValue(fields);
			} catch {
				if (!cancelled) setWorkflows([]);
			} finally {
				if (!cancelled) {
					setWorkflowLoading(false);
					setPriorityLoading(false);
				}
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [form, getOptions, initialWorkflowId, request]);

	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setConfirmLoading(true);
			const created = await request(PROJECT_API, {
				method: 'POST',
				body: {
					name: String(values.name || '').trim(),
					description: values.description || '',
					workflow_template_id: values.workflow_template_id,
					priority: values.priority,
					...toCycleUpdateBody(values),
				},
			});
			if (created == null) return;
			messageApi.success(t('projectList.createSuccess'));
			handleSearch?.();
			setOpen(false);
		} catch {
			// 表单校验失败，或 createRequest 已提示接口错误
		} finally {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width={750}
			title={t('projectList.createTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
			destroyOnHidden
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				initialValues={cycleDefaults}
				className="h-[65vh] overflow-y-auto max-h-[800px]"
			>
				<Form.Item
					name="workflow_template_id"
					label={t('projectList.template')}
					rules={[{ required: true, message: t('projectList.templateRequired') }]}
				>
					<WorkflowTemplatePicker workflows={workflows} loading={workflowLoading} />
				</Form.Item>
				<Form.Item name="name" label={t('projectList.name')} rules={[{ required: true }]}>
					<Input placeholder={t('projectList.namePlaceholder')} />
				</Form.Item>
				<Form.Item name="description" label={t('projectList.description')}>
					<Input.TextArea
						rows={3}
						placeholder={t('projectList.descriptionPlaceholder')}
					/>
				</Form.Item>
				<Form.Item
					name="priority"
					label={t('projectList.priority')}
					rules={[{ required: true }]}
				>
					<Select loading={priorityLoading} options={priorityOptions} />
				</Form.Item>
				<ProjectCycleFields required stacked />
			</Form>
		</NewModal>
	);
}
