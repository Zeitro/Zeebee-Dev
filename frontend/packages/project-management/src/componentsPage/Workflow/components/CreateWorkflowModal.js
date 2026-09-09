'use client';

import { useCallback, useEffect } from 'react';
import { Form, Input } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';

/**
 * 新建工作流弹窗：提交 `{ name, description }`。
 * 请求由父级走 POST `/project_management/workflow`，成功 data 为工作流对象（含 stages / transitions）。
 *
 * @param {{
 * 	open: boolean,
 * 	confirmLoading: boolean,
 * 	onCancel: () => void,
 * 	onSave: (values: {name: string, description: string}) => Promise<void>
 * }} props
 */
export function CreateWorkflowModal({ open, confirmLoading, onCancel, onSave }) {
	const t = useTranslations();
	const [form] = Form.useForm();

	useEffect(() => {
		if (!open) return;
		form.resetFields();
		form.setFieldsValue({ name: '', description: '' });
	}, [form, open]);

	const handleOk = useCallback(async () => {
		const values = await form.validateFields();
		const payload = {
			name: String(values.name || '').trim(),
			description: values.description || '',
		};
		await onSave?.(payload);
	}, [form, onSave]);

	if (!open) return null;

	return (
		<NewModal
			width={520}
			title={t('workflow.createTitle')}
			okText={t('workflow.createOk')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleOk}
			onCancel={onCancel}
			destroyOnHidden
		>
			<Form form={form} layout="vertical" size="large">
				<Form.Item
					name="name"
					label={t('workflow.nameLabel')}
					rules={[{ required: true, message: t('workflow.nameRequired') }]}
				>
					<Input placeholder={t('workflow.namePlaceholder')} />
				</Form.Item>

				<Form.Item name="description" label={t('workflow.descriptionLabel')}>
					<Input.TextArea rows={4} placeholder={t('workflow.descriptionPlaceholder')} />
				</Form.Item>
			</Form>
		</NewModal>
	);
}
