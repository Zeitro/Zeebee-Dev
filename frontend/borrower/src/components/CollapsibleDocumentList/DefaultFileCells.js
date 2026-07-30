'use client';

import React, { useState } from 'react';
import { DeleteOutlined, FolderOpenOutlined, FolderOutlined } from '@ant-design/icons';
import { Space } from 'antd';
import { DeleteBtn, ViewFile } from '@repo/shared-methods';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { getFileTypeIcon } from '@repo/shared-methods/enum';
function getFileSuffix(file, fileKey = 'file_url') {
	const name = file?.[fileKey]?.split('/').pop() ?? '';
	return name.split('.').pop()?.toLowerCase() ?? '';
}

function getFileName(file, fileKey = 'file_url') {
	return file?.[fileKey]?.split('/').pop() ?? '';
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
		<div className="inline-flex h-8 max-w-full items-center gap-1 rounded-md border border-base-border bg-background px-3">
			{getFileTypeIcon(suffix, { className: 'w-[16px] h-[16px]' })}
			<span className="min-w-0 truncate text-sm text-foreground">{displayValue(name)}</span>
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
