/**
 * 新建任务/子任务：NewModal + Form；提交后由父级请求接口并刷新看板。
 */
'use client';

import dynamic from 'next/dynamic';
import { useMemo, useState } from 'react';
import { DatePicker, Form, Input, InputNumber, Select } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { displayValue } from '@repo/shared-methods/utils';
import {
	PROJECT_PRIORITY,
	PROJECT_PRIORITY_FILTER_VALUES,
	PROJECT_PRIORITY_I18N_KEY,
} from '../../../../enum/project';
import { normalizeRichTextValue } from './helpers';

const WangEditor = dynamic(() => import('@repo/shared-methods/components/WangEditor'), {
	ssr: false,
});

/**
 * @param {{
 *   parentTask?: object|null,
 *   taskTypes?: Array,
 *   memberOptions?: Array,
 *   onCancel: () => void,
 *   onSubmit: (values: object) => Promise<void>,
 * }} props
 */
export default function CreateTaskModal({
	parentTask,
	taskTypes = [],
	memberOptions = [],
	onCancel,
	onSubmit,
}) {
	const t = useTranslations();
	const locale = useLocale();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const defaultType = taskTypes[0]?.name;

	const priorityOptions = useMemo(
		() =>
			PROJECT_PRIORITY_FILTER_VALUES.map((value) => ({
				value: Number(value),
				label: t(PROJECT_PRIORITY_I18N_KEY[value]),
			})),
		[t],
	);

	const typeOptions = useMemo(
		() => taskTypes.map((item) => ({ value: item.name, label: item.name })),
		[taskTypes],
	);

	const assigneeOptions = useMemo(
		() =>
			memberOptions.map((member) => ({
				value: member.id,
				label: displayValue(member.name || member.email),
			})),
		[memberOptions],
	);

	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setConfirmLoading(true);
			await onSubmit(values);
			setConfirmLoading(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width="70%"
			title={parentTask ? t('kanban.newSubtask') : t('kanban.newTask')}
			okText={t('common.create')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={onCancel}
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
					disabled={confirmLoading}
					initialValues={{
						priority: PROJECT_PRIORITY.MEDIUM,
						type: defaultType,
					}}
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
						<Form.Item name="type" label={t('kanban.type')}>
							<Select options={typeOptions} allowClear />
						</Form.Item>
						<Form.Item
							name="priority"
							label={t('projectList.priority')}
							rules={[{ required: true }]}
						>
							<Select options={priorityOptions} />
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
						<Form.Item name="assignee_id" label={t('kanban.assignee')}>
							<Select
								allowClear
								showSearch
								optionFilterProp="label"
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
