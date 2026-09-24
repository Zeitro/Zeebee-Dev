/**
 * 项目甘特图：@rsagiev/gantt-task-react-19。
 * 条上显示工期；标题仅 hover；浅主色条对齐主题。
 */
'use client';

import { useCallback, useEffect, useMemo, useState, useSyncExternalStore } from 'react';
import { Segmented, Spin } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { Gantt, ViewMode } from '@rsagiev/gantt-task-react-19';
import '@rsagiev/gantt-task-react-19/dist/index.css';
import { NoData } from '@repo/shared-methods';
import { displayValue, formatDate } from '@repo/shared-methods/utils';
import { useProjectDetailQuery } from '../hooks/useProjectDetailQuery';
import { useProjectTaskCatalog } from '../hooks/useProjectTaskCatalog';
import { mapGanttTasks } from '../lib/mapGanttTasks';
import { useTodayTick } from '../hooks/useTodayTick';
import { ProjectGanttTaskListHeader, ProjectGanttTaskListTable } from './ProjectGanttTaskList';
import { ProjectViewFrame } from './ProjectViewFrame';
import styles from './ProjectTimelineView.module.scss';

const VIEW_MODE = {
	day: ViewMode.Day,
	week: ViewMode.Week,
	month: ViewMode.Month,
};

const COLUMN_WIDTH = { day: 36, week: 84, month: 120 };
/** 插件底部 HorizontalScroll 约占高度，须从 ganttHeight 扣除否则滚动条被裁切 */
const H_SCROLL_RESERVE = 18;
const HEADER_HEIGHT = 50;

/**
 * 订阅容器高度，供 Gantt ganttHeight（不含表头与横向滚动条）。
 * @param {HTMLElement|null} el
 * @returns {number}
 */
function useContainerHeight(el) {
	return useSyncExternalStore(
		(onStoreChange) => {
			if (!el || typeof ResizeObserver === 'undefined') return () => {};
			const observer = new ResizeObserver(onStoreChange);
			observer.observe(el);
			return () => observer.disconnect();
		},
		() => Math.max(0, el?.clientHeight || 0),
		() => 0,
	);
}

/**
 * 插件日视图表头写死「周几, 日」；改成只保留日期数字。
 * @param {HTMLElement|null} container
 * @param {boolean} enabled
 */
function useDayHeaderDigits(container, enabled) {
	useEffect(() => {
		if (!enabled || !container) return undefined;

		const rewrite = () => {
			container.querySelectorAll('text').forEach((node) => {
				const text = (node.textContent || '').trim();
				if (!text.includes(',')) return;
				const day = text.split(',').pop()?.trim();
				if (day && /^\d{1,2}$/.test(day)) node.textContent = day;
			});
		};

		rewrite();
		const observer = new MutationObserver(rewrite);
		observer.observe(container, { childList: true, subtree: true });
		return () => observer.disconnect();
	}, [container, enabled]);
}

/**
 * @param {{ projectId?: number|string, workflow?: object|null }} props
 */
