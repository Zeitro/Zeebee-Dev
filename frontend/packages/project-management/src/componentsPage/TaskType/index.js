/**
 * 任务类型列表页：无页标题；FilterTable 仅放新建按钮，表单走独立 Modal。
 */
'use client';

import { FilterTable } from '@repo/shared-methods';
import { TaskTypeFormModal } from './components/TaskTypeFormModal';
import { useTaskTypeList } from './hooks/useTaskTypeList';

export function ProjectTaskType() {
	const {
		tableFilterRef,
		filterList,
		columns,
		getData,
		messageContextHolder,
		modalContextHolder,
		modalOpen,
		saving,
		editingType,
		closeModal,
		handleSave,
		emptyText,
	} = useTaskTypeList();

	return (
		<>
			{messageContextHolder}
			{modalContextHolder}
			<FilterTable
				filterWrapClassName="overflow-hidden bg-background rounded-lg p-4"
				ref={tableFilterRef}
				columns={columns}
				filterList={filterList}
				hasFilterField
				getData={getData}
				rowProps={{ justify: 'end' }}
				locale={{ emptyText }}
			/>
			<TaskTypeFormModal
				open={modalOpen}
				confirmLoading={saving}
				editingType={editingType}
				onCancel={closeModal}
				onSave={handleSave}
			/>
		</>
	);
}

export default ProjectTaskType;
