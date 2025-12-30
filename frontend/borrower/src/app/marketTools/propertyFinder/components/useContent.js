import { useRef, useState } from 'react';
import { createRequest, formatAmount } from '@repo/shared-methods/utils';
import Agent, { getCopyContent } from './Agent';
const useContent = () => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [list, setList] = useState([
		{
			address: '123 Main St, Anytown, CA 12345',
			agentphone: '5555555555',
			agentname: 'miko',
			id: 1,
		},
	]);
	const [dataSource, setDataSource] = useState([]);
	const [filterValue, setFilterValue] = useState('all'); //price筛选
	const [radius, setRadius] = useState(8);
	const [location, setLocation] = useState(null);

	// 根据price过滤
	const handleFilterUpdate = (value = filterValue) => {
		let filtered = [...list];

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

		setList(filtered);
	};
	const columns = [
		{
			title: 'Address',
			dataIndex: 'address',
			key: 'address',
			render: (v) => <div>{getCopyContent(v)}</div>,
			// render: (address, record) => renderAddressColumn(address, record),
		},
		{
			title: 'Days listed',
			dataIndex: 'daysonmarket',
			key: 'daysonmarket',
			sorter: (a, b) => a.daysonmarket - b.daysonmarket,
			defaultSortOrder: 'ascend',
		},
		{
			title: 'Listed price',
			dataIndex: 'price',
			key: 'price',
			render: (price) => formatAmount(price),
			sorter: (a, b) => a.price - b.price,
		},
		{
			title: 'Distance',
			dataIndex: 'distance',
			key: 'distance',
			render: (distance) => `${distance} mi`,
			sorter: (a, b) => a.distance - b.distance,
		},
		{
			title: 'Agent',
			dataIndex: 'agentname',
			key: 'agentname',
			render: (v, record) => <Agent record={record} />,
			sorter: (a, b) => a.agentname.localeCompare(b.agentname),
		},
	];

	const getList = async () => {
		try {
			setLoading(true);
			const response = await request('/los/openhousesearch', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					Lat: location.lat,
					Lon: location.lng,
					Dist: parseFloat(radius),
					Addr: location.address,
				}),
				isExternalLink: true,
			});
			const arr = response.map((listing, index) => ({
				id: index,
				address: listing.FullAddress,
				city: listing.City,
				county: listing.County,
				state: listing.State,
				zip: listing.ZipCode,
				price: parseFloat(listing.Price.trim()),
				daysonmarket: listing.DaysOnMarket,
				agentname: listing.AgentName,
				agentphone: listing.AgentPhone,
				agentemail: listing.AgentEmail,
				distance: listing.Distance,
			}));
			setList(arr);
			setLoading(false);
		} catch (error) {
			console.log(error, '----list---');
			setLoading(false);
		}
	};

	const handleSearch = async () => {
		// if (!location) {
		// 	return geoInputRef.current.setErrorTips(
		// 		'Please select an option after input your location.',
		// 	);
		// }
		await getList();
		handleFilterUpdate();
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
		setLocation,
		filterValue,
		setFilterValue,
	};
};
export default useContent;
