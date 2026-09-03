/**
 * 工作流编辑器数据与 CRUD：拉取 workflow、状态/流转增删改、拖拽排序、名称更新。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Form, message } from 'antd';
import { useSensor, useSensors, PointerSensor } from '@dnd-kit/core';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import {
	createRequest,
	createWorkflowApi,
	isValidProjectId,
} from '@repo/project-management/api/workflow';
import { useOptionsStore } from '@repo/project-management/store/options';
import {
	applyAddTransitionToStages,
	applyDeleteTransitionFromStages,
	applyStateFormToStages,
	applyUpdateTransitionToStages,
	buildWorkflowUpdateBody,
	mapWorkflowDetail,
} from '../lib/mapWorkflowDetail';
import { formatApiError } from '../../../lib/formatApiError';
import { useWorkflowEditorOptions } from './useWorkflowEditorOptions';

const DEFAULT_CONFIG_DATA = {
	allowedRoles: [],
	requireNote: false,
	requiredFields: [],
	approvalFields: [],
};

/**
 * @param {number|string} workflowId
 * @param {number|string} [projectId] 有值时 get / update 走 `/project_management/project/:id/workflow`
 */
export function useWorkflowEditor(workflowId, projectId) {
	const t = useTranslations();
	const [messageApi, contextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const workflowApi = useMemo(
		() => createWorkflowApi(request, { projectId }),
		[request, projectId],
	);
	const userId = useUserStore((state) => state.userInfo?.id);
	const options = useWorkflowEditorOptions();

	const [workflow, setWorkflow] = useState(null);
	const [loading, setLoading] = useState(true);
	const [selectedStateId, setSelectedStateId] = useState(null);
	const [viewMode, setViewMode] = useState('flow');

	const [stateModalOpen, setStateModalOpen] = useState(false);
	const [editingState, setEditingState] = useState(null);
	const [transitionModalOpen, setTransitionModalOpen] = useState(false);
	const [transFromState, setTransFromState] = useState(null);
	const [configDrawerOpen, setConfigDrawerOpen] = useState(false);
	const [configTransition, setConfigTransition] = useState(null);
	const [configData, setConfigData] = useState(DEFAULT_CONFIG_DATA);

	const [stateForm] = Form.useForm();
	const [transForm] = Form.useForm();

	const [editingName, setEditingName] = useState(false);
	const [nameInputValue, setNameInputValue] = useState('');
	const nameInputRef = useRef(null);
	const selectedStateIdRef = useRef(null);
	const selectedStateNameRef = useRef(null);
	/** 弹窗关闭后短暂忽略画布清空选中（避免点击穿透） */
	const suppressDeselectUntilRef = useRef(0);

	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

	const fetchWorkflow = useCallback(
		async ({ initial = false } = {}) => {
			const hasProject = isValidProjectId(projectId);
			if (!hasProject && (workflowId == null || workflowId === '')) {
				setWorkflow(null);
				setLoading(false);
				return;
			}
			const keepId = selectedStateIdRef.current;
			const keepName = selectedStateNameRef.current;
			suppressDeselectUntilRef.current = Date.now() + 500;
			try {
				if (initial) setLoading(true);
				const data = await workflowApi.get(workflowId);
				const mapped = mapWorkflowDetail(data);
				setWorkflow(mapped);
				const nextStates = mapped?.states || [];
				const restored =
					nextStates.find((s) => String(s.id) === String(keepId)) ||
					(keepName ? nextStates.find((s) => s.name === keepName) : null);
				if (restored) setSelectedStateId(restored.id);
			} catch (error) {
				setWorkflow(null);
				messageApi.error(formatApiError(t, 'workflow.editor.loadFailed', error));
			} finally {
				setLoading(false);
			}
		},
		[workflowApi, workflowId, projectId, t, messageApi],
	);

	useEffect(() => {
		void Promise.resolve().then(() => fetchWorkflow({ initial: true }));
	}, [fetchWorkflow]);

	useEffect(() => {
		const { stageOptions, memberOptions, getOptions } = useOptionsStore.getState();
		if (!stageOptions.length) {
			getOptions('project_management_stage_phase', 'stageOptions').catch(() => {});
		}
		if (!memberOptions.length) {
			getOptions('project_management_member_role', 'memberOptions').catch(() => {});
		}
	}, []);

	const states = useMemo(
		() => [...(workflow?.states || [])].sort((a, b) => a.sort_order - b.sort_order),
		[workflow?.states],
	);
	/** 按 id 从最新 states 取当前选中项，避免编辑后详情仍用旧对象 */
	const selectedState = useMemo(
		() => states.find((s) => String(s.id) === String(selectedStateId)) ?? null,
		[states, selectedStateId],
	);

	useEffect(() => {
		selectedStateIdRef.current = selectedStateId;
		if (selectedState?.name) selectedStateNameRef.current = selectedState.name;
		else if (selectedStateId == null) selectedStateNameRef.current = null;
	}, [selectedStateId, selectedState?.name]);

	const syncSelectionRefs = useCallback(
		(nextId) => {
			selectedStateIdRef.current = nextId;
			if (nextId == null) {
				selectedStateNameRef.current = null;
				return;
			}
			const state = states.find((s) => String(s.id) === String(nextId));
			if (state?.name) selectedStateNameRef.current = state.name;
		},
		[states],
	);

	const setSelectedState = useCallback(
		(next) => {
			const resolveId = (value) => {
				if (typeof value === 'function') {
					const prevState =
						states.find((s) => String(s.id) === String(selectedStateIdRef.current)) ??
						null;
					const result = value(prevState);
					return result == null ? null : result.id;
				}
				return value == null ? null : value.id;
			};
			const nextId = resolveId(next);
			if (nextId == null && Date.now() < suppressDeselectUntilRef.current) return;
			syncSelectionRefs(nextId);
			setSelectedStateId(nextId);
		},
		[states, syncSelectionRefs],
	);
	const transitions = workflow?.transitions || [];
	const flowDiagramKey = states.map((s) => s.id).join(',');

	const handleDragEnd = async (event) => {
		if (!workflow) return;
		const { active, over } = event;
		if (!over || active.id === over.id) return;

		const fromStageId = Number(String(active.id).replace('state-', ''));
		const toStageId = Number(String(over.id).replace('state-', ''));
		const currentStates = [...(workflow.states || [])].sort(
			(a, b) => a.sort_order - b.sort_order,
		);
		const fromState = currentStates.find((s) => s.id === fromStageId);
		const toState = currentStates.find((s) => s.id === toStageId);
		if (!fromState || !toState) return;

		setWorkflow({
			...workflow,
			states: currentStates.map((s) => {
				if (s.id === fromStageId) return { ...s, sort_order: toState.sort_order };
				if (s.id === toStageId) return { ...s, sort_order: fromState.sort_order };
				return s;
			}),
		});

		try {
			await workflowApi.swapStageOrder(workflow.id, {
				from_stage_id: fromStageId,
				to_stage_id: toStageId,
			});
			messageApi.success(t('workflow.editor.reorderSuccess'));
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.reorderFailed', error));
			fetchWorkflow();
		}
	};

	const openAddState = useCallback(() => {
		setEditingState(null);
		setStateModalOpen(true);
	}, []);

	const openEditState = useCallback((state) => {
		setEditingState(state);
		setStateModalOpen(true);
	}, []);

	const handleSaveState = async (values) => {
		if (!workflow) return;
		try {
			const stages = applyStateFormToStages(workflow.stages, { editingState, values });
			await workflowApi.update(workflow.id, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(
				editingState
					? t('workflow.editor.stateUpdateSuccess')
					: t('workflow.editor.stateAddSuccess'),
			);
			setStateModalOpen(false);
			stateForm.resetFields();
			fetchWorkflow();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.stateSaveFailed', error));
		}
	};

	const handleDeleteState = useCallback(
		async (stateId) => {
			if (!workflow) return;
			try {
				await workflowApi.deleteState(workflow.id, stateId, userId || 0);
				messageApi.success(t('workflow.editor.stateDeleteSuccess'));
				const nextId =
					String(selectedStateIdRef.current) === String(stateId)
						? null
						: selectedStateIdRef.current;
				syncSelectionRefs(nextId);
				setSelectedStateId(nextId);
				fetchWorkflow();
			} catch (error) {
				messageApi.error(formatApiError(t, 'workflow.editor.stateDeleteFailed', error));
			}
		},
		[fetchWorkflow, syncSelectionRefs, userId, workflow, workflowApi, t, messageApi],
	);

	useEffect(() => {
		const handleKeyDown = (e) => {
			const target = e.target;
			const isInput =
				target.tagName === 'INPUT' ||
				target.tagName === 'TEXTAREA' ||
				target.isContentEditable;
			const isModalOpen = stateModalOpen || transitionModalOpen || configDrawerOpen;
			if (isInput || isModalOpen) return;

			if (
				(e.key === 'Delete' || e.key === 'Backspace' || e.key === 'Enter') &&
				selectedState
			) {
				e.preventDefault();
				handleDeleteState(selectedState.id);
			}
		};
		window.addEventListener('keydown', handleKeyDown);
		return () => window.removeEventListener('keydown', handleKeyDown);
	}, [selectedState, stateModalOpen, transitionModalOpen, configDrawerOpen, handleDeleteState]);

	const openAddTransition = useCallback((fromState) => {
		setTransFromState(fromState);
		setTransitionModalOpen(true);
	}, []);

	const handleSaveTransition = async (values) => {
		if (!workflow || !transFromState) return;
		try {
			const stages = applyAddTransitionToStages(workflow.stages, {
				fromStageId: transFromState.id,
				values,
			});
			await workflowApi.update(workflow.id, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(t('workflow.editor.transitionAddSuccess'));
			setTransitionModalOpen(false);
			fetchWorkflow();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.transitionAddFailed', error));
		}
	};

	const handleDeleteTransition = async (transitionId) => {
		if (!workflow) return;
		try {
			const stages = applyDeleteTransitionFromStages(workflow.stages, { transitionId });
			await workflowApi.update(workflow.id, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(t('workflow.editor.transitionDeleteSuccess'));
			fetchWorkflow();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.transitionDeleteFailed', error));
		}
	};

	const openConfigDrawer = useCallback((transition) => {
		setConfigTransition(transition);
		setConfigDrawerOpen(true);
		setConfigData({
			...DEFAULT_CONFIG_DATA,
			allowedRoles: Array.isArray(transition?.allowed_actor_types)
				? transition.allowed_actor_types
				: [],
		});
	}, []);

	const closeConfigDrawer = useCallback(() => {
		setConfigDrawerOpen(false);
		setConfigTransition(null);
	}, []);

	const handleSaveConfig = async () => {
		if (!workflow || !configTransition) return;
		try {
			const stages = applyUpdateTransitionToStages(workflow.stages, {
				transitionId: configTransition.id,
				values: { allowed_actor_types: configData.allowedRoles || [] },
			});
			await workflowApi.update(workflow.id, buildWorkflowUpdateBody(workflow, stages));
			messageApi.success(t('workflow.editor.conditionUpdateSuccess'));
			setConfigDrawerOpen(false);
			setConfigTransition(null);
			fetchWorkflow();
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.conditionSaveFailed', error));
		}
	};

	const handleSaveWorkflowName = async () => {
		const trimmed = nameInputValue.trim();
		if (!trimmed || !workflow) return;
		try {
			await workflowApi.update(workflow.id, {
				...buildWorkflowUpdateBody(workflow),
				name: trimmed,
			});
			await fetchWorkflow();
			messageApi.success(t('workflow.editor.nameUpdateSuccess'));
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.editor.nameUpdateFailed', error));
		}
		setEditingName(false);
	};

	const startEditName = () => {
		setNameInputValue(workflow?.name || '');
		setEditingName(true);
		setTimeout(() => nameInputRef.current?.focus(), 0);
	};

	const closeStateModal = () => {
		setStateModalOpen(false);
		stateForm.resetFields();
	};

	const closeTransitionModal = () => {
		setTransitionModalOpen(false);
	};

	return {
		workflow,
		loading,
		states,
		transitions,
		flowDiagramKey,
		selectedState,
		setSelectedState,
		viewMode,
		setViewMode,
		sensors,
		handleDragEnd,
		openAddState,
		openEditState,
		handleSaveState,
		handleDeleteState,
		openAddTransition,
		handleSaveTransition,
		handleDeleteTransition,
		openConfigDrawer,
		closeConfigDrawer,
		handleSaveConfig,
		handleSaveWorkflowName,
		startEditName,
		editingName,
		nameInputValue,
		setNameInputValue,
		nameInputRef,
		stateModalOpen,
		editingState,
		transitionModalOpen,
		transFromState,
		configDrawerOpen,
		configTransition,
		configData,
		setConfigData,
		stateForm,
		transForm,
		closeStateModal,
		closeTransitionModal,
		fetchWorkflow,
		workflowApi,
		userId,
		messageApi,
		contextHolder,
		...options,
	};
}
