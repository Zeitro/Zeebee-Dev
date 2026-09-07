/**
 * 任务依赖：前置/后置依赖的增删改与搜索选择。
 * 变更走 PUT 全量更新：修改详情中的 dependencies 后一并提交。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { DeleteOutlined, LinkOutlined, PlusOutlined } from '@ant-design/icons';
import { Button, Dropdown, Input, Select, Spin, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { DeleteBtn } from '@repo/shared-methods';
import { displayValue } from '@repo/shared-methods/utils';
import {
	TASK_DEPENDENCY_DIRECTION,
	TASK_DEPENDENCY_RULE,
} from '@repo/project-management/enum/task';
import { DEPENDENCY_DIRECTION_OPTIONS, DEPENDENCY_RULE_OPTIONS } from '../constants';
import { splitTaskDependencies } from '../helpers';
import { TaskDetailPanel, TaskDetailPanelAction } from './TaskDetailPanel';

/**
 * @param {{
 *   task: object,
 *   projectId?: number|string,
 *   tasks?: Array,
 *   onTaskClick: (task: object) => void,
 *   fetchCandidates: (params: { project_id: number|string, task_id: number|string, keyword?: string }) => Promise<Array>,
 *   onSaveDependencies: (dependencies: Array) => Promise<object|null>,
 *   initialAddingMode?: typeof TASK_DEPENDENCY_DIRECTION.PRE_TASK | typeof TASK_DEPENDENCY_DIRECTION.POST_TASK | null,
 * }} props
 */
