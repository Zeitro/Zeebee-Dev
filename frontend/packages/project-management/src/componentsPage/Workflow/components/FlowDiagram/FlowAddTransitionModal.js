'use client';

import { Input, Select } from 'antd';
import { NewModal } from '@repo/shared-methods';

/**
 * 流程图连线后添加流转弹窗。
 *
 * @param {object} props
 */
export default function FlowAddTransitionModal({
	open,
	onCancel,
	onOk,
	newTransition,
	states,
	transName,
	onTransNameChange,
	transTrigger,
	onTransTriggerChange,
	t,
}) {
	if (!open) return null;

	const fromName = states.find((s) => s.id === newTransition?.from)?.name;
	const toName = states.find((s) => s.id === newTransition?.to)?.name;

	return (
		<NewModal
			width={520}
			title={t('workflow.editor.addTransitionTitle')}
			okText={t('workflow.editor.addTransitionOk')}
			cancelText={t('common.cancel')}
			onCancel={onCancel}
			onOk={onOk}
			destroyOnHidden
		>
			<div className="mt-4 mb-4">
				<strong>
					{fromName}
					{' → '}
					{toName}
				</strong>
			</div>
			<div className="mb-3">
				<div className="mb-1 text-sm">{t('workflow.editor.transitionNameLabel')}</div>
				<Input
					placeholder={t('workflow.editor.transitionNamePlaceholder')}
					value={transName}
					onChange={(e) => onTransNameChange(e.target.value)}
				/>
			</div>
			<div>
				<div className="mb-1 text-sm">{t('workflow.editor.triggerTypeLabel')}</div>
				<Select value={transTrigger} onChange={onTransTriggerChange} className="w-full">
					<Select.Option value="manual">
						{t('workflow.editor.triggerManual')}
					</Select.Option>
					<Select.Option value="auto">{t('workflow.editor.triggerAuto')}</Select.Option>
				</Select>
			</div>
		</NewModal>
	);
}
