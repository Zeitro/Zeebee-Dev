/**
 * 项目详情视图 key，与 URL `?view=` 一致。看板为默认，不写入 query。
 */
export const PROJECT_DETAIL_VIEW = {
	KANBAN: 'kanban',
	TABLE: 'table',
	TIMELINE: 'timeline',
	HOURS: 'hours',
};

export const PROJECT_DETAIL_VIEWS = Object.values(PROJECT_DETAIL_VIEW);
