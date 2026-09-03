/**
 * 任务类型列表：GET `/project_management/task_type/all` 无分页（量少一次拉全量），
 * 返回 FilterTable 所需 `{ list, total }`；新建 / 编辑 / 删除（按 PROJECT_TASK_TYPE 权限）。
 */
'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { Button, Empty, Modal, Space, Tag, message } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { initConfigFn, PermissionSlot } from '@repo/shared-methods';
import { useAuth } from '@repo/shared-methods/hooks';
import { TASK_TYPE_ALL_API, createTaskTypeApi } from '@repo/project-management/api/taskType';
import { getTaskTypeColorLabel, toTaskTypeColorHex } from '@repo/project-management/enum/taskType';

export function useTaskTypeList(PERMISSIONS) {
	const t = useTranslations();
	const { hasPermission, getActiveByPermission } = useAuth();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const tableFilterRef = useRef();

	const [modalOpen, setModalOpen] = useState(false);
	const [editingType, setEditingType] = useState(null);
	const [saving, setSaving] = useState(false);
	const [messageApi, messageContextHolder] = message.useMessage();
	const [modal, modalContextHolder] = Modal.useModal();

	const refreshList = useCallback(() => {
		cancelMatch(TASK_TYPE_ALL_API);
		tableFilterRef.current?.handleSearch?.();
	}, [cancelMatch]);

	/**
	 * FilterTable 仍会传入 page/limit，此处忽略；全量 list 一次返回。
	 * @returns {Promise<{ list: object[], total: number }>}
	 */
	const getData = useCallback(async () => {
		try {
			const data = await taskTypeApi.list();
			const list = Array.isArray(data?.list) ? data.list : [];
			return { list, total: list.length };
		} catch {
			messageApi.error(t('taskType.loadFailed'));
			return { list: [], total: 0 };
		}
	}, [messageApi, t, taskTypeApi]);

	const openCreateModal = useCallback(() => {
		setEditingType(null);
		setModalOpen(true);
	}, []);

	const openEditModal = useCallback((record) => {
		setEditingType(record);
		setModalOpen(true);
	}, []);

	const closeModal = useCallback(() => {
		setModalOpen(false);
		setEditingType(null);
	}, []);

	const handleSave = useCallback(
		async (values) => {
			const body = {
				name: values.name,
				description: values.description || '',
				color: values.color,
				sort_order: Number(values.sort_order ?? 0),
			};
			try {
				setSaving(true);
				if (editingType) {
					await taskTypeApi.update(editingType.id, body);
					messageApi.success(t('taskType.updateSuccess'));
				} else {
					await taskTypeApi.create(body);
					messageApi.success(t('taskType.createSuccess'));
				}
				closeModal();
				refreshList();
			} catch (err) {
				messageApi.error(
					err?.message ||
						(editingType ? t('taskType.updateFailed') : t('taskType.createFailed')),
				);
			} finally {
				setSaving(false);
			}
		},
		[closeModal, editingType, messageApi, refreshList, t, taskTypeApi],
	);

	const handleDelete = useCallback(
		(record) => {
			modal.confirm({
				...initConfigFn({
					title: t('taskType.deleteConfirmTitle'),
					content: t('taskType.deleteConfirmContent', { name: record.name }),
					okText: t('common.delete'),
					cancelText: t('common.cancel'),
					okButtonProps: { danger: true },
					onOk: async () => {
						try {
							await taskTypeApi.delete(record.id);
							messageApi.success(t('taskType.deleteSuccess'));
							refreshList();
						} catch (err) {
							messageApi.error(err?.message || t('taskType.deleteFailed'));
						}
					},
				}),
			});
		},
		[messageApi, modal, refreshList, t, taskTypeApi],
	);

	/**
	 * 行拖到 toIndex 后，对该行 PUT `/:id/order`，sort_order 用 1-based 新位置。
	 * 表格已乐观更新顺序；成功/失败都 refreshList，失败时拉回服务端顺序。
	 */
	const handleSortEnd = useCallback(
		async (nextList, { toIndex } = {}) => {
			const record = nextList?.[toIndex];
			if (record?.id == null || toIndex == null) return;
			try {
				await taskTypeApi.updateOrder(record.id, { sort_order: toIndex + 1 });
				refreshList();
			} catch (err) {
				messageApi.error(err?.message || t('taskType.orderFailed'));
				refreshList();
			}
		},
		[messageApi, refreshList, t, taskTypeApi],
	);

	const columns = useMemo(
		() => [
			{
				title: t('taskType.nameColumn'),
				dataIndex: 'name',
				key: 'name',
				render: (name) => displayValue(name),
			},
			{
				title: t('taskType.descriptionColumn'),
				dataIndex: 'description',
				key: 'description',
				render: (description) => displayValue(description),
			},
			{
				title: t('taskType.colorColumn'),
				dataIndex: 'color',
				key: 'color',
				width: 140,
				render: (color) => {
					const hex = toTaskTypeColorHex(color);
					const label = getTaskTypeColorLabel(color);
					if (!hex) return displayValue(color);
					return <Tag color={hex}>{label}</Tag>;
				},
			},
			{
				title: t('taskType.sortOrderColumn'),
				dataIndex: 'sort_order',
				key: 'sort_order',
				width: 100,
				render: (sort_order) => displayValue(sort_order),
			},
			getActiveByPermission({
				title: t('taskType.actionColumn'),
				key: 'action',
				width: 180,
				permissionList: [
					PERMISSIONS?.PROJECT_TASK_TYPE?.EDIT,
					PERMISSIONS?.PROJECT_TASK_TYPE?.DELETE,
				],
				render: (_, record) => (
					<Space size="middle">
						<PermissionSlot code={PERMISSIONS?.PROJECT_TASK_TYPE?.EDIT}>
							<Button
								type="link"
								icon={<EditOutlined />}
								onClick={() => openEditModal(record)}
							>
								{t('taskType.edit')}
							</Button>
						</PermissionSlot>
						<PermissionSlot code={PERMISSIONS?.PROJECT_TASK_TYPE?.DELETE}>
							<Button
								type="link"
								danger
								icon={<DeleteOutlined />}
								onClick={() => handleDelete(record)}
							>
								{t('taskType.delete')}
							</Button>
						</PermissionSlot>
					</Space>
				),
			}),
		],
		[getActiveByPermission, handleDelete, openEditModal, PERMISSIONS, t],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'button',
				text: t('taskType.createTitle'),
				icon: <PlusOutlined />,
				onClick: openCreateModal,
				hidden: !hasPermission(PERMISSIONS?.PROJECT_TASK_TYPE?.ADD),
			},
		],
		[hasPermission, openCreateModal, PERMISSIONS, t],
	);

	return {
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
		emptyText: <Empty description={t('taskType.empty')} />,
		handleSortEnd,
	};
}
