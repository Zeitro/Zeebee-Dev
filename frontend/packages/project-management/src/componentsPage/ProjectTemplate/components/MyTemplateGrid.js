/**
 * 我的模板卡片网格。删除确认在卡片内；预览 / 创建回调由页面传入。
 * `templates` 为 POST `/project_management/template/page` 的 list 项，直接读后端字段。
 */
'use client';

import { DeleteOutlined } from '@ant-design/icons';
import { Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { DeleteBtn } from '@repo/shared-methods';
import { displayValue, formatDate } from '@repo/shared-methods/utils';
import { sortedStages, templateSummary } from '../lib/mapTemplatePreview';
import { TemplateCoverCard } from './TemplateCoverCard';

/**
 * @param {{
 *   templates: object[],
 *   canCreate?: boolean,
 *   onPreview: (template: object) => void,
 *   onCreate: (template: object) => void,
 *   onDelete: (id: number|string) => void
 * }} props
 */
export function MyTemplateGrid({ templates, canCreate, onPreview, onCreate, onDelete }) {
	const t = useTranslations();

	return (
		<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{templates.map((template) => {
				const createdAt = template.created_at ? formatDate(template.created_at) : '';
				return (
					<TemplateCoverCard
						key={template.id}
						stages={sortedStages(template)}
						canCreate={canCreate}
						onPreview={() => onPreview(template)}
						onCreate={() => onCreate(template)}
					>
						<div className="flex items-start gap-2">
							<div className="min-w-0 flex-1">
								<div className="mb-1.5 flex items-center gap-2">
									<span className="truncate text-sm font-semibold text-black-dark">
										{displayValue(template.name)}
									</span>
									<Tag color="green" className="m-0 shrink-0">
										{t('projectTemplate.mineTag')}
									</Tag>
								</div>
								<div className="line-clamp-2 text-xs leading-relaxed text-base-gray">
									{templateSummary(template) ||
										t('projectTemplate.customFallback')}
								</div>
								{createdAt ? (
									<div className="mt-2 text-[11px] text-base-label">
										{createdAt}
									</div>
								) : null}
							</div>
							<DeleteBtn
								type="text"
								size="small"
								danger
								isStopPropagation
								title={t('projectTemplate.deleteConfirm')}
								okText={t('common.delete')}
								cancelText={t('common.cancel')}
								onOk={() => onDelete(template.id)}
							>
								<DeleteOutlined />
							</DeleteBtn>
						</div>
					</TemplateCoverCard>
				);
			})}
		</div>
	);
}
