'use client';
import { Button, Radio } from 'antd';
import React, { useCallback, useState } from 'react';
import { FilterTable } from '@repo/shared-methods';
import useColumn from './useColumns';

const DifferentStatusLoans = ({ stage, data }) => {
	const { getColumnsByStage, getDataSourceByStage, titleInfo, totalList } = useColumn(
		data,
		stage,
	);
	return (
		<div className="border border-base-border rounded-lg">
			<div className="flex justify-between items-center">
				<span>{titleInfo[stage]}</span>
				<Button type="text">View all</Button>
			</div>
			<div className="flex gap-4">
				{totalList.map((item, index) => (
					<div key={index} className="border border-base-border rounded-lg">
						<span>{item.label}</span>
						<div>{item?.render?.() || item.value}</div>
					</div>
				))}
			</div>
			<FilterTable
				key={stage}
				// ref={tableFilterRef}
				columns={getColumnsByStage(stage)}
				// scroll={{ x: scrollX }}
				hasFilterField={false}
				pagination={false}
				// filterWrapClassName={filterWrapClassName}
				dataSource={getDataSourceByStage()}
				// onRow={(row) => ({
				// 	onClick: () => {
				// 		handleSelectedRow(row, stage);
				// 	},
				// })}
				// getData={(params) => {
				// 	if (!api) {
				// 		return {
				// 			list: dataSource,
				// 			total: dataSource?.length,
				// 		};
				// 	}
				// 	const filter = hasFilterField
				// 		? {}
				// 		: {
				// 				...(params?.filter || {}),
				// 				stage,
				// 			};
				// 	const payload = {
				// 		method,
				// 	};
				// 	if (method === 'POST') {
				// 		payload.body = {
				// 			...params,
				// 			filter,
				// 		};
				// 	}
				// 	return request(api, payload);
				// }}
			/>
		</div>
	);
};
export default DifferentStatusLoans;
