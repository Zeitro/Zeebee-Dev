'use client';

import React from 'react';
import { Tag } from 'antd';

// Per-file status color
const fileStatusColor = {
	created: 'green',
	replaced: 'blue',
	failed: 'red',
	removed: 'orange',
};

// Job-level status color
const jobStatusColor = {
	Success: 'green',
	PartialFailed: 'orange',
	Failed: 'red',
};

/**
 * Renders the full result of a Google Drive sync, sourced from the transaction
 * operations log metadata:
 * { status, total, created, replaced, failed, removed, folder_web_link, files[], error }
 * Used for the "Synced to Google Drive" event detail in LogTab (replaces the old success modal).
 */
const DriveSyncResult = ({ result }) => {
	if (!result || typeof result !== 'object') {
		return null;
	}
	const files = Array.isArray(result.files) ? result.files : [];
	const status = result.status;

	return (
		<div className="w-full">
			{status && (
				<div className="pb-2">
					<Tag color={jobStatusColor[status] || 'default'} className="m-0">
						{status}
					</Tag>
				</div>
			)}

			{result.error && <div className="pb-2 text-sm text-red-600">{result.error}</div>}

			<div className="flex flex-wrap gap-x-4 gap-y-1 pb-3 text-sm">
				<span>Total: {result.total ?? 0}</span>
				<span className="text-green-600">Created: {result.created ?? 0}</span>
				<span className="text-blue-600">Replaced: {result.replaced ?? 0}</span>
				<span className="text-red-600">Failed: {result.failed ?? 0}</span>
				<span className="text-orange-600">Removed: {result.removed ?? 0}</span>
			</div>

			{result.folder_web_link && (
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
							<Tag color={fileStatusColor[file.status] || 'default'}>
								{file.status}
							</Tag>
						</div>
					))}
				</div>
			)}
		</div>
	);
};

export default DriveSyncResult;
