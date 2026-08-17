/**
 * 看板 UI 常量：卡片高度、列宽、详情抽屉宽度对齐产品 demo。
 */

export const TASK_CARD_HEIGHT_CLASS = 'h-[92px]';
export const KANBAN_COLUMN_MIN_WIDTH_CLASS = 'min-w-[230px]';
export const TASK_CARD_DUE_FORMAT = 'MM-DD';
export const TASK_DETAIL_DRAWER_WIDTH = '66.67%';
export const TASK_LIST_PAGE_SIZE = 200;

/** 任务依赖类型，与后端 `dependency_type` 一致 */
export const DEPENDENCY_TYPE_OPTIONS = [
	{ value: 'finish_to_start', labelKey: 'kanban.dependencyFinishToStart' },
	{ value: 'finish_to_finish', labelKey: 'kanban.dependencyFinishToFinish' },
	{ value: 'start_to_start', labelKey: 'kanban.dependencyStartToStart' },
	{ value: 'start_to_finish', labelKey: 'kanban.dependencyStartToFinish' },
];

export const PROJECT_MEMBER_ROLE_OPTIONS = [
	{ value: 'admin', labelKey: 'kanban.roleAdmin' },
	{ value: 'member', labelKey: 'kanban.roleMember' },
	{ value: 'executor', labelKey: 'kanban.roleExecutor' },
	{ value: 'viewer', labelKey: 'kanban.roleViewer' },
];

export const IMAGE_FILE_EXTENSIONS = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp', 'svg'];
