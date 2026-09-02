/**
 * 看板任务卡：可拖拽；展示标题、子任务数、优先级、截止日期、评论、执行者。
 */
'use client';

import { LinkOutlined, MessageOutlined, PartitionOutlined } from '@ant-design/icons';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useTranslations } from 'next-intl';
import { dayjs, displayValue, formatDate } from '@repo/shared-methods/utils';
import {
	PROJECT_PRIORITY,
	getPriorityI18nKey,
	getPriorityTagClass,
} from '../../../../enum/project';
import { toProjectPriorityLevel } from '../../../../enum/task';
import { isTerminalState } from '../../../../enum/workflow';
import { getAssigneeIds, getMemberInitial, getMemberLabel } from './helpers';

const PILL_CLASS =
	'inline-flex items-center gap-0.5 px-2 py-0.5 rounded bg-base-text2 text-[11px] text-base-gray leading-[18px]';

/**
 * @param {{
 *   task: object,
 *   members?: Array,
 *   parentTitle?: string,
 *   subtaskCount?: number,
 *   workflowState?: object,
 *   onClick: () => void,
 * }} props
 */
export function TaskCard({
	task,
	members = [],
	parentTitle,
	subtaskCount = 0,
	workflowState,
	onClick,
}) {
	const t = useTranslations();
	const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
		id: `task-${task.id}`,
		data: { task, stateId: task.workflow_state_id },
	});
	const isCompleted = isTerminalState(task.workflow_state || workflowState);
	const isOverdue =
		Boolean(task.due_date) &&
		dayjs(task.due_date).isBefore(dayjs(), 'day') &&
		task.status !== 'done' &&
		!isCompleted;
	const commentCount = task.comments?.length || 0;
	const assigneeIds = getAssigneeIds(task);
	const assigneeMembers = assigneeIds
		.map((id) => members.find((member) => member.id === id))
		.filter(Boolean);
	const priorityLevel = toProjectPriorityLevel(task.priority);
	const priorityKey = getPriorityI18nKey(priorityLevel);

	const isCriticalPriority = priorityLevel === PROJECT_PRIORITY.URGENT;

	return (
		<div
			ref={setNodeRef}
			style={{
				transform: CSS.Transform.toString(transform),
				transition,
				opacity: isDragging ? 0.4 : 1,
				zIndex: isDragging ? 999 : 'auto',
			}}
			{...attributes}
			{...listeners}
		>
			<div
				role="button"
				tabIndex={0}
				onClick={onClick}
				onKeyDown={(event) => {
					if (event.key === 'Enter') onClick();
				}}
				className={`w-full min-h-[76px] mb-2.5 rounded-lg bg-white px-3.5 py-3 text-left border border-base-border shadow-[0_1px_4px_rgba(0,0,0,0.06)] cursor-grab overflow-hidden flex flex-col`}
			>
				<div className="flex items-center gap-1.5 h-5">
					<p
						className={`flex-1 text-[13px] font-medium leading-5 truncate ${
							isCompleted ? 'line-through text-base-disabled' : 'text-black-light'
						}`}
					>
						{displayValue(task.title)}
						<span className="ml-1 text-[11px] font-normal text-base-disabled">
							#{task.id}
						</span>
					</p>
					{task.parent_id ? (
						<LinkOutlined
							className="text-xs text-base-label shrink-0"
							title={t('kanban.parentLink', {
								title: parentTitle || `#${task.parent_id}`,
							})}
						/>
					) : null}
				</div>
				<div className="flex items-center gap-1 mt-2 h-[22px] overflow-hidden">
					{subtaskCount > 0 ? (
						<span className={PILL_CLASS}>
							<PartitionOutlined className="text-[10px]" />
							{subtaskCount}
						</span>
					) : null}
					{priorityKey ? (
						<span
							className={
								isCriticalPriority
									? `${PILL_CLASS} !text-base-error`
									: `${PILL_CLASS} ${getPriorityTagClass(priorityLevel)}`
							}
						>
							{t(priorityKey)}
						</span>
					) : null}
					{task.due_date ? (
						<span
							className={`${PILL_CLASS} ${
								isOverdue ? 'text-base-error bg-[var(--error-light)]' : ''
							}`}
						>
							{displayValue(task.due_date, (value) => formatDate(value))}
						</span>
					) : null}
					{commentCount > 0 ? (
						<span className={PILL_CLASS}>
							<MessageOutlined className="text-[10px]" />
							{commentCount}
						</span>
					) : null}
					{assigneeMembers.length > 0 ? (
						<span className="inline-flex items-center ml-auto">
							{assigneeMembers.slice(0, 3).map((member) => (
								<span
									key={member.id}
									title={getMemberLabel(member)}
									className="w-5 h-5 -ml-1 rounded-full bg-blue-medium text-white text-[9px] font-semibold inline-flex items-center justify-center border-[1.5px] border-white"
								>
									{getMemberInitial(member)}
								</span>
							))}
							{assigneeMembers.length > 3 ? (
								<span className="ml-0.5 text-[10px] text-base-label">
									+{assigneeMembers.length - 3}
								</span>
							) : null}
						</span>
					) : null}
				</div>
			</div>
		</div>
	);
}
