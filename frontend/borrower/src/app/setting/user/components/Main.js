import React, { useState, useCallback } from 'react';
import { Row, Col, Form, Select, Input, InputNumber, Spin } from 'antd';
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
			okText={rowData?.id ? 'Save change' : 'Add'}
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
					className={[
						' overflow-hidden overflow-y-auto',
						rowData?.id ? 'h-[64vh]' : 'h-[56vh]',
					].join(' ')}
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
				>
					<Row gutter={16}>
						<Col {...titleCol} className="!text-xl font-medium text-base-gray pb-4">
							User information
						</Col>
						{rowData?.id && (
							<InputNameCombination
								layoutCol={{ ...layoutCol, md: 12 }}
								firstName={['first_name']}
								lastName={['last_name']}
								hideMiddleName={true}
							/>
						)}
						<Col {...titleCol}>
							<Form.Item
								label="Email"
								name={['email']}
								rules={[
									{
										required: true,
									},
									{
										validator: (_, value) =>
											validateRule(value, emailRegex, 'Invalid email format'),
									},
								]}
							>
								<Input placeholder="Enter email" disabled={rowData?.id} />
							</Form.Item>
						</Col>
						<Col {...titleCol}>
							<Form.Item
								label="Phone"
								name={['phone']}
								rules={[
									{
										validator: (_, value) =>
											validateRule(value, phoneRegex, 'Invalid phone number'),
									},
								]}
							>
								<InputPhone />
							</Form.Item>
						</Col>
						<Col {...titleCol}>
							<Form.Item label="Remark" name={['remark']}>
								<TextArea
									placeholder="Enter remark"
									autoSize={{ minRows: 3, maxRows: 5 }}
								/>
							</Form.Item>
						</Col>
						<Col {...titleCol} className="!text-xl font-medium text-base-gray pb-4">
							User role
						</Col>
						<Col {...titleCol}>
							<Form.Item label="Role" name={['role']} rules={[{ required: true }]}>
								<Select options={[]} placeholder="Select role" />
							</Form.Item>
						</Col>
					</Row>
				</Form>
			</Spin>
		</NewModal>
	);
};
export default MainModal;
