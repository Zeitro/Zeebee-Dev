/**
 * 流程图交互状态、事件处理与键盘快捷键。
 */

import { useState, useRef, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { computeLayout, NODE_H, NODE_W } from '../flowLayout';
import { useFlowDiagramKeyboard } from './useFlowDiagramKeyboard';
import { useFlowDiagramInteraction } from './useFlowDiagramInteraction';

const DEFAULT_CONFIG_DATA = {
	allowedRoles: ['project_admin', 'project_owner'],
	requireNote: false,
	requiredFields: [],
	approvalFields: [],
};

/**
 * @param {object} params
 * @param {object} params.workflow
 * @param {object} params.workflowApi
 * @param {() => void} params.onUpdate
 * @param {boolean} params.canEdit
 * @param {number} params.userId
 * @param {import('antd/es/message/interface').MessageInstance} params.messageApi
 * @param {(state: object|null) => void} [params.onSelectState]
 */
export function useFlowDiagram({
	workflow,
	workflowApi,
	messageApi,
	onUpdate,
	canEdit,
	userId,
	onSelectState,
}) {
	const t = useTranslations();

	const states = useMemo(
		() => [...(workflow?.states || [])].sort((a, b) => a.sort_order - b.sort_order),
		[workflow?.states],
	);
	const transitions = useMemo(() => workflow?.transitions || [], [workflow?.transitions]);
	const workflowId = workflow?.id;

	const svgRef = useRef(null);
	const containerRef = useRef(null);

	const [nodePositions, setNodePositions] = useState(() => computeLayout(states));
	const [selectedTransition, setSelectedTransition] = useState(null);
	const [selectedTransitionKey, setSelectedTransitionKey] = useState(null);
	const [addModalOpen, setAddModalOpen] = useState(false);
	const [newTransition, setNewTransition] = useState(null);
	const [transName, setTransName] = useState('');
	const [transTrigger, setTransTrigger] = useState('manual');
	const [configDrawerOpen, setConfigDrawerOpen] = useState(false);
	const [configTransition, setConfigTransition] = useState(null);
	const [configData, setConfigData] = useState(DEFAULT_CONFIG_DATA);

	const positions = nodePositions;

	const maxX = Math.max(...Array.from(positions.values()).map((p) => p.x), 0) + NODE_W + 100;
	const maxY = Math.max(...Array.from(positions.values()).map((p) => p.y), 0) + NODE_H + 200;

	const hasTransition = useCallback(
		(fromId, toId) =>
			transitions.some((item) => item.from_state_id === fromId && item.to_state_id === toId),
		[transitions],
	);

	const handleDeleteTransition = useCallback(async () => {
		if (!selectedTransition) return;
		if (selectedTransition.id === undefined || selectedTransition.id === null) {
			messageApi.error(t('workflow.editor.transitionMissingId'));
			return;
		}
		const transId = selectedTransition.id;
		setSelectedTransition(null);
		setSelectedTransitionKey(null);
		try {
			await workflowApi.deleteTransition(workflowId, transId, userId);
			messageApi.success(t('workflow.editor.transitionDeleteSuccess'));
			onUpdate();
		} catch {
			messageApi.error(t('workflow.editor.transitionDeleteFailed'));
		}
	}, [messageApi, onUpdate, selectedTransition, t, userId, workflowId, workflowApi]);

	const interaction = useFlowDiagramInteraction({
		canEdit,
		t,
		messageApi,
		states,
		positions,
		containerRef,
		onSelectState,
		hasTransition,
		setSelectedTransition,
		setSelectedTransitionKey,
		setAddModalOpen,
		setNewTransition,
		setTransName,
		setTransTrigger,
		setNodePositions,
	});

	useFlowDiagramKeyboard({
		canEdit,
		addModalOpen,
		configDrawerOpen,
		activeNodeId: interaction.activeNodeId,
		selectedTransition,
		onSelectState,
		onUpdate,
		onDeleteTransition: handleDeleteTransition,
		workflowApi,
		workflowId: workflowId,
		userId,
		messageApi,
		t,
		setActiveNodeId: interaction.setActiveNodeId,
	});

	const handleAddTransition = useCallback(async () => {
		if (!newTransition || workflowId == null) return;
		try {
			await workflowApi.addTransition(workflowId, {
				from_state_id: newTransition.from,
				to_state_id: newTransition.to,
				name: transName,
				trigger_type: transTrigger,
				user_id: userId,
			});
			messageApi.success(t('workflow.editor.transitionAddSuccess'));
			setAddModalOpen(false);
			onUpdate();
		} catch {
			messageApi.error(t('workflow.editor.transitionAddFailed'));
		}
	}, [
		messageApi,
		newTransition,
		onUpdate,
		t,
		transName,
		transTrigger,
		userId,
		workflowId,
		workflowApi,
	]);

	const openConfigDrawer = useCallback((transition, e) => {
		e.stopPropagation();
		setConfigTransition(transition);
		setConfigDrawerOpen(true);
		setConfigData(DEFAULT_CONFIG_DATA);
	}, []);

	const handleSaveConfig = useCallback(() => {
		messageApi.success(t('workflow.editor.configSaveSuccess'));
		setConfigDrawerOpen(false);
	}, [messageApi, t]);

	const handleResetLayout = useCallback(() => {
		setNodePositions(computeLayout(states));
		messageApi.success(t('workflow.editor.layoutResetSuccess'));
	}, [messageApi, states, t]);

	return {
		t,
		states,
		transitions,
		svgRef,
		containerRef,
		positions,
		maxX,
		maxY,
		selectedTransition,
		selectedTransitionKey,
		addModalOpen,
		setAddModalOpen,
		newTransition,
		transName,
		setTransName,
		transTrigger,
		setTransTrigger,
		configDrawerOpen,
		setConfigDrawerOpen,
		configTransition,
		setConfigTransition,
		configData,
		setConfigData,
		canEdit,
		handleDeleteTransition,
		handleAddTransition,
		openConfigDrawer,
		handleSaveConfig,
		handleResetLayout,
		...interaction,
	};
}