export function ProjectTimelineView({ projectId }) {
	const t = useTranslations();
	const locale = useLocale();
	const { setTaskQuery } = useProjectDetailQuery();
	const { tasks, loading } = useProjectTaskCatalog(projectId);
	const [scale, setScale] = useState('week');
	const [collapsedIds, setCollapsedIds] = useState(() => new Set());
	const [containerEl, setContainerEl] = useState(null);
	const containerHeight = useContainerHeight(containerEl);
	useDayHeaderDigits(containerEl, scale === 'day');
	useTodayTick(containerEl, scale, t('projectDetail.today'));

	const formatDuration = useCallback((days) => t('projectDetail.spanDays', { count: days }), [t]);

	const ganttTasks = useMemo(
		() => mapGanttTasks(tasks, { collapsedIds, formatDuration }),
		[collapsedIds, formatDuration, tasks],
	);

	const listLabels = useMemo(
		() => ({
			name: t('myTask.titleColumn'),
			duration: t('projectDetail.durationColumn'),
			start: t('kanban.startDate'),
			end: t('kanban.dueDate'),
		}),
		[t],
	);

	const handleExpanderClick = useCallback((task) => {
		setCollapsedIds((prev) => {
			const next = new Set(prev);
			if (task.hideChildren) next.add(task.id);
			else next.delete(task.id);
			return next;
		});
	}, []);

	const openTask = useCallback(
		(task) => {
			const id = Number(task?.id);
			if (!Number.isFinite(id)) return;
			setTaskQuery(id);
		},
		[setTaskQuery],
	);

	const TaskListHeader = useCallback(
		(props) => <ProjectGanttTaskListHeader {...props} labels={listLabels} />,
		[listLabels],
	);

	const TaskListTable = useCallback(
		(props) => (
			<ProjectGanttTaskListTable
				{...props}
				onTaskNameClick={openTask}
				formatDuration={formatDuration}
			/>
		),
		[formatDuration, openTask],
	);

	const TooltipContent = useCallback(
		({ task }) => (
			<div className="min-w-[160px] rounded-lg border border-base-border bg-white px-3 py-2.5 text-xs shadow-lg">
				<div className="mb-1.5 font-medium text-foreground">
					{task.taskTitle || task.name}
				</div>
				<div className="text-foreground">{formatDuration(task.durationDays || 0)}</div>
				<div className="mt-0.5 text-foreground">
					{displayValue(task.displayStart || task.start, formatDate)} –{' '}
					{displayValue(task.displayEnd || task.end, formatDate)}
				</div>
			</div>
		),
		[formatDuration],
	);

	const viewMode = VIEW_MODE[scale];
	const ganttHeight = Math.max(120, containerHeight - HEADER_HEIGHT - H_SCROLL_RESERVE);

	const toolbar = (
		<>
			<span className="text-sm font-medium text-foreground">
				{t('projectDetail.viewTimeline')}
			</span>
			<Segmented
				value={scale}
				onChange={setScale}
				options={[
					{ label: t('projectDetail.scaleDay'), value: 'day' },
					{ label: t('projectDetail.scaleWeek'), value: 'week' },
					{ label: t('projectDetail.scaleMonth'), value: 'month' },
				]}
			/>
		</>
	);

	return (
		<ProjectViewFrame toolbar={toolbar}>
			{loading ? (
				<div className="flex flex-1 items-center justify-center">
					<Spin />
				</div>
			) : ganttTasks.length ? (
				<div ref={setContainerEl} className={styles.ganttHost}>
					{containerHeight > 0 ? (
						<Gantt
							tasks={ganttTasks}
							viewMode={viewMode}
							locale={locale === 'zh' ? 'zh-CN' : 'en-US'}
							preStepsCount={1}
							columnWidth={COLUMN_WIDTH[scale]}
							listCellWidth="120px"
							rowHeight={44}
							headerHeight={HEADER_HEIGHT}
							ganttHeight={ganttHeight}
							barCornerRadius={4}
							barFill={56}
							handleWidth={0}
							fontSize="12"
							barBackgroundColor="var(--base-primary)"
							barBackgroundSelectedColor="var(--blue-dark)"
							barProgressColor="var(--base-primary)"
							barProgressSelectedColor="var(--blue-dark)"
							projectBackgroundColor="var(--base-primary)"
							projectBackgroundSelectedColor="var(--blue-dark)"
							projectProgressColor="var(--base-primary)"
							projectProgressSelectedColor="var(--blue-dark)"
							todayColor="transparent"
							TaskListHeader={TaskListHeader}
							TaskListTable={TaskListTable}
							TooltipContent={TooltipContent}
							onExpanderClick={handleExpanderClick}
							onClick={openTask}
							onDoubleClick={openTask}
						/>
					) : null}
				</div>
			) : (
				<div className="flex flex-1 items-center justify-center">
					<NoData text={t('projectDetail.timelineEmpty')} />
				</div>
			)}
		</ProjectViewFrame>
	);
}
