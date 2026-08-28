'use client';

import { Card, Button, Space, Tag, Popconfirm, Divider } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined, SettingOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';
import {
	getWorkflowStateColorHex,
	toWorkflowStateColorPayload,
} from '@repo/project-management/enum/workflow';

/**
 * 右侧状态详情：基础信息 + 流转规则列表。
 */

/**
 * 详情字段行：左标签、右内容；窄侧栏内标签不换行。
 * @param {{ label: string, children: import('react').ReactNode }} props
 */
function InfoRow({ label, children }) {
	return (
		<div className="grid grid-cols-[max-content_1fr] gap-x-3 items-start">
			<span className="text-[13px] text-base-label leading-6">{label}</span>
			<div className="min-w-0 text-[13px] leading-6">{children}</div>
		</div>
	);
}

export function StateDetailsPanel({
	selectedState,
	states,
	transitions,
	canEdit,
	stageOptions,
	colorOptions,
	onEditState,
	onAddTransition,
	onConfigTransition,
	onDeleteTransition,
}) {
	const t = useTranslations();

	if (!selectedState) {
		return (
			<Card>
				<div className="text-center py-10 text-base-label">
					{t('workflow.editor.selectStateHint')}
				</div>
			</Card>
		);
	}

	const outgoing = transitions.filter((trans) => trans.from_state_id === selectedState.id);
	const colorHex = getWorkflowStateColorHex(selectedState.color);
	const colorKey = toWorkflowStateColorPayload(selectedState.color);
	const colorLabel = colorOptions.find((c) => c.value === colorKey)?.label || colorHex;
	const phaseLabel = stageOptions.find((opt) => opt.value === selectedState.phase)?.label;

	return (
		<Card
			title={t('workflow.editor.stateDetailsTitle')}
			size="small"
			extra={
				canEdit ? (
					<Button
						size="small"
						icon={<EditOutlined />}
						onClick={() => onEditState(selectedState)}
					>
						{t('workflow.editor.edit')}
					</Button>
				) : null
			}
		>
			<div className="mb-4">
				<div className="text-xs font-semibold uppercase tracking-wider text-base-label mb-2">
					{t('workflow.editor.basicInfo')}
				</div>
				<div className="flex rounded-lg border border-base-border overflow-hidden bg-gray-light">
					<div className="w-1 shrink-0" style={{ background: colorHex }} />
					<div className="flex-1 min-w-0 px-3 py-2.5 space-y-2.5">
						<InfoRow label={t('workflow.editor.fieldName')}>
							<div className="flex items-center gap-1.5 flex-wrap">
								<span className="font-medium text-base-text">
									{displayValue(selectedState.name)}
								</span>
								{/* {selectedState.is_initial ? (
									<Tag className="m-0 text-[10px] leading-4">
										{t('workflow.editor.tagInitial')}
									</Tag>
								) : null}
								{selectedState.is_terminal ? (
									<Tag className="m-0 text-[10px] leading-4">
										{t('workflow.editor.tagTerminal')}
									</Tag>
								) : null} */}
							</div>
						</InfoRow>
						<InfoRow label={t('workflow.editor.fieldType')}>
							{phaseLabel ? (
								<Tag className="m-0">{phaseLabel}</Tag>
							) : (
								displayValue(null)
							)}
						</InfoRow>
						<InfoRow label={t('workflow.editor.fieldColor')}>
							<span className="inline-flex items-center gap-1.5 rounded-full border border-base-border bg-white px-2">
								<span
									className="h-2.5 w-2.5 rounded-full shrink-0"
									style={{ background: colorHex }}
								/>
								<span className="text-base-text">{displayValue(colorLabel)}</span>
							</span>
						</InfoRow>
						{selectedState.description ? (
							<InfoRow label={t('workflow.editor.fieldDescription')}>
								<p className="m-0 text-base-gray break-words whitespace-pre-wrap">
									{selectedState.description}
								</p>
							</InfoRow>
						) : null}
					</div>
				</div>
			</div>

			<Divider className="!my-3" />

			<div>
				<div className="text-xs font-semibold uppercase tracking-wider text-base-label mb-2">
					{t('workflow.editor.transitionRules')}
				</div>
				{outgoing.map((trans) => {
					const toState = states.find((s) => s.id === trans.to_state_id);
					return (
						<Card key={trans.id} size="small" className="!mb-2">
							<div className="flex justify-between items-center">
								<span>
									{trans.name || t('workflow.editor.transitionDefaultName')} →{' '}
									<Tag color={getWorkflowStateColorHex(toState?.color)}>
										{toState?.name}
									</Tag>
								</span>
								{canEdit ? (
									<Space size={0}>
										<Button
											size="small"
											type="text"
											icon={<SettingOutlined />}
											onClick={() => onConfigTransition(trans)}
										/>
										<Popconfirm
											title={t('workflow.editor.deleteTransitionConfirm')}
											onConfirm={() => onDeleteTransition(trans.id)}
										>
											<Button
												size="small"
												type="text"
												danger
												icon={<DeleteOutlined />}
											/>
										</Popconfirm>
									</Space>
								) : null}
							</div>
						</Card>
					);
				})}
				{canEdit ? (
					<Button
						type="dashed"
						block
						icon={<PlusOutlined />}
						onClick={() => onAddTransition(selectedState)}
					>
						{t('workflow.editor.addTransition')}
					</Button>
				) : null}
			</div>
		</Card>
	);
}
