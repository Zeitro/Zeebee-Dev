'use client';
import FilterField from '@/components/FilterField';
import Image from 'next/image';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { ReloadOutlined } from '@ant-design/icons';
import { sort_asc, sort_des } from '@/enum/properties';
import { des, asc } from '@/enum/common';
import PropertiesItem from './components/propertiesItem';
import NewPagination from '@/components/NewPagination';
import { createRequest } from '@/lib/fetch';
import { getPagination, transformArrayStructure } from '@/lib/common';
import { Spin } from 'antd';
import NoData from '@/components/NoData';
import { useGetOptions } from '@/hooks/useGetOptions';
const Properties = () => {
	const { request, cancel } = createRequest();
	const { borrowerOptions } = useGetOptions([{ type: 'borrower' }]);
	const filterFieldRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [borrower, setBorrower] = useState([]);
	const pagination = useRef(getPagination());
	const sorter = useRef({
		column: 'total_amount',
		sort: des,
	});
	const [loading, setLoading] = useState(true);
	const [list, setList] = useState([]);
	const scrollRef = useRef(null);

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by any field',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
			},
			{
				filterType: 'treeSelect',
				label: 'Borrower',
				value: borrower,
				paramsKey: 'applicant',
				placeholder: 'Borrower',
				onChange: (value) => setBorrower(value),
				treeData: transformArrayStructure({
					sourceArray: borrowerOptions,
					outputKey: 'title',
				})?.map((i) => ({ ...i, key: i.value })),
			},
			{
				filterType: 'button',
				onClick: () => {
					cancel('/property/page');
					pagination.current.current = 1;
					getList();
				},
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					setSearchValue('');
					setBorrower([]);
					setTimeout(() => {
						getList?.();
					}, 300);
				},
			},
		],
		[searchValue, borrower, borrowerOptions],
	);

	const getList = async () => {
		try {
			const filterParams = filterFieldRef?.current?.getApiParams?.();
			setLoading(true);
			const params = {
				page: pagination.current?.current,
				limit: pagination.current?.pageSize,
				filter: {
					...filterParams,
				},
			};
			if (sorter.current?.column) {
				params.sorter = sorter.current;
			}
			const { total, list } = await request('/property/page', {
				method: 'POST',
				body: params,
			});
			pagination.current.total = total;
			setList(list);
			setLoading(false);
		} catch (error) {
			setLoading(false);
		}
	};

	useEffect(() => {
		getList();
	}, []);

	return (
		<div className={['h-full flex flex-col px-8 pt-10 pb-8'].join(' ')}>
			<FilterField ref={filterFieldRef} list={filterList} className="pb-8" />
			<div className="flex items-center justify-between flex-shrink-0 mb-4 text-base font-medium pb-2 border-b border-base-disabled">
				<p className="text-base font-medium">{pagination.current?.total || 0} Properties</p>
				<div
					className="flex cursor-pointer"
					onClick={() => {
						sorter.current.sort = sorter.current.sort === des ? asc : des;
						pagination.current.current = 1;
						getList();
					}}
				>
					<Image src="/images/sort.svg" alt="" width={24} height={24} />
					Total loan: {sorter.current.sort === des ? sort_des : sort_asc}
				</div>
			</div>
			<div ref={scrollRef} className="flex-1 overflow-hidden overflow-y-auto">
				<Spin spinning={loading}>
					{list?.length ? (
						list.map((item, index) => <PropertiesItem key={index} data={item} />)
					) : loading ? (
						<div className="min-h-[100px]"></div>
					) : (
						<NoData height={200} />
					)}
				</Spin>
			</div>
			{!!list.length && (
				<NewPagination
					className="!mt-4 justify-end"
					onChange={(page) => {
						cancel('/property/page');
						pagination.current.current = page;
						getList();
						scrollRef.current.scrollTo({
							// behavior: 'smooth',
							top: 0,
						});
					}}
					{...pagination.current}
				/>
			)}
		</div>
	);
};
export default Properties;
