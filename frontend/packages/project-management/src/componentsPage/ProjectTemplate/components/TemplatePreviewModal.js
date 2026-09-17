/**
 * 模板预览弹层：看板 / 工作流两个页签。有快照图时优先展示图片，否则用只读视图。
 * 确认走创建，不在此弹层提交项目。
 */
'use client';

import { Button, Spin, Tabs } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { displayValue } from '@repo/shared-methods/utils';
import { TemplateBoardView, TemplateWorkflowView } from './TemplatePreviewViews';

/**
 * @param {{ src?: string, children: import('react').ReactNode }} props
 */
function SnapshotPane({ src, children }) {
	if (src) {
		return (
			// eslint-disable-next-line @next/next/no-img-element -- 快照地址可能是任意域名或 data URL
			<img src={src} alt="" className="block w-full rounded-lg" />
		);
	}
	return children;
}

/**
 * @param {{
 *   loading?: boolean,
 *   preview?: object|null,
 *   onClose: () => void,
 *   onCreate: () => void,
 *   canCreate?: boolean
 * }} props
 */
export function TemplatePreviewModal({
	loading = false,
	preview,
	onClose,
	onCreate,
	canCreate = true,
}) {
	const t = useTranslations();
	const title = preview?.name
		? t('projectTemplate.previewTitle', { name: displayValue(preview.name) })
		: t('projectTemplate.previewFallbackTitle');

	return (
		<NewModal
			width={980}
			title={title}
			okText={t('common.create')}
			cancelText={t('projectTemplate.close')}
			onOk={onCreate}
			onCancel={onClose}
			{...(canCreate
				? {}
				: {
						footer: <Button onClick={onClose}>{t('projectTemplate.close')}</Button>,
					})}
			destroyOnHidden
		>
			{loading || !preview ? (
				<div className="flex justify-center py-16">
					<Spin />
				</div>
			) : (
				<Tabs
					items={[
						{
							key: 'board',
							label: t('projectTemplate.board'),
							children: (
								<SnapshotPane src={preview.snapshot_board}>
									<TemplateBoardView preview={preview} />
								</SnapshotPane>
							),
						},
						{
							key: 'workflow',
							label: t('projectTemplate.workflow'),
							children: (
								<SnapshotPane src={preview.snapshot_workflow}>
									<TemplateWorkflowView preview={preview} />
								</SnapshotPane>
							),
						},
					]}
				/>
			)}
		</NewModal>
	);
}
