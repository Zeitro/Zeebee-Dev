'use client';

import React from 'react';
import { Tag } from 'antd';

const statusColor = {
	created: 'green',
	replaced: 'blue',
	failed: 'red',
	removed: 'orange',
};

const SyncResultContent = ({ result }) => {
	const files = result?.files || [];
	return (
		<div className="pt-2">
			<div className="flex gap-4 pb-3 text-sm">
				<span>Total: {result?.total ?? 0}</span>
				<span className="text-green-600">Created: {result?.created ?? 0}</span>
				<span className="text-blue-600">Replaced: {result?.replaced ?? 0}</span>
				<span className="text-red-600">Failed: {result?.failed ?? 0}</span>
				<span className="text-orange-600">Removed: {result?.removed ?? 0}</span>
			</div>
			{result?.folder_web_link && (
				<div className="pb-3 text-sm">
					<a href={result.folder_web_link} target="_blank" rel="noreferrer">
						Open transaction folder in Google Drive
					</a>
				</div>
			)}
			{files.length > 0 && (
				<div className="max-h-[280px] overflow-auto rounded border border-base-border">
					{files.map((file, index) => (
						<div
							key={`${file.file_name}-${index}`}
							className="flex items-center justify-between gap-2 border-b border-base-border px-3 py-2 text-sm last:border-b-0"
						>
							<span className="truncate" title={file.file_name}>
								{file.drive_url ? (
									<a href={file.drive_url} target="_blank" rel="noreferrer">
										{file.file_name}
									</a>
								) : (
									file.file_name
								)}
							</span>
							<Tag color={statusColor[file.status] || 'default'}>{file.status}</Tag>
						</div>
					))}
				</div>
			)}
		</div>
	);
};

export default SyncResultContent;
