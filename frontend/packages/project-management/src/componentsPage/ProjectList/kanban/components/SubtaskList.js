/**
 * 详情抽屉内子任务列表：展示直属子任务与完成数。
 */
'use client';

import { Button, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';
import { isTerminalState } from '../../../../enum/workflow';
import { getSubtasks } from './helpers';

/**
 * @param {{
 *   parentId: number,
 *   tasks?: Array,
 *   workflow?: object|null,
 *   onTaskClick: (task: object) => void,
 * }} props
 */
export function SubtaskList({ parentId, tasks = [], workflow, onTaskClick }) {
	const t = useTranslations();
	const subtasks = getSubtasks(tasks, parentId);
	if (!subtasks.length) return null;

	const completedCount = subtasks.filter((item) => {
		const state = workflow?.states?.find((row) => row.id === item.workflow_state_id);
		return isTerminalState(state);
	}).length;

	return (
		<div className="mb-4 border border-base-border rounded-lg overflow-hidden">
			<div className="flex items-center gap-2 px-3 py-2 bg-gray-light">
				<span className="font-semibold text-[13px] flex-1">{t('kanban.subtasks')}</span>
				<span className="text-xs text-base-label">
					{completedCount}/{subtasks.length}
				</span>
			</div>
			<div className="py-1">
				{subtasks.map((item) => {
					const state = workflow?.states?.find(
						(row) => row.id === item.workflow_state_id,
					);
					return (
						<Button
							key={item.id}
							type="text"
							block
							className="flex items-center gap-2 px-3 py-1.5 h-auto text-left text-[13px] hover:bg-base-text2 rounded-none"
							onClick={() => onTaskClick(item)}
						>
							<span className="flex-1 truncate">
								{displayValue(item.title)}
								<span className="ml-1 text-[11px] text-base-disabled">
									#{item.id}
								</span>
							</span>
							{state?.name ? <Tag className="m-0">{state.name}</Tag> : null}
						</Button>
					);
				})}
			</div>
		</div>
	);
}
