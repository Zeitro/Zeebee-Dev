'use client';

import { Card, Button, Tag, Popconfirm } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined, HolderOutlined } from '@ant-design/icons';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useTranslations } from 'next-intl';
import { getWorkflowStateColorHex } from '@repo/project-management/enum/workflow';

/**
 * 列表视图：可拖拽排序的状态卡片。
 */
export function SortableStateCard({
	state,
	stateColor,
	isSelected,
	stateTransitions,
	states,
	canEdit,
	stageOptions,
	onSelect,
	onEdit,
	onAddTransition,
	onDelete,
	onDeleteTransition,
}) {
	const t = useTranslations();
	const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
		id: `state-${state.id}`,
	});

	const style = {
		transform: CSS.Transform.toString(transform),
		transition,
		opacity: isDragging ? 0.5 : 1,
		zIndex: isDragging ? 999 : 'auto',
	};

	return (
		<div ref={setNodeRef} style={style}>
			<Card
				size="small"
				hoverable
				style={{
					borderColor: isSelected ? '#1890ff' : stateColor,
					borderWidth: isSelected ? 2 : 1,
					minWidth: 140,
					cursor: isDragging ? 'grabbing' : 'pointer',
				}}
				onClick={onSelect}
				styles={{ body: { padding: '8px 16px' } }}
			>
				<div
					{...attributes}
					{...listeners}
					className="flex items-center gap-2"
					style={{ cursor: isDragging ? 'grabbing' : 'grab' }}
					onClick={(e) => e.stopPropagation()}
				>
					<HolderOutlined className="text-gray-medium text-sm" />
					<span
						className="inline-block w-3 h-3 rounded-full shrink-0"
						style={{ background: stateColor }}
					/>
					<span className="font-medium">{state.name}</span>
				</div>
				<div className="mt-1 text-xs text-base-label">
					{state.is_initial ? (
						<Tag className="text-[10px]">{t('workflow.editor.tagInitial')}</Tag>
					) : null}
					{state.is_terminal ? (
						<Tag className="text-[10px]">{t('workflow.editor.tagTerminal')}</Tag>
					) : null}
					<span>{stageOptions.find((opt) => opt.value === state.phase)?.label}</span>
				</div>
				{stateTransitions.length > 0 ? (
					<div className="mt-2 pt-1.5 border-t border-gray-border">
						<div className="text-[11px] text-base-label mb-1">
							{t('workflow.editor.transitionTo')}
						</div>
						{stateTransitions.map((trans) => {
							const toState = states.find((s) => s.id === trans.to_state_id);
							return (
								<div key={trans.id} className="flex items-center gap-1 mb-0.5">
									<Tag
										className="text-[10px] m-0"
										color={getWorkflowStateColorHex(toState?.color)}
										closable={canEdit}
										onClose={(e) => {
											e.preventDefault();
											e.stopPropagation();
											onDeleteTransition(trans.id);
										}}
									>
										→ {toState?.name || '?'} (
										{trans.name || t('workflow.editor.transitionDefaultName')})
									</Tag>
								</div>
							);
						})}
					</div>
				) : null}
				{canEdit ? (
					<div className="mt-1.5 flex gap-1">
						<Button
							size="small"
							type="text"
							icon={<EditOutlined />}
							onClick={(e) => {
								e.stopPropagation();
								onEdit();
							}}
						/>
						<Button
							size="small"
							type="text"
							icon={<PlusOutlined />}
							onClick={(e) => {
								e.stopPropagation();
								onAddTransition();
							}}
						/>
						<Popconfirm
							title={t('workflow.editor.deleteStateConfirm')}
							onConfirm={onDelete}
							onPopupClick={(e) => e.stopPropagation()}
						>
							<Button
								size="small"
								type="text"
								danger
								icon={<DeleteOutlined />}
								onClick={(e) => e.stopPropagation()}
							/>
						</Popconfirm>
					</div>
				) : null}
			</Card>
		</div>
	);
}
