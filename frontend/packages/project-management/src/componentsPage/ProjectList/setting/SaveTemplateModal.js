/**
 * 将当前项目另存为模板。POST `/project_management/project/:id/template`，不改项目本身。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { SaveOutlined } from '@ant-design/icons';
import { Button, Checkbox, Form, Input } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import { getProjectTemplateApi } from '../../../api/project';

/**
 * @param {{
 *   project: object,
 *   projectId: number|string,
 *   messageApi: { success: (msg: string) => void }
 * }} props
 */
export function SaveAsTemplateButton({ project, projectId, messageApi }) {
	const t = useTranslations();
	const [open, setOpen] = useState(false);

	return (
		<>
			<Button icon={<SaveOutlined />} htmlType="button" onClick={() => setOpen(true)}>
				{t('projectSetting.saveAsTemplate')}
			</Button>
			<SaveTemplateModal
				open={open}
				project={project}
				projectId={projectId}
				messageApi={messageApi}
				onClose={() => setOpen(false)}
			/>
		</>
	);
}

/**
 * @param {{
 *   open: boolean,
 *   project: object,
 *   projectId: number|string,
 *   messageApi: { success: (msg: string) => void },
 *   onClose: () => void
 * }} props
 */
function SaveTemplateModal({ open, project, projectId, messageApi, onClose }) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const [form] = Form.useForm();
	const [saving, setSaving] = useState(false);

	useEffect(() => {
		if (!open) return;
		form.setFieldsValue({
			name: t('projectSetting.templateNameDefault', { name: project?.name || '' }),
			description: project?.description || '',
			include_workflow: true,
			include_tasks: true,
			include_members: true,
		});
	}, [form, open, project?.description, project?.name, t]);

	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setSaving(true);
			const saved = await request(getProjectTemplateApi(projectId), {
				method: 'POST',
				body: {
					name: String(values.name || '').trim(),
					description: values.description || '',
					include_workflow: true,
					include_tasks: !!values.include_tasks,
					include_members: !!values.include_members,
				},
			});
			if (saved == null) return;
			messageApi.success(t('projectSetting.saveTemplateSuccess'));
			onClose();
		} catch {
			// 表单校验失败，或 createRequest 已提示接口错误
		} finally {
			setSaving(false);
		}
	};

	if (!open) return null;

	return (
		<NewModal
			width={520}
			title={t('projectSetting.saveAsTemplateTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			confirmLoading={saving}
			onOk={handleSave}
			onCancel={onClose}
			destroyOnHidden
		>
			<Form form={form} layout="vertical" className="!mt-2">
				<Form.Item
					name="name"
					label={t('projectSetting.templateName')}
					rules={[{ required: true, message: t('projectSetting.templateNameRequired') }]}
				>
					<Input placeholder={t('projectSetting.templateNamePlaceholder')} />
				</Form.Item>
				<Form.Item name="description" label={t('projectSetting.templateDescription')}>
					<Input.TextArea
						rows={2}
						placeholder={t('projectSetting.templateDescriptionPlaceholder')}
					/>
				</Form.Item>
				<div className="mb-2 font-medium">{t('projectSetting.templateInclude')}</div>
				<Form.Item name="include_workflow" valuePropName="checked" className="!mb-2">
					<Checkbox disabled>{t('projectSetting.templateIncludeWorkflow')}</Checkbox>
				</Form.Item>
				<Form.Item name="include_tasks" valuePropName="checked" className="!mb-2">
					<Checkbox>{t('projectSetting.templateIncludeTasks')}</Checkbox>
				</Form.Item>
				<Form.Item name="include_members" valuePropName="checked" className="!mb-0">
					<Checkbox>{t('projectSetting.templateIncludeMembers')}</Checkbox>
				</Form.Item>
			</Form>
		</NewModal>
	);
}
