/**
 * 任务页签右侧操作：看板/表格/时间线/工时、新建任务。
 * 只挂在任务页签，文档页不渲染。
 */
'use client';

import { PlusOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import { useTranslations } from 'next-intl';
import { ProjectViewSwitch } from './ProjectViewSwitch';

/**
 * @param {{
 *   view: string,
 *   onViewChange: (view: string) => void,
 *   showCreate?: boolean,
 *   onCreateTask?: () => void,
 * }} props
 */
export function ProjectTaskToolbar({ view, onViewChange, showCreate = false, onCreateTask }) {
	const t = useTranslations();

	return (
		<div className="flex items-center gap-3">
			<ProjectViewSwitch value={view} onChange={onViewChange} />
			{showCreate ? (
				<Button type="primary" icon={<PlusOutlined />} onClick={onCreateTask}>
					{t('kanban.newTask')}
				</Button>
			) : null}
		</div>
	);
}
