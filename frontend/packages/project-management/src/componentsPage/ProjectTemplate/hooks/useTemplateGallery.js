/**
 * 模板页共用的预览弹层与创建项目弹框开关。不拉列表。
 * 创建走已有 CreateProjectModal；成功后回项目列表。
 */
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { message } from 'antd';
import { useTranslations } from 'next-intl';

export function useTemplateGallery() {
	const t = useTranslations();
	const router = useRouter();
	const [messageApi, messageContextHolder] = message.useMessage();
	const [createOpen, setCreateOpen] = useState(false);
	const [workflowId, setWorkflowId] = useState();
	const [previewOpen, setPreviewOpen] = useState(false);
	const [previewLoading, setPreviewLoading] = useState(false);
	const [preview, setPreview] = useState(null);

	const openCreate = (id) => {
		setWorkflowId(id);
		setCreateOpen(true);
	};

	/**
	 * @param {number|string|undefined} id 创建时预选的工作流模板；没有则走弹框默认项
	 * @param {() => Promise<object|null>} loadPreview
	 */
	const openPreview = async (id, loadPreview) => {
		setWorkflowId(id);
		setPreview(null);
		setPreviewLoading(true);
		setPreviewOpen(true);
		try {
			const data = await loadPreview();
			if (!data) throw new Error('empty');
			setPreview(data);
		} catch {
			messageApi.error(t('projectTemplate.previewFailed'));
			setPreviewOpen(false);
		} finally {
			setPreviewLoading(false);
		}
	};

	const createFromPreview = () => {
		setPreviewOpen(false);
		setCreateOpen(true);
	};

	const goProjectList = () => {
		router.push('/project/projectList');
	};

	return {
		messageApi,
		messageContextHolder,
		createOpen,
		setCreateOpen,
		workflowId,
		previewOpen,
		setPreviewOpen,
		previewLoading,
		preview,
		openCreate,
		openPreview,
		createFromPreview,
		goProjectList,
	};
}
