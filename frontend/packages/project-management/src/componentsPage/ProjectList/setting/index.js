/**
 * 项目设置：封面 / 名称 / 描述 / 优先级更新，以及周期、另存模板、项目工作流和危险区删除。
 * 读路由 id，不写入全局 store。已联调的更新 / 删除请求体保持不变，周期字段附加在同一 PUT。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { CameraOutlined } from '@ant-design/icons';
import { Button, Form, Input, Select, Space, Spin, message } from 'antd';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { NoData, UploadFileButton } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import { getProjectDetailApi } from '../../../api/project';
import {
	getTaskPriorityI18nKey,
	localizeSelectOptions,
	toTaskPriorityValue,
} from '../../../enum/task';
import { useOptionsStore } from '../../../store/options';
import { toCycleFormValues, toCycleUpdateBody } from './cycle';
import { DeleteProjectSection } from './DeleteProjectSection';
import { ProjectCycleFields } from './ProjectCycleFields';
import { ProjectWorkflowSection } from './ProjectWorkflowSection';
import { SaveAsTemplateButton } from './SaveTemplateModal';

const COVER_ACCEPT = '.png,.jpg,.jpeg';

/**
 * @param {{ projectId: number|string }} props
 */
function ProjectSettingContent({ projectId }) {
	const t = useTranslations();
	const router = useRouter();
	const [form] = Form.useForm();
	const { request } = useMemo(() => createRequest(), []);
	const [messageApi, messageContextHolder] = message.useMessage();
	const [saving, setSaving] = useState(false);
	const [loading, setLoading] = useState(() => projectId != null && projectId !== '');
	const [project, setProject] = useState(null);
	const [coverPreview, setCoverPreview] = useState();
	const [reloadToken, setReloadToken] = useState(0);
	const coverUploadRef = useRef(null);

	const priorityOptionsRaw = useOptionsStore((state) => state.priorityOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);
	const [priorityLoading, setPriorityLoading] = useState(
		() => !useOptionsStore.getState().priorityOptions.length,
	);

	const priorityOptions = useMemo(
		() => localizeSelectOptions(priorityOptionsRaw, t, getTaskPriorityI18nKey),
		[priorityOptionsRaw, t],
	);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				await getOptions('project_management_priority', 'priorityOptions');
			} catch {
				// 下拉失败时保留空选项，不阻断设置页
			} finally {
				if (!cancelled) setPriorityLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [getOptions]);

	useEffect(() => {
		if (projectId == null || projectId === '') return undefined;
		let cancelled = false;
		(async () => {
			try {
				const data = await request(getProjectDetailApi(projectId), {
					method: 'GET',
					isHideTips: true,
				});
				if (cancelled) return;
				setProject(data || null);
			} catch {
				if (!cancelled) {
					setProject(null);
					messageApi.error(t('projectSetting.loadFailed'));
				}
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [projectId, request, t, messageApi, reloadToken]);

	const goList = useCallback(() => {
		router.push('/project/projectList');
	}, [router]);

	/** 工作流变更后重拉详情；不把整页打回 loading，避免已填表单被卸载。 */
	const refreshProject = useCallback(() => {
		setReloadToken((token) => token + 1);
	}, []);

	const handleCoverListChange = useCallback((fileList) => {
		const last = fileList[fileList.length - 1];
		setCoverPreview(last || undefined);
	}, []);

	/**
	 * 有新封面时先 triggerUpload 走公共 AWS 预签名上传，再用返回 url 作为 cover_url。
	 */
	const handleUpdate = async (values) => {
		try {
			setSaving(true);
			let cover_url = project?.cover_url ?? '';
			if (coverPreview?.originFileObj) {
				const urls = await coverUploadRef.current?.triggerUpload();
				const uploaded = urls?.[0];
				if (!uploaded) return;
				cover_url = typeof uploaded === 'string' ? uploaded : uploaded.url;
			}
			const next = await request(getProjectDetailApi(projectId), {
				method: 'PUT',
				body: {
					name: String(values.name || '').trim(),
					description: values.description || '',
					cover_url,
					priority: values.priority,
					...toCycleUpdateBody(values),
				},
			});
			if (next == null) return;
			setProject(next);
			setCoverPreview(undefined);
			coverUploadRef.current?.clearFileList();
			messageApi.success(t('projectSetting.updateSuccess'));
		} catch {
			// createRequest 已提示接口错误
		} finally {
			setSaving(false);
		}
	};

	const coverImage = coverPreview?.thumbUrl ?? project?.cover_url;

	return (
		<div className="flex-1 overflow-hidden bg-gray-light px-8 pt-6 pb-10">
			{messageContextHolder}
			<div className="max-w-2xl mx-auto">
				<div className="relative flex items-center justify-center mb-6">
					<h1 className="m-0 text-xl font-semibold text-black-dark">
						{t('projectSetting.title')}
					</h1>
				</div>
				{loading ? (
					<div className="flex items-center justify-center py-24">
						<Spin size="large" />
					</div>
				) : !project ? (
					<NoData text={t('projectSetting.notFound')} />
				) : (
					<>
						<div
							className="relative h-[200px] rounded-xl overflow-hidden bg-gradient-to-r from-[var(--purple-dark)] to-blue-dark bg-cover bg-center"
							style={
								coverImage
									? {
											backgroundImage: `url(${coverImage}), linear-gradient(135deg, var(--purple-dark), var(--blue-dark))`,
											backgroundPosition: 'center',
											backgroundSize: 'cover',
										}
									: undefined
							}
						>
							<div className="absolute top-3 right-3">
								<UploadFileButton
									ref={coverUploadRef}
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
								priority: toTaskPriorityValue(project.priority),
								...toCycleFormValues(project),
							}}
							onFinish={handleUpdate}
							size="large"
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
							<ProjectCycleFields />
							<Form.Item
								name="priority"
								label={t('projectList.priority')}
								rules={[{ required: true }]}
							>
								<Select loading={priorityLoading} options={priorityOptions} />
							</Form.Item>
							<Space>
								<Button type="primary" htmlType="submit" loading={saving}>
									{t('projectSetting.update')}
								</Button>
								<SaveAsTemplateButton
									project={project}
									projectId={projectId}
									messageApi={messageApi}
								/>
							</Space>
						</Form>
						<ProjectWorkflowSection
							project={project}
							projectId={projectId}
							messageApi={messageApi}
							onRefresh={refreshProject}
						/>
						<DeleteProjectSection
							project={project}
							projectId={projectId}
							request={request}
							messageApi={messageApi}
							onDeleted={goList}
						/>
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
