/**
 * 项目工作流列表页：FilterTable 负责筛选/分页，新建走独立 Modal。
 */
'use client';

import { FilterTable } from '@repo/shared-methods';
import { CreateWorkflowModal } from './components/CreateWorkflowModal';
import { useWorkflowList } from './hooks/useWorkflowList';

export { WorkflowEditor } from './WorkflowEditor';

export function ProjectWorkflow() {
	const {
		tableFilterRef,
		filterList,
		columns,
		getData,
		messageContextHolder,
		modalContextHolder,
		modalOpen,
		saving,
		setModalOpen,
		handleSave,
		emptyText,
	} = useWorkflowList();

	return (
		<>
			{messageContextHolder}
			{modalContextHolder}
			<FilterTable
				filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
				ref={tableFilterRef}
				columns={columns}
				filterList={filterList}
				hasFilterField
				getData={getData}
				locale={{ emptyText }}
			/>
			<CreateWorkflowModal
				open={modalOpen}
				confirmLoading={saving}
				onCancel={() => setModalOpen(false)}
				onSave={handleSave}
			/>
		</>
	);
}

export default ProjectWorkflow;
