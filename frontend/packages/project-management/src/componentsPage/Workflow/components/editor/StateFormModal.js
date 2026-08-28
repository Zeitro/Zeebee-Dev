/**
 * 添加 / 编辑状态 Modal。颜色用 ColorPicker，表单存 hex。
 * 状态类型对应接口 `stages.phase`，下拉来自 `stageOptions`（`project_management_stage_phase`）。
 * 提交后由编辑器把表单写入详情 `stages`，再 PUT 整份 workflow。
 */
'use client';

import { useEffect, useMemo } from 'react';
import { ColorPicker, Form, Input, Select } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';
import { localizeSelectOptions } from '@repo/project-management/enum/task';
import {
	WORKFLOW_STATE_COLOR_PRESETS,
	WORKFLOW_STATE_DEFAULT_COLOR,
	getWorkflowPhaseI18nKey,
	toWorkflowStateColorHex,
} from '@repo/project-management/enum/workflow';
import { useOptionsStore } from '@repo/project-management/store/options';

/**
 * ColorPicker 受控为 hex，对接 Form.Item 的 value / onChange。
 * @param {{ value?: string, onChange?: (hex: string) => void, presets?: Array }} props
 */
function HexColorPicker({ value, onChange, presets, ...rest }) {
	return (
		<ColorPicker
			{...rest}
			value={value || WORKFLOW_STATE_DEFAULT_COLOR}
			onChange={(color) => onChange?.(toWorkflowStateColorHex(color))}
			showText
			disabledAlpha
			format="hex"
			size="large"
			presets={presets}
			getPopupContainer={(trigger) => trigger?.parentElement || document.body}
		/>
	);
}

export function StateFormModal({ open, editingState, form, onCancel, onSubmit }) {
	const t = useTranslations();
	const stageOptions = useOptionsStore((state) => state.stageOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);
	const presets = useMemo(
		() => [{ label: t('workflow.editor.presetColors'), colors: WORKFLOW_STATE_COLOR_PRESETS }],
		[t],
	);
	const phaseOptions = useMemo(
		() => localizeSelectOptions(stageOptions, t, getWorkflowPhaseI18nKey),
		[stageOptions, t],
	);

	useEffect(() => {
		if (!open) return;
		if (useOptionsStore.getState().stageOptions.length) return;
		getOptions('project_management_stage_phase', 'stageOptions').catch(() => {});
	}, [open, getOptions]);

	return (
		<NewModal
			open={open}
			width={520}
			title={
				editingState
					? t('workflow.editor.editStateTitle')
					: t('workflow.editor.addStateTitle')
			}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			onCancel={onCancel}
			onOk={() => form.submit()}
			destroyOnHidden
		>
			<Form form={form} layout="vertical" className="mt-4" onFinish={onSubmit} size="large">
				<Form.Item
					name="name"
					label={t('workflow.editor.stateNameLabel')}
					rules={[{ required: true }]}
				>
					<Input placeholder={t('workflow.editor.stateNamePlaceholder')} />
				</Form.Item>
				<Form.Item
					name="phase"
					label={t('workflow.editor.stateTypeLabel')}
					rules={[{ required: true }]}
				>
					<Select
						options={phaseOptions}
						placeholder={t('workflow.editor.selectPlaceholder')}
					/>
				</Form.Item>
				<Form.Item
					name="color"
					label={t('workflow.editor.stateColorLabel')}
					rules={[{ required: true }]}
				>
					<HexColorPicker presets={presets} />
				</Form.Item>
				<Form.Item name="description" label={t('workflow.editor.fieldDescription')}>
					<Input.TextArea rows={4} />
				</Form.Item>
			</Form>
		</NewModal>
	);
}
