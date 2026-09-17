/**
 * 模板封面卡：悬停露出预览 / 创建。纯 UI，不请求、不改外部数据。
 */
'use client';

import { Button } from 'antd';
import { useTranslations } from 'next-intl';
import { getWorkflowStateColor } from '../../../enum/workflow';

const COVER_GRADIENT = 'linear-gradient(135deg, var(--purple-dark), var(--blue-dark))';

/**
 * @param {{
 *   cover?: string,
 *   stages?: object[],
 *   children: import('react').ReactNode,
 *   onPreview: () => void,
 *   onCreate?: () => void,
 *   canCreate?: boolean
 * }} props
 */
export function TemplateCoverCard({
	cover,
	stages = [],
	children,
	onPreview,
	onCreate,
	canCreate = true,
}) {
	const t = useTranslations();
	const bars = (stages.length ? stages : [{ id: 'empty' }]).slice(0, 4);

	return (
		<div className="group h-full overflow-hidden rounded-xl border border-base-border bg-white shadow-sm transition-colors hover:border-base-primary">
			<div
				className="relative h-[140px] bg-cover bg-center"
				style={{
					backgroundImage: cover ? `url(${cover}), ${COVER_GRADIENT}` : undefined,
					background: cover ? undefined : COVER_GRADIENT,
				}}
			>
				<div className="absolute inset-x-3 bottom-3 flex gap-1">
					{bars.map((stage) => (
						<span
							key={stage.id}
							className="h-1.5 flex-1 rounded-full bg-white/80"
							style={
								stage.color
									? { background: getWorkflowStateColor(stage.color) }
									: undefined
							}
						/>
					))}
				</div>
				<div className="absolute inset-0 flex items-end justify-center gap-2.5 bg-gradient-to-t from-black/50 to-transparent pb-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100">
					<Button
						onClick={(event) => {
							event.stopPropagation();
							onPreview();
						}}
					>
						{t('projectTemplate.preview')}
					</Button>
					{canCreate ? (
						<Button
							type="primary"
							onClick={(event) => {
								event.stopPropagation();
								onCreate?.();
							}}
						>
							{t('common.create')}
						</Button>
					) : null}
				</div>
			</div>
			<div className="px-3.5 pt-3 pb-4">{children}</div>
		</div>
	);
}
