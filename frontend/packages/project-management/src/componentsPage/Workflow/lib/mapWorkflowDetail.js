/**
 * 将工作流详情接口 data 转为编辑器内部结构，并组装 PUT body。
 * 接口：`stages` + 每阶段嵌套 `transitions`（from_stage_id / to_stage_id）。
 * 编辑器：扁平 `states` + `transitions`（from_state_id / to_state_id）。
 * 添加 / 编辑 / 删除状态与流转走同一 PUT：改 `stages`（含嵌套 transitions）后整份提交。
 */

import { TASK_PHASE } from '@repo/project-management/enum/task';
import {
	WORKFLOW_STATE_TYPE,
	toWorkflowStateColorPayload,
} from '@repo/project-management/enum/workflow';

const PHASE_TO_STATE_TYPE = {
	[TASK_PHASE.PENDING]: WORKFLOW_STATE_TYPE.INITIAL,
	[TASK_PHASE.IN_PROGRESS]: WORKFLOW_STATE_TYPE.IN_PROGRESS,
	'In Progress': WORKFLOW_STATE_TYPE.IN_PROGRESS,
	[TASK_PHASE.PAUSED]: WORKFLOW_STATE_TYPE.PAUSED,
	[TASK_PHASE.FINISHED]: WORKFLOW_STATE_TYPE.COMPLETED,
	[TASK_PHASE.COMPLETED]: WORKFLOW_STATE_TYPE.COMPLETED,
};

const STATE_TYPE_TO_PHASE = {
	[WORKFLOW_STATE_TYPE.INITIAL]: TASK_PHASE.PENDING,
	[WORKFLOW_STATE_TYPE.IN_PROGRESS]: TASK_PHASE.IN_PROGRESS,
	[WORKFLOW_STATE_TYPE.PAUSED]: TASK_PHASE.PAUSED,
	[WORKFLOW_STATE_TYPE.COMPLETED]: TASK_PHASE.COMPLETED,
};

/**
 * @param {string|null|undefined} phase
 * @returns {string}
 */
export function mapStagePhaseToStateType(phase) {
	return PHASE_TO_STATE_TYPE[phase] || WORKFLOW_STATE_TYPE.IN_PROGRESS;
}

/**
 * 表单 type（initial / in_progress …）转接口 `stages.phase`（Pending / InProgress …）。
 * @param {string|null|undefined} type
 * @returns {string}
 */
export function mapStateTypeToStagePhase(type) {
	if (!type) return TASK_PHASE.IN_PROGRESS;
	if (STATE_TYPE_TO_PHASE[type]) return STATE_TYPE_TO_PHASE[type];
	if (type === 'In Progress') return TASK_PHASE.IN_PROGRESS;
	if (PHASE_TO_STATE_TYPE[type]) return type;
	return TASK_PHASE.IN_PROGRESS;
}

/**
 * 提交给 PUT 的 color：预设 `gray` 对应接口 `grey`。
 * @param {string|null|undefined} color
 * @returns {string}
 */
function toApiStageColor(color) {
	const payload = toWorkflowStateColorPayload(color);
	return payload === 'gray' ? 'grey' : payload;
}

/**
 * 把详情里的 stage 收成 PUT 字段，保留嵌套 transitions。
 * @param {object} stage
 * @returns {object|null}
 */
function toWorkflowStagePayload(stage) {
	if (!stage || typeof stage !== 'object') return null;

	const payload = {
		name: stage.name,
		phase: stage.phase,
		color: stage.color === 'gray' ? 'grey' : stage.color,
		sort_order: Number(stage.sort_order) || 0,
		transitions: Array.isArray(stage.transitions)
			? stage.transitions.map((trans) => {
					const item = {
						name: trans.name,
						from_stage_id: trans.from_stage_id ?? trans.from_state_id,
						to_stage_id: trans.to_stage_id ?? trans.to_state_id,
					};
					if (trans.id != null) item.id = trans.id;
					item.allowed_actor_types = Array.isArray(trans.allowed_actor_types)
						? trans.allowed_actor_types
						: [];
					return item;
				})
			: [],
	};
	if (stage.id != null) payload.id = stage.id;
	if (stage.description != null && stage.description !== '') {
		payload.description = stage.description;
	}
	return payload;
}

/**
 * 按状态弹窗表单改 stages：编辑匹配 id，新增则 append（无 id）。
 * @param {object[]|null|undefined} stages GET 详情的 stages
 * @param {{ editingState?: object|null, values: { name: string, phase: string, color: string, description?: string } }} params
 * @returns {object[]}
 */
export function applyStateFormToStages(stages, { editingState, values }) {
	const current = Array.isArray(stages) ? stages.map((stage) => ({ ...stage })) : [];
	const fields = {
		name: values.name,
		phase: values.phase || mapStateTypeToStagePhase(values.type),
		color: toApiStageColor(values.color),
		description: values.description || '',
	};

	if (editingState?.id != null) {
		const idx = current.findIndex((stage) => String(stage.id) === String(editingState.id));
		if (idx === -1) return current;
		current[idx] = { ...current[idx], ...fields };
		return current;
	}

	const maxOrder = current.reduce(
		(max, stage) => Math.max(max, Number(stage.sort_order) || 0),
		0,
	);
	current.push({
		...fields,
		sort_order: maxOrder + 10,
		transitions: [],
	});
	return current;
}

/**
 * @param {object[]|null|undefined} stages
 * @returns {object[]}
 */
