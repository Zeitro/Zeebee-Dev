'use client';
import { Col, Row } from 'antd';
import React from 'react';
import { generalList } from '@/enum/loan';
const General = ({ rowData = {} }) => {
	return (
		<Row gutter={[40, 24]}>
			{generalList.map((i, j) => {
				return (
					<Col key={j} sm={24} xs={24} md={12}>
						<div className="border border-base-border rounded-2xl pt-6 pb-4 px-6">
							<div className="text-xl font-medium text-black-dark pl-3 pb-4 border-b border-base-border">
								{i.title}
							</div>
							<Row gutter={[24, 24]} className="px-4 py-4">
								{i.fields?.map((item, index) => (
									<Col sm={24} xs={24} md={12} key={index}>
										<div
											className="text-sm/[22px] pb-2"
											style={{ color: 'rgba(0, 0, 0, 0.65)' }}
										>
											{item.label}
										</div>
										<div className="text-xl font-medium text-black-dark">
											{item?.render
												? item?.render?.(rowData[item.value])
												: rowData[item.value]}
										</div>
									</Col>
								))}
							</Row>
						</div>
					</Col>
				);
			})}
		</Row>
	);
};

export default General;
