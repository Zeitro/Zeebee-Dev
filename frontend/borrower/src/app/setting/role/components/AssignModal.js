'use client';
import React, { useState, useCallback, useRef } from 'react';
import { Row, Col, Form, Select, Input, InputNumber, Spin, Checkbox } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import {
	emailRegex,
	phoneRegex,
	creditScoreRegex,
	minCreditScore,
	maxCreditScore,
	validateRule,
	tinOrSSNRegex,
	createRequest,
} from '@repo/shared-methods/utils';
import { layoutCol, titleCol } from '@/enum/common';
import { useGetOptions } from '@repo/shared-methods/hooks';
import {
	CommonFormDate,
	FilterTable,
	InputNameCombination,
	InputPhone,
	InputTIN,
	NewModal,
} from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import useContent from './useContent';
const { TextArea } = Input;
const AssignModal = ({ setOpen, rowData, handleSearch }) => {
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [selectedRowKeys, setSelectedRowKeys] = useState([]);
	const [loading, setLoading] = useState(false);
	const tableFilterRef = useRef();
	const { assignColumns } = useContent();
	// const { maritalStatusOptions } = useGetOptions([
	// 	{
	// 		key: 'marital_status',
	// 		valueName: 'maritalStatusOptions',
	// 	},
	// ]);

	const handleDetails = useCallback(async () => {
		try {
			setLoading(true);
			const res = await request(`/borrower_applicant/${rowData?.id}`, {
				method: 'GET',
			});
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [rowData]);

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

	useDeferredSyncEffect({
		predicate: () => !!rowData?.id,
		deferred: () => {
			handleDetails();
		},
		deps: [handleDetails, rowData?.id],
		isUseMicrotask: true,
	});
	const rowSelection = {
		type: 'checkbox',
		selectedRowKeys,
		onChange: (selectedRowKeys, selectedRows) => {
			setSelectedRowKeys(selectedRowKeys);
			console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
		},
		// getCheckboxProps: (record) => ({
		// 	disabled: record.name === 'Disabled User', // Column configuration not to be checked
		// 	name: record.name,
		// }),
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
			<FilterTable
				filterWrapClassName="!h-[60vh] overflow-hidden"
				ref={tableFilterRef}
				columns={assignColumns}
				scroll={{ x: 1200 }}
				hasFilterField={false}
				// pagination={pagination}
				api="/borrower_applicant/page"
				onRow={() => ({
					onClick: () => {
						// handleSelectedRow(row, stage);
					},
				})}
				rowSelection={rowSelection}
			/>
		</NewModal>
	);
};
export default AssignModal;
