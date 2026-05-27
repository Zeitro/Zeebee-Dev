/** 将 API checklist 转为 DocumentsTab 侧栏/列表所需结构（仅 UI 层，不改 rowData） */

function slugifyCategoryKey(name, documentId) {
	const slug = String(name ?? '')
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '_')
		.replace(/^_|_$/g, '');
	return slug || String(documentId ?? '');
}

function mapChecklistDocument(doc) {
	const files = doc?.files ?? [];
	const fileCount = doc?.file_count ?? files.length;
	const hasFiles = fileCount > 0;
	const updatedAt = files[0]?.updated_at ?? null;

	return {
		id: doc.document_id,
		key: String(doc.document_id),
		name: doc.name,
		required: doc.is_required,
		status: hasFiles ? 'Ready for review' : 'Pending upload',
		updated_at: updatedAt,
		active: true,
		files,
	};
}

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

export function getChecklistProgress(checklist) {
	if (!checklist) {
		return { completed: 0, total: 0 };
	}

	const completed =
		(checklist.total_required_updated_document_count ?? 0) +
		(checklist.total_optional_updated_document_count ?? 0);
	const total = checklist.total_document_count ?? 0;

	return { completed, total };
}
