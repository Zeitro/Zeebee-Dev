/**
 * 创建项目弹窗：NewModal + Form；成功后回调 handleSearch 刷新列表。
 */
'use client';

import { useMemo, useState } from 'react';
import { Form, Input, Select, message } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import { PROJECT_API } from '../../../api/project';
import {
	PROJECT_PRIORITY,
	PROJECT_PRIORITY_FILTER_VALUES,
	PROJECT_PRIORITY_I18N_KEY,
} from '../../../enum/project';

/**
 * @param {{ setOpen: (open: boolean) => void, handleSearch?: () => void }} props
 */
export default function CreateProjectModal({ setOpen, handleSearch }) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);

	const priorityOptions = useMemo(
		() =>
			PROJECT_PRIORITY_FILTER_VALUES.map((value) => ({
				value: Number(value),
				label: t(PROJECT_PRIORITY_I18N_KEY[value]),
			})),
		[t],
	);

	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setConfirmLoading(true);
			await request(PROJECT_API, {
				method: 'POST',
				body: {
					name: values.name,
					description: values.description,
					priority: values.priority,
				},
			});
			message.success(t('projectList.createSuccess'));
			handleSearch?.();
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width={520}
			title={t('projectList.createTitle')}
			okText={t('common.save')}
			cancelText={t('common.cancel')}
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				initialValues={{ priority: PROJECT_PRIORITY.MEDIUM }}
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
					name="priority"
					label={t('projectList.priority')}
					rules={[{ required: true }]}
				>
					<Select options={priorityOptions} />
				</Form.Item>
			</Form>
		</NewModal>
	);
}
