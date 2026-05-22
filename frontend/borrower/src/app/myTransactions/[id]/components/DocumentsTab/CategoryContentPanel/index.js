'use client';

import React from 'react';
import { displayValue } from '@repo/shared-methods/utils';
import { CATEGORY_CONTENT_TYPE, resolveCategoryContentType } from '../constants';
import UploadDocumentList from './UploadDocumentList';

function renderCategoryContent(category) {
	const contentType = resolveCategoryContentType(category);

	switch (contentType) {
		case CATEGORY_CONTENT_TYPE.UPLOAD:
			return <UploadDocumentList category={category} />;
		default:
			return (
				<div className="px-6 py-12 text-center text-base-label">
					This category view is not available yet.
				</div>
			);
	}
}

const CategoryContentPanel = ({ category }) => {
	return (
		<div className="min-w-0 flex-1 overflow-hidden rounded-lg border border-base-border bg-background">
			<div className="border-b border-base-border px-6 py-4">
				<h2 className="m-0 text-lg font-medium leading-7 text-foreground">
					{category?.name ?? displayValue(null)}
				</h2>
			</div>

			{category ? renderCategoryContent(category) : null}
		</div>
	);
};

export default CategoryContentPanel;
