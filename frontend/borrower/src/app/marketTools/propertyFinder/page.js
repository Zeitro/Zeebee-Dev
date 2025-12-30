'use client';
import GeoLocationInput from '@/components/GeoLocationInput/index';
import { DownOutlined, SearchOutlined } from '@ant-design/icons';
import { FilterTable, NewTable, NoData } from '@repo/shared-methods';
import { formatAmount, getPagination } from '@repo/shared-methods/utils';
import { Button, Divider, InputNumber, Select, Tooltip } from 'antd';
import { use, useEffect, useRef, useState } from 'react';
import Agent, { getCopyContent } from './components/Agent';
import { createRequest } from '@repo/shared-methods/utils';
import useContent from './components/useContent';
import Income from './components/Income';
// 根据输入地址，查找周围已售房屋估值
const PropertyFinder = () => {
	const geoInputRef = useRef(null);
	const [errorMessage, setErrorMessage] = useState('');
	const [expandedRowKey, setExpandedRowKey] = useState(null);
	const [options] = useState([
		{ value: 'all', label: 'Price range' },
		{ value: 'less500', label: 'Up to $500K' },
		{ value: '500to1000', label: 'Between $500K - M' },
		{ value: '1000to3000', label: 'Between $1M - $3M' },
		{ value: 'over3000', label: 'Over$0' },
	]);
	// const [pagination, setPagination] = useState(getPagination());
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
	} = useContent();

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
					onChange={({ target: { value } }) => setRadius(value)}
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
					onChange={({ target: { value } }) => {
						setFilterValue(value);
						handleFilterUpdate(value);
					}}
					className="!w-[160px]"
				/>
			</div>
			<NewTable
				className="flex-1 overflow-hidden [&_.ant-spin-nested-loading]:h-full [&_.ant-spin-container]:h-full [&_.ant-table-container]:h-full [&_.ant-table]:h-[calc(100%-64px)]"
				columns={columns}
				scroll={{ x: 800 }}
				dataSource={list}
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
				expandable={{
					expandedRowRender: (record) => <Income record={record} />,
					rowExpandable: () => true,
					expandedRowKeys: expandedRowKey ? [expandedRowKey] : [],
					onExpand: (expanded, record) => {
						setExpandedRowKey(expanded ? record?.id : null);
					},
					fixed: 'right',
					expandIcon: ({ expanded, onExpand, record }) => {
						return (
							<DownOutlined
								onClick={(e) => onExpand(record, e)}
								className={`${
									expanded ? 'rotate-180' : ''
								} transition-all duration-300 ease-in-out`}
							/>
						);
					},
				}}
				// pagination={hasPagination}
			/>
		</div>
	);
};

export default PropertyFinder;
