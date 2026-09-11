'use client';

import { useEffect, useMemo } from 'react';
import { Button, Drawer, Select, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { localizeSelectOptions } from '@repo/project-management/enum/task';
import {
	getWorkflowMemberRoleI18nKey,
	getWorkflowStateColorHex,
} from '@repo/project-management/enum/workflow';
import { useOptionsStore } from '@repo/project-management/store/options';

/**
 * 流程图流转配置抽屉。
 * 可执行成员下来自公共接口 `project_management_member_role`（store.memberOptions），与流转条件 Modal 一致。
 */
export default function FlowConfigDrawer({
	open,
	onClose,
	onSave,
	configTransition,
	states,
	configData,
	onConfigDataChange,
}) {
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

	const fromState = states.find((s) => s.id === configTransition?.from_state_id);
	const toState = states.find((s) => s.id === configTransition?.to_state_id);

	return (
		<Drawer
			title={t('workflow.editor.configDrawerTitle')}
			open={open}
			onClose={onClose}
			width={530}
			extra={
				<Button type="primary" onClick={onSave}>
					{t('common.save')}
				</Button>
			}
		>
			{configTransition ? (
				<div>
					<div style={{ marginBottom: 24 }}>
						<div style={{ fontSize: 13, color: '#666', marginBottom: 8 }}>
							{t('workflow.editor.configPathLabel')}
						</div>
						<div
							style={{
								display: 'flex',
								alignItems: 'center',
								gap: 8,
								padding: '12px 16px',
								background: '#f6f8fa',
								borderRadius: 8,
							}}
						>
							<Tag
								color={getWorkflowStateColorHex(fromState?.color)}
								style={{ margin: 0 }}
							>
								{fromState?.name}
							</Tag>
							<span style={{ color: '#999' }}>→</span>
							<Tag
								color={getWorkflowStateColorHex(toState?.color)}
								style={{ margin: 0 }}
							>
								{toState?.name}
							</Tag>
							<div
								style={{
									width: 20,
									height: 20,
									borderRadius: '50%',
									border: '1.5px solid #1890ff',
									background: '#fff',
									display: 'flex',
									alignItems: 'center',
									justifyContent: 'center',
									marginLeft: 4,
								}}
							>
								<div
									style={{
										width: 12,
										height: 12,
										background:
											'repeating-linear-gradient(45deg, transparent, transparent 2px, #1890ff 2px, #1890ff 3px)',
										borderRadius: '50%',
									}}
								/>
							</div>
						</div>
						{configTransition.name ? (
							<div style={{ marginTop: 6, fontSize: 12, color: '#888' }}>
								{t('workflow.editor.configNameLabel', {
									name: configTransition.name,
								})}
							</div>
						) : null}
						<div style={{ marginTop: 4, fontSize: 11, color: '#bbb' }}>
							{t('workflow.editor.configInheritedHint')}
						</div>
					</div>

					<div style={{ marginBottom: 24 }}>
						<div
							style={{
								fontSize: 13,
								color: '#333',
								fontWeight: 500,
								marginBottom: 8,
							}}
						>
							{t('workflow.editor.configAllowedMembers')}
						</div>
						<Select
							allowClear
							mode="multiple"
							style={{ width: '100%' }}
							value={configData.allowedRoles}
							onChange={(val) =>
								onConfigDataChange({ ...configData, allowedRoles: val })
							}
							options={memberOptions}
							size="large"
							placeholder={t('workflow.editor.selectPlaceholder')}
						/>
					</div>

					{/* <div style={{ marginBottom: 24 }}>
						<div
							style={{
								display: 'flex',
								alignItems: 'center',
								justifyContent: 'space-between',
								marginBottom: 8,
							}}
						>
							<span style={{ fontSize: 13, color: '#333', fontWeight: 500 }}>
								{t('workflow.editor.configRequiredFields')}
							</span>
							<Button size="small" type="text" icon={<EditOutlined />} />
						</div>
						<div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
							<Tag
								color="blue"
								closable
								onClose={() =>
									onConfigDataChange({
										...configData,
										requiredFields: configData.requiredFields.filter(
											(f) => f !== 'title',
										),
									})
								}
							>
								{t('workflow.editor.fieldTitle')}
							</Tag>
							<Tag
								color="blue"
								closable
								onClose={() =>
									onConfigDataChange({
										...configData,
										requiredFields: configData.requiredFields.filter(
											(f) => f !== 'assignee',
										),
									})
								}
							>
								{t('workflow.editor.fieldAssignee')}
							</Tag>
							<Button
								size="small"
								type="dashed"
								icon={<PlusOutlined />}
								onClick={() => {
									onConfigDataChange({
										...configData,
										requiredFields: [...configData.requiredFields, 'new_field'],
									});
								}}
							>
								{t('workflow.editor.addField')}
							</Button>
						</div>
					</div> */}

					{/* <div style={{ marginBottom: 24 }}>
						<div
							style={{
								display: 'flex',
								alignItems: 'center',
								justifyContent: 'space-between',
							}}
						>
							<span style={{ fontSize: 13, color: '#333', fontWeight: 500 }}>
								{t('workflow.editor.configRequireNote')}
							</span>
							<Switch
								checked={configData.requireNote}
								onChange={(val) =>
									onConfigDataChange({ ...configData, requireNote: val })
								}
							/>
						</div>
						{configData.requireNote ? (
							<Input.TextArea
								style={{ marginTop: 8 }}
								rows={2}
								placeholder={t('workflow.editor.notePlaceholder')}
								defaultValue={t('workflow.editor.noteDefaultValue')}
							/>
						) : null}
					</div> */}

					{/* <div style={{ marginBottom: 24 }}>
						<div
							style={{
								fontSize: 13,
								color: '#333',
								fontWeight: 500,
								marginBottom: 8,
							}}
						>
							{t('workflow.editor.configApprovalFields')}
						</div>
						<div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
							{configData.approvalFields.map((field, idx) => (
								<Tag
									key={idx}
									color="orange"
									closable
									onClose={() =>
										onConfigDataChange({
											...configData,
											approvalFields: configData.approvalFields.filter(
												(_, i) => i !== idx,
											),
										})
									}
								>
									{field}
								</Tag>
							))}
							<Button
								size="small"
								type="dashed"
								icon={<PlusOutlined />}
								onClick={() => {
									onConfigDataChange({
										...configData,
										approvalFields: [
											...configData.approvalFields,
											t('workflow.editor.approvalFieldDefault', {
												index: configData.approvalFields.length + 1,
											}),
										],
									});
								}}
							>
								{t('workflow.editor.addApprovalField')}
							</Button>
						</div>
					</div> */}
				</div>
			) : null}
		</Drawer>
	);
}
