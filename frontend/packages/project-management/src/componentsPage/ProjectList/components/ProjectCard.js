/**
 * 项目卡片：对齐 pm-frontend 封面 + 标题/描述 + 底栏操作；纯 UI，不改 rowData。
 */
'use client';

import { AppstoreOutlined, SettingOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { displayValue } from '@repo/shared-methods/utils';

/**
 * @param {{ row: object }} props
 */
export function ProjectCard({ row }) {
	const t = useTranslations('projectList');
	const router = useRouter();
	const projectId = row?.id;
	const coverImage = row?.cover_image;

	const goKanban = () => {
		if (projectId == null) return;
		router.push(`/project/${projectId}/kanban`);
	};

	const goSettings = () => {
		if (projectId == null) return;
		router.push(`/project/${projectId}/settings`);
	};

	return (
		<div className="flex flex-col h-full min-h-[280px]">
			<div className="cursor-pointer flex-1" onClick={goKanban}>
				<div
					className="h-[140px] bg-gradient-to-r from-[var(--purple-dark)] to-blue-dark bg-cover bg-center"
					style={
						coverImage
							? {
									backgroundImage: `url(${coverImage}), linear-gradient(135deg, var(--purple-dark), var(--blue-dark))`,
								}
							: undefined
					}
				/>
				<div className="px-4 pt-4 pb-3">
					<p className="font-medium text-base text-black-dark truncate">
						{displayValue(row?.name)}
					</p>
					<p className="text-sm text-base-gray mt-1 line-clamp-2 min-h-[40px]">
						{displayValue(row?.description)}
					</p>
				</div>
			</div>
			<div className="mt-auto flex items-stretch border-t border-base-border">
				<button
					type="button"
					className="flex-1 py-3 text-base-gray hover:text-base-primary"
					aria-label={t('openKanban')}
					onClick={goKanban}
				>
					<AppstoreOutlined />
				</button>
				<span className="w-px bg-base-border" />
				<button
					type="button"
					className="flex-1 py-3 text-base-gray hover:text-base-primary"
					aria-label={t('openSettings')}
					onClick={goSettings}
				>
					<SettingOutlined />
				</button>
			</div>
		</div>
	);
}
