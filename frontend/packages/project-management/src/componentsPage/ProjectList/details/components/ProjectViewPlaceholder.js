/**
 * 列表 / 时间线 / 工时空态。真实页面接入后由 details 入口替换本组件。
 */
'use client';

import { useTranslations } from 'next-intl';
import { NoData } from '@repo/shared-methods';
import { PROJECT_DETAIL_VIEW } from '../constants';

const VIEW_LABEL_KEY = {
	[PROJECT_DETAIL_VIEW.TABLE]: 'projectDetail.viewTable',
	[PROJECT_DETAIL_VIEW.TIMELINE]: 'projectDetail.viewTimeline',
	[PROJECT_DETAIL_VIEW.HOURS]: 'projectDetail.viewHours',
};

/**
 * @param {{ view: string }} props
 */
export function ProjectViewPlaceholder({ view }) {
	const t = useTranslations();
	const labelKey = VIEW_LABEL_KEY[view] || 'projectDetail.viewTable';

	return (
		<div className="flex flex-1 items-center justify-center bg-white min-h-0">
			<NoData text={t('projectDetail.viewEmpty', { view: t(labelKey) })} />
		</div>
	);
}
