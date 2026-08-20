/**
 * 工作流编辑器（独立编辑模式）：编排 hooks 与子组件，不含业务逻辑。
 */
'use client';

import { Spin, Empty } from 'antd';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useWorkflowEditor } from './hooks/useWorkflowEditor';
import FlowDiagram from './components/FlowDiagram';
import { WorkflowEditorHeader } from './components/editor/WorkflowEditorHeader';
import { StateDetailsPanel } from './components/editor/StateDetailsPanel';
import { StateListView } from './components/editor/StateListView';
import { StateFormModal } from './components/editor/StateFormModal';
import { TransitionFormModal } from './components/editor/TransitionFormModal';
import { TransitionConditionModal } from './components/editor/TransitionConditionModal';

/**
 * @param {{ workflowId: number|string }} props
 */
export function WorkflowEditor({ workflowId }) {
	const t = useTranslations();
	const router = useRouter();
	const editor = useWorkflowEditor(workflowId);

	if (editor.loading) {
		return (
			<div className="text-center py-24">
				{editor.contextHolder}
				<Spin size="large" />
			</div>
		);
	}

	if (!editor.workflow) {
		return (
			<div className="py-24">
				{editor.contextHolder}
				<Empty description={t('workflow.editor.loadFailed')} />
			</div>
		);
	}

	const detailsPanel = (
		<StateDetailsPanel
			selectedState={editor.selectedState}
			states={editor.states}
			transitions={editor.transitions}
			canEdit
			typeOptions={editor.typeOptions}
			colorOptions={editor.colorOptions}
			colorMap={editor.colorMap}
			onEditState={editor.openEditState}
			onAddTransition={editor.openAddTransition}
			onEditTransition={editor.openEditTransition}
			onDeleteTransition={editor.handleDeleteTransition}
		/>
	);

	return (
		<div className="flex-1 overflow-hidden rounded-lg p-6">
			{editor.contextHolder}
			<WorkflowEditorHeader
				workflowName={editor.workflow?.name}
				editingName={editor.editingName}
				nameInputValue={editor.nameInputValue}
				nameInputRef={editor.nameInputRef}
				viewMode={editor.viewMode}
				canEdit
				onBack={() => router.push('/project/workflow')}
				onStartEditName={editor.startEditName}
				onNameChange={editor.setNameInputValue}
				onSaveName={editor.handleSaveWorkflowName}
				onViewModeChange={editor.setViewMode}
				onAddState={editor.openAddState}
			/>

			{editor.viewMode === 'flow' ? (
				<div className="flex gap-6 flex-1">
					<div className="flex-1 min-w-0">
						<FlowDiagram
							layoutKey={editor.flowDiagramKey}
							workflow={editor.workflow}
							workflowApi={editor.workflowApi}
							messageApi={editor.messageApi}
							onUpdate={editor.fetchWorkflow}
							canEdit
							userId={editor.userId || 0}
							onSelectState={editor.setSelectedState}
						/>
					</div>
					<div className="w-80 shrink-0">{detailsPanel}</div>
				</div>
			) : (
				<div className="flex gap-6 flex-1">
					<div className="flex-1">
						<StateListView
							states={editor.states}
							transitions={editor.transitions}
							selectedState={editor.selectedState}
							canEdit
							sensors={editor.sensors}
							typeOptions={editor.typeOptions}
							colorMap={editor.colorMap}
							onDragEnd={editor.handleDragEnd}
							onSelectState={editor.setSelectedState}
							onEditState={editor.openEditState}
							onAddTransition={editor.openAddTransition}
							onDeleteState={editor.handleDeleteState}
							onDeleteTransition={editor.handleDeleteTransition}
						/>
					</div>
					<div className="w-80">{detailsPanel}</div>
				</div>
			)}

			<StateFormModal
				open={editor.stateModalOpen}
				editingState={editor.editingState}
				form={editor.stateForm}
				colorOptions={editor.colorOptions}
				typeOptions={editor.typeOptions}
				onCancel={editor.closeStateModal}
				onSubmit={editor.handleSaveState}
			/>

			<TransitionFormModal
				open={editor.transitionModalOpen}
				fromState={editor.transFromState}
				states={editor.states}
				form={editor.transForm}
				colorMap={editor.colorMap}
				triggerOptions={editor.triggerOptions}
				onCancel={editor.closeTransitionModal}
				onSubmit={editor.handleSaveTransition}
			/>

			<TransitionConditionModal
				open={editor.conditionModalOpen}
				form={editor.conditionForm}
				memberOptions={editor.memberOptions}
				onCancel={editor.closeConditionModal}
				onSubmit={editor.handleSaveTransitionConditions}
			/>
		</div>
	);
}

export default WorkflowEditor;
