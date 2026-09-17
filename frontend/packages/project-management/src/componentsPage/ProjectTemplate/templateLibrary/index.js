/**
 * 产品模板库：组织工作流模板画廊。悬停预览看板 / 工作流，创建走已有 CreateProjectModal 并预选该模板。
 * 列表 GET `/project_management/workflow/all`；预览再取 GET `/project_management/workflow/:id`（列表不含流转）。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { Spin, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { useAuth } from '@repo/shared-methods/hooks';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { getWorkflowApi, WORKFLOW_ALL_API } from '../../../api/workflow';
import CreateProjectModal from '../../ProjectList/components/CreateProjectModal';
import { TemplateCoverCard } from '../components/TemplateCoverCard';
import { TemplatePreviewModal } from '../components/TemplatePreviewModal';
import { useTemplateGallery } from '../hooks/useTemplateGallery';
import { mapTemplatePreview, sortedStages } from '../lib/mapTemplatePreview';

/**
 * @param {{ PERMISSIONS?: { PROJECT_LIST?: { ADD?: string } } }} [props]
 */
export function ProjectTemplateLibrary({ PERMISSIONS } = {}) {
	const t = useTranslations();
	const { hasPermission } = useAuth();
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(true);
	const [templates, setTemplates] = useState([]);
	const gallery = useTemplateGallery();

	const canCreate = !PERMISSIONS || hasPermission(PERMISSIONS?.PROJECT_LIST?.ADD);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const data = await request(WORKFLOW_ALL_API, { method: 'GET', isHideTips: true });
				if (cancelled) return;
				setTemplates(Array.isArray(data?.list) ? data.list : []);
			} catch {
				if (!cancelled) {
					gallery.messageApi.error(t('projectTemplate.loadFailed'));
					setTemplates([]);
				}
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [gallery.messageApi, request, t]);

	const openPreview = (workflow) => {
		gallery.openPreview(workflow.id, async () => {
			const data = await request(getWorkflowApi(workflow.id), {
				method: 'GET',
				isHideTips: true,
			});
			return mapTemplatePreview(data);
		});
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
					{templates.map((workflow) => {
						const stages = sortedStages(workflow);
						const hint = stages
							.map((stage) => stage?.name)
							.filter(Boolean)
							.join(' → ');
						const description = String(workflow.description || '').trim();
						return (
							<TemplateCoverCard
								key={workflow.id}
								stages={stages}
								canCreate={canCreate}
								onPreview={() => openPreview(workflow)}
								onCreate={() => gallery.openCreate(workflow.id)}
							>
								<div className="mb-1.5 flex items-center gap-2">
									<span className="truncate text-sm font-semibold text-black-dark">
										{displayValue(workflow.name)}
									</span>
									{workflow.is_default ? (
										<Tag color="blue" className="m-0 shrink-0">
											{t('projectSetting.workflowDefault')}
										</Tag>
									) : null}
								</div>
								<div className="line-clamp-2 text-xs leading-relaxed text-base-gray">
									{description || hint || t('projectList.templateEmptyHint')}
								</div>
							</TemplateCoverCard>
						);
					})}
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
