/**
 * 我的任务 mock：后端 `/tasks` 未就绪时替代分页响应。
 * 仅返回当前用户为创建者或执行者的任务；支持 title / priority / statusType 筛选。
 */

import { PROJECT_PRIORITY } from '../../enum/project';
import { WORKFLOW_STATE_TYPE } from '../../enum/workflow';
import { MOCK_PROJECT_LIST } from '../ProjectList/mockData';

/** @type {Record<number, { id: number, type: string, name: string }>} */
export const MOCK_WORKFLOW_STATE_MAP = {
	101: { id: 101, type: WORKFLOW_STATE_TYPE.INITIAL, name: '待办' },
	102: { id: 102, type: WORKFLOW_STATE_TYPE.IN_PROGRESS, name: '进行中' },
	103: { id: 103, type: WORKFLOW_STATE_TYPE.IN_PROGRESS, name: '评审' },
	104: { id: 104, type: WORKFLOW_STATE_TYPE.COMPLETED, name: '完成' },
};

/** @type {Record<number, string>} */
export const MOCK_MEMBER_MAP = {
	1001: '张明',
	1002: '李华',
	1003: '王芳',
};

/**
 * @param {number|string|undefined} userId
 * @param {string|undefined} userName
 * @returns {object[]}
 */
function buildMockMyTasks(userId, userName) {
	const uid = Number(userId) || 1001;
	const creatorName = userName || MOCK_MEMBER_MAP[uid] || 'guanlu';

	return [
		{
			id: 2001,
			project_id: 1,
			title: '子任务',
			status: 'todo',
			workflow_state_id: 101,
			priority: PROJECT_PRIORITY.LOW,
			type: 'task',
			type_label: '任务',
			creator_id: uid,
			assignee_id: uid,
			due_date: '2026-08-17',
		},
		{
			id: 2002,
			project_id: 1,
			title: '222',
			status: 'in_progress',
			workflow_state_id: 102,
			priority: PROJECT_PRIORITY.LOW,
			type: 'requirement',
			type_label: '需求',
			creator_id: uid,
			assignee_id: uid,
			due_date: '2026-08-25',
		},
		{
			id: 2003,
			project_id: 1,
			title: '1',
			status: 'todo',
			workflow_state_id: 101,
			priority: PROJECT_PRIORITY.LOW,
			type: 'task',
			type_label: '任务',
			creator_id: uid,
			assignee_id: 1002,
			due_date: null,
		},
		{
			id: 2004,
			project_id: 2,
			title: '梳理 onboarding 流程',
			status: 'in_progress',
			workflow_state_id: 102,
			priority: PROJECT_PRIORITY.MEDIUM,
			type: 'feature',
			type_label: '功能',
			creator_id: 1002,
			assignee_id: uid,
			due_date: '2026-09-01',
		},
	].map((task) => ({
		...task,
		creator_name: task.creator_id === uid ? creatorName : MOCK_MEMBER_MAP[task.creator_id],
		assignee_name: task.assignee_id === uid ? creatorName : MOCK_MEMBER_MAP[task.assignee_id],
	}));
}

/**
 * @param {{ page?: number, limit?: number, filter?: { title?: string, priority?: number|string, statusType?: string } }} params
 * @param {number|string|undefined} userId
 * @param {string|undefined} userName
 * @returns {{ list: object[], total: number }}
 */
export function getMockMyTaskPage(params, userId, userName) {
	const page = Number(params?.page) || 1;
	const limit = Number(params?.limit) || 15;
	const filter = params?.filter || {};
	const { title, priority, statusType } = filter;
	const uid = Number(userId) || 1001;

	let list = buildMockMyTasks(uid, userName).filter(
		(task) => task.creator_id === uid || task.assignee_id === uid,
	);

	if (title?.trim()) {
		const kw = title.trim().toLowerCase();
		list = list.filter((item) => item.title?.toLowerCase().includes(kw));
	}

	if (priority !== undefined && priority !== '') {
		list = list.filter((item) => item.priority === Number(priority));
	}

	if (statusType) {
		list = list.filter((item) => {
			const state = MOCK_WORKFLOW_STATE_MAP[item.workflow_state_id];
			return state?.type === statusType;
		});
	}

	const total = list.length;
	const start = (page - 1) * limit;

	return {
		list: list.slice(start, start + limit),
		total,
	};
}

/**
 * @param {number|string|undefined} projectId
 * @returns {string}
 */
export function getMockProjectName(projectId) {
	if (projectId == null) return '';
	const project = MOCK_PROJECT_LIST.find((item) => String(item.id) === String(projectId));
	if (String(projectId) === '1') return 'zb';
	return project?.name || '';
}
