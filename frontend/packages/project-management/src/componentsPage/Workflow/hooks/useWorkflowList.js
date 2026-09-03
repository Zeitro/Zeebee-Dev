/**
 * 工作流列表：GET `/project_management/workflow/all` 无分页（量少一次拉全量），
 * 返回 FilterTable 所需 `{ list, total }`；新建 / 编辑 / 删除 / 设默认（按 PROJECT_WORKFLOW 权限）。
 */
'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Empty, Modal, Space, Tag, message } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { createRequest, displayValue, EMPTY_DISPLAY } from '@repo/shared-methods/utils';
import { initConfigFn, PermissionSlot } from '@repo/shared-methods';
import { useAuth } from '@repo/shared-methods/hooks';
import { WORKFLOW_ALL_API, createWorkflowApi } from '@repo/project-management/api/workflow';
import { getWorkflowStateColor } from '@repo/project-management/enum/workflow';
import { formatApiError } from '../../../lib/formatApiError';

/**
 * 按 sort_order 排 stages，供状态列展示。
 * @param {object} wf
 * @returns {object[]}
 */
function getWorkflowStages(wf) {
	const stages = Array.isArray(wf?.stages) ? [...wf.stages] : [];
	return stages.sort((a, b) => (Number(a?.sort_order) || 0) - (Number(b?.sort_order) || 0));
}

export function useWorkflowList(PERMISSIONS) {
	const t = useTranslations();
	const router = useRouter();
	const { hasPermission, getActiveByPermission } = useAuth();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const workflowApi = useMemo(() => createWorkflowApi(request), [request]);
	const tableFilterRef = useRef();

	const [modalOpen, setModalOpen] = useState(false);
	const [saving, setSaving] = useState(false);
	const [messageApi, messageContextHolder] = message.useMessage();
	const [modal, modalContextHolder] = Modal.useModal();

	const refreshList = useCallback(() => {
		cancelMatch(WORKFLOW_ALL_API);
		tableFilterRef.current?.handleSearch?.();
	}, [cancelMatch]);

	/**
	 * FilterTable 仍会传入 page/limit，此处忽略；全量 list 一次返回。
	 * @returns {Promise<{ list: object[], total: number }>}
	 */
	const getData = useCallback(async () => {
		try {
			const data = await workflowApi.list();
			const list = Array.isArray(data?.list) ? data.list : [];
			return { list, total: list.length };
		} catch (error) {
			messageApi.error(formatApiError(t, 'workflow.loadFailed', error));
			return { list: [], total: 0 };
		}
	}, [messageApi, t, workflowApi]);

	const openCreateModal = useCallback(() => {
		setModalOpen(true);
	}, []);

	const handleSave = useCallback(
		async ({ name: workflowName, description }) => {
			try {
				setSaving(true);
				await workflowApi.create({
					name: String(workflowName || '').trim(),
					description: description || '',
				});
				messageApi.success(t('workflow.createSuccess'));
				setModalOpen(false);
				refreshList();
			} catch (error) {
				messageApi.error(formatApiError(t, 'workflow.createFailed', error));
			} finally {
				setSaving(false);
			}
		},
		[messageApi, refreshList, t, workflowApi],
	);

	const handleEdit = useCallback(
		(wf) => {
			router.push(`/project/workflow/${wf.id}`);
		},
		[router],
	);

	const handleSetDefault = useCallback(
		async (wf) => {
			try {
				await workflowApi.setDefault(wf.id);
				messageApi.success(t('workflow.setDefaultSuccess', { name: wf.name }));
				refreshList();
			} catch (error) {
				messageApi.error(formatApiError(t, 'workflow.setDefaultFailed', error));
			}
		},
		[messageApi, refreshList, t, workflowApi],
	);

	const handleDelete = useCallback(
		(wf) => {
			modal.confirm({
				...initConfigFn({
					title: t('workflow.deleteConfirmTitle'),
					content: t('workflow.deleteConfirmContent', { name: wf.name }),
					okText: t('common.delete'),
					cancelText: t('common.cancel'),
					okButtonProps: { danger: true },
					onOk: async () => {
						try {
							await workflowApi.delete(wf.id);
							messageApi.success(t('workflow.deleteSuccess'));
							refreshList();
						} catch (error) {
							messageApi.error(formatApiError(t, 'workflow.deleteFailed', error));
						}
					},
				}),
			});
		},
		[messageApi, modal, refreshList, t, workflowApi],
	);

	const columns = useMemo(
		() => [
			{
				title: t('workflow.nameColumn'),
				key: 'name',
				render: (_, wf) => (
					<Space size="middle">
						<span className="font-medium">{displayValue(wf.name)}</span>
						{wf.is_default ? <Tag color="blue">{t('workflow.setDefault')}</Tag> : null}
						{wf.is_system ? <Tag color="orange">{t('kanban.system')}</Tag> : null}
					</Space>
				),
			},
			{
				title: t('workflow.statusColumn'),
				key: 'stages',
				render: (_, wf) => {
					const list = getWorkflowStages(wf);
					if (list.length === 0) return EMPTY_DISPLAY;
					return (
						<Space wrap size={8}>
							{list.map((s) => (
								<Tag key={String(s.id)} className="rounded-[10px]">
									<span
										className="inline-block w-2.5 h-2.5 rounded-full mr-2 align-middle"
										style={{ background: getWorkflowStateColor(s.color) }}
									/>
									{displayValue(s.name)}
								</Tag>
							))}
						</Space>
					);
				},
			},
			getActiveByPermission({
				title: t('workflow.actionColumn'),
				key: 'action',
				width: 280,
				permissionList: [
					PERMISSIONS?.PROJECT_WORKFLOW?.EDIT,
					PERMISSIONS?.PROJECT_WORKFLOW?.DELETE,
					PERMISSIONS?.PROJECT_WORKFLOW?.SET_DEFAULT,
				],
				render: (_, wf) => (
					<Space size="middle">
						<PermissionSlot code={PERMISSIONS?.PROJECT_WORKFLOW?.EDIT}>
							<Button
								type="link"
								icon={<EditOutlined />}
								onClick={() => handleEdit(wf)}
							>
								{t('workflow.edit')}
							</Button>
						</PermissionSlot>
						<PermissionSlot code={PERMISSIONS?.PROJECT_WORKFLOW?.DELETE}>
							<Button
								type="link"
								danger
								icon={<DeleteOutlined />}
								onClick={() => handleDelete(wf)}
							>
								{t('workflow.delete')}
							</Button>
						</PermissionSlot>
						{!wf.is_default ? (
							<PermissionSlot code={PERMISSIONS?.PROJECT_WORKFLOW?.SET_DEFAULT}>
								<Button type="link" onClick={() => handleSetDefault(wf)}>
									{t('workflow.setDefault')}
								</Button>
							</PermissionSlot>
						) : null}
					</Space>
				),
			}),
		],
		[getActiveByPermission, handleDelete, handleEdit, handleSetDefault, PERMISSIONS, t],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'button',
				text: t('workflow.createTitle'),
				icon: <PlusOutlined />,
				onClick: openCreateModal,
				hidden: !hasPermission(PERMISSIONS?.PROJECT_WORKFLOW?.ADD),
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
		setModalOpen,
		handleSave,
		emptyText: <Empty description={t('workflow.empty')} />,
	};
}
