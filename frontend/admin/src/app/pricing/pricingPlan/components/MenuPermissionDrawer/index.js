'use client';

import React, { useCallback, useMemo, useState } from 'react';
import { Button, Checkbox, Drawer, Form, Space, Spin, message } from 'antd';
import MenuPermission from '@repo/shared-methods/componentsUi/SettingRole/components/MenuPermission';
import { createRequest } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

/**
 * 定价计划菜单权限抽屉：树勾选与角色配置一致，保存 POST `/pricing_plan/{id}/bind_resources`。
 *
 * @param {object} props
 * @param {boolean} props.open
 * @param {(open: boolean) => void} props.setOpen
 * @param {object} props.rowData 当前定价计划行，需含 id
 * @param {'Borrower'|'Lender'|'Admin'} props.systemType 拉取 `/resource/tree` 的 type 参数
 * @param {() => void} [props.successCallback]
 */

/** 标记公共权限并收集默认勾选（与 SettingRole Main 一致） */
const updateMenuList = (list, initCheckedArr = []) => {
	const menus = list.map((item) => {
		if (item.is_public) {
			item.disableCheckbox = true;
			initCheckedArr.push({
				resource_id: item.id,
				data_scope: item.enable_data_scope ? 'All' : undefined,
			});
		}
		if (item.children) {
			const obj = updateMenuList(item.children, initCheckedArr);
			item.children = obj.menus;
		}
		return item;
	});
	return { menus, initCheckedArr };
};

/**
 * 将定价计划详情中的已绑资源转为 MenuPermission 表单值。
 * 兼容 `resources`（对象数组）与 `resource_ids`（id 数组）两种后端形态。
 *
 * @param {object} detail
 * @returns {Array<{ resource_id: number, data_scope?: string }>}
 */
const mapPlanResourcesToForm = (detail) => {
	if (detail?.resources?.length) {
		return detail.resources.map((item) => ({
			resource_id: item.id ?? item.resource_id,
			data_scope: item.data_scope || 'All',
		}));
	}
	return [];
};

const MenuPermissionDrawer = ({ open, setOpen, rowData, systemType, successCallback }) => {
	const { request } = useMemo(() => createRequest(), []);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [isFold, setIsFold] = useState(false);
	const [isSelectAll, setIsSelectAll] = useState(null);
	const [menuList, setMenuList] = useState([]);
	const [messageApi, contextHolder] = message.useMessage();

	const handleMenuDetails = useCallback(async () => {
		const payload = new URLSearchParams({ type: systemType });
		const res = await request(`/resource/tree?${payload}`);
		const menuInfo = updateMenuList(res?.list || []);
		setMenuList(menuInfo?.menus || []);
		return menuInfo?.initCheckedArr || [];
	}, [request, systemType]);

	const handlePlanDetails = useCallback(async () => {
		const res = await request(`/pricing_plan/${rowData?.id}`);
		return mapPlanResourcesToForm(res);
	}, [request, rowData?.id]);

	const handleClose = () => {
		setOpen(false);
		form.resetFields();
		setIsFold(false);
		setIsSelectAll(null);
	};

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const { resources } = await form.validateFields();
			const resourceIds = resources?.map((item) => item.resource_id) || [];
			await request(`/pricing_plan/${rowData.id}/bind_resources`, {
				method: 'POST',
				body: { resource_ids: resourceIds },
			});
			messageApi.success('Menu permission saved successfully.');
			successCallback?.();
			handleClose();
		} catch {
			// 校验失败或请求异常
		} finally {
			setConfirmLoading(false);
		}
	};

	useDeferredSyncEffect({
		predicate: () => open && !!rowData?.id && !!systemType,
		deferred: async () => {
			try {
				setLoading(true);
				const [defaultChecked, boundResources] = await Promise.all([
					handleMenuDetails(),
					handlePlanDetails(),
				]);
				form.setFieldsValue({
					resources: boundResources.length ? boundResources : defaultChecked,
				});
			} catch {
				// 请求失败时保留空表单
			} finally {
				setLoading(false);
			}
		},
		deps: [open, rowData?.id, systemType, handleMenuDetails, handlePlanDetails, form],
		delay: 0,
	});

	return (
		<Drawer
			title="Menu permission"
			placement="right"
			width={700}
			open={open}
			onClose={handleClose}
			destroyOnHidden
			extra={
				<Space>
					<Button onClick={handleClose}>Cancel</Button>
					<Button type="primary" onClick={handleSave} loading={confirmLoading}>
						Save
					</Button>
				</Space>
			}
		>
			{contextHolder}
			<Spin spinning={loading}>
				<Form
					form={form}
					name="pricing-plan-menu-permission"
					layout="vertical"
					autoComplete="off"
					size="large"
					disabled={confirmLoading || loading}
					scrollToFirstError
					className="h-full overflow-hidden overflow-y-auto"
				>
					<div className="relative">
						<Space className="absolute top-0 right-0 z-10">
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
							label="Menu permission"
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
								blockNode={true}
							/>
						</Form.Item>
					</div>
				</Form>
			</Spin>
		</Drawer>
	);
};

export default MenuPermissionDrawer;
