/**
 * 任务依赖：前置/后置依赖的增删改与搜索选择。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { DeleteOutlined, LinkOutlined, PlusOutlined } from '@ant-design/icons';
import { Avatar, Button, Dropdown, Input, Select, Tag, Tooltip, message } from 'antd';
import { useTranslations } from 'next-intl';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { displayValue } from '@repo/shared-methods/utils';
import { DEPENDENCY_TYPE_OPTIONS } from '../constants';
import { getMemberInitial, getMemberLabel } from '../helpers';
import { TaskDetailPanel, TaskDetailPanelAction } from './TaskDetailPanel';

/**
 * @param {{
 *   task: object,
 *   tasks?: Array,
 *   workflow?: object|null,
 *   members?: Array,
 *   onTaskClick: (task: object) => void,
 *   fetchDependencies: (taskId: number) => Promise<Array>,
 *   fetchDependents: (taskId: number) => Promise<Array>,
 *   onAddDependency: (taskId: number, dependsOnTaskId: number) => Promise<void>,
 *   onAddDependent: (taskId: number, dependentTaskId: number) => Promise<void>,
 *   onUpdateDependency: (taskId: number, dependencyId: number, type: string) => Promise<void>,
 *   onUpdateDependent: (taskId: number, dependencyId: number, type: string) => Promise<void>,
 *   onDeleteDependency: (taskId: number, dependencyId: number) => Promise<void>,
 *   onDeleteDependent: (taskId: number, dependencyId: number) => Promise<void>,
 *   initialAddingMode?: 'predecessor' | 'successor' | null,
 * }} props
 */
