import React, { useState, useCallback } from 'react';
import { Row, Spin, Col, Space } from 'antd';
import CardItem from './Item';
import { BtnPopUpModal } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import ViewMore from './ViewMore';
import AddModal from './Add';
import { RightCircleTwoTone } from '@ant-design/icons';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import _ from 'lodash';

// 选择property (用于loan申请中)
const MoreCard = ({ disabled = false, className, id, changeFormValue = () => {} }) => {
	const { request } = createRequest();
	const [list, setList] = useState([]);
	const [rowData, setRowData] = useState([]);
	const [loading, setLoading] = useState(true);

	const moveItemToFirst = (array, predicate) => {
		// 找到并移除匹配的项
		const removedItems = _.remove(array, predicate);
		// 将移除的项添加到数组开头
		if (removedItems.length > 0) {
			array.unshift(...removedItems);
		}
		return array;
	};

	const getList = useCallback(async (propertyId) => {
		try {
			setLoading(true);
			const res = await request('/property/all', {
				method: 'GET',
			});
			let lists = [...(res?.lists || [])];
			if (propertyId) {
				lists = moveItemToFirst([...(res?.lists || [])], (item) => item.id === propertyId);
			}
			setList(lists);
			setLoading(false);
			return lists;
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	const handleClickItem = (propertyItem, isNoChange) => {
		let arr = [...list];
		if (!isNoChange) {
			arr = moveItemToFirst([...list], (item) => item.id === propertyItem?.id);
		}
		setList(arr);
		setRowData(propertyItem);
	};

	useDeferredSyncEffect({
		deferred: () => {
			getList(id);
		},
		deps: [id],
		delay: 500,
	});

	useDeferredSyncEffect({
		predicate: () => rowData?.id,
		deferred: () => {
			changeFormValue(rowData);
		},
		deps: [rowData],
		delay: 0,
	});
	return (
		<Col className={['w-full', className].join(' ')}>
			<div className="text-base font-medium text-black-dark pb-2">choose</div>
			<Row gutter={[16, 16]} className="w-full">
				{loading ? (
					<Spin className="h-[100px] w-[100px] !flex items-center justify-center" />
				) : (
					list.slice(0, 3).map((item, index) => {
						return (
							<CardItem
								key={index}
								rowData={rowData}
								disabled={disabled}
								item={item}
								handleClickItem={(item) => handleClickItem(item, true)}
							/>
						);
					})
				)}
				<Col xl={4} xs={24} sm={24}>
					<Space direction="vertical" className="w-full">
						<BtnPopUpModal
							getContent={(setOpen) => (
								<ViewMore
									setOpen={setOpen}
									rowData={rowData}
									disabled={disabled}
									list={list}
									handleClickItem={handleClickItem}
								/>
							)}
							type="default"
							size="large"
							className={[
								'w-full max-w-40 block !h-[49px]',
								!disabled && '!text-base-primary',
							].join(' ')}
							disabled={disabled}
						>
							View more
							<RightCircleTwoTone />
						</BtnPopUpModal>
						<BtnPopUpModal
							getContent={(setOpen) => (
								<AddModal
									setOpen={setOpen}
									handleSearch={async (id) => {
										const arr = await getList(id);
										const changeItem = arr.find((item) => item.id === id);
										setRowData(changeItem);
									}}
								/>
							)}
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
					</Space>
				</Col>
			</Row>
		</Col>
	);
};
export default MoreCard;
