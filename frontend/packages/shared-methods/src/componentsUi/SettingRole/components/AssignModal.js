'use client';
import React, { useState, useCallback, useRef } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { FilterTable, NewModal } from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import useContent from './useContent';
const AssignModal = ({ setOpen, rowData, handleSearch, systemType }) => {
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [selectedRowKeys, setSelectedRowKeys] = useState([]);
	const [loading, setLoading] = useState(false);
	const [list, setList] = useState([]);
	const tableFilterRef = useRef();
	const { assignColumns } = useContent({ handleSearch, systemType });

	const handleSave = async () => {
		try {
			// setConfirmLoading(true);
			// const value = await form.validateFields();
			// let url = '/borrower_applicant';
			// let method = 'POST';
			// if (rowData?.id) {
			// 	url = `/borrower_applicant/${rowData?.id}`;
			// 	method = 'PUT';
			// }
			// const res = await request(url, {
			// 	method,
			// 	body: { ...value },
			// });
			// handleSearch?.(res?.id);
			// setConfirmLoading(false);
			// setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	const rowSelection = {
		type: 'checkbox',
		selectedRowKeys,
		onChange: (selectedRowKeys, selectedRows) => {
			setSelectedRowKeys(selectedRowKeys);
		},
		// getCheckboxProps: (record) => {
		// 	console.log(record);
		// 	const checkedItem = record.roles?.find((i) => i.id === rowData?.id);
		// 	return {
		// 		checked: !!checkedItem,
		// 		disabled: !!checkedItem, // Column configuration not to be checked
		// 		name: record.name,
		// 	};
		// },
	};
	return (
		<NewModal
			width={700}
			title="Assign"
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<div className="text-xl font-medium text-base-gray pb-4">
				{rowData?.name}
				<span className="text-sm font-normal pl-2">
					(<span className="text-base-error pr-1">*</span>Please add users within this
					permission range.)
				</span>
			</div>
			<FilterTable
				filterWrapClassName="!h-[60vh] overflow-hidden"
				ref={tableFilterRef}
				columns={assignColumns}
				scroll={{ x: 1100 }}
				hasFilterField={false}
				pagination={false}
				rowSelection={rowSelection}
				settingPagination={{ pageSize: 999 }}
				getData={async (params) => {
					const res = await request('/organization_user/page', {
						method: 'POST',
						body: { ...params },
					});
					const checkList = res?.list
						?.filter((i) => i.roles?.some((j) => j.id === rowData?.id))
						?.map((i) => i.id);

					setList(res?.list);
					setSelectedRowKeys(checkList);
					return res;
				}}
			/>
		</NewModal>
	);
};
export default AssignModal;
