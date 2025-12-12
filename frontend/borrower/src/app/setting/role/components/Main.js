'use client';
import React, { useState, useCallback } from 'react';
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
	InputNameCombination,
	InputPhone,
	InputTIN,
	NewModal,
} from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import MenuPermission from './MenuPermission';
const { TextArea } = Input;
const MainModal = ({ setOpen, rowData, handleSearch }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
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
			form.setFieldsValue(res);
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [form, rowData]);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/borrower_applicant';
			let method = 'POST';
			if (rowData?.id) {
				url = `/borrower_applicant/${rowData?.id}`;
				method = 'PUT';
			}
			const res = await request(url, {
				method,
				body: { ...value },
			});
			handleSearch?.(res?.id);
			setConfirmLoading(false);
			setOpen(false);
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
	return (
		<NewModal
			width={600}
			title={`${rowData?.id ? 'Edit' : 'Add'} user`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Spin spinning={loading}>
				<Form
					form={form}
					name="user-form"
					layout="vertical"
					autoComplete="off"
					className="h-[40vh] overflow-hidden overflow-y-auto"
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
				>
					<Row gutter={16}>
						<Col {...layoutCol}>
							<Form.Item
								label="Role name"
								name={['name']}
								rules={[{ required: true }]}
							>
								<Input placeholder="Enter name" />
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Date permission"
								name={['permission']}
								rules={[
									{
										required: true,
									},
								]}
							>
								<Select options={[]} placeholder="Select role" />
							</Form.Item>
						</Col>
						<Col {...layoutCol}>
							<Form.Item
								label="Menu permission"
								name={['menu']}
								rules={[{ required: true }]}
								valuePropName="checked"
							>
								<MenuPermission
									treeData={[
										{
											label: 'my loans',
											value: 'loans',
											children: [{ label: 'view', value: 'view' }],
										},
										{ label: 'my properties', value: 'properties' },
									]}
								/>
							</Form.Item>
						</Col>
					</Row>
				</Form>
			</Spin>
		</NewModal>
	);
};
export default MainModal;