export function TaskDependencies({
	task,
	tasks = [],
	workflow,
	members = [],
	onTaskClick,
	fetchDependencies,
	fetchDependents,
	onAddDependency,
	onAddDependent,
	onUpdateDependency,
	onUpdateDependent,
	onDeleteDependency,
	onDeleteDependent,
	initialAddingMode = null,
}) {
	const t = useTranslations();
	const [dependencies, setDependencies] = useState([]);
	const [dependents, setDependents] = useState([]);
	const [addingMode, setAddingMode] = useState(initialAddingMode);
	const [search, setSearch] = useState('');

	const dependencyTypeOptions = useMemo(
		() =>
			DEPENDENCY_TYPE_OPTIONS.map((item) => ({
				value: item.value,
				label: t(item.labelKey),
			})),
		[t],
	);

	const reload = useCallback(async () => {
		if (!task?.id) return;
		const [depRows, dependentRows] = await Promise.all([
			fetchDependencies(task.id),
			fetchDependents(task.id),
		]);
		setDependencies(depRows);
		setDependents(dependentRows);
	}, [fetchDependencies, fetchDependents, task]);

	// 父级 key 随 task / initialAddingMode 变化 remount；此处仅异步拉依赖列表
	useDeferredSyncEffect({
		deferred: () => {
			reload();
		},
		deps: [reload, task?.id],
		isUseMicrotask: true,
	});

	const candidateTasks = useMemo(() => {
		const keyword = search.trim().toLowerCase();
		const isPredecessor = addingMode === 'predecessor';
		return tasks.filter((item) => {
			if (item.id === task.id) return false;
			if (isPredecessor) {
				if (dependencies.some((dep) => dep.depends_on_task_id === item.id)) return false;
			} else if (dependents.some((dep) => dep.task_id === item.id)) {
				return false;
			}
			if (!keyword) return true;
			return item.title.toLowerCase().includes(keyword) || String(item.id) === keyword;
		});
	}, [addingMode, dependencies, dependents, search, task.id, tasks]);

	const handleAdd = async (targetTask) => {
		const alreadyRelated =
			dependencies.some((dep) => dep.depends_on_task_id === targetTask.id) ||
			dependents.some((dep) => dep.task_id === targetTask.id);
		if (alreadyRelated) {
			message.warning(t('kanban.dependencyExists'));
			return;
		}
		try {
			if (addingMode === 'predecessor') {
				await onAddDependency(task.id, targetTask.id);
			} else {
				await onAddDependent(task.id, targetTask.id);
			}
			setSearch('');
			setAddingMode(null);
			await reload();
		} catch (error) {
			message.error(error?.message || t('kanban.dependencyAddFailed'));
		}
	};

	const hasDependencies = dependencies.length > 0 || dependents.length > 0;
	const showEmpty = !hasDependencies && !addingMode;

	return (
		<TaskDetailPanel
			title={t('kanban.dependencies')}
			isEmpty={showEmpty}
			emptyText={t('kanban.dependenciesEmpty')}
			action={
				<Dropdown
					menu={{
						items: [
							{
								key: 'predecessor',
								label: t('kanban.predecessor'),
								icon: <LinkOutlined />,
							},
							{
								key: 'successor',
								label: t('kanban.successor'),
								icon: <LinkOutlined />,
							},
						],
						onClick: ({ key }) => {
							setAddingMode(key);
							setSearch('');
						},
					}}
					trigger={['click']}
				>
					<TaskDetailPanelAction
						title={t('kanban.addDependencyTask')}
						icon={<PlusOutlined />}
					/>
				</Dropdown>
			}
		>
			{addingMode ? (
				<div className="px-4 py-3 border-b border-base-border/60">
					<div className="text-xs text-base-label mb-2">
						{addingMode === 'predecessor'
							? t('kanban.predecessor')
							: t('kanban.successor')}
					</div>
					<Input
						size="small"
						prefix={<LinkOutlined className="text-base-disabled" />}
						placeholder={t('kanban.dependencySearchPlaceholder')}
						value={search}
						onChange={(event) => setSearch(event.target.value)}
						allowClear
						className="mb-2"
					/>
					<div className="border border-base-border rounded-lg bg-base-text2 max-h-60 overflow-y-auto">
						<div className="px-3 py-2 text-xs text-base-label border-b border-base-border bg-white rounded-t-lg">
							{search.trim()
								? t('kanban.searchResults')
								: t('kanban.selectableTasks')}
						</div>
						{candidateTasks.length === 0 ? (
							<div className="px-3 py-3 text-[13px] text-base-disabled">
								{t('kanban.noMatchingTasks')}
							</div>
						) : (
							candidateTasks.map((item) => {
								const stateObj = workflow?.states?.find(
									(state) => state.id === item.workflow_state_id,
								);
								const assignee = members.find(
									(member) => member.id === item.assignee_id,
								);
								return (
									<div
										key={item.id}
										className="flex items-center gap-2.5 px-3 py-2 text-[13px] border-b border-base-border last:border-b-0 cursor-pointer hover:bg-base-active transition-colors"
										onClick={() => handleAdd(item)}
									>
										<span className="text-base-label text-xs">#{item.id}</span>
										<span className="flex-1 truncate">
											{displayValue(item.title)}
										</span>
										{stateObj ? (
											<Tag className="m-0 text-[11px] leading-[18px] px-1.5">
												{stateObj.name}
											</Tag>
										) : null}
										{assignee ? (
											<Tooltip title={getMemberLabel(assignee)}>
												<Avatar
													size={20}
													style={{
														backgroundColor: 'var(--blue-medium)',
													}}
													className="text-[10px] shrink-0"
												>
													{getMemberInitial(assignee)}
												</Avatar>
											</Tooltip>
										) : null}
									</div>
								);
							})
						)}
					</div>
					<Button
						size="small"
						type="text"
						className="text-base-label mt-1 !px-0"
						onClick={() => {
							setAddingMode(null);
							setSearch('');
						}}
					>
						{t('common.cancel')}
					</Button>
				</div>
			) : null}

			{dependencies.length > 0 ? (
				<div className="px-4 py-3 border-b border-base-border/60 last:border-b-0">
					<div className="text-xs font-medium text-base-label mb-2">
						{t('kanban.predecessor')}
					</div>
					{dependencies.map((dep) => (
						<div
							key={dep.id}
							className="flex items-center gap-2 px-3 py-2 bg-base-text2 rounded-lg mb-2 last:mb-0"
						>
							<LinkOutlined className="text-base-disabled text-xs" />
							<Button
								type="link"
								className="flex-1 text-left text-[13px] truncate h-auto p-0 justify-start"
								onClick={() => {
									const depTask = tasks.find(
										(item) => item.id === dep.depends_on_task_id,
									);
									if (depTask) onTaskClick(depTask);
								}}
							>
								#{dep.depends_on_task_id} {displayValue(dep.depends_on_task?.title)}
							</Button>
							<Select
								size="small"
								value={dep.dependency_type}
								className="w-[110px]"
								options={dependencyTypeOptions}
								onChange={async (value) => {
									try {
										await onUpdateDependency(task.id, dep.id, value);
										await reload();
									} catch {
										message.error(t('kanban.dependencyUpdateFailed'));
									}
								}}
							/>
							<Button
								size="small"
								type="text"
								danger
								icon={<DeleteOutlined />}
								onClick={async () => {
									try {
										await onDeleteDependency(task.id, dep.id);
										await reload();
									} catch {
										message.error(t('kanban.dependencyRemoveFailed'));
									}
								}}
							/>
						</div>
					))}
				</div>
			) : null}

			{dependents.length > 0 ? (
				<div className="px-4 py-3">
					<div className="text-xs font-medium text-base-label mb-2">
						{t('kanban.successor')}
					</div>
					{dependents.map((dep) => (
						<div
							key={dep.id}
							className="flex items-center gap-2 px-3 py-2 bg-base-text2 rounded-lg mb-2 last:mb-0"
						>
							<LinkOutlined className="text-base-disabled text-xs" />
							<Select
								size="small"
								value={dep.dependency_type}
								className="w-[110px]"
								options={dependencyTypeOptions}
								onChange={async (value) => {
									try {
										await onUpdateDependent(task.id, dep.id, value);
										await reload();
									} catch {
										message.error(t('kanban.dependencyUpdateFailed'));
									}
								}}
							/>
							<Button
								type="link"
								className="flex-1 text-left text-[13px] truncate h-auto p-0 justify-start"
								onClick={() => {
									const depTask = tasks.find((item) => item.id === dep.task_id);
									if (depTask) onTaskClick(depTask);
								}}
							>
								#{dep.task_id} {displayValue(dep.task?.title)}
							</Button>
							<Button
								size="small"
								type="text"
								danger
								icon={<DeleteOutlined />}
								onClick={async () => {
									try {
										await onDeleteDependent(task.id, dep.id);
										await reload();
									} catch {
										message.error(t('kanban.dependencyRemoveFailed'));
									}
								}}
							/>
						</div>
					))}
				</div>
			) : null}
		</TaskDetailPanel>
	);
}
