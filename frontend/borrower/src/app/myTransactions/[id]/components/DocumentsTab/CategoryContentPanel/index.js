'use client';

import React, { useState } from 'react';
import { displayValue } from '@repo/shared-methods/utils';
import { CATEGORY_CONTENT_TYPE, resolveCategoryContentType } from '../constants';
import UploadDocumentList from './UploadDocumentList';
import ParseStatusBar from './ParseStatusBar';
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
 * 分类内容面板：标题 + Parse file + 解析步骤条 + 文档列表。
 * 解析请求与状态放在本组件；ParseFileButton / ParseStatusBar 仅负责展示。
 */
const CategoryContentPanel = ({ category, onRefresh, transactionId }) => {
	/** @type {['idle' | 'uploaded' | 'parsing' | 'complete', Function]} */
	const [parseStatus, setParseStatus] = useState('idle');

	const handleParseFile = () => {
		if (parseStatus === 'parsing') return;
		// 先进入「解析中」；真实解析接口 / 轮询后续写在此处
		setParseStatus('parsing');
	};

	return (
		<div className="min-w-0 flex-1 overflow-hidden rounded-lg border border-base-border bg-background">
			<div className="flex items-center justify-between gap-4 border-b border-base-border px-6 py-4">
				<h2 className="m-0 min-w-0 truncate text-lg font-medium leading-7 text-foreground">
					{category?.name ?? displayValue(null)}
				</h2>
				<ParseFileButton
					categoryName={category?.name}
					parseStatus={parseStatus}
					onClick={handleParseFile}
				/>
			</div>

			<ParseStatusBar status={parseStatus} />

			{category ? renderCategoryContent(category, onRefresh, transactionId) : null}
		</div>
	);
};

export default CategoryContentPanel;
