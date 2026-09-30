/**
 * 产品模板库：公共模板画廊。悬停预览看板 / 工作流，创建走已有 CreateProjectModal 并预选该模板。
 * 列表 POST `/project_management/template/page`，`template_type = Public`；
 * 预览 GET `/project_management/template/:id`。公共模板只读，不提供删除。
 */
'use client';

import { Spin, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { useAuth } from '@repo/shared-methods/hooks';
import { displayValue } from '@repo/shared-methods/utils';
import { TEMPLATE_TYPE } from '../../../enum/template';
import CreateProjectModal from '../../ProjectList/components/CreateProjectModal';
import { TemplateCoverCard } from '../components/TemplateCoverCard';
import { TemplatePreviewModal } from '../components/TemplatePreviewModal';
import { useTemplateGallery } from '../hooks/useTemplateGallery';
import { useTemplatePage } from '../hooks/useTemplatePage';
import { mapTemplatePreview, sortedStages, templateSummary } from '../lib/mapTemplatePreview';

/**
 * @param {{ PERMISSIONS?: { PROJECT_LIST?: { ADD?: string } } }} [props]
 */
export function ProjectTemplateLibrary({ PERMISSIONS } = {}) {
	const t = useTranslations();
	const { hasPermission } = useAuth();
	const gallery = useTemplateGallery();
	const { loading, templates, loadDetail } = useTemplatePage({
		templateType: TEMPLATE_TYPE.PUBLIC,
		messageApi: gallery.messageApi,
	});

	const canCreate = !PERMISSIONS || hasPermission(PERMISSIONS?.PROJECT_LIST?.ADD);

	const openPreview = (template) => {
		gallery.openPreview(template.id, async () =>
			mapTemplatePreview(await loadDetail(template.id)),
		);
	};

	return (
		<div className="flex-1 overflow-auto px-8 pt-10 pb-8">
			{gallery.messageContextHolder}
			<div className="mb-5">
				<div className="text-lg font-semibold text-black-dark">
					{t('projectTemplate.libraryTitle')}
				</div>
				<div className="mt-1 text-[13px] text-base-gray">
					{t('projectTemplate.libraryHint')}
				</div>
			</div>
			{loading ? (
				<div className="flex justify-center py-12">
					<Spin />
				</div>
			) : templates.length === 0 ? (
				<NoData text={t('projectTemplate.emptyLibrary')} className="min-h-[240px]" />
			) : (
				<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					{templates.map((template) => (
						<TemplateCoverCard
							key={template.id}
							stages={sortedStages(template)}
							canCreate={canCreate}
							onPreview={() => openPreview(template)}
							onCreate={() => gallery.openCreate(template.id)}
						>
							<div className="mb-1.5 flex items-center gap-2">
								<span className="truncate text-sm font-semibold text-black-dark">
									{displayValue(template.name)}
								</span>
								{template.is_default ? (
									<Tag color="blue" className="m-0 shrink-0">
										{t('projectSetting.workflowDefault')}
									</Tag>
								) : null}
							</div>
							<div className="line-clamp-2 text-xs leading-relaxed text-base-gray">
								{templateSummary(template) || t('projectList.templateEmptyHint')}
							</div>
						</TemplateCoverCard>
					))}
				</div>
			)}
			{gallery.previewOpen ? (
				<TemplatePreviewModal
					loading={gallery.previewLoading}
					preview={gallery.preview}
					canCreate={canCreate}
					onClose={() => gallery.setPreviewOpen(false)}
					onCreate={gallery.createFromPreview}
				/>
			) : null}
			{gallery.createOpen ? (
				<CreateProjectModal
					messageApi={gallery.messageApi}
					setOpen={gallery.setCreateOpen}
					handleSearch={gallery.goProjectList}
					initialWorkflowId={gallery.workflowId}
				/>
			) : null}
		</div>
	);
}

export default ProjectTemplateLibrary;
