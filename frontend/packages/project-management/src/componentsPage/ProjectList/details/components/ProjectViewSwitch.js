/**
 * 项目详情视图切换：看板 / 列表 / 时间线 / 工时。只改 view，不拉数据。
 */
'use client';

import {
	AppstoreOutlined,
	ClockCircleOutlined,
	FieldTimeOutlined,
	UnorderedListOutlined,
} from '@ant-design/icons';
import { Segmented } from 'antd';
import { useTranslations } from 'next-intl';
import { PROJECT_DETAIL_VIEW } from '../constants';

/**
 * @param {{ value: string, onChange: (view: string) => void }} props
 */
export function ProjectViewSwitch({ value, onChange }) {
	const t = useTranslations();

	return (
		<Segmented
			value={value}
			onChange={onChange}
			options={[
				{
					label: t('projectDetail.viewKanban'),
					value: PROJECT_DETAIL_VIEW.KANBAN,
					icon: <AppstoreOutlined />,
				},
				{
					label: t('projectDetail.viewTable'),
					value: PROJECT_DETAIL_VIEW.TABLE,
					icon: <UnorderedListOutlined />,
				},
				{
					label: t('projectDetail.viewTimeline'),
					value: PROJECT_DETAIL_VIEW.TIMELINE,
					icon: <FieldTimeOutlined />,
				},
				{
					label: t('projectDetail.viewHours'),
					value: PROJECT_DETAIL_VIEW.HOURS,
					icon: <ClockCircleOutlined />,
				},
			]}
		/>
	);
}
