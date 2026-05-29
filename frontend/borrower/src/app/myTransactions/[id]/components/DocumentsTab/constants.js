/** 右侧内容区展示类型；非 upload 时可在此扩展并在 CategoryContentPanel 中分支 */
export const CATEGORY_CONTENT_TYPE = {
	UPLOAD: 'upload',
};

export function resolveCategoryContentType(category) {
	return category?.contentType ?? CATEGORY_CONTENT_TYPE.UPLOAD;
}
