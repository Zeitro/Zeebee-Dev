/**
 * 创建项目卡：始终插在网格首位，点击打开 CreateProjectModal。
 */
'use client';

import { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import CreateProjectModal from './CreateProjectModal';
import { message } from 'antd';

/**
 * @param {{ handleSearch?: () => void }} props
 */
export function CreateProjectCard({ handleSearch }) {
	const t = useTranslations('projectList');
	const [open, setOpen] = useState(false);
	const [messageApi, messageContextHolder] = message.useMessage();

	return (
		<div className="h-full min-h-[280px] bg-gray-light">
			{messageContextHolder}
			<div
				className="flex flex-col items-center justify-center h-full min-h-[280px] cursor-pointer"
				onClick={() => setOpen(true)}
				onKeyDown={(e) => {
					if (e.key === 'Enter') setOpen(true);
				}}
				role="button"
				tabIndex={0}
			>
				<PlusOutlined className="text-[48px] text-base-disabled" />
				<p className="mt-3 text-base text-base-gray">{t('createTitle')}</p>
			</div>
			{open && (
				<CreateProjectModal
					messageApi={messageApi}
					setOpen={setOpen}
					handleSearch={handleSearch}
				/>
			)}
		</div>
	);
}
