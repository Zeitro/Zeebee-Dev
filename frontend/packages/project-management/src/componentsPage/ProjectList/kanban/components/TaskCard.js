/**
 * 看板任务卡：可拖拽；展示标题、子任务数、优先级、截止日期、评论、执行者。
 * DragOverlay 用 TaskCardDragPreview，避免重复挂 useSortable。
 */
'use client';

import {
	MessageOutlined,
	PaperClipOutlined,
	PartitionOutlined,
	SubnodeOutlined,
} from '@ant-design/icons';
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
 * 任务卡视觉内容（列内卡与拖拽预览共用）。
 * @param {{
 *   task: object,
 *   members?: Array,
 *   parentTitle?: string,
 *   subtaskCount?: number,
 *   workflowState?: object,
 *   onClick?: () => void,
 *   className?: string,
 * }} props
 */
function TaskCardBody({
	task,
	members = [],
	parentTitle,
	subtaskCount = 0,
	workflowState,
	onClick,
	className = '',
}) {
	const t = useTranslations();
	const isCompleted = isTerminalState(task.workflow_state || workflowState);
	const isOverdue =
		Boolean(task.due_date) &&
		dayjs(task.due_date).isBefore(dayjs(), 'day') &&
		task.status !== 'done' &&
		!isCompleted;
	const commentCount = task.comments?.length || 0;
	const attachmentCount =
		Number(task.attachment_count) ||
		(Array.isArray(task.attachments) ? task.attachments.length : 0);
	const assigneeIds = getAssigneeIds(task);
	const assigneeMembers = assigneeIds
		.map((id) => members.find((member) => member.id === id))
		.filter(Boolean);
	const priorityLevel = toProjectPriorityLevel(task.priority);
	const priorityKey = getPriorityI18nKey(priorityLevel);
	const isCriticalPriority = priorityLevel === PROJECT_PRIORITY.URGENT;

	return (
		<div
			role={onClick ? 'button' : undefined}
			tabIndex={onClick ? 0 : undefined}
			onClick={onClick}
			onKeyDown={
				onClick
					? (event) => {
							if (event.key === 'Enter') onClick();
						}
					: undefined
			}
			className={`w-full min-h-[76px] mb-2.5 rounded-lg bg-white px-3.5 py-3 text-left border border-base-border shadow-[0_1px_4px_rgba(0,0,0,0.06)] cursor-grab flex flex-col transition-shadow duration-150 hover:shadow-[0_0_12px_rgba(0,0,0,0.14),0_2px_8px_rgba(0,0,0,0.08)] ${className}`}
		>
			<div className="flex items-center gap-1.5 h-5">
				<p
					className={`flex-1 text-[13px] font-medium leading-5 truncate ${
						isCompleted ? 'line-through text-base-disabled' : 'text-black-light'
					}`}
				>
					{displayValue(task.title)}
				</p>
				{task.parent_id ? (
					<SubnodeOutlined
						className="text-xs text-base-label shrink-0"
						title={t('kanban.parentLink', {
							title: parentTitle || `#${task.parent_id}`,
						})}
					/>
				) : null}
			</div>
			<div className="flex items-center gap-1 my-2 h-[22px] overflow-hidden">
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
				{subtaskCount > 0 ? (
					<span
						className={PILL_CLASS}
						title={t('kanban.subtaskCount', { count: subtaskCount })}
					>
						<PartitionOutlined className="text-[10px]" />
						{subtaskCount}
					</span>
				) : null}
				{commentCount > 0 ? (
					<span className={PILL_CLASS}>
						<MessageOutlined className="text-[10px]" />
						{commentCount}
					</span>
				) : null}
				{attachmentCount > 0 ? (
					<span
						className={PILL_CLASS}
						title={t('kanban.attachmentCount', { count: attachmentCount })}
					>
						<PaperClipOutlined className="text-[10px]" />
						{attachmentCount}
					</span>
				) : null}
			</div>
			<div className="flex items-center justify-between">
				<span className="flex items-center text-xs text-base-gray">#{task.id}</span>
				{assigneeMembers.length > 0 ? (
					<span className="inline-flex items-center ml-auto">
						{assigneeMembers.slice(0, 3).map((member) => (
							<span
								key={member.id}
								title={getMemberLabel(member)}
								className="w-5 h-5 -ml-1 pt-[2px] rounded-full bg-blue-medium text-white text-[9px] font-semibold inline-flex items-center justify-center border-[1.5px] border-white"
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
	);
}

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
	const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
		id: `task-${task.id}`,
		data: { task, stateId: task.workflow_state_id },
	});

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
			<TaskCardBody
				task={task}
				members={members}
				parentTitle={parentTitle}
				subtaskCount={subtaskCount}
				workflowState={workflowState}
				onClick={onClick}
			/>
		</div>
	);
}

/**
 * DragOverlay 完整预览：与列内卡同结构，不挂 sortable。
 * @param {{
 *   task: object,
 *   members?: Array,
 *   parentTitle?: string,
 *   subtaskCount?: number,
 *   workflowState?: object,
 * }} props
 */
export function TaskCardDragPreview({
	task,
	members = [],
	parentTitle,
	subtaskCount = 0,
	workflowState,
}) {
	return (
		<div className="w-[230px] pointer-events-none opacity-90 shadow-md rounded-lg">
			<TaskCardBody
				task={task}
				members={members}
				parentTitle={parentTitle}
				subtaskCount={subtaskCount}
				workflowState={workflowState}
				className="!mb-0 !shadow-none cursor-grabbing"
			/>
		</div>
	);
}