export function TaskDependencies({
	task,
	projectId,
	tasks = [],
	onTaskClick,
	fetchCandidates,
	onSaveDependencies,
	initialAddingMode = null,
	messageApi,
}) {
	const t = useTranslations();
	const [addingMode, setAddingMode] = useState(initialAddingMode);
	const [search, setSearch] = useState('');
	const [candidateTasks, setCandidateTasks] = useState([]);
	const [candidatesLoading, setCandidatesLoading] = useState(false);
	const [saving, setSaving] = useState(false);
	const [removingDepKey, setRemovingDepKey] = useState(null);

	const { predecessors, successors } = useMemo(
		() => splitTaskDependencies(task?.dependencies),
		[task?.dependencies],
	);

	const dependencyDirectionOptions = useMemo(
		() =>
			DEPENDENCY_DIRECTION_OPTIONS.map((item) => ({
				key: item.value,
				label: t(item.labelKey),
				icon: <LinkOutlined />,
			})),
		[t],
	);

	const dependencyDirectionLabelKey = useMemo(() => {
		const option = DEPENDENCY_DIRECTION_OPTIONS.find((item) => item.value === addingMode);
		return option?.labelKey ?? 'kanban.predecessor';
	}, [addingMode]);

	const dependencyRuleOptions = useMemo(
		() =>
			DEPENDENCY_RULE_OPTIONS.map((item) => ({
				value: item.value,
				label: t(item.labelKey),
			})),
		[t],
	);

	const relatedTaskIds = useMemo(() => {
		const ids = new Set();
		predecessors.forEach((dep) => {
			if (dep.dependency_task_id != null) ids.add(dep.dependency_task_id);
		});
		successors.forEach((dep) => {
			if (dep.dependency_task_id != null) ids.add(dep.dependency_task_id);
		});
		return ids;
	}, [predecessors, successors]);

	useEffect(() => {
		if (!addingMode || !projectId || !task?.id) return undefined;

		let cancelled = false;
		const timer = setTimeout(async () => {
			setCandidatesLoading(true);
			try {
				const list = await fetchCandidates({
					project_id: projectId,
					task_id: task.id,
					keyword: search.trim(),
				});
				if (!cancelled) setCandidateTasks(list);
			} catch {
				if (!cancelled) setCandidateTasks([]);
			} finally {
				if (!cancelled) setCandidatesLoading(false);
			}
		}, 300);

		return () => {
			cancelled = true;
			clearTimeout(timer);
		};
	}, [addingMode, fetchCandidates, projectId, search, task?.id]);

	const visibleCandidates = useMemo(() => {
		return candidateTasks.filter((item) => !relatedTaskIds.has(item.id));
	}, [candidateTasks, relatedTaskIds]);

	const findRelatedTask = (dependencyTaskId) =>
		tasks.find((item) => item.id === dependencyTaskId);

	/** 依赖项点击跳转：不依赖看板 tasks，有 id 即可导航 */
	const handleDependencyClick = (dep) => {
		const dependencyTaskId = Number(dep.dependency_task_id);
		if (!dependencyTaskId) return;
		const relatedTask = findRelatedTask(dependencyTaskId);
		onTaskClick({
			...(dep.dependency_task || {}),
			...(relatedTask || {}),
			id: dependencyTaskId,
			title: dep?.dependency_task_title,
		});
	};

	const persistDependencies = async (nextDependencies, successKey) => {
		setSaving(true);
		try {
			await onSaveDependencies(nextDependencies);
			messageApi.success(t(successKey));
		} catch (error) {
			messageApi.error(error?.message || t('kanban.dependencyAddFailed'));
			throw error;
		} finally {
			setSaving(false);
		}
	};

	const handleAdd = async (targetTask) => {
		if (relatedTaskIds.has(targetTask.id)) {
			messageApi.warning(t('kanban.dependencyExists'));
			return;
		}
		if (
			addingMode !== TASK_DEPENDENCY_DIRECTION.PRE_TASK &&
			addingMode !== TASK_DEPENDENCY_DIRECTION.POST_TASK
		) {
			return;
		}
		const nextDependencies = [
			...(task.dependencies || []),
			{
				id: 0,
				task_id: task.id,
				dependency_task_id: targetTask.id,
				dependency_type: addingMode,
				dependency_rule: TASK_DEPENDENCY_RULE.FINISH_TO_START,
			},
		];
		try {
			await persistDependencies(nextDependencies, 'kanban.dependencyAdded');
			setSearch('');
			setAddingMode(null);
			setCandidateTasks([]);
		} catch {
			/* message 已在 persistDependencies 展示 */
		}
	};

	const isSameDependency = (a, b) =>
		a === b ||
		(a?.id != null && b?.id != null && a.id !== 0 && a.id === b.id) ||
		(a?.dependency_task_id === b?.dependency_task_id &&
			a?.dependency_type === b?.dependency_type);

	const getDepKey = (dep) => dep.id ?? `${dep.dependency_type}-${dep.dependency_task_id}`;

	const handleUpdateRule = async (dep, rule) => {
		const nextDependencies = (task.dependencies || []).map((item) =>
			isSameDependency(item, dep) ? { ...item, dependency_rule: rule } : item,
		);
		try {
			await persistDependencies(nextDependencies, 'kanban.saveSuccess');
		} catch {
			messageApi.error(t('kanban.dependencyUpdateFailed'));
		}
	};

	const handleRemove = async (dep) => {
		const depKey = getDepKey(dep);
		setRemovingDepKey(depKey);
		const nextDependencies = (task.dependencies || []).filter(
			(item) => !isSameDependency(item, dep),
		);
		try {
			await persistDependencies(nextDependencies, 'kanban.dependencyRemoved');
		} catch {
			messageApi.error(t('kanban.dependencyRemoveFailed'));
		} finally {
			setRemovingDepKey(null);
		}
	};

	const renderDeleteButton = (dep) => (
		<DeleteBtn
			type="text"
			size="small"
			danger
			icon={<DeleteOutlined />}
			title={t('kanban.dependencyDeleteConfirm')}
			okText={t('common.confirm')}
			cancelText={t('common.cancel')}
			onOk={() => handleRemove(dep)}
			disabled={saving}
			confirmLoading={removingDepKey === getDepKey(dep)}
			className="!px-0 !gap-0"
		>
			<></>
		</DeleteBtn>
	);

	const hasDependencies = predecessors.length > 0 || successors.length > 0;
	const showEmpty = !hasDependencies && !addingMode;

	return (
		<TaskDetailPanel
			title={t('kanban.dependencies')}
			isEmpty={showEmpty}
			emptyText={t('kanban.dependenciesEmpty')}
			action={
				<Dropdown
					menu={{
						items: dependencyDirectionOptions,
						onClick: ({ key }) => {
							setAddingMode(key);
							setSearch('');
							setCandidateTasks([]);
						},
					}}
					trigger={['click']}
					disabled={saving}
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
						{t(dependencyDirectionLabelKey)}
					</div>
					<Input
						size="small"
						prefix={<LinkOutlined className="text-base-disabled" />}
						placeholder={t('kanban.dependencySearchPlaceholder')}
						value={search}
						onChange={(event) => setSearch(event.target.value)}
						allowClear
						disabled={saving}
						className="mb-2"
					/>
					<div className="border border-base-border rounded-lg bg-base-text2 max-h-60 overflow-y-auto">
						<div className="px-3 py-2 text-xs text-base-label border-b border-base-border bg-white rounded-t-lg">
							{search.trim()
								? t('kanban.searchResults')
								: t('kanban.selectableTasks')}
						</div>
						{candidatesLoading ? (
							<div className="flex justify-center py-6">
								<Spin size="small" />
							</div>
						) : visibleCandidates.length === 0 ? (
							<div className="px-3 py-3 text-[13px] text-base-disabled">
								{t('kanban.noMatchingTasks')}
							</div>
						) : (
							visibleCandidates.map((item) => (
								<div
									key={item.id}
									className="flex items-center gap-2.5 px-3 py-2 text-[13px] border-b border-base-border last:border-b-0 cursor-pointer hover:bg-base-active transition-colors"
									onClick={() => handleAdd(item)}
								>
									<span className="text-base-label text-xs">#{item.id}</span>
									<span className="flex-1 truncate">
										{displayValue(item.title)}
									</span>
									{item.phase ? (
										<Tag className="m-0 text-[11px] leading-[18px] px-1.5">
											{displayValue(item.phase)}
										</Tag>
									) : null}
								</div>
							))
						)}
					</div>
					<Button
						size="small"
						type="text"
						className="text-base-label mt-1 !px-0"
						disabled={saving}
						onClick={() => {
							setAddingMode(null);
							setSearch('');
							setCandidateTasks([]);
						}}
					>
						{t('common.cancel')}
					</Button>
				</div>
			) : null}

			{predecessors.length > 0 ? (
				<div className="px-4 py-3 border-b border-base-border/60 last:border-b-0">
					<div className="text-xs font-medium text-base-label mb-2">
						{t('kanban.predecessor')}
					</div>
					{predecessors.map((dep) => (
						<div
							key={dep.id ?? `${dep.dependency_type}-${dep.dependency_task_id}`}
							className="flex items-center gap-2 px-3 py-2 bg-base-text2 rounded-lg mb-2 last:mb-0"
						>
							<LinkOutlined className="text-base-disabled text-xs" />
							<Button
								type="link"
								className="flex-1 text-[13px] truncate h-auto !p-0 !justify-start"
								onClick={() => handleDependencyClick(dep)}
							>
								{displayValue(dep?.dependency_task_title)}
								<span className="text-[11px]/3 pt-1">
									#{dep.dependency_task_id}
								</span>
							</Button>
							<Select
								size="small"
								value={dep.dependency_rule}
								className="w-[110px]"
								options={dependencyRuleOptions}
								disabled={saving}
								onChange={(value) => handleUpdateRule(dep, value)}
							/>
							{renderDeleteButton(dep)}
						</div>
					))}
				</div>
			) : null}

			{successors.length > 0 ? (
				<div className="px-4 py-3">
					<div className="text-xs font-medium text-base-label mb-2">
						{t('kanban.successor')}
					</div>
					{successors.map((dep) => (
						<div
							key={dep.id ?? `${dep.dependency_type}-${dep.dependency_task_id}`}
							className="flex items-center gap-2 px-3 py-2 bg-base-text2 rounded-lg mb-2 last:mb-0"
						>
							<LinkOutlined className="text-base-disabled text-xs" />
							<Button
								type="link"
								className="flex-1 text-[13px] truncate h-auto !p-0 !justify-start"
								onClick={() => handleDependencyClick(dep)}
							>
								{displayValue(dep?.dependency_task_title)}
								<span className="text-[11px]/3 pt-1">
									#{dep.dependency_task_id}
								</span>
							</Button>
							<Select
								size="small"
								value={dep.dependency_rule}
								className="w-[110px]"
								options={dependencyRuleOptions}
								disabled={saving}
								onChange={(value) => handleUpdateRule(dep, value)}
							/>
							{renderDeleteButton(dep)}
						</div>
					))}
				</div>
			) : null}
		</TaskDetailPanel>
	);
}
