/**
 * 看板 UI 常量：卡片高度、列宽、详情抽屉宽度对齐产品 demo。
 */

// export const TASK_CARD_HEIGHT_CLASS = 'h-[92px]';
export const KANBAN_COLUMN_MIN_WIDTH_CLASS = 'min-w-[230px]';
// export const TASK_CARD_DUE_FORMAT = 'MM-DD';
export const TASK_DETAIL_DRAWER_WIDTH = '66.67%';
export const TASK_LIST_PAGE_SIZE = 200;

/** 任务依赖规则选项，与详情 `dependency_rule` 一致 */
export const DEPENDENCY_RULE_OPTIONS = [
	{ value: 'FS', labelKey: 'kanban.dependencyFinishToStart' },
	{ value: 'FF', labelKey: 'kanban.dependencyFinishToFinish' },
	{ value: 'SS', labelKey: 'kanban.dependencyStartToStart' },
	{ value: 'SF', labelKey: 'kanban.dependencyStartToFinish' },
];

import { TASK_DEPENDENCY_DIRECTION } from '@repo/project-management/enum/task';
import { WORKFLOW_TRANSITION_MEMBER_OPTIONS } from '@repo/project-management/enum/workflow';

/** 添加依赖方向选项，value 与后端 `dependency_type` 一致 */
export const DEPENDENCY_DIRECTION_OPTIONS = [
	{ value: TASK_DEPENDENCY_DIRECTION.PRE_TASK, labelKey: 'kanban.predecessor' },
	{ value: TASK_DEPENDENCY_DIRECTION.POST_TASK, labelKey: 'kanban.successor' },
];

/** 与后端 `project_management_member_role` 一致，排除 ProjectOwner / TaskCreator */
export const PROJECT_MEMBER_ROLE_OPTIONS = WORKFLOW_TRANSITION_MEMBER_OPTIONS.filter(
	(option) => option.value === 'ProjectAdmin' || option.value === 'TaskExecutor',
);

export const IMAGE_FILE_EXTENSIONS = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp', 'svg'];
