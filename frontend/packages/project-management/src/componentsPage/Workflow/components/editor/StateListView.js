'use client';

import { Card } from 'antd';
import { DndContext, pointerWithin } from '@dnd-kit/core';
import { SortableContext, horizontalListSortingStrategy } from '@dnd-kit/sortable';
import { useTranslations } from 'next-intl';
import { getWorkflowStateColorHex } from '@repo/project-management/enum/workflow';
import { SortableStateCard } from './SortableStateCard';

/**
 * 列表视图：横向可拖拽状态卡片网格。
 */
export function StateListView({
	states,
	transitions,
	selectedState,
	canEdit,
	sensors,
	typeOptions,
	colorMap,
	onDragEnd,
	onSelectState,
	onEditState,
	onAddTransition,
	onDeleteState,
	onDeleteTransition,
}) {
	const t = useTranslations();

	return (
		<Card
			title={t('workflow.editor.stateNodesTitle')}
			size="small"
			extra={
				<span className="text-xs text-base-label">
					{t('workflow.editor.dragReorderHint')}
				</span>
			}
		>
			<DndContext sensors={sensors} collisionDetection={pointerWithin} onDragEnd={onDragEnd}>
				<SortableContext
					items={states.map((s) => `state-${s.id}`)}
					strategy={horizontalListSortingStrategy}
				>
					<div className="flex flex-wrap gap-4 justify-center">
						{states.map((state) => (
							<SortableStateCard
								key={state.id}
								state={state}
								stateColor={getWorkflowStateColorHex(state.color)}
								isSelected={selectedState?.id === state.id}
								stateTransitions={transitions.filter(
									(trans) => trans.from_state_id === state.id,
								)}
								states={states}
								canEdit={canEdit}
								typeOptions={typeOptions}
								colorMap={colorMap}
								onSelect={() => onSelectState(state)}
								onEdit={() => onEditState(state)}
								onAddTransition={() => onAddTransition(state)}
								onDelete={() => onDeleteState(state.id)}
								onDeleteTransition={onDeleteTransition}
							/>
						))}
					</div>
				</SortableContext>
			</DndContext>
		</Card>
	);
}
