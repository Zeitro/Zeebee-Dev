'use client';

import React, { useState } from 'react';
import { DeleteOutlined, FolderOpenOutlined, FolderOutlined } from '@ant-design/icons';
import { Space, Typography } from 'antd';
import { DeleteBtn, ViewFile } from '@repo/shared-methods';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { getFileTypeIcon } from '@repo/shared-methods/enum';

/** 中间省略时保留的尾部长度（含扩展名），过短文件名不截断 */
const FILENAME_ELLIPSIS_SUFFIX_COUNT = 14;

function getFileSuffix(file, fileKey = 'file_url') {
	const name = file?.[fileKey]?.split('/').pop() ?? '';
	return name.split('.').pop()?.toLowerCase() ?? '';
}

function getFileName(file, fileKey = 'file_url') {
	return file?.[fileKey]?.split('/').pop() ?? '';
}

/**
 * 文件名中间省略：Typography.Text + ellipsis.suffix（antd 官方中间省略写法）。
 * @param {{ children: string, suffixCount?: number }} props
 */
function EllipsisMiddleFilename({ children, suffixCount = FILENAME_ELLIPSIS_SUFFIX_COUNT }) {
	const text = children ?? '';
	if (!text) {
		return <span className="text-sm text-foreground">{displayValue(null)}</span>;
	}
	if (text.length <= suffixCount) {
		return (
			<span className="min-w-0 truncate text-sm text-foreground" title={text}>
				{text}
			</span>
		);
	}
	const start = text.slice(0, text.length - suffixCount);
	const suffix = text.slice(-suffixCount).trim();
	return (
		<Typography.Text
			className="!mb-0 min-w-0 text-sm text-foreground"
			style={{ maxWidth: '100%' }}
			ellipsis={{ suffix, tooltip: text }}
		>
			{start}
		</Typography.Text>
	);
}

export function DefaultParentNameCell({ record, fieldKey = 'document_type', open }) {
	return (
		<>
			{open ? <FolderOpenOutlined /> : <FolderOutlined />}
			<span className="pl-2">{record[fieldKey]}</span>
		</>
	);
}

export function DefaultFileNameCell({ file, fileKey = 'file_url' }) {
	const name = file?.file_name ?? getFileName(file, fileKey);
	const suffix = getFileSuffix(file, fileKey);

	return (
		<div className="inline-flex h-8 max-w-full min-w-0 items-center gap-1 overflow-hidden rounded-md border border-base-border bg-background px-3">
			{getFileTypeIcon(suffix, { className: 'h-[16px] w-[16px] shrink-0' })}
			<EllipsisMiddleFilename>{name}</EllipsisMiddleFilename>
		</div>
	);
}

export function DefaultFileActions({
	file,
	hideDelete,
	onDeleteFile,
	onAfterDelete,
	fileKey = 'file_url',
}) {
	const [loading, setLoading] = useState(false);
	const { request } = createRequest();

	const handleDelete = async () => {
		if (onDeleteFile) {
			await onDeleteFile(file);
			onAfterDelete?.();
			return;
		}

		try {
			setLoading(true);
			await request(`/file/${file.id}`, { method: 'DELETE' });
			onAfterDelete?.();
		} finally {
			setLoading(false);
		}
	};

	return (
		<Space size={0}>
			<ViewFile type="link" file={file} fileKey={fileKey} size="small" />
			{!hideDelete && (
				<DeleteBtn
					icon={<DeleteOutlined />}
					onOk={handleDelete}
					confirmLoading={loading}
					isStopPropagation
					content={file?.file_name ?? file?.name}
				>
					Delete
				</DeleteBtn>
			)}
		</Space>
	);
}
