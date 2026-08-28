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

import { WORKFLOW_TRANSITION_MEMBER_OPTIONS } from '@repo/project-management/enum/workflow';

/** 与后端 `project_management_member_role` 一致，排除 ProjectOwner / TaskCreator */
export const PROJECT_MEMBER_ROLE_OPTIONS = WORKFLOW_TRANSITION_MEMBER_OPTIONS.filter(
	(option) => option.value === 'ProjectAdmin' || option.value === 'TaskExecutor',
);

export const IMAGE_FILE_EXTENSIONS = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp', 'svg'];
