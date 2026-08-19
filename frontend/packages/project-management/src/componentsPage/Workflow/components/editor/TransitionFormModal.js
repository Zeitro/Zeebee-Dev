'use client';

import { Modal, Form, Input, Select } from 'antd';
import { useTranslations } from 'next-intl';

/**
 * 添加流转 Modal。
 */
export function TransitionFormModal({
	open,
	fromState,
	states,
	form,
	colorMap,
	triggerOptions,
	onCancel,
	onSubmit,
}) {
	const t = useTranslations();

	return (
		<Modal
			title={t('workflow.editor.addTransitionTitle', { from: fromState?.name || '' })}
			open={open}
			onCancel={onCancel}
			onOk={() => form.submit()}
		>
			<Form form={form} layout="vertical" onFinish={onSubmit}>
				<Form.Item
					name="to_state_id"
					label={t('workflow.editor.targetStateLabel')}
					rules={[{ required: true }]}
				>
					<Select>
						{states
							.filter((s) => s.id !== fromState?.id)
							.map((s) => (
								<Select.Option key={s.id} value={s.id}>
									<span
										className="inline-block w-2.5 h-2.5 rounded-full mr-2 align-middle"
										style={{ background: colorMap[s.color] }}
									/>
									{s.name}
								</Select.Option>
							))}
					</Select>
				</Form.Item>
				<Form.Item name="name" label={t('workflow.editor.transitionNameLabel')}>
					<Input placeholder={t('workflow.editor.transitionNamePlaceholder')} />
				</Form.Item>
				<Form.Item name="trigger_type" label={t('workflow.editor.triggerTypeLabel')}>
					<Select options={triggerOptions} />
				</Form.Item>
			</Form>
		</Modal>
	);
}
