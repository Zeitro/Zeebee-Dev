/**
 * 创建项目的工作流模板横滑选择。同一弹框内始终展示，不切步骤、不隐藏。
 * 受 Form.Item 的 value / onChange 控制，不发请求、不改 rowData。
 */
'use client';

import { CheckCircleFilled } from '@ant-design/icons';
import { Spin } from 'antd';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';
import { getWorkflowStateColor } from '../../../enum/workflow';

/**
 * @param {object|null|undefined} workflow
 * @returns {object[]}
 */
function sortedStages(workflow) {
	const stages = Array.isArray(workflow?.stages) ? [...workflow.stages] : [];
	return stages.sort((a, b) => (Number(a?.sort_order) || 0) - (Number(b?.sort_order) || 0));
}

/**
 * @param {{
 *   value?: number,
 *   onChange?: (id: number) => void,
 *   workflows?: object[],
 *   loading?: boolean,
 *   disabled?: boolean,
 *   id?: string
 * }} props
 */
export function WorkflowTemplatePicker({
	value,
	onChange,
	workflows = [],
	loading = false,
	disabled = false,
	id,
}) {
	const t = useTranslations();

	const handleSelect = (workflowId) => {
		if (disabled) return;
		onChange?.(workflowId);
	};

	if (loading) {
		return (
			<div className="flex h-24 items-center justify-center">
				<Spin />
			</div>
		);
	}

	if (!workflows.length) {
		return <div className="py-2 text-sm text-base-gray">{t('projectList.templateEmpty')}</div>;
	}

	return (
		<div id={id} className="flex gap-3 overflow-x-auto px-0.5 pt-1 pb-2 [scrollbar-width:thin]">
			{workflows.map((workflow) => {
				const selected = value === workflow.id;
				const stages = sortedStages(workflow);
				const hint = stages
					.map((stage) => stage?.name)
					.filter(Boolean)
					.join(' → ');
				return (
					<div
						key={workflow.id}
						role="button"
						tabIndex={disabled ? -1 : 0}
						aria-pressed={selected}
						className={[
							'w-[148px] shrink-0 cursor-pointer overflow-hidden rounded-xl border bg-white text-left transition-all duration-150',
							'hover:-translate-y-0.5 hover:border-base-primary hover:shadow-md',
							selected
								? 'border-base-primary bg-blue-light-1 shadow-md'
								: 'border-base-border',
							disabled ? 'pointer-events-none opacity-60' : '',
						].join(' ')}
						onClick={() => handleSelect(workflow.id)}
						onKeyDown={(event) => {
							if (event.key === 'Enter' || event.key === ' ') {
								event.preventDefault();
								handleSelect(workflow.id);
							}
						}}
					>
						<div
							className="relative h-24"
							style={{
								background:
									'linear-gradient(135deg, var(--purple-dark), var(--blue-dark))',
							}}
						>
							{workflow.is_default ? (
								<span className="absolute top-2 left-2 rounded bg-white/90 px-1.5 py-0.5 text-[10px] leading-4 text-base-primary">
									{t('projectSetting.workflowDefault')}
								</span>
							) : null}
							{selected ? (
								<CheckCircleFilled className="absolute top-2 right-2 text-base text-white" />
							) : null}
							<div className="absolute inset-x-3 bottom-3 flex gap-1">
								{(stages.length ? stages : [{ id: 'empty' }])
									.slice(0, 4)
									.map((stage) => (
										<span
											key={stage.id}
											className="h-1.5 flex-1 rounded-full bg-white/80"
											style={
												stage.color
													? {
															background: getWorkflowStateColor(
																stage.color,
															),
														}
													: undefined
											}
										/>
									))}
							</div>
						</div>
						<div className="px-3 pt-2.5 pb-3">
							<div className="truncate text-[13px] font-medium text-black-dark">
								{displayValue(workflow.name)}
							</div>
							<div className="mt-1 line-clamp-2 text-[11px] leading-snug text-base-gray">
								{hint || t('projectList.templateEmptyHint')}
							</div>
						</div>
					</div>
				);
			})}
		</div>
	);
}
