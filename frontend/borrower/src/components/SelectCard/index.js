import React, { useState, useCallback, useRef } from 'react';
import { Row, Spin, Col, Space, Form } from 'antd';
import { BtnPopUpModal } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import ViewMore from './ViewMore';
import { RightCircleTwoTone } from '@ant-design/icons';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import _ from 'lodash';
import useDom from './useDom';
import { useEffect } from 'react';
// 选择property (用于loan申请中)
const SelectCard = ({
	disabled = false,
	className,
	id,
	changeFormValue = () => {},
	type = 'property',
	api = '/property/all',
	title = 'Choose a property',
	addTitle = 'Add a property',
	name,
	required = true,
	xl = 6,
	sm = 24,
	xs = 24,
	btnCol = { xl: 6, sm: 24, xs: 24 },
	sliceIndex = 3,
}) => {
	const isInit = useRef(true);
	const { request } = createRequest();
	const [list, setList] = useState([]);
	const [rowData, setRowData] = useState([]);
	const [loading, setLoading] = useState(true);
	const { getCardItem, getAddDom } = useDom();
	const moveItemToFirst = (array, predicate) => {
		// 找到并移除匹配的项
		const removedItems = _.remove(array, predicate);
		// 将移除的项添加到数组开头
		if (removedItems.length > 0) {
			array.unshift(...removedItems);
		}
		return array;
	};

	const getList = useCallback(
		async (propertyId) => {
			try {
				setLoading(true);
				const res = await request(api, {
					method: 'GET',
				});
				let lists = [...(res?.lists || [])];
				if (propertyId) {
					lists = moveItemToFirst(
						[...(res?.lists || [])],
						(item) => item.id === propertyId,
					);
				}
				setList(lists);
				setLoading(false);
				return lists;
			} catch (error) {
				console.log(error);
				setLoading(false);
			}
		},
		// eslint-disable-next-line react-hooks/exhaustive-deps
		[api],
	);

	const handleClickItem = useCallback(
		(propertyItem, isNoChange, newList = list) => {
			let arr = [...newList];
			if (!isNoChange) {
				arr = moveItemToFirst([...arr], (item) => item.id === propertyItem?.id);
			}
			setList(arr);
			setRowData(propertyItem);
		},
		[list],
	);

	useDeferredSyncEffect({
		deferred: () => {
			getList();
		},
		deps: [api],
		delay: 0,
	});

	useEffect(() => {
		if (id && list?.length && isInit?.current) {
			const findItem = list.find((i) => i.id === id);
			handleClickItem(findItem);
			isInit.current = false;
		}
	}, [handleClickItem, id, list]);

	useDeferredSyncEffect({
		predicate: () => rowData?.id,
		deferred: () => {
			changeFormValue(rowData);
		},
		deps: [rowData],
		delay: 0,
	});

	const handleSelect = async (id) => {
		const arr = await getList(id);
		const changeItem = arr.find((item) => item.id === id);
		setRowData(changeItem);
	};

	return (
		<Col className={['w-full', className].join(' ')}>
			<div className="text-base font-medium text-black-dark pb-2">{title}</div>
			<Form.Item noStyle name={name} rules={[{ required }]}></Form.Item>
			<Row gutter={[16, 16]} className="w-full">
				{loading ? (
					<Spin className="h-[100px] w-[100px] !flex items-center justify-center" />
				) : (
					list.slice(0, sliceIndex).map((item, index) => {
						return getCardItem(type, {
							index,
							rowData,
							item,
							disabled,
							handleClickItem: (item, newList) => {
								if (disabled) return;
								handleClickItem(item, false, newList);
							},
							getList,
							xl,
							xs,
							sm,
						});
					})
				)}
				<Col xl={6} xs={24} sm={24} {...btnCol}>
					<Space direction="vertical" className="w-full">
						<BtnPopUpModal
							getContent={(setOpen) => (
								<ViewMore
									setOpen={setOpen}
									rowData={rowData}
									disabled={disabled}
									list={list}
									handleClickItem={handleClickItem}
									title={title}
								>
									{({ item, index, setData, data }) => {
										return getCardItem(type, {
											index,
											md: 12,
											xl: 12,
											isShowOperate: false,
											rowData: data,
											item,
											disabled,
											handleClickItem: (item) => setData(item),
										});
									}}
								</ViewMore>
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
							<RightCircleTwoTone
								twoToneColor={
									disabled
										? ['var(--base-disabled)', '#fff']
										: ['var(--base-primary)', '#E8F2FF']
								}
							/>
						</BtnPopUpModal>
						<BtnPopUpModal
							getContent={(setOpen) =>
								getAddDom(type, { setOpen, handleSearch: handleSelect })
							}
							type="default"
							size="large"
							className={[
								'w-full max-w-40 block !h-[49px]',
								!disabled && '!text-base-primary',
							].join(' ')}
							disabled={disabled}
						>
							{addTitle}
						</BtnPopUpModal>
					</Space>
				</Col>
			</Row>
		</Col>
	);
};
export default SelectCard;
