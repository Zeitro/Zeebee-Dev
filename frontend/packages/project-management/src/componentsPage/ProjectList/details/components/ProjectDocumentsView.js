/**
 * 项目文档：汇总各任务附件。列表与详情请求都在本组件内。
 * 删除/上传仍走任务详情，避免再写一套附件保存。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { Button, Image, Input, Spin } from 'antd';
import { DownloadOutlined, FileOutlined, SearchOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { createRequest, displayValue, formatTime } from '@repo/shared-methods/utils';
import { createTaskApi } from '@repo/project-management/api/task';
import { IMAGE_FILE_EXTENSIONS } from '../../kanban/components/constants';
import { getFileExtension } from '../../kanban/components/helpers';
import { useProjectDetailQuery } from '../hooks/useProjectDetailQuery';
import { useProjectTaskCatalog } from '../hooks/useProjectTaskCatalog';
import { ProjectViewFrame } from './ProjectViewFrame';

const DETAIL_CONCURRENCY = 5;

/**
 * @param {Array} tasks
 * @param {{ detail: Function }} taskApi
 * @param {number|string} projectId
 * @returns {Promise<Array>}
 */
async function collectAttachments(tasks, taskApi, projectId) {
	const queue = [...tasks];
	const files = [];

	const worker = async () => {
		while (queue.length) {
			const task = queue.shift();
			try {
				const detail = await taskApi.detail(task.id, projectId);
				const attachments = Array.isArray(detail?.attachments) ? detail.attachments : [];
				attachments.forEach((file) => {
					if (!file?.file_name && !file?.file_url) return;
					files.push({
						...file,
						taskId: task.id,
						taskTitle: task.title,
					});
				});
			} catch {
				// 单条详情失败不阻断其余文档
			}
		}
	};

	await Promise.all(Array.from({ length: DETAIL_CONCURRENCY }, worker));
	return files;
}

/**
 * @param {{ projectId?: number|string }} props
 */
export function ProjectDocumentsView({ projectId }) {
	const t = useTranslations();
	const { setTaskQuery } = useProjectDetailQuery();
	const { tasks, loading: listLoading, truncated } = useProjectTaskCatalog(projectId);
	const { request } = useMemo(() => createRequest(), []);
	const taskApi = useMemo(() => createTaskApi(request), [request]);
	const [files, setFiles] = useState(null);
	const [keyword, setKeyword] = useState('');
	const [preview, setPreview] = useState(null);

	useEffect(() => {
		if (!projectId || listLoading) return undefined;
		let cancelled = false;
		collectAttachments(tasks, taskApi, projectId)
			.then((next) => {
				if (!cancelled) setFiles(next);
			})
			.catch(() => {
				if (!cancelled) setFiles([]);
			});
		return () => {
			cancelled = true;
		};
	}, [listLoading, projectId, taskApi, tasks]);

	const visible = useMemo(() => {
		const query = keyword.trim().toLowerCase();
		const list = files || [];
		if (!query) return list;
		return list.filter((file) =>
			`${file.file_name || ''} ${file.taskTitle || ''}`.toLowerCase().includes(query),
		);
	}, [files, keyword]);

	const openFile = (file) => {
		const ext = getFileExtension(file.file_name || '');
		if (IMAGE_FILE_EXTENSIONS.includes(ext)) {
			setPreview(file);
			return;
		}
		if (file.file_url) window.open(file.file_url, '_blank', 'noopener,noreferrer');
	};

	const loading = listLoading || files == null;

	return (
		<ProjectViewFrame
			toolbar={
				<>
					<Input
						allowClear
						prefix={<SearchOutlined className="text-base-label" />}
						placeholder={t('projectDetail.docsSearch')}
						value={keyword}
						onChange={(event) => setKeyword(event.target.value)}
						className="max-w-xs"
					/>
					<span className="text-xs text-base-label">
						{truncated
							? t('projectDetail.catalogTruncated')
							: t('projectDetail.docsCount', { count: files?.length || 0 })}
					</span>
				</>
			}
		>
			{loading ? (
				<div className="flex flex-1 items-center justify-center">
					<Spin />
				</div>
			) : visible.length ? (
				<div className="min-h-0 flex-1 overflow-auto p-4">
					<div className="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3">
						{visible.map((file) => (
							<div
								key={`${file.taskId}-${file.id ?? file.file_url}`}
								className="flex min-w-0 flex-col gap-2 rounded-lg border border-base-border px-4 py-3 transition-colors hover:border-base-primary"
							>
								<div className="flex items-start gap-2">
									<FileOutlined className="mt-1 text-blue-medium" />
									<Button
										type="link"
										className="h-auto min-w-0 flex-1 truncate p-0 text-left text-sm font-medium"
										onClick={() => openFile(file)}
									>
										{displayValue(file.file_name)}
									</Button>
									{file.file_url ? (
										<Button
											type="text"
											size="small"
											icon={<DownloadOutlined />}
											href={file.file_url}
											download={file.file_name}
										/>
									) : null}
								</div>
								<Button
									type="link"
									className="h-auto max-w-full justify-start truncate p-0 text-xs"
									onClick={() => setTaskQuery(file.taskId)}
								>
									{displayValue(file.taskTitle)}
								</Button>
								<div className="text-xs text-base-label">
									{displayValue(file.created_at, formatTime)}
								</div>
							</div>
						))}
					</div>
				</div>
			) : (
				<div className="flex flex-1 items-center justify-center">
					<NoData text={t('projectDetail.docsEmpty')} />
				</div>
			)}
			<Image
				width={0}
				height={0}
				alt=""
				preview={{
					visible: Boolean(preview),
					src: preview?.file_url,
					onVisibleChange: (visibleImage) => {
						if (!visibleImage) setPreview(null);
					},
				}}
			/>
		</ProjectViewFrame>
	);
}
