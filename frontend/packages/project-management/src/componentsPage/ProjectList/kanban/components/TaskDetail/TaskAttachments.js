/**
 * 任务附件：UploadFileButton 直传 AWS，成功后走 saveTask 统一 PUT 更新。
 */
'use client';

import { useCallback, useState } from 'react';
import { DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { Button, Image, Spin } from 'antd';
import { useTranslations } from 'next-intl';
import { DeleteBtn, NewModal, UploadFileButton } from '@repo/shared-methods';
import { displayValue, formatTime } from '@repo/shared-methods/utils';
import { IMAGE_FILE_EXTENSIONS } from '../constants';
import { getFileExtension } from '../helpers';
import { TaskDetailPanel, TaskDetailPanelAction } from './TaskDetailPanel';

/**
 * @param {{
 *   task: object,
 *   saveTask: (sourceTask: object, values: object) => Promise<object|null>,
 *   updateTaskDetail: (sourceTask: object, overrides: object) => Promise<object|null>,
 *   messageApi?: import('antd/es/message/interface').MessageInstance,
 *   onTaskUpdated?: (task: object|null) => void,
 * }} props
 */
export function TaskAttachments({ task, saveTask, updateTaskDetail, messageApi, onTaskUpdated }) {
	const t = useTranslations();
	const items = Array.isArray(task?.attachments) ? task.attachments : [];
	const [previewItem, setPreviewItem] = useState(null);
	const [saving, setSaving] = useState(false);
	const [deletingId, setDeletingId] = useState(null);

	const handleUploadSuccess = useCallback(
		async (fileUrl, file) => {
			if (!fileUrl || !file?.name) return;
			const currentAttachments = Array.isArray(task?.attachments) ? task.attachments : [];
			setSaving(true);
			try {
				const updated = await saveTask(task, {
					attachments: [
						...currentAttachments,
						{
							file_name: file.name,
							file_url: fileUrl,
							storage_type: 'Aws S3',
							upload_type: 'Manual',
						},
					],
				});
				onTaskUpdated?.(updated);
			} finally {
				setSaving(false);
			}
		},
		[onTaskUpdated, saveTask, task],
	);

	const handleDelete = async (fileId) => {
		setDeletingId(fileId);
		try {
			const updated = await updateTaskDetail(task, {
				attachments: items.filter((item) => item.id !== fileId),
			});
			if (!updated) {
				messageApi?.error(t('kanban.attachmentDeleteFailed'));
				return;
			}
			messageApi?.success(t('kanban.attachmentDeleteSuccess'));
			onTaskUpdated?.(updated);
		} catch {
			messageApi?.error(t('kanban.attachmentDeleteFailed'));
		} finally {
			setDeletingId(null);
		}
	};

	const handleOpen = (item) => {
		const ext = getFileExtension(item.file_name);
		if (IMAGE_FILE_EXTENSIONS.includes(ext) || ext === 'pdf') {
			setPreviewItem(item);
		} else {
			window.open(item.file_url, '_blank', 'noopener,noreferrer');
		}
	};

	const uploadButton = (
		<UploadFileButton
			showUploadList={false}
			immediateUpload={true}
			disabled={saving}
			onSuccess={handleUploadSuccess}
		>
			<Spin spinning={saving} size="small">
				<TaskDetailPanelAction
					title={t('kanban.uploadAttachment')}
					icon={<UploadOutlined />}
				/>
			</Spin>
		</UploadFileButton>
	);

	return (
		<>
			<TaskDetailPanel
				title={t('kanban.attachments')}
				action={uploadButton}
				isEmpty={items.length === 0}
				emptyText={t('kanban.attachmentsEmpty')}
			>
				{items.map((item) => (
					<div
						key={item.id ?? item.file_url}
						className="flex items-center gap-2.5 px-4 py-2.5 border-b border-base-border/60 last:border-b-0 text-xs hover:bg-base-active transition-colors"
					>
						<Button
							type="link"
							size="small"
							className="font-medium truncate max-w-[220px] !px-0"
							title={t('kanban.attachmentPreview', { name: item.file_name })}
							onClick={() => handleOpen(item)}
						>
							{item.file_name}
						</Button>
						<Button
							type="link"
							size="small"
							icon={<DownloadOutlined className="text-blue-medium" />}
							href={item.file_url}
							download={item.file_name}
							title={t('kanban.download')}
						/>
						<span className="text-base-label flex-1">
							{displayValue(item.created_at, formatTime)}
						</span>
						<DeleteBtn
							type="text"
							size="small"
							danger
							icon={<DeleteOutlined />}
							title={t('kanban.attachmentDeleteConfirm')}
							okText={t('common.confirm')}
							cancelText={t('common.cancel')}
							onOk={() => handleDelete(item.id)}
							confirmLoading={deletingId === item.id}
							className="!px-0 !gap-0"
						>
							<></>
						</DeleteBtn>
					</div>
				))}
			</TaskDetailPanel>

			{previewItem ? (
				<NewModal
					title={previewItem.file_name}
					footer={null}
					width={800}
					onCancel={() => setPreviewItem(null)}
				>
					{getFileExtension(previewItem.file_name) === 'pdf' ? (
						<iframe
							src={previewItem.file_url}
							title={previewItem.file_name}
							className="w-full h-[70vh] border-0"
						/>
					) : (
						<Image
							src={previewItem.file_url}
							alt={previewItem.file_name}
							className="max-w-full block mx-auto"
						/>
					)}
				</NewModal>
			) : null}
		</>
	);
}
