/**
 * 任务详情统一区块卡片：标题栏 + 固定尺寸操作按钮 + 空态。
 */

import { forwardRef } from 'react';
import { Button, Tooltip } from 'antd';

/** 区块右上角操作按钮统一尺寸 */
export const PANEL_ACTION_CLASS = '!w-6 !h-6 !min-w-6 !p-0 inline-flex items-center justify-center';

/**
 * @param {{
 *   title: string,
 *   icon: import('react').ReactNode,
 *   onClick?: () => void,
 * } & import('antd/es/button').ButtonProps} props
 */
export const TaskDetailPanelAction = forwardRef(function TaskDetailPanelAction(
	{ title, icon, ...props },
	ref,
) {
	return (
		<Tooltip title={title}>
			<Button
				ref={ref}
				type="default"
				size="small"
				shape="circle"
				icon={icon}
				className={PANEL_ACTION_CLASS}
				{...props}
			/>
		</Tooltip>
	);
});

/**
 * @param {{
 *   title: string,
 *   titleSuffix?: import('react').ReactNode,
 *   action?: import('react').ReactNode,
 *   children?: import('react').ReactNode,
 *   emptyText?: string,
 *   isEmpty?: boolean,
 *   className?: string,
 * }} props
 */
export function TaskDetailPanel({
	title,
	titleSuffix,
	action,
	children,
	emptyText,
	isEmpty = false,
	className = '',
}) {
	return (
		<section
			className={`mb-6 rounded-xl border border-base-border bg-white overflow-hidden ${className}`.trim()}
		>
			<div className="flex items-center justify-between gap-3 px-4 py-2.5 border-b border-base-border/60">
				<h3 className="text-xs font-semibold uppercase tracking-wider text-base-label m-0 inline-flex items-center gap-2">
					{title}
					{titleSuffix ? (
						<span className="font-normal normal-case tracking-normal tabular-nums text-base-disabled">
							{titleSuffix}
						</span>
					) : null}
				</h3>
				<div className="w-6 h-6 shrink-0 flex items-center justify-center">{action}</div>
			</div>
			{isEmpty && emptyText ? (
				<div className="px-4 py-5 text-center text-[13px] text-base-disabled">
					{emptyText}
				</div>
			) : (
				children
			)}
		</section>
	);
}
