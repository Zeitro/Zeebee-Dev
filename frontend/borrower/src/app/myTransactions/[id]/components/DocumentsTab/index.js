'use client';

import React, { useMemo, useState } from 'react';
import { Button, Progress, message } from 'antd';
import { CloudSyncOutlined, DiffOutlined, ImportOutlined, UploadOutlined } from '@ant-design/icons';
import { displayValue } from '@repo/shared-methods/utils';
import { MOCK_DOCUMENT_CATEGORIES } from './constants';
import DocumentCategoryButton from './DocumentCategoryButton';

function DocumentRow({ doc, onUpload }) {
	return (
		<div className="flex items-center justify-between border-b border-base-border px-6 py-4 last:border-b-0">
			<div className="flex min-w-0 items-center gap-2">
				<span className="flex size-6 shrink-0 items-center justify-center rounded-md bg-base-text2">
					<DiffOutlined className="text-sm text-base-gray" />
				</span>
				<span className="text-base font-medium leading-6 text-base-gray">{doc.name}</span>
			</div>
			<Button icon={<UploadOutlined />} onClick={() => onUpload(doc)}>
				Upload
			</Button>
		</div>
	);
}

const DocumentsTab = ({ rowData }) => {
	const categories = rowData?.document_categories ?? MOCK_DOCUMENT_CATEGORIES;
	const completed = rowData?.document_completed_count ?? rowData?.documents_completed ?? 16;
	const total = rowData?.document_total_count ?? rowData?.documents_total ?? 24;
	const progressPercent = total > 0 ? Math.round((completed / total) * 100) : 0;

	const defaultCategoryKey = useMemo(() => {
		const firstPending = categories.find((c) => c.pendingCount > 0 && !c.completed);
		return firstPending?.key ?? categories[0]?.key;
	}, [categories]);

	const [activeKey, setActiveKey] = useState(defaultCategoryKey);

	const activeCategory = categories.find((c) => c.key === activeKey) ?? categories[0];
	const requiredDocs = activeCategory?.documents?.filter((d) => d.required) ?? [];
	const optionalDocs = activeCategory?.documents?.filter((d) => !d.required) ?? [];

	const requiredCategories = categories.filter((c) => c.required);
	const optionalCategories = categories.filter((c) => !c.required);

	const handleUpload = (doc) => {
		message.info(`Upload is not wired yet for "${doc.name}".`);
	};

	return (
		<div className="flex gap-4 p-4">
			<div className="flex w-[311px] shrink-0 flex-col gap-2">
				<div className="flex gap-2">
					<Button
						className="flex-1"
						icon={<ImportOutlined />}
						onClick={() => message.info('Import is not wired yet.')}
					>
						Import
					</Button>
					<Button
						icon={<CloudSyncOutlined />}
						onClick={() => message.info('Sync to Google Drive is not wired yet.')}
					>
						Sync to Google Drive
					</Button>
				</div>

				<div className="rounded-lg border border-base-border bg-background p-6">
					<div className="mb-5 flex items-center gap-8">
						<p className="m-0 text-base font-medium leading-6 text-foreground">
							Overall progress
						</p>
						<p className="m-0 text-sm leading-[22px] text-base-gray">
							{total > 0 ? `${completed}/${total}` : displayValue(null)}
						</p>
					</div>
					<Progress
						percent={progressPercent}
						showInfo={false}
						strokeColor="var(--blue-dark)"
						trailColor="var(--blue-light-1)"
						size="small"
					/>
				</div>

				<div className="rounded-lg border border-base-border bg-background py-4">
					<p className="m-0 px-4 pb-2 text-sm font-medium leading-[22px] text-base-gray">
						Required
					</p>
					{requiredCategories.map((cat) => (
						<DocumentCategoryButton
							key={cat.key}
							category={cat}
							isActive={cat.key === activeKey}
							onSelect={setActiveKey}
						/>
					))}

					<div className="my-2 border-y border-base-border px-4 py-2">
						<p className="m-0 text-sm font-medium leading-[22px] text-base-gray">
							Optional
						</p>
					</div>

					{optionalCategories.map((cat) => (
						<DocumentCategoryButton
							key={cat.key}
							category={cat}
							isActive={cat.key === activeKey}
							onSelect={setActiveKey}
						/>
					))}
				</div>
			</div>

			<div className="min-w-0 flex-1 rounded-lg border border-base-border bg-background">
				<div className="border-b border-base-border px-6 py-4">
					<h2 className="m-0 text-2xl font-medium leading-8 text-foreground">
						{activeCategory?.name ?? displayValue(null)}
					</h2>
				</div>

				{requiredDocs.map((doc) => (
					<DocumentRow key={doc.key} doc={doc} onUpload={handleUpload} />
				))}

				{optionalDocs.length > 0 && (
					<>
						<div className="border-y border-base-border px-6 py-2">
							<p className="m-0 text-sm font-medium leading-[22px] text-base-gray">
								Optional
							</p>
						</div>
						{optionalDocs.map((doc) => (
							<DocumentRow key={doc.key} doc={doc} onUpload={handleUpload} />
						))}
					</>
				)}

				{!requiredDocs.length && !optionalDocs.length && (
					<div className="px-6 py-12 text-center text-base-label">
						No documents in this category yet.
					</div>
				)}
			</div>
		</div>
	);
};

export default DocumentsTab;
