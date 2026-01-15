'use client';
import { SearchOutlined } from '@ant-design/icons';
import { NewTable, NoData, GeoLocationInput } from '@repo/shared-methods';
import { Button, InputNumber, Select } from 'antd';
import { useRef, useState } from 'react';
import useContent from './components/useContent';
import Image from 'next/image';
// 根据输入地址，查找周围已售房屋估值
const PropertyFinder = () => {
	const geoInputRef = useRef(null);
	const [errorMessage, setErrorMessage] = useState('');
	// const [expandedRowKey, setExpandedRowKey] = useState(null);
	const [options] = useState([
		{ value: 'all', label: 'Price range' },
		{ value: 'less500', label: 'Up to $500K' },
		{ value: '500to1000', label: 'Between $500K - M' },
		{ value: '1000to3000', label: 'Between $1M - $3M' },
		{ value: 'over3000', label: 'Over$0' },
	]);
	const {
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
		onChange,
		pagination,
		isSearch,
	} = useContent({
		setErrorTips: (message) => {
			if (geoInputRef.current && geoInputRef.current?.setErrorTips) {
				geoInputRef.current.setErrorTips(message);
			}
		},
	});

	return (
		<div className="flex-1 overflow-hidden pt-8 px-8 flex flex-col">
			<div className="flex items-center gap-4 shrink-0">
				<GeoLocationInput
					ref={geoInputRef}
					onLocationSelect={(loc) => setLocation(loc)}
					placeholder="Search by city, county, or zip code..."
					className="!w-[30%]"
					onError={(error) => setErrorMessage(error)}
				/>
				<InputNumber
					suffix="mi"
					placeholder="Distance"
					min={1}
					size="large"
					className="!w-[160px]"
					controls={false}
					value={radius}
					onChange={(value) => setRadius(value)}
				/>
				<Button
					type="primary"
					loading={loading}
					disabled={!!errorMessage}
					size="large"
					onClick={handleSearch}
					icon={<SearchOutlined />}
				>
					Search
				</Button>
			</div>
			<div className="text-base-error h-6 pt-1 text-sm shrink-0">{errorMessage}</div>
			<div className="text-base-label text-sm shrink-0">
				*California listings only, data updated weekly.
			</div>
			<div className="flex items-center justify-between flex-shrink-0 pt-3 mb-4 text-base font-medium pb-2 border-b border-base-disabled">
				<p className="text-base font-medium">{list?.length || 0} Property finder</p>
				<Select
					options={options}
					value={filterValue}
					onChange={(value) => {
						setFilterValue(value);
						handleFilterUpdate(list, value);
					}}
					className="!w-[170px]"
				/>
			</div>
			{!isSearch ? (
				<div className="flex items-center flex-col flex-1 pt-[100px]">
					<Image
						src="/images/searchOpenHouse.png"
						alt="searchOpenHouse"
						width={191}
						height={109}
					/>
					<div className="text-lg font-medium text-base-gray pt-4">
						Enter a city, zip code, or state to find open houses near you.
					</div>
				</div>
			) : (
				<NewTable
					className="flex-1 overflow-hidden [&_.ant-spin-nested-loading]:h-full [&_.ant-spin-container]:h-full [&_.ant-table-container]:h-full [&_.ant-table]:h-[calc(100%-64px)]"
					columns={columns}
					scroll={{ x: 1200 }}
					dataSource={dataSource}
					locale={{
						emptyText: (
							<NoData
								text={
									<>
										<div className="text-lg font-medium pb-2">
											No open houses found in this area.
										</div>
										<div className="text-sm">
											Try adjusting your search radius or explore other nearby
											cities.
										</div>
									</>
								}
							/>
						),
					}}
					onChange={onChange}
					pagination={pagination}
					loading={loading}
					// expandable={{
					// 	expandedRowRender: (record) => <Income record={record} />,
					// 	rowExpandable: () => true,
					// 	expandedRowKeys: expandedRowKey ? [expandedRowKey] : [],
					// 	onExpand: (expanded, record) => {
					// 		setExpandedRowKey(expanded ? record?.id : null);
					// 	},
					// 	fixed: 'right',
					// 	expandIcon: ({ expanded, onExpand, record }) => {
					// 		return (
					// 			<DownOutlined
					// 				onClick={(e) => onExpand(record, e)}
					// 				className={`${
					// 					expanded ? 'rotate-180' : ''
					// 				} transition-all duration-300 ease-in-out`}
					// 			/>
					// 		);
					// 	},
					// }}
					// pagination={hasPagination}
				/>
			)}
		</div>
	);
};

export default PropertyFinder;
