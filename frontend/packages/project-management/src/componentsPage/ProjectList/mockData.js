/**
 * 项目 mock：看板 / 我的任务在对应接口未就绪时共用 MOCK_PROJECT_LIST。
 */

import { PROJECT_PRIORITY, PROJECT_STATUS } from '../../enum/project';

/** @type {object[]} */
export const MOCK_PROJECT_LIST = [
	{
		id: 1,
		name: 'Zeebee 官网重构',
		description: 'Next.js 15 + Tailwind，对齐 borrower 设计 token 与项目管理入口。',
		cover_image:
			'https://images.unsplash.com/photo-1557683316-973673baf926?w=400&h=200&fit=crop',
		priority: PROJECT_PRIORITY.HIGH,
		status: PROJECT_STATUS.ACTIVE,
		workflow_id: 1,
		owner_id: 1001,
		start_date: '2026-03-01',
		end_date: '2026-06-30',
		created_at: '2026-02-10T08:00:00Z',
		updated_at: '2026-08-01T12:00:00Z',
	},
	{
		id: 2,
		name: '借款方 onboarding 优化',
		description: '简化注册与组织加入流程，补齐空状态与错误提示。',
		cover_image:
			'https://images.unsplash.com/photo-1558618666-fcd25c85f82e?w=400&h=200&fit=crop',
		priority: PROJECT_PRIORITY.MEDIUM,
		status: PROJECT_STATUS.PLANNING,
		workflow_id: null,
		owner_id: 1002,
		start_date: null,
		end_date: null,
		created_at: '2026-07-01T09:30:00Z',
		updated_at: '2026-07-15T16:20:00Z',
	},
	{
		id: 3,
		name: '看板拖拽体验',
		description: '任务卡片跨列拖拽、详情抽屉与成员指派联调。',
		cover_image: '',
		priority: PROJECT_PRIORITY.URGENT,
		status: PROJECT_STATUS.ACTIVE,
		workflow_id: 2,
		owner_id: 1001,
		start_date: '2026-08-01',
		end_date: '2026-09-15',
		created_at: '2026-07-20T10:00:00Z',
		updated_at: '2026-08-10T14:45:00Z',
	},
];
