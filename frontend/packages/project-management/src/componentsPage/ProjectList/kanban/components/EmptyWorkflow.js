/**
 * 项目尚未配置工作流时的空态；入口跳转项目设置。
 */
'use client';

import { ApartmentOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { KanbanHeader } from './KanbanHeader';

/**
 * @param {{ project?: object|null }} props
 */
export function EmptyWorkflow({ project }) {
	const t = useTranslations();
	const router = useRouter();
	const projectId = project?.id;

	return (
		<div className="flex flex-col h-[calc(100vh-130px)] px-8 pt-6 pb-6">
			<KanbanHeader project={project} taskCount={0} showCreate={false} />
			<div className="flex-1 flex flex-col items-center justify-center">
				<NoData text={t('kanban.emptyWorkflow')} />
				{projectId != null ? (
					<Button
						type="primary"
						icon={<ApartmentOutlined />}
						onClick={() => router.push(`/project/${projectId}/settings`)}
					>
						{t('kanban.configureWorkflow')}
					</Button>
				) : null}
			</div>
		</div>
	);
}
