/**
 * 将工作流详情 / 模板分页 list 项收成预览弹层结构。纯映射，不改入参、不发请求。
 * 工作流 / 模板：`stages` + 嵌套 `transitions`（from_stage_id / to_stage_id）。
 * 另存模板：可含 `snapshot_board` / `snapshot_workflow`，或 `workflows[].states` / `tasks`。
 */

import { mapWorkflowDetail } from '../../Workflow/lib/mapWorkflowDetail';

/**
 * @param {object[]|null|undefined} tasks
 * @returns {object[]}
 */
function mapTasks(tasks) {
	if (!Array.isArray(tasks)) return [];
	return tasks.map((task) => ({
		id: task.id,
		title: task.title,
		parent_id: task.parent_id ?? task.parent_task_id ?? null,
		state_name: task.state_name,
		workflow_state_id: task.workflow_state_id ?? task.stage_id ?? null,
	}));
}

/**
 * @param {object} data
 * @returns {object}
 */
function fromWorkflow(data) {
	const mapped = mapWorkflowDetail(data);
	return {
		name: data.name,
		description: data.description,
		snapshot_board: data.snapshot_board,
		snapshot_workflow: data.snapshot_workflow,
		states: mapped?.states || [],
		transitions: mapped?.transitions || [],
		tasks: mapTasks(data.tasks),
	};
}

/**
 * @param {object[]} workflows
 * @returns {object[]}
 */
function mergeWorkflowStates(workflows) {
	const seen = new Map();
	workflows.forEach((workflow) => {
		const states = Array.isArray(workflow?.states) ? [...workflow.states] : [];
		states.sort((a, b) => (Number(a?.sort_order) || 0) - (Number(b?.sort_order) || 0));
		states.forEach((state) => {
			if (state?.name && !seen.has(state.name)) seen.set(state.name, state);
		});
	});
	return [...seen.values()];
}

/**
 * @param {object[]} workflows
 * @param {object[]} states
 * @returns {object[]}
 */
function mergeTransitions(workflows, states) {
	const idByName = new Map(states.map((state) => [state.name, state.id]));
	const out = [];
	const seen = new Set();
	let index = 1;
	workflows.forEach((workflow) => {
		const nameById = new Map((workflow?.states || []).map((state) => [state.id, state.name]));
		(workflow?.transitions || []).forEach((trans) => {
			const from = nameById.get(trans.from_state_id);
			const to = nameById.get(trans.to_state_id);
			if (!from || !to) return;
			const key = `${from}->${to}`;
			if (seen.has(key)) return;
			seen.add(key);
			out.push({
				id: trans.id || index++,
				from_state_id: idByName.get(from),
				to_state_id: idByName.get(to),
				name: trans.name,
			});
		});
	});
	return out;
}

/**
 * @param {object|null|undefined} data
 * @returns {object|null}
 */
export function mapTemplatePreview(data) {
	if (!data || typeof data !== 'object') return null;

	const workflows = Array.isArray(data.workflows) ? data.workflows : [];
	if (workflows.length) {
		const states = mergeWorkflowStates(workflows);
		return {
			name: data.name,
			description: data.description,
			snapshot_board: data.snapshot_board,
			snapshot_workflow: data.snapshot_workflow,
			states,
			transitions: mergeTransitions(workflows, states),
			tasks: mapTasks(data.tasks),
		};
	}

	if (Array.isArray(data.stages)) return fromWorkflow(data);

	return {
		name: data.name,
		description: data.description,
		snapshot_board: data.snapshot_board,
		snapshot_workflow: data.snapshot_workflow,
		states: Array.isArray(data.states) ? data.states : [],
		transitions: Array.isArray(data.transitions) ? data.transitions : [],
		tasks: mapTasks(data.tasks),
	};
}

/**
 * 封面状态色条用，按 sort_order 排。
 * @param {object|null|undefined} workflow
 * @returns {object[]}
 */
export function sortedStages(workflow) {
	const stages = Array.isArray(workflow?.stages) ? [...workflow.stages] : [];
	return stages.sort((a, b) => (Number(a?.sort_order) || 0) - (Number(b?.sort_order) || 0));
}

/**
 * 卡片摘要：`description` 优先，否则按阶段名拼 `A → B`；都没有返回空串，由调用方给兜底文案。
 * @param {object|null|undefined} template 模板分页 list 项
 * @returns {string}
 */
export function templateSummary(template) {
	const description = String(template?.description || '').trim();
	if (description) return description;
	return sortedStages(template)
		.map((stage) => stage?.name)
		.filter(Boolean)
		.join(' → ');
}
