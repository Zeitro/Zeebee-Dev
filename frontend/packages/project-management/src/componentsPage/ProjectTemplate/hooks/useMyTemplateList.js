/**
 * 我的模板列表：拉取、删除。不渲染 UI。
 * GET `/project_management/project_template/all`；DELETE `/project_management/project_template/:id`。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { useTranslations } from 'next-intl';
import { createRequest } from '@repo/shared-methods/utils';
import {
	getProjectTemplateDetailApi,
	PROJECT_TEMPLATE_ALL_API,
} from '../../../api/projectTemplate';

/**
 * @param {unknown} data
 * @returns {object[]}
 */
function toTemplateList(data) {
	if (Array.isArray(data)) return data;
	return Array.isArray(data?.list) ? data.list : [];
}

/**
 * @param {{ error: (msg: string) => void, success: (msg: string) => void }} messageApi
 */
export function useMyTemplateList(messageApi) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(true);
	const [templates, setTemplates] = useState([]);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const data = await request(PROJECT_TEMPLATE_ALL_API, {
					method: 'GET',
					isHideTips: true,
				});
				if (cancelled) return;
				setTemplates(toTemplateList(data));
			} catch {
				if (!cancelled) {
					messageApi.error(t('projectTemplate.loadFailed'));
					setTemplates([]);
				}
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [messageApi, request, t]);

	const handleDelete = async (id) => {
		try {
			await request(getProjectTemplateDetailApi(id), { method: 'DELETE', isHideTips: true });
			messageApi.success(t('projectTemplate.deleteSuccess'));
			setTemplates((list) => list.filter((item) => item.id !== id));
		} catch {
			messageApi.error(t('projectTemplate.deleteFailed'));
		}
	};

	const loadPreview = async (id) => {
		return request(getProjectTemplateDetailApi(id), { method: 'GET', isHideTips: true });
	};

	return { loading, templates, handleDelete, loadPreview };
}
