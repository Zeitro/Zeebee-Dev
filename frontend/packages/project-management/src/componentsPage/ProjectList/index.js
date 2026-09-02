/**
 * 项目列表：FilterTableCard 负责筛选/分页；卡片视觉对齐 pm-frontend，
 * 创建卡通过占位数据插在当前页网格首位。
 */
'use client';

import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { ReloadOutlined } from '@ant-design/icons';
import { useTranslations } from 'next-intl';
import { FilterTableCard } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import { PROJECT_PAGE_API } from '../../api/project';
import { getTaskPriorityI18nKey, localizeSelectOptions } from '../../enum/task';
import { useOptionsStore } from '../../store/options';
import { CREATE_PROJECT_ITEM, isCreateProjectItem } from './components/constants';
import { CreateProjectCard } from './components/CreateProjectCard';
import { ProjectCard } from './components/ProjectCard';

const PROJECT_CARD_ITEM_CLASS =
	'!p-0 w-[280px] min-h-[280px] overflow-hidden rounded-xl hover:border-base-primary';

export const ProjectList = () => {
	const t = useTranslations();
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const tableFilterRef = useRef();
	const [keyword, setKeyword] = useState('');
	const [priority, setPriority] = useState();

	const priorityOptionsRaw = useOptionsStore((state) => state.priorityOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);

	const priorityOptions = useMemo(
		() => localizeSelectOptions(priorityOptionsRaw, t, getTaskPriorityI18nKey),
		[priorityOptionsRaw, t],
	);

	useEffect(() => {
		const { priorityOptions: cached } = useOptionsStore.getState();
		if (!cached.length) {
			getOptions('project_management_priority', 'priorityOptions').catch(() => {});
		}
	}, [getOptions]);

	const handleRefresh = useCallback(() => {
		cancelMatch(PROJECT_PAGE_API);
		tableFilterRef?.current?.handleSearch();
	}, [cancelMatch]);

	/**
	 * POST 分页；filter 对齐后端 `{ entity_id, entity_type, keyword }`，
	 * 列表首位再插入创建卡占位（不计入 total）。
	 * @param {{ page?: number, limit?: number, filter?: object }} params
	 * @returns {Promise<{ list: object[], total: number }>}
	 */
	const getData = useCallback(
		async (params) => {
			const filter = {
				// entity_id: 0,
				// entity_type: '',
				...(params?.filter || {}),
			};
			if (typeof filter.priority !== 'string' || !filter.priority) {
				delete filter.priority;
			}
			try {
				const res = await request(PROJECT_PAGE_API, {
					method: 'POST',
					body: {
						...params,
						filter,
					},
				});
				const list = res?.list ?? [];
				return {
					list: [CREATE_PROJECT_ITEM, ...list],
					total: res?.total ?? 0,
				};
			} catch {
				return {
					list: [CREATE_PROJECT_ITEM],
					total: 0,
				};
			}
		},
		[request],
	);

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
			getData={getData}
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
