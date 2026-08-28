'use client';

import { useEffect, useMemo } from 'react';
import { Form, Select } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { useTranslations } from 'next-intl';
import { localizeSelectOptions } from '@repo/project-management/enum/task';
import { getWorkflowMemberRoleI18nKey } from '@repo/project-management/enum/workflow';
import { useOptionsStore } from '@repo/project-management/store/options';

/**
 * 流转条件 Modal（成员权限等）。
 * 可执行成员下来自公共接口 `project_management_member_role`（store.memberOptions）。
 * 当前编辑器入口改用 FlowConfigDrawer，本组件暂保留以便切回。
 */
export function TransitionConditionModal({ open, form, onCancel, onSubmit }) {
	const t = useTranslations();
	const memberOptionsRaw = useOptionsStore((state) => state.memberOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);
	const memberOptions = useMemo(
		() => localizeSelectOptions(memberOptionsRaw, t, getWorkflowMemberRoleI18nKey),
		[memberOptionsRaw, t],
	);

	useEffect(() => {
		if (!open) return;
		if (useOptionsStore.getState().memberOptions.length) return;
		getOptions('project_management_member_role', 'memberOptions').catch(() => {});
	}, [open, getOptions]);

	return (
		<NewModal
			open={open}
			width={600}
			title={t('workflow.editor.conditionModalTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			onCancel={onCancel}
			onOk={() => form.submit()}
			destroyOnHidden
		>
			<Form form={form} layout="vertical" className="mt-4" onFinish={onSubmit} size="large">
				<Form.Item
					name="allowed_actor_types"
					label={t('workflow.editor.allowedMembersLabel')}
				>
					<Select
						mode="multiple"
						placeholder={t('workflow.editor.selectPlaceholder')}
						options={memberOptions}
						allowClear
					/>
				</Form.Item>
				{/* <Form.Item label={t('workflow.editor.requiredFieldsLabel')}>
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
				</Form.Item> */}
			</Form>
		</NewModal>
	);
}
