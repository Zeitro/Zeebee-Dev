/**
 * 任务附件：无后端 API 时用 localStorage 持久化（与 demo 一致）。
 */
'use client';

import { useEffect, useRef, useState } from 'react';
import {
	DeleteOutlined,
	DownloadOutlined,
	PaperClipOutlined,
	UploadOutlined,
} from '@ant-design/icons';
import { Button, Input, Modal, Popconfirm, message } from 'antd';
import { useTranslations } from 'next-intl';
import { dayjs } from '@repo/shared-methods/utils';
import { useUserStore } from '@repo/shared-methods/store';
import { IMAGE_FILE_EXTENSIONS } from './constants';
import { dataUrlToBlob, formatFileSize, getFileExtension, getMemberLabel } from './helpers';

/**
 * @param {{
 *   taskId: number,
 *   projectId?: number|string,
 *   members?: Array,
 * }} props
 */
export function TaskAttachments({ taskId, projectId, members = [] }) {
	const t = useTranslations();
	const userInfo = useUserStore((state) => state.userInfo);
	const storageKey = `pm-attachments-${projectId || 'default'}-${taskId}`;
	const [items, setItems] = useState(() => {
		if (typeof window === 'undefined') return [];
		try {
			return JSON.parse(localStorage.getItem(storageKey) || '[]');
		} catch {
			return [];
		}
	});
	const fileInputRef = useRef(null);
	const [previewItem, setPreviewItem] = useState(null);

	useEffect(() => {
		if (typeof window === 'undefined') return;
		try {
			setItems(JSON.parse(localStorage.getItem(storageKey) || '[]'));
		} catch {
			setItems([]);
		}
	}, [storageKey]);

	const persist = (next) => {
		setItems(next);
		try {
			localStorage.setItem(storageKey, JSON.stringify(next));
		} catch {
			message.error(t('kanban.attachmentSaveFailed'));
		}
	};

	const handleOpen = (item) => {
		const ext = getFileExtension(item.name);
		if (IMAGE_FILE_EXTENSIONS.includes(ext) || ext === 'pdf') {
			setPreviewItem(item);
		} else {
			const url = URL.createObjectURL(dataUrlToBlob(item.dataUrl));
			window.open(url, '_blank');
		}
	};

	const handleFileChange = (event) => {
		const file = event.target.files?.[0];
		event.target.value = '';
		if (!file) return;
		const reader = new FileReader();
		reader.onload = () => {
			persist([
				{
					id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
					name: file.name,
					dataUrl: String(reader.result),
					size: file.size,
					uploaderId: userInfo?.id || 0,
					uploaderName: userInfo?.name || userInfo?.email || t('kanban.unknownUser'),
					uploadedAt: new Date().toISOString(),
					note: '',
				},
				...items,
			]);
		};
		reader.readAsDataURL(file);
	};

	const handleNoteChange = (id, note) => {
		const next = items.map((item) => (item.id === id ? { ...item, note } : item));
		setItems(next);
		try {
			localStorage.setItem(storageKey, JSON.stringify(next));
		} catch {
			/* ignore */
		}
	};

	return (
		<div className="mb-3">
			<div className="flex items-center gap-2">
				<PaperClipOutlined className="text-base-label text-xs" />
				<span className="font-medium text-[13px]">{t('kanban.attachments')}</span>
				<Button
					size="small"
					icon={<UploadOutlined />}
					onClick={() => fileInputRef.current?.click()}
				>
					{t('kanban.uploadAttachment')}
				</Button>
				{/* <input
					ref={fileInputRef}
					type="file"
					className="hidden"
					onChange={handleFileChange}
				/> */}
			</div>
			{items.length > 0 ? (
				<div className="mt-2 border border-base-border rounded-md overflow-hidden">
					{items.map((item) => {
						const uploader = members.find((member) => member.id === item.uploaderId);
						return (
							<div
								key={item.id}
								className="flex items-center gap-2.5 px-2.5 py-1.5 border-b border-base-border last:border-b-0 text-xs"
							>
								<button
									type="button"
									className="text-base-primary font-medium truncate max-w-[220px] text-left"
									title={t('kanban.attachmentPreview', { name: item.name })}
									onClick={() => handleOpen(item)}
								>
									{item.name}
								</button>
								<a
									href={item.dataUrl}
									download={item.name}
									title={t('kanban.download')}
								>
									<DownloadOutlined className="text-blue-medium" />
								</a>
								<span className="text-base-label shrink-0">
									{dayjs(item.uploadedAt).format('YYYY-MM-DD HH:mm')}{' '}
									{t('kanban.uploadedBy', {
										name:
											getMemberLabel(uploader) ||
											item.uploaderName ||
											t('kanban.unknownUser'),
									})}
								</span>
								<span className="text-base-disabled shrink-0">
									{formatFileSize(item.size)}
								</span>
								<Input
									size="small"
									variant="borderless"
									placeholder={t('kanban.attachmentNote')}
									value={item.note}
									onChange={(event) =>
										handleNoteChange(item.id, event.target.value)
									}
									className="flex-1 min-w-[60px] !px-1"
								/>
								<Popconfirm
									title={t('kanban.attachmentDeleteConfirm')}
									okText={t('common.confirm')}
									cancelText={t('common.cancel')}
									onConfirm={() =>
										persist(items.filter((row) => row.id !== item.id))
									}
								>
									<Button
										type="text"
										size="small"
										danger
										icon={<DeleteOutlined />}
									/>
								</Popconfirm>
							</div>
						);
					})}
				</div>
			) : null}
			<Modal
				open={Boolean(previewItem)}
				title={previewItem?.name}
				footer={null}
				width={800}
				onCancel={() => setPreviewItem(null)}
			>
				{previewItem ? (
					getFileExtension(previewItem.name) === 'pdf' ? (
						<iframe
							src={previewItem.dataUrl}
							title={previewItem.name}
							className="w-full h-[70vh] border-0"
						/>
					) : (
						<img
							src={previewItem.dataUrl}
							alt={previewItem.name}
							className="max-w-full block mx-auto"
						/>
					)
				) : null}
			</Modal>
		</div>
	);
}
