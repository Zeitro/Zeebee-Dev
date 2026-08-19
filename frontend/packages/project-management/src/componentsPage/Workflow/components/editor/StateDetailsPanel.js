'use client';

import { Card, Button, Space, Tag, Popconfirm, Divider } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';

/**
 * 右侧状态详情：基础信息 + 流转规则列表。
 */
export function StateDetailsPanel({
	selectedState,
	states,
	transitions,
	canEdit,
	typeOptions,
	colorOptions,
	colorMap,
	onEditState,
	onAddTransition,
	onEditTransition,
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
			<div className="mb-3">
				<div className="font-medium mb-1">{t('workflow.editor.basicInfo')}</div>
				<p>
					<strong>{t('workflow.editor.fieldName')}:</strong> {selectedState.name}
				</p>
				<p>
					<strong>{t('workflow.editor.fieldType')}:</strong>{' '}
					{typeOptions.find((opt) => opt.value === selectedState.type)?.label}
				</p>
				<p>
					<strong>{t('workflow.editor.fieldColor')}:</strong>{' '}
					<span
						className="inline-block w-4 h-4 rounded-full align-middle mr-1"
						style={{ background: colorMap[selectedState.color] }}
					/>
					{colorOptions.find((c) => c.value === selectedState.color)?.label}
				</p>
				{selectedState.description ? (
					<p>
						<strong>{t('workflow.editor.fieldDescription')}:</strong>{' '}
						{selectedState.description}
					</p>
				) : null}
			</div>

			<Divider className="!my-3" />

			<div>
				<div className="font-medium mb-2">{t('workflow.editor.transitionRules')}</div>
				{outgoing.map((trans) => {
					const toState = states.find((s) => s.id === trans.to_state_id);
					return (
						<Card key={trans.id} size="small" className="!mb-2">
							<div className="flex justify-between items-center">
								<span>
									{trans.name || t('workflow.editor.transitionDefaultName')} →{' '}
									<Tag color={colorMap[toState?.color || 'gray']}>
										{toState?.name}
									</Tag>
								</span>
								{canEdit ? (
									<Space size={0}>
										<Button
											size="small"
											type="text"
											icon={<EditOutlined />}
											onClick={() => onEditTransition(trans)}
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
