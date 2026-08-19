/**
 * 流程图键盘快捷键：Delete/Backspace/Enter 删除节点或流转。
 */

import { useEffect } from 'react';

/**
 * @param {object} params
 * @param {import('antd/es/message/interface').MessageInstance} params.messageApi
 * @param {boolean} params.canEdit
 * @param {boolean} params.addModalOpen
 * @param {boolean} params.configDrawerOpen
 * @param {number|null} params.activeNodeId
 * @param {object|null} params.selectedTransition
 * @param {() => void} params.onSelectState
 * @param {() => void} params.onUpdate
 * @param {() => Promise<void>} params.onDeleteTransition
 * @param {object} params.workflowApi
 * @param {number} params.workflowId
 * @param {number} params.userId
 * @param {(key: string) => string} params.t
 * @param {React.Dispatch<React.SetStateAction<number|null>>} params.setActiveNodeId
 */
export function useFlowDiagramKeyboard({
	canEdit,
	addModalOpen,
	configDrawerOpen,
	activeNodeId,
	selectedTransition,
	onSelectState,
	onUpdate,
	onDeleteTransition,
	workflowApi,
	workflowId,
	userId,
	messageApi,
	t,
	setActiveNodeId,
}) {
	useEffect(() => {
		if (!canEdit) return;
		const handleKeyDown = (e) => {
			const target = e.target;
			const isInput =
				target.tagName === 'INPUT' ||
				target.tagName === 'TEXTAREA' ||
				target.isContentEditable;
			if (isInput || addModalOpen || configDrawerOpen) return;

			if (e.key === 'Delete' || e.key === 'Backspace' || e.key === 'Enter') {
				if (activeNodeId !== null) {
					e.preventDefault();
					const stateId = activeNodeId;
					setActiveNodeId(null);
					onSelectState?.(null);
					workflowApi
						.deleteState(workflowId, stateId, userId)
						.then(() => {
							messageApi.success(t('workflow.editor.stateDeleteSuccess'));
							onUpdate();
						})
						.catch(() => messageApi.error(t('workflow.editor.stateDeleteFailed')));
				} else if (selectedTransition?.id !== undefined && selectedTransition.id !== null) {
					e.preventDefault();
					onDeleteTransition();
				}
			}
		};
		window.addEventListener('keydown', handleKeyDown);
		return () => window.removeEventListener('keydown', handleKeyDown);
	}, [
		activeNodeId,
		addModalOpen,
		canEdit,
		configDrawerOpen,
		onDeleteTransition,
		onSelectState,
		onUpdate,
		selectedTransition,
		setActiveNodeId,
		t,
		userId,
		messageApi,
		workflowApi,
		workflowId,
	]);
}
