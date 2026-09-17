/**
 * 我的模板：当前用户另存的项目模板。悬停预览，创建走已有 CreateProjectModal。
 */
'use client';

import { Spin } from 'antd';
import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { useAuth } from '@repo/shared-methods/hooks';
import CreateProjectModal from '../../ProjectList/components/CreateProjectModal';
import { MyTemplateGrid } from '../components/MyTemplateGrid';
import { TemplatePreviewModal } from '../components/TemplatePreviewModal';
import { useMyTemplateList } from '../hooks/useMyTemplateList';
import { useTemplateGallery } from '../hooks/useTemplateGallery';
import { mapTemplatePreview } from '../lib/mapTemplatePreview';

/**
 * @param {{ PERMISSIONS?: { PROJECT_LIST?: { ADD?: string } } }} [props]
 */
export function ProjectTemplateMine({ PERMISSIONS } = {}) {
	const t = useTranslations();
	const { hasPermission } = useAuth();
	const gallery = useTemplateGallery();
	const { loading, templates, handleDelete, loadPreview } = useMyTemplateList(gallery.messageApi);
	const canCreate = !PERMISSIONS || hasPermission(PERMISSIONS?.PROJECT_LIST?.ADD);

	const openPreview = (template) => {
		gallery.openPreview(template.workflow_template_id || template.workflow_id, async () =>
			mapTemplatePreview(await loadPreview(template.id)),
		);
	};

	return (
		<div className="flex-1 overflow-auto px-8 pt-10 pb-8">
			{gallery.messageContextHolder}
			<div className="mb-5">
				<div className="text-lg font-semibold text-black-dark">
					{t('projectTemplate.mineTitle')}
				</div>
				<div className="mt-1 text-[13px] text-base-gray">
					{t('projectTemplate.mineHint')}
				</div>
			</div>
			{loading ? (
				<div className="flex justify-center py-12">
					<Spin />
				</div>
			) : templates.length === 0 ? (
				<NoData text={t('projectTemplate.emptyMine')} className="min-h-[240px]" />
			) : (
				<MyTemplateGrid
					templates={templates}
					canCreate={canCreate}
					onPreview={openPreview}
					onCreate={(template) =>
						gallery.openCreate(template.workflow_template_id || template.workflow_id)
					}
					onDelete={handleDelete}
				/>
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

export default ProjectTemplateMine;
