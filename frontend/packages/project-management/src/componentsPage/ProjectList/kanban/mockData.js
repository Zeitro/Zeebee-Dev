/**
 * 看板 mock：后端工作流/任务/成员接口未就绪时替代 GET 响应。
 * 项目 2（无 workflow_id）返回空工作流，用于 EmptyWorkflow；其余项目展示示例列与任务。
 */

import { PROJECT_PRIORITY } from '../../../enum/project';
import { WORKFLOW_STATE_TYPE } from '../../../enum/workflow';
import { MOCK_PROJECT_LIST } from '../mockData';

const MOCK_WORKFLOW_ID = 1;

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
		description: '',
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
		description: '',
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
		description: '',
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
		description: '',
	},
];

/** @type {object[]} */
const MOCK_MEMBERS = [
	{
		id: 1001,
		name: '张明',
		email: 'zhang.ming@example.com',
		role: 'admin',
		status: 'active',
	},
	{
		id: 1002,
		name: '李华',
		email: 'li.hua@example.com',
		role: 'member',
		status: 'active',
	},
	{
		id: 1003,
		name: '王芳',
		email: 'wang.fang@example.com',
		role: 'member',
		status: 'active',
	},
];

/** @type {object[]} */
const MOCK_TASK_TYPES = [
	{ id: 1, name: 'feature', label: '功能' },
	{ id: 2, name: 'task', label: '任务' },
	{ id: 3, name: 'bug', label: '缺陷' },
];

/**
 * @param {number|string} projectId
 * @returns {object[]}
 */
function buildMockTasks(projectId) {
	const pid = Number(projectId) || 1;
	return [
		{
			id: 1001,
			project_id: pid,
			parent_id: null,
			title: '设计首页信息架构',
			description: '梳理 borrower 项目管理入口与路由结构。',
			workflow_state_id: 101,
			priority: PROJECT_PRIORITY.HIGH,
			type: 'feature',
			assignee_ids: [1001, 1002],
			assignee_id: 1001,
			due_date: '2026-08-25',
			status: 'todo',
			comments: [{ id: 1, content: '需要确认菜单权限' }],
		},
		{
			id: 1002,
			project_id: pid,
			parent_id: 1001,
			title: '梳理导航与路由',
			description: '',
			workflow_state_id: 101,
			priority: PROJECT_PRIORITY.MEDIUM,
			type: 'task',
			assignee_ids: [1002],
			assignee_id: 1002,
			due_date: null,
			status: 'todo',
		},
		{
			id: 1003,
			project_id: pid,
			parent_id: null,
			title: '实现 FilterTableCard 集成',
			description: '项目列表筛选与卡片布局。',
			workflow_state_id: 102,
			priority: PROJECT_PRIORITY.MEDIUM,
			type: 'feature',
			assignee_ids: [1001],
			assignee_id: 1001,
			due_date: '2026-08-20',
			status: 'in_progress',
		},
		{
			id: 1004,
			project_id: pid,
			parent_id: null,
			title: '看板拖拽联调',
			description: 'DnD 跨列 transition 与详情抽屉。',
			workflow_state_id: 103,
			priority: PROJECT_PRIORITY.URGENT,
			type: 'bug',
			assignee_ids: [1003],
			assignee_id: 1003,
			due_date: '2026-08-18',
			status: 'in_review',
		},
		{
			id: 1005,
			project_id: pid,
			parent_id: null,
			title: 'i18n 文案补齐',
			description: '',
			workflow_state_id: 104,
			priority: PROJECT_PRIORITY.LOW,
			type: 'task',
			assignee_ids: [1002],
			assignee_id: 1002,
			due_date: '2026-08-10',
			status: 'done',
		},
	];
}

/**
 * @param {number|string|undefined} projectId
 * @returns {{ project: object|null, workflow: object|null, tasks: object[] }}
 */
export function getMockKanbanBoard(projectId) {
	const project =
		MOCK_PROJECT_LIST.find((item) => String(item.id) === String(projectId)) ||
		MOCK_PROJECT_LIST[0];
	const pid = Number(projectId) || project.id;

	if (!project.workflow_id) {
		return {
			project: { ...project, id: pid },
			workflow: null,
			tasks: [],
		};
	}

	return {
		project: { ...project, id: pid },
		workflow: {
			id: MOCK_WORKFLOW_ID,
			project_id: pid,
			name: '默认看板工作流',
			description: '',
			is_default: true,
			states: MOCK_WORKFLOW_STATES,
		},
		tasks: buildMockTasks(pid),
	};
}

/**
 * @param {number|string|undefined} projectId
 * @returns {{ members: object[], projectMembers: object[], taskTypes: object[] }}
 */
export function getMockKanbanMembers(projectId) {
	const project =
		MOCK_PROJECT_LIST.find((item) => String(item.id) === String(projectId)) ||
		MOCK_PROJECT_LIST[0];
	const pid = Number(projectId) || project.id;
	const ownerId = project.owner_id || 1001;

	return {
		members: MOCK_MEMBERS,
		projectMembers: [
			{ id: 1, project_id: pid, user_id: ownerId, role: 'owner' },
			{ id: 2, project_id: pid, user_id: 1002, role: 'member' },
			{ id: 3, project_id: pid, user_id: 1003, role: 'member' },
		],
		taskTypes: MOCK_TASK_TYPES,
	};
}
