/**
 * 任务视图外壳：统一圆角面板与工具条，各视图自己填 toolbar / 内容。
 */
'use client';

/**
 * @param {{ toolbar?: import('react').ReactNode, children: import('react').ReactNode }} props
 */
export function ProjectViewFrame({ toolbar, children }) {
	return (
		<div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-lg border border-base-border bg-white">
			{toolbar ? (
				<div className="flex shrink-0 items-center justify-between gap-3 border-b border-base-border px-4 py-3">
					{toolbar}
				</div>
			) : null}
			<div className="flex min-h-0 flex-1 flex-col overflow-hidden">{children}</div>
		</div>
	);
}
