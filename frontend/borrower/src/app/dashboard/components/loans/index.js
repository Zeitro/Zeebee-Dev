'use client';
import { Button } from 'antd';
import React from 'react';
import { FilterTable } from '@repo/shared-methods';
import useColumn from './useColumns';
import { ArrowRightOutlined } from '@ant-design/icons';

const DifferentStatusLoans = ({
	stage,
	data,
	handleGotoLoans,
	handleGotoLoanDetails,
	className,
	dataType,
}) => {
	const { getColumnsByStage, getDataSourceByStage, titleInfo, totalList } = useColumn(
		data,
		stage,
	);

	return (
		<div className="border border-base-border rounded-lg mt-6 bg-background">
			<div className="flex justify-between items-center p-4">
				<span className="text-xl font-medium text-base-text">{titleInfo[stage]}</span>
				<Button
					type="text"
					size="small"
					className="!text-base-gray"
					onClick={() => handleGotoLoans(stage, dataType)}
				>
					View all
					<ArrowRightOutlined />
				</Button>
			</div>
			<div
				className={`px-4 pb-6 grid grid-cols-1 lg:grid-cols-3 xl:grid-cols-5 gap-4 ${className}`}
			>
				{totalList.map((item, index) => (
					<div key={index} className="border border-base-border rounded-md py-4 px-6">
						<span className="font-bold text-base-gray pb-1">{item.label}</span>
						<div className="text-2xl font-medium text-base-text">
							{item?.render?.() || item.value}
						</div>
					</div>
				))}
			</div>
			<FilterTable
				key={stage}
				columns={getColumnsByStage(stage)}
				hasFilterField={false}
				pagination={false}
				className={'rounded-lg'}
				dataSource={getDataSourceByStage()}
				onRow={(row) => ({
					onClick: () => {
						handleGotoLoanDetails(row);
					},
				})}
			/>
		</div>
	);
};
export default DifferentStatusLoans;
