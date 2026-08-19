/**
 * 工作流编辑器数据与 CRUD：拉取 workflow、状态/流转增删改、拖拽排序、名称更新。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Form, message } from 'antd';
import { useSensor, useSensors, PointerSensor } from '@dnd-kit/core';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest, createWorkflowApi } from '@repo/project-management/api/workflow';
import { useWorkflowEditorOptions } from './useWorkflowEditorOptions';

/**
 * @param {number|string} workflowId
 */
export function useWorkflowEditor(workflowId) {
	const t = useTranslations();
	const [messageApi, contextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const workflowApi = useMemo(() => createWorkflowApi(request), [request]);
	const userId = useUserStore((state) => state.userInfo?.id);
	const options = useWorkflowEditorOptions();

	const [workflow, setWorkflow] = useState(null);
	const [loading, setLoading] = useState(true);
	const [selectedState, setSelectedState] = useState(null);
	const [viewMode, setViewMode] = useState('flow');

	const [stateModalOpen, setStateModalOpen] = useState(false);
	const [editingState, setEditingState] = useState(null);
	const [transitionModalOpen, setTransitionModalOpen] = useState(false);
	const [transFromState, setTransFromState] = useState(null);
	const [conditionModalOpen, setConditionModalOpen] = useState(false);
	const [editingTransition, setEditingTransition] = useState(null);

	const [stateForm] = Form.useForm();
	const [transForm] = Form.useForm();
	const [conditionForm] = Form.useForm();

	const [editingName, setEditingName] = useState(false);
	const [nameInputValue, setNameInputValue] = useState('');
	const nameInputRef = useRef(null);

	const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

	const fetchWorkflow = useCallback(async () => {
		try {
			setLoading(true);
			const wf = await workflowApi.get(workflowId);
			setWorkflow(wf || null);
		} catch {
			messageApi.error(t('workflow.editor.loadFailed'));
		} finally {
			setLoading(false);
		}
	}, [workflowApi, workflowId, t, messageApi]);

	useEffect(() => {
		void Promise.resolve().then(() => fetchWorkflow());
	}, [fetchWorkflow]);

	const states = useMemo(
		() => [...(workflow?.states || [])].sort((a, b) => a.sort_order - b.sort_order),
		[workflow?.states],
	);
	const transitions = workflow?.transitions || [];
	const flowDiagramKey = states.map((s) => s.id).join(',');

	const handleDragEnd = async (event) => {
		if (!workflow) return;
		const { active, over } = event;
		if (!over || active.id === over.id) return;

		const activeId = Number(String(active.id).replace('state-', ''));
		const overId = Number(String(over.id).replace('state-', ''));
		const currentStates = [...(workflow.states || [])].sort(
			(a, b) => a.sort_order - b.sort_order,
		);
		const activeIdx = currentStates.findIndex((s) => s.id === activeId);
		const overIdx = currentStates.findIndex((s) => s.id === overId);
		if (activeIdx === -1 || overIdx === -1) return;

		const reordered = [...currentStates];
		const [moved] = reordered.splice(activeIdx, 1);
		reordered.splice(overIdx, 0, moved);

		setWorkflow({
			...workflow,
			states: reordered.map((s, i) => ({ ...s, sort_order: i + 1 })),
		});

		try {
			await workflowApi.reorderStates(
				workflow.id,
				reordered.map((s) => s.id),
				userId || 0,
			);
			messageApi.success(t('workflow.editor.reorderSuccess'));
		} catch {
			messageApi.error(t('workflow.editor.reorderFailed'));
			fetchWorkflow();
		}
	};

	const openAddState = useCallback(() => {
		setEditingState(null);
		stateForm.resetFields();
		stateForm.setFieldsValue({ type: 'in_progress', color: 'blue' });
		setStateModalOpen(true);
	}, [stateForm]);

	const openEditState = useCallback(
		(state) => {
			setEditingState(state);
			stateForm.setFieldsValue({
				name: state.name,
				type: state.type,
				color: state.color,
				description: state.description,
			});
			setStateModalOpen(true);
		},
		[stateForm],
	);

	const handleSaveState = async (values) => {
		if (!workflow) return;
		try {
			if (editingState) {
				await workflowApi.updateState(workflow.id, editingState.id, {
					...values,
					user_id: userId || 0,
				});
				messageApi.success(t('workflow.editor.stateUpdateSuccess'));
			} else {
				const maxOrder = Math.max(0, ...(workflow.states?.map((s) => s.sort_order) || []));
				await workflowApi.addState(workflow.id, {
					...values,
					sort_order: maxOrder + 1,
					is_initial: values.type === 'initial',
					is_terminal: values.type === 'completed',
					user_id: userId || 0,
				});
				messageApi.success(t('workflow.editor.stateAddSuccess'));
			}
			setStateModalOpen(false);
			stateForm.resetFields();
			fetchWorkflow();
		} catch {
			messageApi.error(t('workflow.editor.stateSaveFailed'));
		}
	};

	const handleDeleteState = useCallback(
		async (stateId) => {
			if (!workflow) return;
			try {
				await workflowApi.deleteState(workflow.id, stateId, userId || 0);
				messageApi.success(t('workflow.editor.stateDeleteSuccess'));
				setSelectedState((prev) => (prev?.id === stateId ? null : prev));
				fetchWorkflow();
			} catch (err) {
				messageApi.error(err?.message || t('workflow.editor.stateDeleteFailed'));
			}
		},
		[fetchWorkflow, userId, workflow, workflowApi, t, messageApi],
	);

	useEffect(() => {
		const handleKeyDown = (e) => {
			const target = e.target;
			const isInput =
				target.tagName === 'INPUT' ||
				target.tagName === 'TEXTAREA' ||
				target.isContentEditable;
			const isModalOpen = stateModalOpen || transitionModalOpen || conditionModalOpen;
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
	}, [selectedState, stateModalOpen, transitionModalOpen, conditionModalOpen, handleDeleteState]);

	const openAddTransition = useCallback(
		(fromState) => {
			setTransFromState(fromState);
			transForm.resetFields();
			transForm.setFieldsValue({ trigger_type: 'manual' });
			setTransitionModalOpen(true);
		},
		[transForm],
	);

	const handleSaveTransition = async (values) => {
		if (!workflow || !transFromState) return;
		try {
			await workflowApi.addTransition(workflow.id, {
				from_state_id: transFromState.id,
				to_state_id: Number(values.to_state_id),
				name: values.name,
				trigger_type: values.trigger_type,
				user_id: userId || 0,
			});
			messageApi.success(t('workflow.editor.transitionAddSuccess'));
			setTransitionModalOpen(false);
			transForm.resetFields();
			fetchWorkflow();
		} catch {
			messageApi.error(t('workflow.editor.transitionAddFailed'));
		}
	};

	const handleDeleteTransition = async (transitionId) => {
		if (!workflow) return;
		try {
			await workflowApi.deleteTransition(workflow.id, transitionId, userId || 0);
			messageApi.success(t('workflow.editor.transitionDeleteSuccess'));
			fetchWorkflow();
		} catch {
			messageApi.error(t('workflow.editor.transitionDeleteFailed'));
		}
	};

	const parseConditions = (conditions) => {
		if (!conditions) return {};
		try {
			return JSON.parse(conditions);
		} catch {
			return {};
		}
	};

	const openEditTransition = useCallback(
		(transition) => {
			setEditingTransition(transition);
			conditionForm.resetFields();
			const cond = parseConditions(transition.conditions);
			const defaultMembers = options.memberOptions.map((m) => m.value);
			conditionForm.setFieldsValue({
				allowed_members: cond.allowed_members?.length
					? cond.allowed_members
					: defaultMembers,
				require_comment: false,
			});
			setConditionModalOpen(true);
		},
		[conditionForm, options.memberOptions],
	);

	const handleSaveTransitionConditions = async (values) => {
		if (!workflow || !editingTransition) return;
		try {
			await workflowApi.updateTransition(workflow.id, editingTransition.id, {
				conditions: JSON.stringify({ allowed_members: values.allowed_members }),
				user_id: userId || 0,
			});
			messageApi.success(t('workflow.editor.conditionUpdateSuccess'));
			setConditionModalOpen(false);
			conditionForm.resetFields();
			setEditingTransition(null);
			fetchWorkflow();
		} catch {
			messageApi.error(t('workflow.editor.conditionSaveFailed'));
		}
	};

	const handleSaveWorkflowName = async () => {
		const trimmed = nameInputValue.trim();
		if (!trimmed || !workflow) return;
		try {
			await workflowApi.update(workflow.id, { name: trimmed, user_id: userId || 0 });
			await fetchWorkflow();
			messageApi.success(t('workflow.editor.nameUpdateSuccess'));
		} catch (err) {
			messageApi.error(err?.message || t('workflow.editor.nameUpdateFailed'));
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
		transForm.resetFields();
	};

	const closeConditionModal = () => {
		setConditionModalOpen(false);
		conditionForm.resetFields();
		setEditingTransition(null);
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
		openEditTransition,
		handleSaveTransitionConditions,
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
		conditionModalOpen,
		stateForm,
		transForm,
		conditionForm,
		closeStateModal,
		closeTransitionModal,
		closeConditionModal,
		fetchWorkflow,
		workflowApi,
		userId,
		messageApi,
		contextHolder,
		...options,
	};
}
