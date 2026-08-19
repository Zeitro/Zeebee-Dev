/**
 * 任务类型表单弹窗：新建 / 编辑共用；颜色用 ColorPicker，提交 hex。
 */
'use client';

import { useCallback, useEffect, useMemo } from 'react';
import { ColorPicker, Form, Input } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';
import {
	TASK_TYPE_COLOR_PRESETS,
	TASK_TYPE_DEFAULT_COLOR,
	toTaskTypeColorHex,
} from '@repo/project-management/enum/taskType';

/**
 * ColorPicker 受控为 hex，对接 Form.Item 的 value / onChange。
 * @param {{ value?: string, onChange?: (hex: string) => void, presets?: Array }} props
 */
function HexColorPicker({ value, onChange, presets, ...rest }) {
	return (
		<ColorPicker
			{...rest}
			value={value || TASK_TYPE_DEFAULT_COLOR}
			onChange={(color) => onChange?.(toTaskTypeColorHex(color))}
			showText
			disabledAlpha
			format="hex"
			presets={presets}
		/>
	);
}

/**
 * @param {{
 * 	open: boolean,
 * 	confirmLoading: boolean,
 * 	editingType: object|null,
 * 	onCancel: () => void,
 * 	onSave: (values: {name: string, description: string, color: string}) => Promise<void>
 * }} props
 */
export function TaskTypeFormModal({ open, confirmLoading, editingType, onCancel, onSave }) {
	const t = useTranslations();
	const [form] = Form.useForm();
	const isEdit = Boolean(editingType);

	const presets = useMemo(
		() => [{ label: t('taskType.presetColors'), colors: TASK_TYPE_COLOR_PRESETS }],
		[t],
	);

	useEffect(() => {
		if (!open) return;
		if (editingType) {
			form.setFieldsValue({
				name: editingType.name,
				description: editingType.description || '',
				color: toTaskTypeColorHex(editingType.color) || TASK_TYPE_DEFAULT_COLOR,
			});
			return;
		}
		form.resetFields();
		form.setFieldsValue({
			name: '',
			description: '',
			color: TASK_TYPE_DEFAULT_COLOR,
		});
	}, [editingType, form, open]);

	const handleOk = useCallback(async () => {
		const values = await form.validateFields();
		await onSave?.({
			name: String(values.name || '').trim(),
			description: values.description || '',
			color: toTaskTypeColorHex(values.color) || TASK_TYPE_DEFAULT_COLOR,
		});
	}, [form, onSave]);

	if (!open) return null;

	return (
		<NewModal
			width={520}
			title={isEdit ? t('taskType.editTitle') : t('taskType.createTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleOk}
			onCancel={onCancel}
			destroyOnHidden
		>
			<Form form={form} layout="vertical" style={{ marginTop: 16 }}>
				<Form.Item
					name="name"
					label={t('taskType.nameLabel')}
					rules={[{ required: true, message: t('taskType.nameRequired') }]}
				>
					<Input placeholder={t('taskType.namePlaceholder')} />
				</Form.Item>
				<Form.Item name="description" label={t('taskType.descriptionLabel')}>
					<Input.TextArea rows={3} placeholder={t('taskType.descriptionPlaceholder')} />
				</Form.Item>
				<Form.Item
					name="color"
					label={t('taskType.colorLabel')}
					rules={[{ required: true, message: t('taskType.colorRequired') }]}
				>
					<HexColorPicker presets={presets} />
				</Form.Item>
			</Form>
		</NewModal>
	);
}
