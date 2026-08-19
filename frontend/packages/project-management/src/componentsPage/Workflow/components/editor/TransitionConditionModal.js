'use client';

import { Modal, Form, Select, Switch } from 'antd';
import { useTranslations } from 'next-intl';

/**
 * 流转条件 Modal（成员权限等）。
 */
export function TransitionConditionModal({ open, form, memberOptions, onCancel, onSubmit }) {
	const t = useTranslations();

	return (
		<Modal
			title={t('workflow.editor.conditionModalTitle')}
			open={open}
			onCancel={onCancel}
			onOk={() => form.submit()}
		>
			<Form form={form} layout="vertical" onFinish={onSubmit}>
				<Form.Item name="allowed_members" label={t('workflow.editor.allowedMembersLabel')}>
					<Select
						mode="multiple"
						placeholder={t('workflow.editor.selectPlaceholder')}
						options={memberOptions}
						allowClear
					/>
				</Form.Item>
				<Form.Item label={t('workflow.editor.requiredFieldsLabel')}>
					<span className="text-xs text-base-label">
						{t('workflow.editor.requiredFieldsHint')}
					</span>
				</Form.Item>
				<Form.Item
					label={t('workflow.editor.requireCommentLabel')}
					valuePropName="checked"
					name="require_comment"
				>
					<Switch
						checkedChildren={t('workflow.editor.switchOn')}
						unCheckedChildren={t('workflow.editor.switchOff')}
					/>
				</Form.Item>
				<Form.Item label={t('workflow.editor.approvalFieldsLabel')}>
					<span className="text-blue-medium cursor-pointer">
						{t('workflow.editor.addApprovalField')}
					</span>
				</Form.Item>
			</Form>
		</Modal>
	);
}
