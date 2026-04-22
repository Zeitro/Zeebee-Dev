'use client';
import { BtnPopUpModal, FilterField, NewPagination, NoData } from '@repo/shared-methods';
import Image from 'next/image';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { ReloadOutlined } from '@ant-design/icons';
import { sort_asc, sort_des } from '@/enum/properties';
import { des, asc } from '@repo/shared-methods/enum';
import PropertiesItem from './components/propertiesItem';
import { createRequest, getPagination, transformArrayStructure } from '@repo/shared-methods/utils';
import { Spin } from 'antd';
import { useGetOptions, useAuth } from '@repo/shared-methods/hooks';
import AddModal from '@repo/shared-methods/componentsUi/SelectCard/Add/property';
import { PlusOutlined } from '@ant-design/icons';
import { PERMISSIONS } from '@config/permissions';

const Properties = () => {
	const { hasPermission } = useAuth();
	const { request, cancel } = createRequest();
	const { borrowerOptions } = useGetOptions([{ type: 'borrower' }]);
	const filterFieldRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	const [borrower, setBorrower] = useState([]);

	const [pagination, setPagination] = useState(getPagination());
	const paginationRef = useRef(pagination);

	const [sorter, setSorter] = useState({
		column: 'total_amount',
		sort: des,
	});
	const sorterRef = useRef(sorter);

	const [loading, setLoading] = useState(true);
	const [list, setList] = useState([]);
	const scrollRef = useRef(null);

	useEffect(() => {
		paginationRef.current = pagination;
	}, [pagination]);

	useEffect(() => {
		sorterRef.current = sorter;
	}, [sorter]);

	const getList = useCallback(async () => {
		try {
			const filterParams = filterFieldRef?.current?.getApiParams?.();
			setLoading(true);
			const params = {
				page: paginationRef.current?.current,
				limit: paginationRef.current?.pageSize,
				filter: {
					...filterParams,
				},
			};
			if (sorterRef.current?.column) {
				params.sorter = sorterRef.current;
			}
			const { total, list } = await request('/property/page', {
				method: 'POST',
				body: params,
			});
			setPagination((v) => ({ ...v, total }));
			setList(list);
			setLoading(false);
			return list;
		} catch {
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const handleSearch = (isDelete) => {
		if (isDelete) {
			paginationRef.current.current = 1;
		}
		getList();
	};

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
					setPagination((v) => {
						const next = { ...v, current: 1 };
						paginationRef.current = next;
						return next;
					});
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
						setPagination((v) => {
							const next = { ...v, current: 1 };
							paginationRef.current = next;
							return next;
						});
						getList?.();
					}, 300);
				},
				colProps: {
					className: 'flex-1',
				},
			},
			{
				filterType: 'button',
				hidden: !hasPermission(PERMISSIONS?.PROPERTY?.ADD),
				children: (
					<BtnPopUpModal
						getContent={(setOpen) => (
							<AddModal setOpen={setOpen} handleSearch={() => getList()} />
						)}
						type="primary"
						icon={<PlusOutlined />}
						size="large"
					>
						Add a property
					</BtnPopUpModal>
				),
			},
		],
		[searchValue, borrower, borrowerOptions, cancel, getList, hasPermission],
	);

	useEffect(() => {
		const t = setTimeout(() => {
			getList();
		}, 0);
		return () => clearTimeout(t);
	}, [getList]);

	return (
		<div className={['flex-1 flex flex-col px-8 pt-10 pb-8 overflow-hidden'].join(' ')}>
			<FilterField ref={filterFieldRef} list={filterList} className="pb-8" />
			<div className="flex items-center justify-between flex-shrink-0 mb-4 text-base font-medium pb-2 border-b border-base-disabled">
				<p className="text-base font-medium">{pagination?.total || 0} Properties</p>
				<div
					className="flex cursor-pointer"
					onClick={() => {
						setSorter((v) => {
							const next = { ...v, sort: v.sort === des ? asc : des };
							sorterRef.current = next;
							return next;
						});
						setPagination((v) => {
							const next = { ...v, current: 1 };
							paginationRef.current = next;
							return next;
						});
						getList();
					}}
				>
					<Image src="/images/sort.svg" alt="" width={24} height={24} />
					Total loan: {sorter.sort === des ? sort_des : sort_asc}
				</div>
			</div>
			<div ref={scrollRef} className="flex-1 overflow-hidden overflow-y-auto">
				<Spin spinning={loading}>
					{list?.length ? (
						list.map((item, index) => (
							<PropertiesItem key={index} data={item} handleSearch={handleSearch} />
						))
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
						setPagination((v) => {
							const next = { ...v, current: page };
							paginationRef.current = next;
							return next;
						});
						getList();
						scrollRef.current.scrollTo({
							// behavior: 'smooth',
							top: 0,
						});
					}}
					{...pagination}
				/>
			)}
		</div>
	);
};
export default Properties;
