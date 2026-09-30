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
import { Tooltip } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue, formatDate } from '@repo/shared-methods/utils';
import { getPriorityI18nKey, getPriorityTagClass } from '../../../../enum/project';
import { toProjectPriorityLevel } from '../../../../enum/task';
import { getAssigneeIds, getMemberInitial, getMemberLabel } from './helpers';
import { TASK_PHASE } from '@repo/project-management/enum/task';

/** 标签统一带 1px 边框，保证有色边框的优先级标签与灰底标签等高；颜色类由各变体单独给，避免 Tailwind 同属性冲突 */
const PILL_BASE =
	'inline-flex items-center gap-0.5 px-2 py-px rounded border text-[11px] leading-[18px] whitespace-nowrap';
const PILL_CLASS = `${PILL_BASE} border-transparent bg-base-text2 text-base-gray`;
const PILL_OVERDUE_CLASS = `${PILL_BASE} border-transparent bg-[var(--error-light)] text-base-error`;

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
	const isCompleted = task?.phase === TASK_PHASE.FINISHED;
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
	const hasPills = Boolean(
		priorityKey || task.due_date || subtaskCount > 0 || commentCount > 0 || attachmentCount > 0,
	);

	return (
		<div
			onClick={onClick}
			className={`w-full min-h-[76px] rounded-lg bg-white px-3.5 py-3 text-left border border-base-border shadow-[0_1px_4px_rgba(0,0,0,0.06)] cursor-pointer flex flex-col transition-[box-shadow,border-color] duration-150 hover:border-blue-light hover:shadow-[0_2px_8px_rgba(0,0,0,0.1)] ${
				isCompleted ? 'opacity-70' : ''
			} ${className}`}
		>
			<div className="flex items-start gap-1.5">
				<p
					className={`flex-1 min-w-0 text-[13px] font-medium leading-5 line-clamp-2 break-words ${
						isCompleted ? 'line-through text-base-label' : 'text-black-light'
					}`}
				>
					{displayValue(task.title)}
				</p>
				{task.parent_id ? (
					<Tooltip
						title={t('kanban.parentLink', {
							title: parentTitle || `#${task.parent_id}`,
						})}
					>
						<SubnodeOutlined className="mt-1 text-xs text-base-label shrink-0" />
					</Tooltip>
				) : null}
			</div>
			{hasPills ? (
				<div className="flex flex-wrap items-center gap-1 mt-2">
					{priorityKey ? (
						<span className={`${PILL_BASE} ${getPriorityTagClass(priorityLevel)}`}>
							{t(priorityKey)}
						</span>
					) : null}
					{task.due_date ? (
						<span className={task?.overdue ? PILL_OVERDUE_CLASS : PILL_CLASS}>
							{displayValue(task.due_date, (value) => formatDate(value))}
						</span>
					) : null}
					{subtaskCount > 0 ? (
						<Tooltip title={t('kanban.subtaskCount', { count: subtaskCount })}>
							<span className={PILL_CLASS}>
								<PartitionOutlined className="text-[10px]" />
								{subtaskCount}
							</span>
						</Tooltip>
					) : null}
					{commentCount > 0 ? (
						<Tooltip title={t('kanban.commentCount', { count: commentCount })}>
							<span className={PILL_CLASS}>
								<MessageOutlined className="text-[10px]" />
								{commentCount}
							</span>
						</Tooltip>
					) : null}
					{attachmentCount > 0 ? (
						<Tooltip title={t('kanban.attachmentCount', { count: attachmentCount })}>
							<span className={PILL_CLASS}>
								<PaperClipOutlined className="text-[10px]" />
								{attachmentCount}
							</span>
						</Tooltip>
					) : null}
				</div>
			) : null}
			<div className="flex items-center justify-between mt-2">
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

	/** 外层已由 dnd-kit 挂 role / tabIndex；这里只补键盘打开，并保留 listeners 自带的 onKeyDown */
	const handleKeyDown = (event) => {
		listeners?.onKeyDown?.(event);
		if (event.target !== event.currentTarget) return;
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			onClick?.();
		}
	};

	return (
		<div
			ref={setNodeRef}
			style={{
				transform: CSS.Transform.toString(transform),
				transition,
				opacity: isDragging ? 0.4 : 1,
				zIndex: isDragging ? 999 : 'auto',
			}}
			className="mb-2.5 rounded-lg outline-none focus-visible:ring-2 focus-visible:ring-base-primary focus-visible:ring-offset-1"
			{...attributes}
			{...listeners}
			onKeyDown={handleKeyDown}
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
				className="!shadow-none !cursor-grabbing"
			/>
		</div>
	);
}
