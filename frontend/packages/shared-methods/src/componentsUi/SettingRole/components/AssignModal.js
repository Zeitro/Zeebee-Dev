'use client';
import React, { useState, useRef } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { FilterTable, NewModal } from '@repo/shared-methods';
import useContent from './useContent';
import { message } from 'antd';
const AssignModal = ({ setOpen, rowData, handleSearch, systemType }) => {
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [selectedRowKeys, setSelectedRowKeys] = useState([]);
	const tableFilterRef = useRef();
	const { assignColumns } = useContent({ handleSearch, systemType });

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const res = await request(`/role/${rowData?.id}/assign_organization_users`, {
				method: 'POST',
				body: {
					organization_user_ids: selectedRowKeys,
				},
			});
			handleSearch?.(res?.id);
			setConfirmLoading(false);
			message.open({
				type: 'success',
				content: 'Assign successful.',
			});
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	const rowSelection = {
		type: 'checkbox',
		selectedRowKeys,
		onChange: (selectedRowKeys) => {
			setSelectedRowKeys(selectedRowKeys);
		},
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

					setSelectedRowKeys(checkList);
					return res;
				}}
			/>
		</NewModal>
	);
};
export default AssignModal;
