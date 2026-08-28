/**
 * 流程图交互状态、事件处理与键盘快捷键。
 */

import { useState, useRef, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import {
	applyAddTransitionToStages,
	applyDeleteTransitionFromStages,
	applyUpdateTransitionToStages,
	buildWorkflowUpdateBody,
} from '../../../lib/mapWorkflowDetail';
import { formatApiError } from '../../../../../lib/formatApiError';
import { computeLayout } from '../flowLayout';
import { computeDiagramExtents } from '../computeTransitionPath';
import { useFlowDiagramKeyboard } from './useFlowDiagramKeyboard';
import { useFlowDiagramInteraction } from './useFlowDiagramInteraction';

const DEFAULT_CONFIG_DATA = {
	allowedRoles: [],
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
 * @param {number|string|null} [params.selectedStateId]
 */
export function useFlowDiagram({
	workflow,
	workflowApi,
	messageApi,
	onUpdate,
	canEdit,
	userId,
	selectedStateId,
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

	const diagramExtents = useMemo(
		() => computeDiagramExtents(positions, states, transitions),
		[positions, states, transitions],
	);
	const maxX = diagramExtents.width;
	const maxY = diagramExtents.height;
	const canvasOffsetY = diagramExtents.offsetY;

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
			const stages = applyDeleteTransitionFromStages(workflow?.stages, {
				transitionId: transId,
			});
			await workflowApi.update(workflowId, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(t('workflow.editor.transitionDeleteSuccess'));
			onUpdate();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.transitionDeleteFailed', error));
		}
	}, [messageApi, onUpdate, selectedTransition, t, workflow, workflowId, workflowApi]);

	const interaction = useFlowDiagramInteraction({
		canEdit,
		t,
		messageApi,
		states,
		positions,
		containerRef,
		canvasOffsetY,
		onSelectState,
		selectedStateId,
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
			const stages = applyAddTransitionToStages(workflow?.stages, {
				fromStageId: newTransition.from,
				values: {
					name: transName,
					to_state_id: newTransition.to,
				},
			});
			await workflowApi.update(workflowId, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(t('workflow.editor.transitionAddSuccess'));
			setAddModalOpen(false);
			setTransName('');
			setTransTrigger('manual');
			onUpdate();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.transitionAddFailed', error));
		}
	}, [messageApi, newTransition, onUpdate, t, transName, workflow, workflowId, workflowApi]);

	const closeAddTransitionModal = useCallback(() => {
		setAddModalOpen(false);
		setTransName('');
		setTransTrigger('manual');
	}, []);

	const openConfigDrawer = useCallback((transition, e) => {
		e?.stopPropagation();
		setConfigTransition(transition);
		setConfigDrawerOpen(true);
		setConfigData({
			...DEFAULT_CONFIG_DATA,
			allowedRoles: Array.isArray(transition?.allowed_actor_types)
				? transition.allowed_actor_types
				: [],
		});
	}, []);

	const handleSaveConfig = useCallback(async () => {
		if (!workflow || !configTransition) return;
		try {
			const stages = applyUpdateTransitionToStages(workflow.stages, {
				transitionId: configTransition.id,
				values: { allowed_actor_types: configData.allowedRoles || [] },
			});
			await workflowApi.update(workflowId, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(t('workflow.editor.conditionUpdateSuccess'));
			setConfigDrawerOpen(false);
			setConfigTransition(null);
			onUpdate();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.conditionSaveFailed', error));
		}
	}, [
		configData.allowedRoles,
		configTransition,
		messageApi,
		onUpdate,
		t,
		workflow,
		workflowId,
		workflowApi,
	]);

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
		canvasOffsetY,
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
		closeAddTransitionModal,
		openConfigDrawer,
		handleSaveConfig,
		handleResetLayout,
		...interaction,
	};
}
