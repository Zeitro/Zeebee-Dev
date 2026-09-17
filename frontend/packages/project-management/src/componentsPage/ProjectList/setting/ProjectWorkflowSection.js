/**
 * 项目工作流：展示详情里的 workflow（或 workflows），可编辑、设默认、删除，并从组织工作流复制挂入。
 * 编辑跳转已有编辑器；挂入 / 设默认走项目维度接口，不改组织工作流列表。
 */
'use client';

import { useCallback, useMemo, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Select, Space, Tag } from 'antd';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { DeleteBtn, NewModal } from '@repo/shared-methods';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import {
	WORKFLOW_ALL_API,
	getProjectDefaultWorkflowApi,
	getProjectWorkflowAttachApi,
	getWorkflowApi,
} from '../../../api/workflow';

/**
 * 详情目前是单个 workflow；若接口以后返回 workflows 数组则优先用它。
 * @param {object|null|undefined} project
 * @returns {object[]}
 */
function listProjectWorkflows(project) {
	if (Array.isArray(project?.workflows) && project.workflows.length) return project.workflows;
	return project?.workflow ? [project.workflow] : [];
}

/**
 * @param {object} workflow
 * @returns {number}
 */
function getStateCount(workflow) {
	if (Array.isArray(workflow?.stages)) return workflow.stages.length;
	if (Array.isArray(workflow?.states)) return workflow.states.length;
	return 0;
}

/**
 * @param {{
 *   project: object,
 *   projectId: number|string,
 *   messageApi: { success: (msg: string) => void, warning: (msg: string) => void },
 *   onRefresh: () => void
 * }} props
 */
export function ProjectWorkflowSection({ project, projectId, messageApi, onRefresh }) {
	const t = useTranslations();
	const router = useRouter();
	const { request } = useMemo(() => createRequest(), []);
	const [modalOpen, setModalOpen] = useState(false);
	const [selectedIds, setSelectedIds] = useState([]);
	const [tenantWorkflows, setTenantWorkflows] = useState([]);
	const [optionsLoading, setOptionsLoading] = useState(false);
	const [saving, setSaving] = useState(false);

	const projectWorkflows = listProjectWorkflows(project);
	const defaultWorkflowId =
		project?.default_workflow_id ?? projectWorkflows.find((row) => row.is_default)?.id;

	const openAddModal = useCallback(async () => {
		setModalOpen(true);
		setSelectedIds([]);
		setOptionsLoading(true);
		try {
			const data = await request(WORKFLOW_ALL_API, { method: 'GET', isHideTips: true });
			setTenantWorkflows(Array.isArray(data?.list) ? data.list : []);
		} catch {
			setTenantWorkflows([]);
		} finally {
			setOptionsLoading(false);
		}
	}, [request]);

	const closeAddModal = () => {
		setModalOpen(false);
		setSelectedIds([]);
	};

	const handleAdd = async () => {
		if (!selectedIds.length) {
			messageApi.warning(t('projectSetting.addWorkflowRequired'));
			return;
		}
		try {
			setSaving(true);
			const added = await request(getProjectWorkflowAttachApi(projectId), {
				method: 'POST',
				body: { workflow_ids: selectedIds },
			});
			if (added == null) return;
			messageApi.success(t('projectSetting.addWorkflowSuccess'));
			closeAddModal();
			onRefresh();
		} catch {
			// createRequest 已提示接口错误
		} finally {
			setSaving(false);
		}
	};

	const handleSetDefault = async (workflowId) => {
		try {
			const updated = await request(getProjectDefaultWorkflowApi(projectId), {
				method: 'PUT',
				body: { workflow_id: workflowId },
			});
			if (updated == null) return;
			messageApi.success(t('projectSetting.setDefaultWorkflowSuccess'));
			onRefresh();
		} catch {
			// createRequest 已提示接口错误
		}
	};

	const handleDelete = async (workflowId) => {
		try {
			await request(getWorkflowApi(workflowId), { method: 'DELETE' });
			messageApi.success(t('workflow.deleteSuccess'));
			onRefresh();
			const remaining = projectWorkflows.filter((row) => row.id !== workflowId);
			if (workflowId === defaultWorkflowId && remaining[0]?.id != null) {
				await request(getProjectDefaultWorkflowApi(projectId), {
					method: 'PUT',
					body: { workflow_id: remaining[0].id },
				});
				onRefresh();
			}
		} catch {
			// createRequest 已提示接口错误
		}
	};

	const options = tenantWorkflows
		.filter(
			(row) =>
				!projectWorkflows.some(
					(attached) => attached.id === row.id || attached.name === row.name,
				),
		)
		.map((row) => ({
			value: row.id,
			label: `${row.name}${row.is_default ? t('projectSetting.workflowOrgDefault') : ''}`,
		}));

	return (
		<div className="mt-10">
			<div className="mb-1 flex items-center justify-between">
				<h4 className="m-0 text-base font-medium text-black-dark">
					{t('projectSetting.workflowsTitle')}
				</h4>
				<Button icon={<PlusOutlined />} onClick={openAddModal}>
					{t('projectSetting.addWorkflow')}
				</Button>
			</div>
			<p className="mt-1 mb-3 text-xs text-base-gray">{t('projectSetting.workflowsHint')}</p>
			{projectWorkflows.length === 0 ? (
				<div className="py-2 text-[13px] text-base-gray">
					{t('projectSetting.workflowsEmpty')}
				</div>
			) : (
				projectWorkflows.map((row) => (
					<div
						key={row.id}
						className="mb-2 flex items-center justify-between rounded-lg border border-base-border px-3 py-2.5"
					>
						<div>
							<span className="font-medium">{displayValue(row.name)}</span>
							{row.id === defaultWorkflowId ? (
								<Tag color="blue" className="!ml-2">
									{t('projectSetting.workflowDefault')}
								</Tag>
							) : null}
							<span className="ml-2 text-xs text-base-gray">
								{t('projectSetting.workflowStateCount', { n: getStateCount(row) })}
							</span>
						</div>
						<Space size={4}>
							{/* <Button
								size="small"
								onClick={() =>
									router.push(
										`/project/workflow/${row.id}?projectId=${projectId}`,
									)
								}
							>
								{t('workflow.edit')}
							</Button> */}
							<DeleteBtn
								type="default"
								size="small"
								title={t('workflow.deleteConfirmTitle')}
								content={t('workflow.deleteConfirmContent', { name: row.name })}
								okText={t('workflow.delete')}
								cancelText={t('common.cancel')}
								onOk={() => handleDelete(row.id)}
							>
								{t('workflow.delete')}
							</DeleteBtn>
						</Space>
					</div>
				))
			)}
			{modalOpen ? (
				<NewModal
					title={t('projectSetting.addWorkflowTitle')}
					okText={t('projectSetting.addWorkflowOk')}
					cancelText={t('common.cancel')}
					confirmLoading={saving}
					onOk={handleAdd}
					onCancel={closeAddModal}
					destroyOnHidden
				>
					<Select
						mode="multiple"
						allowClear
						showSearch
						optionFilterProp="label"
						loading={optionsLoading}
						placeholder={t('projectSetting.workflowPlaceholder')}
						className="w-full"
						value={selectedIds}
						onChange={(value) => setSelectedIds(value)}
						options={options}
						size="large"
					/>
					<p className="mt-2 mb-4 text-xs text-base-gray">
						{t('projectSetting.addWorkflowHint')}
					</p>
				</NewModal>
			) : null}
		</div>
	);
}