function cloneStagesWithTransitions(stages) {
	return (Array.isArray(stages) ? stages : []).map((stage) => ({
		...stage,
		transitions: Array.isArray(stage.transitions)
			? stage.transitions.map((trans) => ({ ...trans }))
			: [],
	}));
}

/**
 * 在 from 阶段的 transitions 末尾追加一条（无 id，由后端生成）。
 * @param {object[]|null|undefined} stages
 * @param {{ fromStageId: number|string, values: { name?: string, to_state_id?: number|string, to_stage_id?: number|string, allowed_actor_types?: string[], allowed_members?: string[] } }} params
 * @returns {object[]}
 */
export function applyAddTransitionToStages(stages, { fromStageId, values }) {
	const current = cloneStagesWithTransitions(stages);
	const idx = current.findIndex((stage) => String(stage.id) === String(fromStageId));
	if (idx === -1) return current;
	current[idx].transitions.push({
		name: values.name || '',
		from_stage_id: fromStageId,
		to_stage_id: Number(values.to_state_id ?? values.to_stage_id),
		allowed_actor_types: values.allowed_actor_types || values.allowed_members || [],
	});
	return current;
}

/**
 * 按 id 更新某条 transition（name / 目标阶段 / 允许角色）。
 * @param {object[]|null|undefined} stages
 * @param {{ transitionId: number|string, values: { name?: string, to_state_id?: number|string, to_stage_id?: number|string, allowed_actor_types?: string[], allowed_members?: string[] } }} params
 * @returns {object[]}
 */
export function applyUpdateTransitionToStages(stages, { transitionId, values }) {
	const current = cloneStagesWithTransitions(stages);
	for (const stage of current) {
		const transIdx = stage.transitions.findIndex(
			(trans) => String(trans.id) === String(transitionId),
		);
		if (transIdx === -1) continue;
		const prev = stage.transitions[transIdx];
		const next = { ...prev };
		if (values.name != null) next.name = values.name;
		if (values.to_state_id != null || values.to_stage_id != null) {
			next.to_stage_id = Number(values.to_state_id ?? values.to_stage_id);
		}
		if (values.allowed_actor_types != null || values.allowed_members != null) {
			next.allowed_actor_types = values.allowed_actor_types || values.allowed_members || [];
		}
		stage.transitions[transIdx] = next;
		break;
	}
	return current;
}

/**
 * 按 id 从各 stage.transitions 中删除流转。
 * @param {object[]|null|undefined} stages
 * @param {{ transitionId: number|string }} params
 * @returns {object[]}
 */
export function applyDeleteTransitionFromStages(stages, { transitionId }) {
	const current = cloneStagesWithTransitions(stages);
	current.forEach((stage) => {
		stage.transitions = stage.transitions.filter(
			(trans) => String(trans.id) !== String(transitionId),
		);
	});
	return current;
}

/**
 * 组装 PUT `/project_management/workflow/:id` 的 body。
 * @param {object|null|undefined} workflow GET 详情（map 后仍含 stages）
 * @param {object[]|undefined} stages 覆盖 stages；默认用 workflow.stages
 * @returns {{ name: string, description: string, project_id: number, stages: object[] }}
 */
export function buildWorkflowUpdateBody(workflow, stages) {
	const list = Array.isArray(stages)
		? stages
		: Array.isArray(workflow?.stages)
			? workflow.stages
			: [];
	return {
		name: workflow?.name || '',
		description: workflow?.description || '',
		project_id: workflow?.project_id ?? 0,
		stages: list.map(toWorkflowStagePayload).filter(Boolean),
	};
}

/**
 * @param {string|null|undefined} color
 * @returns {string|undefined}
 */
function normalizeStateColor(color) {
	return color === 'grey' ? 'gray' : color;
}

/**
 * 从 stages 嵌套 / 顶层 transitions 收集流转，按 id 去重。
 * @param {object} data
 * @returns {object[]}
 */
function collectTransitions(data) {
	const seen = new Set();
	const list = [];

	const push = (trans) => {
		if (!trans || typeof trans !== 'object') return;
		if (trans.id != null) {
			if (seen.has(trans.id)) return;
			seen.add(trans.id);
		}
		list.push({
			...trans,
			from_state_id: trans.from_stage_id ?? trans.from_state_id,
			to_state_id: trans.to_stage_id ?? trans.to_state_id,
		});
	};

	(data.stages || []).forEach((stage) => {
		(stage.transitions || []).forEach(push);
	});
	if (list.length === 0 && Array.isArray(data.transitions)) {
		data.transitions.forEach(push);
	}
	return list;
}

/**
 * @param {object|null|undefined} data GET `/project_management/workflow/:id` 解包后的对象
 * @returns {object|null}
 */
export function mapWorkflowDetail(data) {
	if (!data || typeof data !== 'object') return null;

	const stages = Array.isArray(data.stages) ? [...data.stages] : [];
	stages.sort((a, b) => (Number(a?.sort_order) || 0) - (Number(b?.sort_order) || 0));

	const states = stages.map((stage, index) => {
		const rest = { ...stage };
		delete rest.transitions;
		const type = mapStagePhaseToStateType(stage.phase);
		return {
			...rest,
			color: normalizeStateColor(stage.color),
			type,
			is_initial: index === 0,
			is_terminal: type === WORKFLOW_STATE_TYPE.COMPLETED,
		};
	});

	return {
		...data,
		states,
		transitions: collectTransitions(data),
	};
}
