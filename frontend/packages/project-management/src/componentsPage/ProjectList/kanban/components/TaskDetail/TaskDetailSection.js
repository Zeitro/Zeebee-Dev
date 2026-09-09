/**
 * 任务详情抽屉区块标题与间距；纯展示容器，不含业务状态。
 */

/**
 * @param {{
 *   title?: string,
 *   extra?: import('react').ReactNode,
 *   children: import('react').ReactNode,
 *   className?: string,
 * }} props
 */
export function TaskDetailSection({ title, extra, children, className = '' }) {
	if (!title && !extra) {
		return <section className={`mb-6 ${className}`.trim()}>{children}</section>;
	}

	return (
		<section className={`mb-6 ${className}`.trim()}>
			<div className="flex items-center justify-between gap-3 mb-3">
				{title ? (
					<h3 className="text-xs font-semibold uppercase tracking-wider text-base-label m-0">
						{title}
					</h3>
				) : (
					<span />
				)}
				{extra}
			</div>
			{children}
		</section>
	);
}
