/**
 * 工作流编辑器 mock：GET /workflows/:id 未就绪时替代详情响应。
 * 状态 id 与看板 mock（kanban/mockData.js）对齐，便于联调。
 */

import { WORKFLOW_STATE_TYPE } from '../../enum/workflow';

const MOCK_WORKFLOW_ID = 1;

/** @type {object[]} */
const MOCK_WORKFLOW_STATES = [
	{
		id: 101,
		workflow_id: MOCK_WORKFLOW_ID,
		name: '待办',
		type: WORKFLOW_STATE_TYPE.INITIAL,
		color: 'gray',
		sort_order: 0,
		is_initial: true,
		is_terminal: false,
		description: '任务创建后的初始状态',
	},
	{
		id: 102,
		workflow_id: MOCK_WORKFLOW_ID,
		name: '进行中',
		type: WORKFLOW_STATE_TYPE.IN_PROGRESS,
		color: 'blue',
		sort_order: 1,
		is_initial: false,
		is_terminal: false,
		description: '执行中的任务',
	},
	{
		id: 103,
		workflow_id: MOCK_WORKFLOW_ID,
		name: '评审',
		type: WORKFLOW_STATE_TYPE.IN_PROGRESS,
		color: 'purple',
		sort_order: 2,
		is_initial: false,
		is_terminal: false,
		description: '等待评审确认',
	},
	{
		id: 104,
		workflow_id: MOCK_WORKFLOW_ID,
		name: '完成',
		type: WORKFLOW_STATE_TYPE.COMPLETED,
		color: 'green',
		sort_order: 3,
		is_initial: false,
		is_terminal: true,
		description: '任务已完成',
	},
];

/** @type {object[]} */
const MOCK_WORKFLOW_TRANSITIONS = [
	{
		id: 201,
		workflow_id: MOCK_WORKFLOW_ID,
		from_state_id: 101,
		to_state_id: 102,
		name: '开始处理',
		trigger_type: 'manual',
		conditions: JSON.stringify({ allowed_members: ['assignee', 'creator'] }),
	},
	{
		id: 202,
		workflow_id: MOCK_WORKFLOW_ID,
		from_state_id: 102,
		to_state_id: 103,
		name: '提交评审',
		trigger_type: 'manual',
		conditions: JSON.stringify({ allowed_members: ['assignee'] }),
	},
	{
		id: 203,
		workflow_id: MOCK_WORKFLOW_ID,
		from_state_id: 103,
		to_state_id: 102,
		name: '退回修改',
		trigger_type: 'manual',
		conditions: JSON.stringify({ allowed_members: ['owner', 'admin'] }),
	},
	{
		id: 204,
		workflow_id: MOCK_WORKFLOW_ID,
		from_state_id: 103,
		to_state_id: 104,
		name: '评审通过',
		trigger_type: 'manual',
		conditions: JSON.stringify({ allowed_members: ['owner', 'admin'] }),
	},
	{
		id: 205,
		workflow_id: MOCK_WORKFLOW_ID,
		from_state_id: 102,
		to_state_id: 104,
		name: '直接完成',
		trigger_type: 'automatic',
		conditions: null,
	},
];

/** @type {object} */
const MOCK_WORKFLOW_DETAIL = {
	id: MOCK_WORKFLOW_ID,
	project_id: 1,
	name: '默认看板工作流',
	description: '项目管理默认任务流转模板',
	is_default: true,
	is_system: false,
	states: MOCK_WORKFLOW_STATES,
	transitions: MOCK_WORKFLOW_TRANSITIONS,
};

/**
 * @param {number|string|undefined} workflowId
 * @returns {object|null}
 */
export function getMockWorkflowDetail(workflowId) {
	const id = Number(workflowId) || MOCK_WORKFLOW_ID;
	if (id !== MOCK_WORKFLOW_ID) {
		return {
			...MOCK_WORKFLOW_DETAIL,
			id,
			name: `工作流 #${id}`,
			states: MOCK_WORKFLOW_STATES.map((state) => ({ ...state, workflow_id: id })),
			transitions: MOCK_WORKFLOW_TRANSITIONS.map((trans) => ({ ...trans, workflow_id: id })),
		};
	}
	return MOCK_WORKFLOW_DETAIL;
}
