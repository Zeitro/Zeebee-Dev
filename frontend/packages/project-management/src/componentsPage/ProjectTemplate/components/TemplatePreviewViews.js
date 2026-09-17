/**
 * 模板预览的看板 / 工作流只读视图。纯展示，不改 preview、不发请求。
 */
'use client';

import { HolderOutlined, PlusOutlined } from '@ant-design/icons';
import { Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { displayValue } from '@repo/shared-methods/utils';
import {
	WORKFLOW_STATE_TYPE_OPTIONS,
	getWorkflowPhaseI18nKey,
	getWorkflowStateColor,
	getWorkflowStateColorHex,
	isInitialState,
} from '../../../enum/workflow';

/**
 * @param {object} state
 * @param {(key: string) => string} t
 * @returns {string}
 */
function stateTypeLabel(state, t) {
	const phaseKey = getWorkflowPhaseI18nKey(state?.phase);
	if (phaseKey) return t(phaseKey);
	const typeKey = WORKFLOW_STATE_TYPE_OPTIONS.find(
		(item) => item.value === state?.type,
	)?.labelKey;
	return typeKey ? t(typeKey) : displayValue(state?.type);
}

/**
 * @param {object} state
 * @param {object[]} tasks
 * @returns {object[]}
 */
function tasksInState(state, tasks) {
	return (tasks || []).filter((task) => {
		if (task.parent_id) return false;
		if (task.state_name && task.state_name === state.name) return true;
		if (task.workflow_state_id != null && task.workflow_state_id === state.id) return true;
		return false;
	});
}

/**
 * @param {{ preview: object }} props
 */
export function TemplateBoardView({ preview }) {
	const t = useTranslations();
	const states = preview?.states || [];
	if (!states.length) {
		return <NoData text={t('projectTemplate.emptyWorkflow')} className="min-h-[240px]" />;
	}

	return (
		<div className="flex gap-3 overflow-x-auto pb-2 min-h-[280px]">
			{states.map((state) => {
				const tasks = tasksInState(state, preview.tasks);
				return (
					<div
						key={state.id}
						className="flex min-w-[180px] flex-1 flex-col rounded-[10px] bg-gray-light p-3"
					>
						<div className="mb-3 flex items-center justify-between gap-2">
							<div className="flex min-w-0 items-center gap-2">
								<HolderOutlined className="text-base-disabled" />
								<span
									className="h-2.5 w-2.5 shrink-0 rounded-full"
									style={{ background: getWorkflowStateColor(state.color) }}
								/>
								<span className="truncate text-sm font-semibold">{state.name}</span>
								<Tag className="m-0 rounded-full">{tasks.length}</Tag>
							</div>
							{isInitialState(state) ? (
								<PlusOutlined className="text-base-disabled" />
							) : null}
						</div>
						<div className="min-h-[180px] flex-1">
							{tasks.map((task) => (
								<div
									key={task.id}
									className="mb-2 rounded-lg bg-white px-3 py-2.5 text-[13px] shadow-sm"
								>
									{displayValue(task.title)}
								</div>
							))}
						</div>
					</div>
				);
			})}
		</div>
	);
}

/**
 * @param {{ preview: object }} props
 */
export function TemplateWorkflowView({ preview }) {
	const t = useTranslations();
	const states = preview?.states || [];
	const transitions = preview?.transitions || [];
	if (!states.length) {
		return <NoData text={t('projectTemplate.emptyWorkflow')} className="min-h-[240px]" />;
	}

	return (
		<div>
			<div className="mb-3 font-semibold">{t('workflow.editor.stateNodesTitle')}</div>
			<div className="flex flex-wrap gap-4">
				{states.map((state) => {
					const color = getWorkflowStateColorHex(state.color);
					const outgoing = transitions.filter((item) => item.from_state_id === state.id);
					return (
						<div
							key={state.id}
							className="min-w-[168px] max-w-[240px] flex-1 rounded-lg border bg-white px-3.5 py-3"
							style={{ borderColor: color }}
						>
							<div className="flex items-center gap-2">
								<HolderOutlined className="text-base-disabled" />
								<span
									className="h-3 w-3 shrink-0 rounded-full"
									style={{ background: color }}
								/>
								<span className="font-medium">{state.name}</span>
							</div>
							<div className="mt-1.5 text-xs text-base-label">
								{state.is_initial ? (
									<Tag className="text-[10px]">
										{t('workflow.editor.tagInitial')}
									</Tag>
								) : null}
								{state.is_terminal ? (
									<Tag className="text-[10px]">
										{t('workflow.editor.tagTerminal')}
									</Tag>
								) : null}
								<span>{stateTypeLabel(state, t)}</span>
							</div>
							{outgoing.length ? (
								<div className="mt-2 border-t border-base-border pt-1.5">
									<div className="mb-1 text-[11px] text-base-label">
										{t('workflow.editor.transitionTo')}
									</div>
									{outgoing.map((trans) => {
										const to = states.find(
											(item) => item.id === trans.to_state_id,
										);
										return (
											<div
												key={`${trans.from_state_id}-${trans.to_state_id}`}
												className="mb-1"
											>
												<Tag
													className="m-0 text-[10px]"
													color={getWorkflowStateColorHex(to?.color)}
												>
													→ {to?.name || '?'} (
													{trans.name ||
														t('workflow.editor.transitionDefaultName')}
													)
												</Tag>
											</div>
										);
									})}
								</div>
							) : null}
						</div>
					);
				})}
			</div>
		</div>
	);
}
