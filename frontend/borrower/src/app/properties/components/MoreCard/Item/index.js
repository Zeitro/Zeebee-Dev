import React from 'react';
import PropertySvg from '@/assets/images/property.svg';
import { Row } from 'antd';
import { Col } from 'antd';
import { Button } from 'antd';
import { Space } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { DeleteOutlined } from '@ant-design/icons';
import { useState } from 'react';
import HomeSvg from '@/assets/images/home2.svg';
const CardItem = ({ isShowOperate = true, setRowData, rowData, item, disabled, ...props }) => {
	const [list, setList] = useState([{}, {}, {}]);
	return (
		<Col md={6} {...props} onClick={() => setRowData({ ...item })}>
			<div
				className={[
					'border border-base-border rounded-md pl-2 py-3 pr-4',
					rowData?.id && rowData?.id === item?.id && 'border-[#0958D9]',
					!disabled &&
						'cursor-pointer hover:border-[#0958D9] [&:hover_.common-fill]:[color:#0958D9]',
				].join(' ')}
			>
				<div className="flex items-center gap-2">
					<div
						className={[
							'common-fill text-[rgba(0,0,0,0.45)]',
							rowData?.id && rowData?.id === item?.id && 'text-[#0958D9]',
						].join(' ')}
					>
						<HomeSvg />
					</div>
					<p className="common-ellipsis text-base font-medium">123 dd</p>
				</div>
				<Row className="px-2 pt-1 text-[rgba(0, 0, 0, 0.65)]">
					<Col md={8}>3units</Col>
					<Col md={8}>5 units</Col>
					<Col md={8}>1500 sqft</Col>
				</Row>
				{isShowOperate && (
					<Space className="w-full justify-end pt-2">
						<Button
							size="small"
							icon={<EditOutlined className="text-base-primary" />}
						></Button>
						<Button
							size="small"
							icon={<DeleteOutlined className="text-base-error" />}
						></Button>
					</Space>
				)}
			</div>
		</Col>
	);
};
export default CardItem;
