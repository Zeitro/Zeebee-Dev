/**
 * 项目工时甘特：仅按执行者分行统计工时；缩放日 / 周 / 月 / 年。
 * 成员与任务在本组件内请求；空工时成员仍占行。
 */
'use client';

import { useCallback, useEffect, useMemo, useState, useSyncExternalStore } from 'react';
import { Segmented, Spin } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { Gantt, ViewMode } from '@rsagiev/gantt-task-react-19';
import '@rsagiev/gantt-task-react-19/dist/index.css';
import { NoData } from '@repo/shared-methods';
import { createRequest, displayValue, formatDate } from '@repo/shared-methods/utils';
import {
	PROJECT_MEMBER_LIST_PAGE_SIZE,
	createProjectMemberApi,
} from '@repo/project-management/api/projectMember';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { unwrapList } from '../../kanban/components/helpers';
import { useProjectTaskCatalog } from '../hooks/useProjectTaskCatalog';
import { useTodayTick } from '../hooks/useTodayTick';
import { groupTasksByAssignee } from '../lib/taskSchedule';
import { mapHoursGanttTasks } from '../lib/mapHoursGanttTasks';
import { ProjectHoursGanttListHeader, ProjectHoursGanttListTable } from './ProjectHoursGanttList';
import { ProjectViewFrame } from './ProjectViewFrame';
import styles from './ProjectTimelineView.module.scss';

const VIEW_MODE = {
	day: ViewMode.Day,
	week: ViewMode.Week,
	month: ViewMode.Month,
	year: ViewMode.Year,
};

const COLUMN_WIDTH = { day: 36, week: 84, month: 120, year: 72 };
const H_SCROLL_RESERVE = 18;
const HEADER_HEIGHT = 50;

/**
 * @param {Array} tasks
 * @param {Array} members
 * @returns {Array}
 */
function mergeAssigneeRows(tasks, members) {
	const grouped = groupTasksByAssignee(tasks);
	const rows = new Map(
		grouped.filter((row) => row.id !== 'unassigned').map((row) => [row.id, row]),
	);
	members.forEach((member) => {
		const id = member.organization_user_id ?? member.user_id;
		if (id == null) return;
		const current = rows.get(id);
		if (current) {
			if (!current.name) current.name = member.name || member.full_name;
			return;
		}
		rows.set(id, {
			id,
			name: member.name || member.full_name || '',
			tasks: [],
			hours: 0,
		});
	});
	const result = [...rows.values()];
	const unassigned = grouped.find((row) => row.id === 'unassigned');
	if (unassigned) result.push(unassigned);
	return result;
}

/**
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
 * @param {{ projectId?: number|string }} props
 */
export function ProjectHoursView({ projectId }) {
	const t = useTranslations();
	const locale = useLocale();
	const { tasks, loading: taskLoading } = useProjectTaskCatalog(projectId);
	const [scale, setScale] = useState('week');
	const [members, setMembers] = useState([]);
	const [membersReady, setMembersReady] = useState(false);
	const [containerEl, setContainerEl] = useState(null);
	const containerHeight = useContainerHeight(containerEl);
	useDayHeaderDigits(containerEl, scale === 'day');
	useTodayTick(containerEl, scale, t('projectDetail.today'));

	const { request } = useMemo(() => createRequest(), []);
	const memberApi = useMemo(() => createProjectMemberApi(request), [request]);

	useEffect(() => {
		if (!projectId) return undefined;
		let cancelled = false;
		memberApi
			.page({
				page: 1,
				limit: PROJECT_MEMBER_LIST_PAGE_SIZE,
				filter: { project_id: Number(projectId), keyword: '', role: '' },
			})
			.then((res) => {
				if (cancelled) return;
				setMembers(unwrapList(res).map(mapProjectMember).filter(Boolean));
				setMembersReady(true);
			})
			.catch(() => {
				if (cancelled) return;
				setMembers([]);
				setMembersReady(true);
			});
		return () => {
			cancelled = true;
		};
	}, [memberApi, projectId]);

	const formatHours = useCallback(
		(hours) => t('projectDetail.hoursValue', { hours: hours || 0 }),
		[t],
	);

	const assigneeRows = useMemo(() => mergeAssigneeRows(tasks, members), [members, tasks]);

	const ganttTasks = useMemo(
		() =>
			mapHoursGanttTasks(assigneeRows, {
				formatHours,
				unassignedLabel: t('projectDetail.unassigned'),
			}),
		[assigneeRows, formatHours, t],
	);

	const listLabels = useMemo(
		() => ({
			assignee: t('kanban.assignee'),
			hours: t('projectDetail.hoursColumn'),
		}),
		[t],
	);

	const TaskListHeader = useCallback(
		(props) => <ProjectHoursGanttListHeader {...props} labels={listLabels} />,
		[listLabels],
	);

	const TaskListTable = useCallback(
		(props) => <ProjectHoursGanttListTable {...props} formatHours={formatHours} />,
		[formatHours],
	);

	const TooltipContent = useCallback(
		({ task }) => (
			<div className="min-w-[160px] rounded-lg border border-base-border bg-white px-3 py-2.5 text-xs shadow-lg">
				<div className="mb-1.5 font-medium text-foreground">
					{task.taskTitle || task.name}
				</div>
				<div className="text-foreground">{formatHours(task.hours || 0)}</div>
				<div className="mt-0.5 text-foreground">
					{displayValue(task.displayStart || task.start, formatDate)} –{' '}
					{displayValue(task.displayEnd || task.end, formatDate)}
				</div>
			</div>
		),
		[formatHours],
	);

	const viewMode = VIEW_MODE[scale];
	const ganttHeight = Math.max(120, containerHeight - HEADER_HEIGHT - H_SCROLL_RESERVE);

	const toolbar = (
		<>
			<span className="text-sm font-medium text-foreground">
				{t('projectDetail.viewHours')}
			</span>
			<Segmented
				value={scale}
				onChange={setScale}
				options={[
					{ label: t('projectDetail.scaleDay'), value: 'day' },
					{ label: t('projectDetail.scaleWeek'), value: 'week' },
					{ label: t('projectDetail.scaleMonth'), value: 'month' },
					{ label: t('projectDetail.scaleYear'), value: 'year' },
				]}
			/>
		</>
	);

	return (
		<ProjectViewFrame toolbar={toolbar}>
			{taskLoading || !membersReady ? (
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
						/>
					) : null}
				</div>
			) : (
				<div className="flex flex-1 items-center justify-center">
					<NoData text={t('projectDetail.hoursEmpty')} />
				</div>
			)}
		</ProjectViewFrame>
	);
}
