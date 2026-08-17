/**
 * 项目列表：FilterTableCard 负责筛选/分页；卡片视觉对齐 pm-frontend，
 * 创建卡通过占位数据插在当前页网格首位。
 */
'use client';

import React, { useCallback, useMemo, useRef, useState } from 'react';
import { ReloadOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { FilterTableCard } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import { PROJECT_PAGE_API } from '../../api/project';
import { PROJECT_PRIORITY_FILTER_VALUES, PROJECT_PRIORITY_I18N_KEY } from '../../enum/project';
import { CREATE_PROJECT_ITEM, isCreateProjectItem } from './components/constants';
import { CreateProjectCard } from './components/CreateProjectCard';
import { ProjectCard } from './components/ProjectCard';
import { getMockProjectPage } from './mockData';

const PROJECT_CARD_ITEM_CLASS =
	'!p-0 w-[280px] min-h-[280px] overflow-hidden rounded-xl hover:border-base-primary';

export const ProjectList = () => {
	const t = useTranslations();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const tableFilterRef = useRef();
	const [keyword, setKeyword] = useState('');
	const [priority, setPriority] = useState();

	const priorityOptions = useMemo(
		() =>
			PROJECT_PRIORITY_FILTER_VALUES.map((value) => ({
				value,
				label: t(PROJECT_PRIORITY_I18N_KEY[value]),
			})),
		[t],
	);

	const handleRefresh = useCallback(() => {
		cancelMatch(PROJECT_PAGE_API);
		tableFilterRef?.current?.handleSearch();
	}, [cancelMatch]);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: t('projectList.searchPlaceholder'),
				paramsKey: 'keyword',
				value: keyword,
				onChange: ({ target: { value } }) => setKeyword(value),
			},
			{
				filterType: 'select',
				placeholder: t('projectList.priority'),
				paramsKey: 'priority',
				value: priority,
				onChange: (value) => setPriority(value),
				options: priorityOptions,
				style: { width: 160 },
			},
			{
				filterType: 'button',
				text: t('common.search'),
				onClick: handleRefresh,
			},
			{
				filterType: 'button',
				type: 'default',
				text: t('common.reset'),
				icon: <ReloadOutlined />,
				onClick: () => {
					cancelMatch(PROJECT_PAGE_API);
					setKeyword('');
					setPriority(undefined);
					tableFilterRef?.current?.handleSearch({
						filter: {
							keyword: '',
							priority: undefined,
						},
					});
				},
			},
		],
		[keyword, priority, priorityOptions, t, cancelMatch, handleRefresh],
	);

	return (
		<FilterTableCard
			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
			ref={tableFilterRef}
			filterList={filterList}
			hasFilterField={true}
			typeKey="card"
			cardTitle={t('projectList.cardTitle')}
			itemClassName={PROJECT_CARD_ITEM_CLASS}
			getData={async (params) => {
				const filter = { ...(params?.filter || {}) };
				if (filter.priority !== undefined && filter.priority !== '') {
					filter.priority = Number(filter.priority);
				} else {
					delete filter.priority;
				}
				// try {
				// 	const res = await request(PROJECT_PAGE_API, {
				// 		method: 'POST',
				// 		body: {
				// 			...params,
				// 			filter,
				// 		},
				// 	});
				// 	const list = res?.list ?? res?.data ?? [];
				// 	return {
				// 		list: [CREATE_PROJECT_ITEM, ...list],
				// 		total: res?.total ?? 0,
				// 	};
				// } catch {
				// 	return {
				// 		list: [CREATE_PROJECT_ITEM],
				// 		total: 0,
				// 	};
				// }
				const { list, total } = getMockProjectPage({ ...params, filter });
				return {
					list: [CREATE_PROJECT_ITEM, ...list],
					total,
				};
			}}
			renderItem={(row) =>
				isCreateProjectItem(row) ? (
					<CreateProjectCard handleSearch={handleRefresh} />
				) : (
					<ProjectCard row={row} />
				)
			}
		/>
	);
};
