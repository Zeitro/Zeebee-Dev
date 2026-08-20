'use client';

import { Form, Input, Select } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';

/**
 * 添加 / 编辑状态 Modal。
 */
export function StateFormModal({
	open,
	editingState,
	form,
	colorOptions,
	typeOptions,
	onCancel,
	onSubmit,
}) {
	const t = useTranslations();

	if (!open) return null;

	return (
		<NewModal
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
			<Form form={form} layout="vertical" className="mt-4" onFinish={onSubmit}>
				<Form.Item
					name="name"
					label={t('workflow.editor.stateNameLabel')}
					rules={[{ required: true }]}
				>
					<Input placeholder={t('workflow.editor.stateNamePlaceholder')} />
				</Form.Item>
				<Form.Item
					name="type"
					label={t('workflow.editor.stateTypeLabel')}
					rules={[{ required: true }]}
				>
					<Select options={typeOptions} />
				</Form.Item>
				<Form.Item
					name="color"
					label={t('workflow.editor.stateColorLabel')}
					rules={[{ required: true }]}
				>
					<Form.Item noStyle shouldUpdate={(prev, cur) => prev.color !== cur.color}>
						{() => (
							<div className="flex gap-3 flex-wrap">
								{colorOptions.map((c) => (
									<div
										key={c.value}
										role="button"
										tabIndex={0}
										onClick={() => form.setFieldValue('color', c.value)}
										onKeyDown={(e) => {
											if (e.key === 'Enter' || e.key === ' ') {
												form.setFieldValue('color', c.value);
											}
										}}
										className="rounded-full cursor-pointer transition-[border]"
										style={{
											width: 32,
											height: 32,
											background: c.hex,
											border:
												form.getFieldValue('color') === c.value
													? '3px solid #333'
													: '3px solid transparent',
										}}
										title={c.label}
									/>
								))}
							</div>
						)}
					</Form.Item>
				</Form.Item>
				<Form.Item name="description" label={t('workflow.editor.fieldDescription')}>
					<Input.TextArea rows={2} />
				</Form.Item>
			</Form>
		</NewModal>
	);
}
