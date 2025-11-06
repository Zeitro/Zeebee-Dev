import React from 'react';
import PropertySvg from '@/assets/images/property.svg';
import { Row } from 'antd';
import { Col } from 'antd';
import { Button } from 'antd';
import { Space } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { DeleteOutlined } from '@ant-design/icons';
import { useState } from 'react';
import CardItem from './Item';
import { BtnPopUpModal } from '@repo/shared-methods';
import ViewMore from './ViewMore';
import AddModal from './Add';
import { RightCircleTwoTone } from '@ant-design/icons';

// 选择property (用于loan申请中)
const MoreCard = ({ disabled = false, className }) => {
	const [list, setList] = useState([{}, {}, {}]);
	const [rowData, setRowData] = useState([]);
	return (
		<Col className={['w-full', className].join(' ')}>
			<div className="text-base font-medium text-black-dark pb-2">choose</div>
			<Row gutter={16} className="w-full">
				{list.map((item, index) => {
					return (
						<CardItem
							key={index}
							setRowData={setRowData}
							rowData={rowData}
							disabled={disabled}
						/>
					);
				})}
				<Col md={4}>
					<BtnPopUpModal
						getContent={(setOpen) => (
							<ViewMore
								setOpen={setOpen}
								setRowData={setRowData}
								rowData={rowData}
								disabled={disabled}
							/>
						)}
						type="default"
						size="large"
						className={[
							'w-full max-w-40 block !h-[49px] mb-2 ',
							!disabled && '!text-base-primary',
						].join(' ')}
						disabled={disabled}
					>
						View more
						<RightCircleTwoTone />
					</BtnPopUpModal>
					<BtnPopUpModal
						getContent={(setOpen) => <AddModal setOpen={setOpen} />}
						type="default"
						size="large"
						className={[
							'w-full max-w-40 block !h-[49px]',
							!disabled && '!text-base-primary',
						].join(' ')}
						disabled={disabled}
					>
						Add a property
					</BtnPopUpModal>
				</Col>
			</Row>
		</Col>
	);
};
export default MoreCard;
