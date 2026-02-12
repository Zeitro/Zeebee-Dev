import { useState } from 'react';
import { createRequest, formatAmount, getPagination } from '@repo/shared-methods/utils';
import Agent, { getCopyContent } from './Agent';
const useContent = ({ setErrorTips }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [list, setList] = useState([]);
	const [dataSource, setDataSource] = useState([]);
	const [filterValue, setFilterValue] = useState('all'); //price筛选
	const [radius, setRadius] = useState(8);
	const [location, setLocation] = useState(null);
	const [pagination, setPagination] = useState(getPagination());
	const [isSearch, setIsSearch] = useState(false);

	// 根据price过滤
	const handleFilterUpdate = (arr = list, value = filterValue) => {
		let filtered = [...arr];

		if (value === 'less500') {
			filtered = filtered.filter((listing) => listing.price <= 500000);
		} else if (value === '500to1000') {
			filtered = filtered.filter(
				(listing) => listing.price > 500000 && listing.price <= 1000000,
			);
		} else if (value === '1000to3000') {
			filtered = filtered.filter(
				(listing) => listing.price > 1000000 && listing.price <= 3000000,
			);
		} else if (value === 'over3000') {
			filtered = filtered.filter((listing) => listing.price > 3000000);
		}
		setPagination((v) => ({ ...v, total: filtered.length, current: 1 }));
		setDataSource(filtered);
	};
	const columns = [
		{
			title: 'Address',
			dataIndex: 'full_address',
			key: 'full_address',
			render: (v) =>
				getCopyContent(v, {
					className: '!font-roboto',
				}),
			ellipsis: true,
			width: 360,
		},
		{
			title: 'Days listed',
			dataIndex: 'days_on_market',
			key: 'days_on_market',
			sorter: (a, b) => a.days_on_market - b.days_on_market,
			defaultSortOrder: 'ascend',
			ellipsis: true,
		},
		{
			title: 'Listed price',
			dataIndex: 'price',
			key: 'price',
			render: (price) => formatAmount(price),
			sorter: (a, b) => a.price - b.price,
			ellipsis: true,
		},
		{
			title: 'Distance',
			dataIndex: 'distance',
			key: 'distance',
			render: (distance) => `${distance} mi`,
			sorter: (a, b) => a.distance - b.distance,
			ellipsis: true,
		},
		{
			title: 'Agent',
			dataIndex: 'agent_name',
			key: 'agent_name',
			render: (v, record) => <Agent record={record} />,
			sorter: (a, b) => a.agent_name.localeCompare(b.agent_name),
			ellipsis: true,
		},
	];

	const getList = async () => {
		try {
			setLoading(true);
			setIsSearch(true);
			const params = new URLSearchParams({
				lat: location.lat,
				lon: location.lng,
				dist: parseFloat(radius),
				addr: location.address,
			});
			const { list } = await request(`/zeitro/listings/search?${params}`);
			const arr = list.map((listing, index) => ({
				...listing,
				id: index,
				price: parseFloat(listing.price?.trim?.()),
			}));
			setList(arr);
			handleFilterUpdate(arr);
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};

	const handleSearch = async () => {
		if (!location) {
			return setErrorTips('Please select an option after input your location.');
		}
		await getList();
	};

	const onChange = (paginationObj) => {
		setPagination((v) => {
			if (paginationObj.current !== v.current) {
				v.current = paginationObj.current;
			} else {
				v.current = 1;
			}
			return { ...v };
		});
	};
	return {
		columns,
		handleFilterUpdate,
		handleSearch,
		loading,
		list,
		dataSource,
		radius,
		setRadius,
		location,
		setLocation,
		filterValue,
		setFilterValue,
		onChange,
		pagination,
		isSearch,
	};
};
export default useContent;
