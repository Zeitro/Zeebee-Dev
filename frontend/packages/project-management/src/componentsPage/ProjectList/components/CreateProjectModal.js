/**
 * 创建项目弹窗：NewModal + Form；POST `/project_management/project`，
 * 工作流下拉来自 GET `/project_management/workflow/all`，成功后回调 handleSearch 刷新列表。
 */
'use client';

import { useEffect, useMemo, useState } from 'react';
import { Form, Input, Select } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { DEFAULT_WORKFLOW_TEMPLATE_ID, PROJECT_API } from '../../../api/project';
import { WORKFLOW_ALL_API } from '../../../api/workflow';
// import {
// 	PROJECT_PRIORITY,
// 	PROJECT_PRIORITY_FILTER_VALUES,
// 	PROJECT_PRIORITY_I18N_KEY,
// } from '../../../enum/project';

/**
 * 下拉文案：工作流名称 + stages.name（按 sort_order），便于区分模板。
 * @param {object} wf
 * @returns {string}
 */
function formatWorkflowLabel(wf) {
	const title = wf?.name || '';
	const stages = Array.isArray(wf?.stages) ? [...wf.stages] : [];
	const stageNames = stages
		.sort((a, b) => (Number(a?.sort_order) || 0) - (Number(b?.sort_order) || 0))
		.map((stage) => stage?.name)
		.filter(Boolean);
	if (title && stageNames.length) {
		return `${title} (${stageNames.join(' → ')})`;
	}
	if (stageNames.length) return stageNames.join(' → ');
	return displayValue(title);
}

/**
 * GET workflow/all 的 data → Select options；默认选 is_default，否则第一项。
 * @param {unknown} data
 * @returns {{ options: Array<{value: number, label: string}>, defaultId: number|undefined }}
 */
function toWorkflowSelect(data) {
	const list = Array.isArray(data?.list) ? data.list : [];
	const options = list.map((wf) => ({
		value: wf.id,
		label: formatWorkflowLabel(wf),
	}));
	const defaultId = list.find((wf) => wf.is_default)?.id ?? list[0]?.id;
	return { options, defaultId };
}

/**
 * @param {{ setOpen: (open: boolean) => void, handleSearch?: () => void }} props
 */
export default function CreateProjectModal({ messageApi, setOpen, handleSearch }) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [workflowOptions, setWorkflowOptions] = useState([]);
	const [workflowLoading, setWorkflowLoading] = useState(true);

	// const priorityOptions = useMemo(
	// 	() =>
	// 		PROJECT_PRIORITY_FILTER_VALUES.map((value) => ({
	// 			value: Number(value),
	// 			label: t(PROJECT_PRIORITY_I18N_KEY[value]),
	// 		})),
	// 	[t],
	// );

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const data = await request(WORKFLOW_ALL_API, {
					method: 'GET',
					isHideTips: true,
				});
				if (cancelled) return;
				const { options, defaultId } = toWorkflowSelect(data);
				setWorkflowOptions(options);
				if (defaultId != null) {
					form.setFieldsValue({ workflow_template_id: defaultId });
				}
			} catch {
				if (!cancelled) setWorkflowOptions([]);
			} finally {
				if (!cancelled) setWorkflowLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [form, request]);

	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setConfirmLoading(true);
			const created = await request(PROJECT_API, {
				method: 'POST',
				body: {
					name: String(values.name || '').trim(),
					description: values.description || '',
					cover_url: '',
					entity_id: 0,
					entity_type: '',
					workflow_template_id:
						values.workflow_template_id ?? DEFAULT_WORKFLOW_TEMPLATE_ID,
					// priority: values.priority,
				},
			});
			if (created == null) return;
			messageApi.success(t('projectList.createSuccess'));
			handleSearch?.();
			setOpen(false);
		} catch {
			// 表单校验失败，或 createRequest 已提示接口错误
		} finally {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width={750}
			title={t('projectList.createTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
			destroyOnHidden
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				// initialValues={{ priority: PROJECT_PRIORITY.MEDIUM }}
			>
				<Form.Item name="name" label={t('projectList.name')} rules={[{ required: true }]}>
					<Input placeholder={t('projectList.namePlaceholder')} />
				</Form.Item>
				<Form.Item name="description" label={t('projectList.description')}>
					<Input.TextArea
						rows={3}
						placeholder={t('projectList.descriptionPlaceholder')}
					/>
				</Form.Item>
				<Form.Item
					name="workflow_template_id"
					label={t('projectList.workflow')}
					rules={[{ required: true }]}
				>
					<Select
						showSearch
						loading={workflowLoading}
						optionFilterProp="label"
						placeholder={t('projectList.workflowPlaceholder')}
						options={workflowOptions}
					/>
				</Form.Item>
				{/* <Form.Item
						name="priority"
						label={t('projectList.priority')}
						rules={[{ required: true }]}
					>
						<Select options={priorityOptions} />
					</Form.Item> */}
			</Form>
		</NewModal>
	);
}
