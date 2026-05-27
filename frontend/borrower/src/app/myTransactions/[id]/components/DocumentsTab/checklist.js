/**
 * 将交易详情 API 返回的 `rowData.checklist` 转为 DocumentsTab 侧栏与右侧列表所用结构。
 * 纯 UI 映射层：不修改 rowData，不发起请求。
 */

/**
 * 由分类名称生成稳定的 React key；无有效 slug 时回退为 document_id。
 * @param {string} name - API 分类名（如 "Purchase Contract"）
 * @param {number|string} documentId - 父级 document_id，作 slug 为空时的兜底
 */
function slugifyCategoryKey(name, documentId) {
	const slug = String(name ?? '')
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '_')
		.replace(/^_|_$/g, '');
	return slug || String(documentId ?? '');
}

/**
 * 将 API 子文档（children 中的一项）映射为 CollapsibleDocumentList 的单条 document。
 * document_type 对应 API `name`；status 由是否已有上传文件推导；updated_at 取首个文件的更新时间。
 */
function mapChecklistDocument(doc) {
	const files = doc?.files ?? [];
	const fileCount = doc?.file_count ?? files.length;
	const hasFiles = fileCount > 0;
	const updatedAt = files[0]?.updated_at ?? null;

	return {
		id: doc.document_id,
		key: String(doc.document_id),
		document_type: doc.name,
		required: doc.is_required,
		status: hasFiles ? 'Ready for review' : 'Pending upload',
		updated_at: updatedAt,
		active: true,
		files,
	};
}

/**
 * 将 API 一级分类（required_documents / optional_documents 中的一项）映射为侧栏 category。
 * completed：子文档已全部上传；pendingCount：尚未上传的子文档数（用于角标）。
 */
function mapChecklistCategory(doc, required) {
	const children = doc?.children ?? [];
	const subCount = doc?.sub_document_count ?? children.length;
	const uploaded = doc?.uploaded_document_count ?? 0;
	const pendingCount = Math.max(0, subCount - uploaded);

	return {
		key: slugifyCategoryKey(doc.name, doc.document_id),
		name: doc.name,
		required,
		completed: subCount > 0 && uploaded >= subCount,
		pendingCount,
		documents: children.map(mapChecklistDocument),
		contentType: 'upload',
	};
}

/**
 * 合并必填与选填分类，得到 DocumentsTab 左侧列表的 categories 数组。
 * @param {object|null|undefined} checklist - `rowData.checklist`
 * @returns {Array<{ key, name, required, completed, pendingCount, documents, contentType }>}
 */
export function buildDocumentCategoriesFromChecklist(checklist) {
	if (!checklist) return [];

	const required = (checklist.required_documents ?? []).map((doc) =>
		mapChecklistCategory(doc, true),
	);
	const optional = (checklist.optional_documents ?? []).map((doc) =>
		mapChecklistCategory(doc, false),
	);

	return [...required, ...optional];
}

/**
 * 从 checklist 汇总字段计算整体进度（侧栏顶部 Progress 与 x/y 文案）。
 * @param {object|null|undefined} checklist
 * @returns {{ completed: number, total: number }}
 */
export function getChecklistProgress(checklist) {
	if (!checklist) {
		return { completed: 0, total: 0 };
	}

	const completed = checklist.total_required_updated_document_count ?? 0;
	const total = checklist.total_required_document_count ?? 0;

	return { completed, total };
}
