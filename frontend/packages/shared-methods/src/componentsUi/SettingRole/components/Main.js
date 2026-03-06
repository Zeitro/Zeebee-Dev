'use client';
import React, { useState, useCallback, useRef } from 'react';
import { Row, Col, Form, Select, Input, Spin, Checkbox, Radio, message, Space } from 'antd';
import { createRequest } from '@repo/shared-methods/utils';
import { titleCol } from '@repo/shared-methods/enum';
import { NewModal, FieldTips } from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import MenuPermission from './MenuPermission';
const { TextArea } = Input;
const MainModal = ({ setOpen, rowData, handleSearch, systemType }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [operationsValue, setOperationsValue] = useState([]);
	const [menuList, setMenuList] = useState([]);

	const handleMenuDetails = useCallback(async () => {
		const payload = new URLSearchParams({ type: systemType });
		const url = `/resource/tree?${payload}`;
		const res = await request(url);
		setMenuList(res?.list);
		setLoading(false);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [systemType]);

	const handleDetails = useCallback(async () => {
		const url = `/role/${rowData?.id}`;
		const res = await request(url);
		res.resources = res.resources?.map((i) => ({
			resource_id: i.id,
			data_scope: i.data_scope,
		}));
		form.setFieldsValue({ ...res });
		setLoading(false);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [rowData]);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/role';
			let method = 'POST';
			if (rowData?.id) {
				url = `/role/${rowData?.id}`;
				method = 'PUT';
			}
			const res = await request(url, {
				method,
				body: { ...value },
			});
			handleSearch?.(res?.id);
			message.open({
				content: 'Save successfully.',
				type: 'success',
			});
			setConfirmLoading(false);

			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	useDeferredSyncEffect({
		predicate: () => !!systemType,
		deferred: async () => {
			try {
				setLoading(true);
				await handleMenuDetails();
				if (rowData?.id) {
					handleDetails();
				}
				setLoading(false);
			} catch (error) {
				setLoading(false);
				throw error;
			}
		},
		deps: [handleMenuDetails, systemType, rowData?.id],
		isUseMicrotask: true,
	});
	return (
		<NewModal
			width={700}
			title={`${rowData?.id ? 'Edit' : 'Add'} role`}
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
					className="h-[60vh] overflow-hidden overflow-y-auto"
					size="large"
					disabled={confirmLoading}
					scrollToFirstError
					// initialValues={{ ...rowData }}
				>
					<Row gutter={16}>
						{/* <Col md={24} className="!text-xl font-medium pb-5 text-base-gray">
							Role information
						</Col> */}
						<Col md={18}>
							<Form.Item
								label="Role name"
								name={['name']}
								rules={[{ required: true }]}
							>
								<Input placeholder="Enter name" />
							</Form.Item>
						</Col>
						<Col md={6}>
							<Form.Item
								label="Status"
								name={['status']}
								rules={[
									{
										required: true,
									},
								]}
							>
								<Select
									options={[
										{ value: 'Active', label: 'Active' },
										{ value: 'Disable', label: 'Disable' },
									]}
									placeholder="Select status"
								/>
							</Form.Item>
						</Col>
						<Col {...titleCol}>
							<Form.Item label="Description" name={['description']}>
								<TextArea
									placeholder="Enter description"
									autoSize={{ minRows: 3, maxRows: 3 }}
								/>
							</Form.Item>
						</Col>
						{/* <Col {...titleCol} className="!text-xl font-medium pb-5 text-base-gray">
							Role permissions
						</Col> */}
						{/* <Col md={24}>
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
						</Col> */}
						<Col md={24} className="relative">
							{/* <Checkbox.Group
								options={[
									{
										label: 'Fold',
										value: 'fold',
									},
									{
										label: 'Select all',
										value: 'all',
									},
								]}
								value={operationsValue}
								onChange={(settingValue) => {
									setOperationsValue(settingValue);
									// menuPermissionRef.current.handleChangeSetting({
									// 	settingValue,
									// 	checkedValue: resources,
									// 	isFold: operationsValue.includes('fold'),
									// 	isSelectAll: operationsValue.includes('all'),
									// });
								}}
								
							/> */}
							<Space className="absolute top-0 right-0">
								<Checkbox
									checked={operationsValue?.includes('Fold')}
									onChange={(e) => {
										console.log(e);

										setOperationsValue(e.target.checked ? ['Fold'] : []);
									}}
								>
									Fold
								</Checkbox>
								<Checkbox
									checked={operationsValue?.includes('All')}
									onChange={(e) => {
										setOperationsValue(e.target.checked ? ['All'] : []);
									}}
								>
									Select all
								</Checkbox>
							</Space>

							<Form.Item
								label={
									<FieldTips
										rootClassName="!max-w-[300px]"
										className="ml-1"
										isShowDescription
										title={
											<>
												Add a menu to the character.
												<p>Note:</p>
												1. Selecting the Menu class will display it in the
												side menu bar. Only selecting the Button class will
												display the corresponding action button.
												<br />
												2. When you have data permissions, please select
												either "Personal data only" or "All data".
											</>
										}
									>
										Menu permission
									</FieldTips>
								}
								name={['resources']}
								rules={[{ required: true }]}
							>
								<MenuPermission
									treeData={menuList}
									isFold={operationsValue.includes('Fold')}
									isSelectAll={operationsValue.includes('All')}
									className="border border-base-border rounded-lg !py-4 !px-4 [&_.ant-tree-treenode]:!mb-2"
									titleName="name"
									keyName="id"
									height={300}
									blockNode={true}
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
