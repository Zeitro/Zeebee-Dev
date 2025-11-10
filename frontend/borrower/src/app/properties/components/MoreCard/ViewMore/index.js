import React, { useState } from 'react';
import { Row } from 'antd';
import { NewModal } from '@repo/shared-methods';
import CardItem from '../Item';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const ViewMore = ({ setOpen, rowData, disabled, list, handleClickItem }) => {
	const [data, setData] = useState({});
	const handleOk = () => {
		handleClickItem(data);
		setOpen(false);
	};

	useDeferredSyncEffect({
		deferred: () => {
			setData(rowData);
		},
		deps: [rowData],
		isUseMicrotask: true,
	});
	return (
		<NewModal
			title="Choose a property"
			width={750}
			onCancel={() => setOpen(false)}
			onOk={handleOk}
		>
			<Row
				gutter={[16, 16]}
				className="min-h-[100px] max-h-[60vh] overflow-hidden overflow-y-auto"
			>
				{list.map((item, index) => {
					return (
						<CardItem
							key={index}
							md={12}
							xl={12}
							isShowOperate={false}
							rowData={data}
							disabled={disabled}
							item={item}
							handleClickItem={(item) => setData(item)}
						/>
					);
				})}
			</Row>
		</NewModal>
	);
};
export default ViewMore;
