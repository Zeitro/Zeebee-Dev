/**
 * 项目设置：封面 / 名称 / 描述更新，以及危险区删除；读路由 id，不写入全局 store。
 */
'use client';

import { useCallback, useState } from 'react';
import { CameraOutlined } from '@ant-design/icons';
import { Button, Form, Input, Modal, message } from 'antd';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { NoData, UploadFileButton } from '@repo/shared-methods';
import { deleteMockProject, getMockProject, updateMockProject } from '../mockData';

const COVER_ACCEPT = '.png,.jpg,.jpeg,.webp,.gif';

/**
 * @param {{ projectId: number|string }} props
 */
function ProjectSettingContent({ projectId }) {
	const t = useTranslations();
	const router = useRouter();
	const [form] = Form.useForm();
	const [saving, setSaving] = useState(false);
	const [project, setProject] = useState(() => getMockProject(projectId));
	const [coverPreview, setCoverPreview] = useState();

	const goList = useCallback(() => {
		router.push('/project/projectList');
	}, [router]);

	const handleCoverListChange = useCallback((fileList) => {
		const last = fileList[fileList.length - 1];
		if (last?.thumbUrl) setCoverPreview(last.thumbUrl);
	}, []);

	const handleUpdate = async (values) => {
		const cover_image = coverPreview ?? project?.cover_image ?? '';
		try {
			setSaving(true);
			// MOCK: start — 接口就绪后改为 PUT getProjectApi(projectId)
			const next = updateMockProject(projectId, {
				name: values.name,
				description: values.description,
				cover_image,
			});
			if (!next) throw new Error('not found');
			setProject(next);
			message.success(t('projectSetting.updateSuccess'));
		} catch {
			message.error(t('projectSetting.updateFailed'));
		} finally {
			setSaving(false);
		}
	};

	const handleDelete = () => {
		Modal.confirm({
			title: t('projectSetting.deleteConfirmTitle'),
			content: t('projectSetting.deleteConfirmContent', { name: project?.name }),
			okText: t('projectSetting.delete'),
			okType: 'danger',
			cancelText: t('common.cancel'),
			onOk: async () => {
				try {
					// MOCK: start — 接口就绪后改为 DELETE getProjectApi(projectId)
					if (!deleteMockProject(projectId)) throw new Error('not found');
					message.success(t('projectSetting.deleteSuccess'));
					goList();
				} catch {
					message.error(t('projectSetting.deleteFailed'));
					return Promise.reject();
				}
			},
		});
	};

	const coverImage = coverPreview ?? project?.cover_image;

	return (
		<div className="flex-1 overflow-hidden bg-gray-light px-8 pt-6 pb-10">
			<div className="max-w-2xl mx-auto">
				<div className="relative flex items-center justify-center mb-6">
					<h1 className="m-0 text-xl font-semibold text-black-dark">
						{t('projectSetting.title')}
					</h1>
				</div>
				{!project ? (
					<NoData text={t('projectSetting.notFound')} />
				) : (
					<>
						<div
							className="relative h-[200px] rounded-xl overflow-hidden bg-gradient-to-r from-[var(--purple-dark)] to-blue-dark bg-cover bg-center"
							style={
								coverImage
									? {
											backgroundImage: `url(${coverImage}), linear-gradient(135deg, var(--purple-dark), var(--blue-dark))`,
											backgroundPosition: 'top',
										}
									: undefined
							}
						>
							<div className="absolute top-3 right-3">
								<UploadFileButton
									accept={COVER_ACCEPT}
									showUploadList={false}
									immediateUpload={false}
									maxCount={1}
									onGetUploadList={handleCoverListChange}
								>
									<Button size="small" icon={<CameraOutlined />}>
										{t('projectSetting.changeCover')}
									</Button>
								</UploadFileButton>
							</div>
						</div>
						<Form
							form={form}
							layout="vertical"
							colon={false}
							className="!pt-8"
							disabled={saving}
							initialValues={{
								name: project.name,
								description: project.description,
							}}
							onFinish={handleUpdate}
						>
							<Form.Item
								name="name"
								label={t('projectList.name')}
								rules={[{ required: true }]}
							>
								<Input placeholder={t('projectList.namePlaceholder')} />
							</Form.Item>
							<Form.Item name="description" label={t('projectList.description')}>
								<Input.TextArea
									rows={4}
									placeholder={t('projectList.descriptionPlaceholder')}
								/>
							</Form.Item>
							<Button type="primary" htmlType="submit" loading={saving}>
								{t('projectSetting.update')}
							</Button>
						</Form>
						<div className="mt-10 rounded-lg border border-[var(--error-medium)] bg-[var(--error-light)] px-4 py-4">
							<p className="m-0 font-medium text-base-error">
								{t('projectSetting.deleteTitle')}
							</p>
							<p className="mt-2 mb-4 text-sm text-base-error">
								{t('projectSetting.deleteWarning')}
							</p>
							<Button
								color="danger"
								variant="outlined"
								onClick={handleDelete}
								disabled={saving}
							>
								{t('projectSetting.delete')}
							</Button>
						</div>
					</>
				)}
			</div>
		</div>
	);
}

/**
 * 项目设置页：由 `projectId` 或路由 params 的 `id` 驱动，内部 key remount 重置表单。
 */
export function ProjectSetting({ projectId: projectIdProp } = {}) {
	const searchParams = useSearchParams();
	const id = searchParams.get('id');
	const projectId = projectIdProp ?? id ?? null;
	return <ProjectSettingContent key={projectId} projectId={projectId} />;
}
