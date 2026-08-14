'use client';
import React from 'react';
import { Spin } from 'antd';
import { NoData } from '../NoData/index.js';
import { useResolvedPaginationComponent } from '../FilterTable/FilterTableUiContext';

/**
 * FilterTableCard 的卡片列表 + 分页；分页组件由应用注入。
 * @param {React.ComponentType} [PaginationComponent] 覆盖 Context 中的分页组件
 */
export const Cards = ({
	PaginationComponent,
	renderItem = () => {},
	itemClassName,
	pagination,
	loading,
	list,
	scrollRef,
	handleChangePagination,
	cardTitle,
}) => {
	const Pagination = useResolvedPaginationComponent(PaginationComponent);

	return (
		<>
			<div className="flex items-center justify-between flex-shrink-0 pt-4 text-base font-medium pb-2 border-b border-base-disabled">
				<p className="text-base font-medium">
					{pagination?.total || 0} {cardTitle}
				</p>
				{/* <div
						className="flex cursor-pointer"
						onClick={() => {
							setSorter((v) => {
								const next = { ...v, sort: v.sort === des ? asc : des };
								return next;
							});
							setPagination((v) => {
								const next = { ...v, current: 1 };
								return next;
							});
							getList();
						}}
					>
						<Image src="/images/sort.svg" alt="" width={24} height={24} />
						Total loan: {sorter.sort === des ? 'Descending' : 'Ascending'}
					</div> */}
			</div>
			<Spin spinning={loading} wrapperClassName="flex-1 overflow-hidden overflow-y-auto">
				<div ref={scrollRef} className="flex flex-wrap gap-6">
					{list?.length ? (
						list.map((item, index) => (
							<div
								key={index}
								className={[
									'border border-base-border rounded-lg px-6 py-6 flex flex-col',
									itemClassName,
								].join(' ')}
								style={{ boxShadow: '0 2px 8px 0 rgba(0,0,0,0.15)' }}
							>
								{renderItem(item, index)}
							</div>
						))
					) : loading ? (
						<div className="min-h-[100px]"></div>
					) : (
						<NoData height={200} />
					)}
				</div>
			</Spin>
			{list.length > 0 && (
				<Pagination
					className="justify-end"
					onChange={handleChangePagination}
					{...pagination}
				/>
			)}
		</>
	);
};

Cards.displayName = 'Cards';
