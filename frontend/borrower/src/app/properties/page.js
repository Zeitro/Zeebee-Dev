'use client';
import FilterField from '@/components/FilterField';
import Image from 'next/image';
import React, { useEffect, useMemo, useState } from 'react';
import { sort_asc, sort_des } from '@/enum/properties';
import PropertiesItem from './components/propertiesItem';
import NewPagination from '@/components/NewPagination';
const Properties = () => {
	//    const [type, setType] = useState('');
	const [searchValue, setSearchValue] = useState('');
	const [borrower, setBorrower] = useState(null);
	const [list, setList] = useState([{ totalLoan: 10 }, {}]);
	const [sortOrder, setSortOrder] = useState(sort_des);
	const filterList = useMemo(
		() => [
			// {
			//     type: 'select-input',
			//     selectValue: type,
			//     value: searchValue,
			//     selectProps: {
			//         onChange: (e) => setType(e),
			//     },
			//     onChange: ({ target: { value } }) => setSearchValue(value),
			//     options: [
			//         {
			//             label: 'Address',
			//             value: 'address',
			//         },
			//         {
			//             label: 'State',
			//             value: 'state',
			//         },
			//         {
			//             label: 'City',
			//             value: 'city',
			//         },
			//         {
			//             label: 'County',
			//             value: 'county',
			//         },
			//         {
			//             label: 'Zip Code',
			//             value: 'zipCode',
			//         },
			//     ],
			// },
			{
				type: 'input-search',
				placeholder: 'Search by any field',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
				onSearch: () => {},
			},
			{
				type: 'select',
				placeholder: 'Borrower',
				value: borrower,
				onChange: (value) => setBorrower(value),
				options: [
					{
						label: 'Address',
						value: 'address',
					},
					{
						label: 'State',
						value: 'state',
					},
					{
						label: 'City',
						value: 'city',
					},
				],
			},
		],
		[searchValue, borrower],
	);

	return (
		<div className={['h-full flex flex-col px-8 pt-10 pb-8'].join(' ')}>
			<FilterField list={filterList} className="pb-8" />
			<div className="flex items-center justify-between flex-shrink-0 mb-4 text-base font-medium pb-2 border-b border-base-disabled">
				<p className="text-base font-medium">{list.length} Properties</p>
				<div
					className="flex cursor-pointer"
					onClick={() => setSortOrder(sortOrder === sort_des ? sort_asc : sort_des)}
				>
					<Image src={require('@/assets/images/sort.svg')} alt="" />
					Total loan: {sortOrder}
				</div>
			</div>
			<div className="flex-1 overflow-hidden overflow-y-auto">
				{list.map((item, index) => (
					<PropertiesItem key={index} data={item} />
				))}
			</div>
			<NewPagination defaultCurrent={2} total={30} className="!mt-4 justify-end"/>
		</div>
	);
};
export default Properties;
