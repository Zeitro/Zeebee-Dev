/**
 * 任务类型列表：拉取/分页、新建/编辑/删除。
 * GET /task-types 无分页协议，在 getData 内切页，返回 FilterTable 所需 { list, total }。
 */
'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { Button, Empty, Modal, Space, Tag, message } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { TASK_TYPES_API, createTaskTypeApi } from '@repo/project-management/api/taskType';
import { getTaskTypeColorLabel, toTaskTypeColorHex } from '@repo/project-management/enum/taskType';

/**
 * @param {unknown} data
 * @param {number} [page]
 * @param {number} [limit]
 * @returns {{ list: object[], total: number }}
 */
function toTaskTypePage(data, page = 1, limit = 20) {
	const all = Array.isArray(data) ? data : Array.isArray(data?.list) ? data.list : [];
	const start = (Math.max(page, 1) - 1) * limit;
	return {
		list: all.slice(start, start + limit),
		total: all.length,
	};
}

export function useTaskTypeList() {
	const t = useTranslations();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const taskTypeApi = useMemo(() => createTaskTypeApi(request), [request]);
	const tableFilterRef = useRef();

	const [modalOpen, setModalOpen] = useState(false);
	const [editingType, setEditingType] = useState(null);
	const [saving, setSaving] = useState(false);
	const [messageApi, messageContextHolder] = message.useMessage();
	const [modal, modalContextHolder] = Modal.useModal();

	const refreshList = useCallback(() => {
		cancelMatch(TASK_TYPES_API);
		tableFilterRef.current?.handleSearch?.();
	}, [cancelMatch]);

	const getData = useCallback(
		async (params) => {
			try {
				const data = await taskTypeApi.list();
				return toTaskTypePage(data, params?.page, params?.limit);
			} catch {
				messageApi.error(t('taskType.loadFailed'));
				return { list: [], total: 0 };
			}
		},
		[messageApi, t, taskTypeApi],
	);

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
			try {
				setSaving(true);
				if (editingType) {
					await taskTypeApi.update(editingType.id, values);
					messageApi.success(t('taskType.updateSuccess'));
				} else {
					await taskTypeApi.create(values);
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
			});
		},
		[messageApi, modal, refreshList, t, taskTypeApi],
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
				title: t('taskType.actionColumn'),
				key: 'action',
				width: 180,
				render: (_, record) => (
					<Space size="middle">
						<Button
							type="link"
							icon={<EditOutlined />}
							onClick={() => openEditModal(record)}
						>
							{t('taskType.edit')}
						</Button>
						<Button
							type="link"
							danger
							icon={<DeleteOutlined />}
							onClick={() => handleDelete(record)}
						>
							{t('taskType.delete')}
						</Button>
					</Space>
				),
			},
		],
		[handleDelete, openEditModal, t],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'button',
				text: t('taskType.createTitle'),
				icon: <PlusOutlined />,
				onClick: openCreateModal,
			},
		],
		[openCreateModal, t],
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
	};
}
