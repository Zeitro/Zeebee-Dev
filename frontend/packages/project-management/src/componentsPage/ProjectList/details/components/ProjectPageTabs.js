/**
 * 项目详情页签：任务 / 文档。任务专属操作走 extra，文档页不传。
 */
'use client';

import { Tabs } from 'antd';
import { useTranslations } from 'next-intl';
import { PROJECT_DETAIL_TAB } from '../constants';

/**
 * @param {{
 *   value: string,
 *   onChange: (tab: string) => void,
 *   extra?: import('react').ReactNode,
 * }} props
 */
export function ProjectPageTabs({ value, onChange, extra = null }) {
	const t = useTranslations();

	return (
		<Tabs
			activeKey={value}
			onChange={onChange}
			tabBarExtraContent={extra}
			className="[&_.ant-tabs-nav]:mb-3"
			items={[
				{ key: PROJECT_DETAIL_TAB.TASKS, label: t('projectDetail.tabTasks') },
				{ key: PROJECT_DETAIL_TAB.DOCS, label: t('projectDetail.tabDocs') },
			]}
		/>
	);
}
