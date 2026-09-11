'use client';

import { useEffect } from 'react';
import { Form, Input, Select } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';
import { getWorkflowStateColorHex } from '@repo/project-management/enum/workflow';

/**
 * 添加流转 Modal。
 */
export function TransitionFormModal({
	open,
	fromState,
	states,
	form,
	triggerOptions,
	onCancel,
	onSubmit,
}) {
	const t = useTranslations();

	useEffect(() => {
		if (!open) return;
		form.resetFields();
		form.setFieldsValue({ trigger_type: 'manual' });
	}, [form, open, fromState?.id]);

	return (
		<NewModal
			open={open}
			width={520}
			title={t('workflow.editor.addTransitionTitle', { from: fromState?.name || '' })}
			okText={t('workflow.editor.addTransitionOk')}
			cancelText={t('common.cancel')}
			onCancel={onCancel}
			onOk={() => form.submit()}
			destroyOnHidden
		>
			<Form
				form={form}
				layout="vertical"
				className="mt-4"
				onFinish={onSubmit}
				size="large"
				initialValues={{ trigger_type: 'manual' }}
			>
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
										style={{ background: getWorkflowStateColorHex(s.color) }}
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
		</NewModal>
	);
}
