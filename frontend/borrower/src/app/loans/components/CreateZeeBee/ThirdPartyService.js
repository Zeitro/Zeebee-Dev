'use client';
import React, { useState } from 'react';
import NewTable from '@/components/NewTable';
import { useColumnByZeeBeeLoan } from './components/column';
import { AddEditBtn } from './components/AddEditBtn';
const ThirdPartyService = () => {
	const [list, setList] = useState([{ id: 1 }]);
	const { thirdPartyServicesColumns } = useColumnByZeeBeeLoan();

	return (
		<div className="gap-y-8 flex flex-col w-full">
			<div className="text-2xl font-medium">Third party services</div>
			{list?.length > 0 && (
				<NewTable
					dataSource={list}
					columns={thirdPartyServicesColumns}
					pagination={false}
					rowKey="id"
				/>
			)}
			<AddEditBtn>Add third party service</AddEditBtn>
			<div className="pb-6"></div>
		</div>
	);
};

export default ThirdPartyService;
