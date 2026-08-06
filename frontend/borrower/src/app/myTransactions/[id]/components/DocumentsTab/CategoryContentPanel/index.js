'use client';

import React from 'react';
import { displayValue } from '@repo/shared-methods/utils';
import { CATEGORY_CONTENT_TYPE, resolveCategoryContentType } from '../constants';
import UploadDocumentList from './UploadDocumentList';
import ParseFileButton from './ParseFileButton';

function renderCategoryContent(category, onRefresh, transactionId) {
	const contentType = resolveCategoryContentType(category);

	switch (contentType) {
		case CATEGORY_CONTENT_TYPE.UPLOAD:
			return (
				<UploadDocumentList
					category={category}
					onRefresh={onRefresh}
					transactionId={transactionId}
				/>
			);
		default:
			return (
				<div className="px-6 py-12 text-center text-base-label">
					This category view is not available yet.
				</div>
			);
	}
}

/**
 * 分类内容面板：标题 + Parse file + 文档列表。
 * 解析请求与状态由 ParseFileButton 自行管理；pollParseStatus 由 DocumentsTab 注入。
 */
const CategoryContentPanel = ({ category, onRefresh, transactionId, pollParseStatus }) => {
	const hasParseAllFilesButton = category?.file_parse_config?.rules?.some?.(
		(rule) => rule.triggered_by === 'Manual',
	);

	return (
		<div className="min-w-0 flex-1 overflow-hidden rounded-lg border border-base-border bg-background">
			<div className="flex items-center justify-between gap-4 border-b border-base-border px-6 py-4">
				<h2 className="m-0 min-w-0 truncate text-lg font-medium leading-7 text-foreground">
					{category?.name ?? displayValue(null)}
				</h2>
				{hasParseAllFilesButton && (
					<ParseFileButton
						documentId={category?.document_id}
						transactionId={transactionId}
						category={category}
						pollParseStatus={pollParseStatus}
					/>
				)}
			</div>

			{/* ParseStatusBar：有 file_parse_config 时展示解析步骤，后续接轮询后打开 */}

			{category ? renderCategoryContent(category, onRefresh, transactionId) : null}
		</div>
	);
};

export default CategoryContentPanel;
