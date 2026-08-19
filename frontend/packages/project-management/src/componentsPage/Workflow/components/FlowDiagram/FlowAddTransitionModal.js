'use client';

import { Input, Modal, Select } from 'antd';

/**
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
	return (
		<Modal
			title={t('workflow.editor.addTransitionTitle')}
			open={open}
			onCancel={onCancel}
			onOk={onOk}
			okText={t('workflow.editor.addTransitionOk')}
			cancelText={t('common.cancel')}
		>
			<div style={{ marginBottom: 16 }}>
				<div style={{ marginBottom: 8 }}>
					<strong>
						{states.find((s) => s.id === newTransition?.from)?.name}
						{' → '}
						{states.find((s) => s.id === newTransition?.to)?.name}
					</strong>
				</div>
			</div>
			<div style={{ marginBottom: 12 }}>
				<div style={{ marginBottom: 4, fontSize: 13 }}>
					{t('workflow.editor.transitionNameLabel')}
				</div>
				<Input
					placeholder={t('workflow.editor.transitionNamePlaceholder')}
					value={transName}
					onChange={(e) => onTransNameChange(e.target.value)}
				/>
			</div>
			<div>
				<div style={{ marginBottom: 4, fontSize: 13 }}>
					{t('workflow.editor.triggerTypeLabel')}
				</div>
				<Select
					value={transTrigger}
					onChange={onTransTriggerChange}
					style={{ width: '100%' }}
				>
					<Select.Option value="manual">
						{t('workflow.editor.triggerManual')}
					</Select.Option>
					<Select.Option value="auto">{t('workflow.editor.triggerAuto')}</Select.Option>
				</Select>
			</div>
		</Modal>
	);
}
