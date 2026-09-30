/**
 * 模板库 / 我的模板共用：按 `template_type` 拉取模板分页、详情、删除。不渲染 UI。
 * 分页 POST `/project_management/template/page`；
 * 详情 GET / 删除 DELETE `/project_management/template/:id`（详情 data 与 list 项同结构，含 `stages`）。
 * 返回的 `templates` 为接口原始 list 项，字段名与后端一致（id / template_type / name / description / is_default / stages …）。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { useTranslations } from 'next-intl';
import { createRequest } from '@repo/shared-methods/utils';
import { getTemplateApi, TEMPLATE_PAGE_API } from '../../../api/template';

/**
 * @param {{
 *   templateType: string,
 *   messageApi: { error: (msg: string) => void, success: (msg: string) => void },
 *   keyword?: string,
 *   page?: number,
 *   limit?: number
 * }} params `templateType` 取 `TEMPLATE_TYPE.PUBLIC` / `TEMPLATE_TYPE.ORGANIZATION`
 * @returns {{
 *   loading: boolean,
 *   templates: object[],
 *   total: number,
 *   loadDetail: (id: number|string) => Promise<object|null>,
 *   handleDelete: (id: number|string) => Promise<void>
 * }}
 */
export function useTemplatePage({ templateType, messageApi, keyword = '', page = 1, limit = 20 }) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(true);
	const [templates, setTemplates] = useState([]);
	const [total, setTotal] = useState(0);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const data = await request(TEMPLATE_PAGE_API, {
					method: 'POST',
					isHideTips: true,
					body: {
						page,
						limit,
						filter: { keyword, template_type: templateType },
					},
				});
				if (cancelled) return;
				setTemplates(Array.isArray(data?.list) ? data.list : []);
				setTotal(Number(data?.total) || 0);
			} catch {
				if (!cancelled) {
					messageApi.error(t('projectTemplate.loadFailed'));
					setTemplates([]);
					setTotal(0);
				}
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [keyword, limit, messageApi, page, request, t, templateType]);

	const loadDetail = (id) => request(getTemplateApi(id), { method: 'GET', isHideTips: true });

	const handleDelete = async (id) => {
		try {
			await request(getTemplateApi(id), { method: 'DELETE', isHideTips: true });
			messageApi.success(t('projectTemplate.deleteSuccess'));
			setTemplates((list) => list.filter((item) => item.id !== id));
			setTotal((count) => Math.max(0, count - 1));
		} catch {
			messageApi.error(t('projectTemplate.deleteFailed'));
		}
	};

	return { loading, templates, total, loadDetail, handleDelete };
}
