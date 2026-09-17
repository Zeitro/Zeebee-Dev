/**
 * 项目设置危险区：确认后 DELETE `/project_management/project/:id`，成功回项目列表。
 */
'use client';

import { useState } from 'react';
import { Button, Modal } from 'antd';
import { useTranslations } from 'next-intl';
import { getProjectDetailApi } from '../../../api/project';

/**
 * @param {{
 *   project: object,
 *   projectId: number|string,
 *   request: (url: string, options?: object) => Promise<unknown>,
 *   messageApi: { success: (msg: string) => void },
 *   onDeleted: () => void
 * }} props
 */
export function DeleteProjectSection({ project, projectId, request, messageApi, onDeleted }) {
	const t = useTranslations();
	const [modal, modalContextHolder] = Modal.useModal();
	const [deleting, setDeleting] = useState(false);

	const handleDelete = () => {
		modal.confirm({
			title: t('projectSetting.deleteConfirmTitle'),
			content: t('projectSetting.deleteConfirmContent', { name: project?.name }),
			okText: t('projectSetting.delete'),
			okType: 'danger',
			cancelText: t('common.cancel'),
			centered: true,
			onOk: async () => {
				try {
					setDeleting(true);
					await request(getProjectDetailApi(projectId), { method: 'DELETE' });
					messageApi.success(t('projectSetting.deleteSuccess'));
					onDeleted();
					setDeleting(false);
				} catch {
					setDeleting(false);
					// createRequest 已提示接口错误
					return Promise.reject();
				}
			},
		});
	};

	return (
		<div className="mt-10 rounded-lg border border-[var(--error-medium)] bg-[var(--error-light)] px-4 py-4">
			{modalContextHolder}
			<p className="m-0 font-medium text-base-error">{t('projectSetting.deleteTitle')}</p>
			<p className="mt-2 mb-4 text-sm text-base-error">{t('projectSetting.deleteWarning')}</p>
			<Button color="danger" variant="outlined" onClick={handleDelete} disabled={deleting}>
				{t('projectSetting.delete')}
			</Button>
		</div>
	);
}
