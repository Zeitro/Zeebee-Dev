/**
 * 任务页内容区：看板沿用已接好的列拖拽；表格 / 时间线 / 工时各自拉数、各自打开任务。
 */
'use client';

import { EmptyWorkflow } from '../../kanban/components/EmptyWorkflow';
import { ProjectKanbanBoard } from '../../kanban';
import { PROJECT_DETAIL_VIEW } from '../constants';
import { ProjectHoursView } from './ProjectHoursView';
import { ProjectTableView } from './ProjectTableView';
import { ProjectTimelineView } from './ProjectTimelineView';

/**
 * @param {{
 *   view: string,
 *   projectId?: number|string,
 *   project?: object|null,
 *   workflow?: object|null,
 *   tasks?: Array,
 *   stageTaskMeta?: object,
 *   members?: Array,
 *   loadMoreStageTasks?: (stageId: number|string) => void,
 *   transitionTask?: (task: object, stateId: number) => Promise<void>|void,
 *   catalogRevision?: number,
 *   onCreateTask?: () => void,
 *   onTaskClick?: (task: object) => void,
 * }} props
 */
export function ProjectTaskViews({
	view,
	projectId,
	project,
	workflow,
	tasks,
	stageTaskMeta,
	members,
	loadMoreStageTasks,
	transitionTask,
	catalogRevision = 0,
	onCreateTask,
	onTaskClick,
}) {
	const catalogKey = `${projectId}-${catalogRevision}`;
	if (view === PROJECT_DETAIL_VIEW.TABLE) {
		return <ProjectTableView key={catalogKey} projectId={projectId} workflow={workflow} />;
	}
	if (view === PROJECT_DETAIL_VIEW.TIMELINE) {
		return <ProjectTimelineView key={catalogKey} projectId={projectId} workflow={workflow} />;
	}
	if (view === PROJECT_DETAIL_VIEW.HOURS) {
		return <ProjectHoursView key={catalogKey} projectId={projectId} />;
	}

	if (!workflow?.states?.length) {
		return <EmptyWorkflow project={project} />;
	}

	return (
		<ProjectKanbanBoard
			workflow={workflow}
			tasks={tasks}
			stageTaskMeta={stageTaskMeta}
			members={members}
			loadMoreStageTasks={loadMoreStageTasks}
			transitionTask={transitionTask}
			onCreateTask={onCreateTask}
			onTaskClick={onTaskClick}
		/>
	);
}
