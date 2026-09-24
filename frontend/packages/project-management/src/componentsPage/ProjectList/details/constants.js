/**
 * 项目详情 URL：`?tab=` 任务/文档，`?view=` 任务内看板/表格/时间线/工时。
 * 默认任务 + 看板，不写入 query。
 */
export const PROJECT_DETAIL_TAB = {
	TASKS: 'tasks',
	DOCS: 'docs',
};

export const PROJECT_DETAIL_TABS = Object.values(PROJECT_DETAIL_TAB);

export const PROJECT_DETAIL_VIEW = {
	KANBAN: 'kanban',
	TABLE: 'table',
	TIMELINE: 'timeline',
	HOURS: 'hours',
};

export const PROJECT_DETAIL_VIEWS = Object.values(PROJECT_DETAIL_VIEW);

/** 表格/时间线/工时/文档自拉任务，避免只用看板已加载的列分页。 */
export const PROJECT_TASK_CATALOG_PAGE_SIZE = 50;
export const PROJECT_TASK_CATALOG_MAX_PAGES = 8;
