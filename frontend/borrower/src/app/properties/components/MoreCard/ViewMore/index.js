import React from 'react';
import PropertySvg from '@/assets/images/property.svg';
import { Row } from 'antd';
import { Col } from 'antd';
import { Button } from 'antd';
import { Space } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { DeleteOutlined } from '@ant-design/icons';
import { useState } from 'react';
import { NewModal } from '@repo/shared-methods';
import CardItem from '../Item';
const ViewMore = ({ setOpen, setRowData, rowData, disabled }) => {
	const [list, setList] = useState([{}, {}, {}]);
	return (
		<NewModal title="Choose a property" onCancel={() => setOpen(false)}>
			<Row
				gutter={[16, 16]}
				className="min-h-[100px] max-h-[300px] overflow-hidden overflow-y-auto"
			>
				{list.map((item, index) => {
					return (
						<CardItem
							key={index}
							md={12}
							isShowOperate={false}
							setRowData={setRowData}
							rowData={rowData}
							disabled={disabled}
						/>
					);
				})}
			</Row>
		</NewModal>
	);
};
export default ViewMore;
