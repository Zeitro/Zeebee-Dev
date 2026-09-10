/**
 * 任务类型表单弹窗：新建 / 编辑共用。
 * 提交 body 对齐接口 `{ name, description, color, sort_order }`；颜色预设色用英文名。
 */
'use client';

import { useCallback, useEffect, useMemo } from 'react';
import { ColorPicker, Form, Input, InputNumber } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';
import {
	TASK_TYPE_COLOR_PRESETS,
	TASK_TYPE_DEFAULT_COLOR,
	toTaskTypeColorHex,
	toTaskTypeColorPayload,
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
 * 	onSave: (values: {name: string, description: string, color: string, sort_order: number}) => Promise<void>
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
				sort_order: Number(editingType.sort_order ?? 0),
			});
			return;
		}
		form.resetFields();
		form.setFieldsValue({
			name: '',
			description: '',
			color: TASK_TYPE_DEFAULT_COLOR,
			sort_order: 0,
		});
	}, [editingType, form, open]);

	const handleOk = useCallback(async () => {
		const values = await form.validateFields();
		await onSave?.({
			name: String(values.name || '').trim(),
			description: values.description || '',
			color: toTaskTypeColorPayload(values.color),
			sort_order: Number(values.sort_order ?? 0),
		});
	}, [form, onSave]);

	if (!open) return null;

	return (
		<NewModal
			width={600}
			title={isEdit ? t('taskType.editTitle') : t('taskType.createTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleOk}
			onCancel={onCancel}
			destroyOnHidden
		>
			<Form form={form} layout="vertical" style={{ marginTop: 16 }} size="large">
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
				<Form.Item
					name="sort_order"
					label={t('taskType.sortOrderLabel')}
					rules={[{ required: true, message: t('taskType.sortOrderRequired') }]}
				>
					<InputNumber className="!w-full" min={0} precision={0} />
				</Form.Item>
			</Form>
		</NewModal>
	);
}
