import React from 'react';
import { Row, Col, Button, Space } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import HomeSvg from '@/assets/images/home2.svg';
import { assemblyAddress, formatNumber } from '@repo/shared-methods/utils';
const CardItem = ({ isShowOperate = true, rowData, item, disabled, handleClickItem, ...props }) => {
	return (
		<Col
			xl={6}
			// lg={list.length > 2 ? 24 : 6}
			sm={24}
			xs={24}
			{...props}
		>
			<div
				onClick={() => {
					handleClickItem(item);
				}}
				style={{
					borderColor:
						rowData?.id && rowData?.id === item?.id ? '#0958D9' : 'var(--base-border)',
				}}
				className={[
					'border border-base-border rounded-md pl-2 py-3 pr-4',
					!disabled &&
						'cursor-pointer hover:!border-[#0958D9] [&:hover_.common-fill]:![color:#0958D9]',
				].join(' ')}
			>
				<div className="flex items-center gap-2">
					<div
						className={['common-fill text-[rgba(0,0,0,0.45)]'].join(' ')}
						style={{
							color:
								rowData?.id && rowData?.id === item?.id
									? '#0958D9'
									: 'rgba(0,0,0,0.45)',
						}}
					>
						<HomeSvg />
					</div>
					<p className="common-ellipsis text-base font-medium">
						{assemblyAddress(item?.address)}
					</p>
				</div>
				<Row className="pl-2 pt-1 text-[rgba(0, 0, 0, 0.65)]">
					<Col md={10} className="common-ellipsis">
						{item?.property_type || '--'}
					</Col>
					<Col md={7} className="common-ellipsis">
						{item?.unit ? `${item?.unit} units` : '--'}
					</Col>
					<Col md={7} className="common-ellipsis">
						{formatNumber(item?.living_area_sqft)}
					</Col>
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
