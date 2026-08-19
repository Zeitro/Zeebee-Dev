/**
 * 工作流列表：名称筛选、拉取/分页、新建/设默认/删除。
 * GET /workflows 无分页协议，在 getData 内按名称过滤后再切页，返回 FilterTable 所需 { list, total }。
 */
'use client';

import { useCallback, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Empty, Modal, Space, Tag, message } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest, displayValue, EMPTY_DISPLAY } from '@repo/shared-methods/utils';
import { WORKFLOW_API, createWorkflowApi } from '@repo/project-management/api/workflow';
import { getWorkflowStateColor } from '@repo/project-management/enum/workflow';

/**
 * @param {unknown} data
 * @param {string} [name]
 * @param {number} [page]
 * @param {number} [limit]
 * @returns {{ list: object[], total: number }}
 */
function toWorkflowPage(data, name, page = 1, limit = 20) {
	const all = Array.isArray(data) ? data : [];
	const keyword = String(name || '')
		.trim()
		.toLowerCase();
	const filtered = keyword
		? all.filter((wf) =>
				String(wf.name || '')
					.toLowerCase()
					.includes(keyword),
			)
		: all;
	const start = (Math.max(page, 1) - 1) * limit;
	return {
		list: filtered.slice(start, start + limit),
		total: filtered.length,
	};
}

export function useWorkflowList() {
	const t = useTranslations();
	const router = useRouter();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const workflowApi = useMemo(() => createWorkflowApi(request), [request]);
	const userId = useUserStore((state) => state.userInfo?.id);
	const tableFilterRef = useRef();

	const [name, setName] = useState('');
	const [modalOpen, setModalOpen] = useState(false);
	const [saving, setSaving] = useState(false);
	const [messageApi, messageContextHolder] = message.useMessage();
	const [modal, modalContextHolder] = Modal.useModal();

	const refreshList = useCallback(() => {
		cancelMatch(WORKFLOW_API);
		tableFilterRef.current?.handleSearch?.();
	}, [cancelMatch]);

	const getData = useCallback(
		async (params) => {
			try {
				const data = await workflowApi.list();
				return toWorkflowPage(data, params?.filter?.name, params?.page, params?.limit);
			} catch {
				messageApi.error(t('workflow.loadFailed'));
				return { list: [], total: 0 };
			}
		},
		[messageApi, t, workflowApi],
	);

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
					use_template: true,
					user_id: userId || 0,
				});
				messageApi.success(t('workflow.createSuccess'));
				setModalOpen(false);
				refreshList();
			} catch (err) {
				messageApi.error(err?.message || t('workflow.createFailed'));
			} finally {
				setSaving(false);
			}
		},
		[messageApi, refreshList, t, userId, workflowApi],
	);

	const handleEdit = useCallback(
		(wf) => {
			router.push(`/project/workflow/${wf.id}/edit`);
		},
		[router],
	);

	const handleSetDefault = useCallback(
		async (wf) => {
			try {
				await workflowApi.update(wf.id, {
					name: wf.name,
					description: wf.description || '',
					is_default: true,
					user_id: userId || 0,
				});
				messageApi.success(t('workflow.setDefaultSuccess', { name: wf.name }));
				refreshList();
			} catch (err) {
				messageApi.error(err?.message || t('workflow.setDefaultFailed'));
			}
		},
		[messageApi, refreshList, t, userId, workflowApi],
	);

	const handleDelete = useCallback(
		(wf) => {
			modal.confirm({
				title: t('workflow.deleteConfirmTitle'),
				content: t('workflow.deleteConfirmContent', { name: wf.name }),
				okText: t('common.delete'),
				cancelText: t('common.cancel'),
				okButtonProps: { danger: true },
				onOk: async () => {
					try {
						await workflowApi.delete(wf.id, userId || 0);
						messageApi.success(t('workflow.deleteSuccess'));
						refreshList();
					} catch (err) {
						messageApi.error(err?.message || t('workflow.deleteFailed'));
					}
				},
			});
		},
		[messageApi, modal, refreshList, t, userId, workflowApi],
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
				key: 'states',
				render: (_, wf) => {
					const list = Array.isArray(wf.states) ? wf.states : [];
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
			{
				title: t('workflow.actionColumn'),
				key: 'action',
				width: 280,
				render: (_, wf) => (
					<Space size="middle">
						{!wf.is_default ? (
							<Button type="link" onClick={() => handleSetDefault(wf)}>
								{t('workflow.setDefault')}
							</Button>
						) : null}
						<Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(wf)}>
							{t('workflow.edit')}
						</Button>
						<Button
							type="link"
							danger
							icon={<DeleteOutlined />}
							onClick={() => handleDelete(wf)}
						>
							{t('workflow.delete')}
						</Button>
					</Space>
				),
			},
		],
		[handleDelete, handleEdit, handleSetDefault, t],
	);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: t('workflow.searchPlaceholder'),
				paramsKey: 'name',
				value: name,
				onChange: ({ target: { value } }) => setName(value),
			},
			{
				filterType: 'button',
				text: t('common.search'),
				onClick: refreshList,
			},
			{
				filterType: 'button',
				type: 'default',
				text: t('common.reset'),
				icon: <ReloadOutlined />,
				onClick: () => {
					cancelMatch(WORKFLOW_API);
					setName('');
					tableFilterRef.current?.handleSearch?.({
						filter: { name: '' },
					});
				},
				colProps: { className: 'flex-1' },
			},
			{
				filterType: 'button',
				text: t('workflow.createTitle'),
				icon: <PlusOutlined />,
				onClick: openCreateModal,
			},
		],
		[cancelMatch, name, openCreateModal, refreshList, t],
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
