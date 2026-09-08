'use client';
import React, { useState, useCallback } from 'react';
import { Row, Col, Form, Select, Input, Spin, Checkbox, message, Space } from 'antd';
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
	const [isFold, setIsFold] = useState(false);
	const [isSelectAll, setIsSelectAll] = useState(null); //必须为null，避免初始化时触发了（只能点击按钮触发）
	const [menuList, setMenuList] = useState([]);

	const updateMenuList = (list, initCheckedArr = []) => {
		const menus = list.map((item) => {
			if (item.is_public) {
				item.disableCheckbox = true;
				initCheckedArr.push({
					resource_id: item.id,
					data_scope: item.enable_data_scope,
				});
			}
			if (item.children) {
				const obj = updateMenuList(item.children, initCheckedArr);
				item.children = obj.menus;
			}
			return item;
		});
		return {
			menus,
			initCheckedArr,
		};
	};

	const handleMenuDetails = useCallback(async () => {
		let params = { type: systemType };
		// admin端单独处理权限数据范围
		if (systemType === 'Admin') {
			params.scope = 'Organization';
		}
		const payload = new URLSearchParams(params);
		const url = `/resource/tree?${payload}`;
		const res = await request(url);
		const menuInfo = updateMenuList(res?.list || []);
		if (!rowData?.id) {
			form.setFieldsValue({
				resources: menuInfo?.initCheckedArr,
			});
		}
		setMenuList(menuInfo?.menus);
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
					initialValues={{ status: 'Active' }}
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
						<Col md={24} className="relative">
							<Space className="absolute top-0 right-0">
								<Checkbox
									checked={isFold}
									onChange={(e) => setIsFold(e.target.checked)}
								>
									Fold
								</Checkbox>
								<Checkbox
									checked={isSelectAll}
									onChange={(e) => setIsSelectAll(e.target.checked)}
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
												Add menu to the character.
												<p>Note:</p>
												1. Permission categories: 'Menu' is a page-level
												permission, and 'Button' is an operation button
												permission.
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
									isFold={isFold}
									isSelectAll={isSelectAll}
									className="border border-base-border rounded-lg !py-4 !px-4 [&_.ant-tree-treenode]:!mb-2"
									titleName="name"
									keyName="id"
									// height={300}
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
