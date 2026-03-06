'use client';
import { useDeferredSyncEffect, useDisUpdateEffect } from '@repo/shared-methods/hooks';
import { Radio, Tree } from 'antd';
import { useEffect, useState, useMemo } from 'react';

const MenuPermission = ({
	value,
	onChange,
	treeData = [],
	titleName = 'label',
	keyName = 'value',
	childrenName = 'children',
	isSelectAll, // 是否全选，默认不选
	isFold, // 是否折叠，默认展开
	...options
}) => {
	const [checkedKeys, setCheckedKeys] = useState();
	const [expandedKeys, setExpandedKeys] = useState([]);
	const [dataScope, setDataScope] = useState({});

	// valueArr获取所有选项的值，方便全选；expandedArr获取所有有子级的选项的值，方便展开；
	const getAllValueAndExpandedKeys = (list) => {
		let valueArr = [];
		let expandedArr = [];
		list.map((item) => {
			valueArr.push(item[keyName]);
			if (item?.[childrenName]?.length) {
				expandedArr.push(item[keyName]);
				const childrenInfo = getAllValueAndExpandedKeys(item[childrenName]);
				if (childrenInfo?.valueArr?.length) {
					valueArr = valueArr.concat(childrenInfo?.valueArr);
				}
				if (childrenInfo?.expandedArr?.length) {
					expandedArr = expandedArr.concat(childrenInfo?.expandedArr);
				}
			}
		});
		return {
			valueArr,
			expandedArr,
		};
	};
	// eslint-disable-next-line react-hooks/preserve-manual-memoization
	const info = useMemo(() => getAllValueAndExpandedKeys(treeData), [treeData]);

	// 详情时的初始化数据
	useEffect(() => {
		const updateValues = () => {
			let dataScope = {};
			const checkedList =
				value?.map((i) => {
					dataScope[i.resource_id] = i.data_scope;
					return i.resource_id;
				}) || [];
			setCheckedKeys(checkedList);
			setDataScope(dataScope);
		};
		updateValues();
	}, [value]);

	// 展开/收缩
	useDeferredSyncEffect({
		deferred: () => {
			if (isFold) {
				setExpandedKeys([]);
			} else {
				setExpandedKeys(info.expandedArr);
			}
		},
		deps: [isFold, info],
		delay: 0,
	});

	// 全选/取消全选
	useDisUpdateEffect(() => {
		const updateValues = () => {
			let arr = [];
			if (isSelectAll) {
				arr = info?.valueArr || [];
				const checkedList =
					arr?.map((i) => {
						return {
							resource_id: i,
							data_scope: 'All',
						};
					}) || [];
				setCheckedKeys(arr);
				onChange?.(checkedList);
			} else {
				arr = [];
				setCheckedKeys([]);
				onChange?.([]);
			}
		};
		updateValues();
	}, [isSelectAll, info]);

	return (
		<Tree
			checkable
			multiple
			expandedKeys={expandedKeys}
			checkedKeys={checkedKeys}
			fieldNames={{ title: titleName, key: keyName, children: childrenName }}
			treeData={treeData}
			onExpand={(expandedKeys) => {
				setExpandedKeys(expandedKeys);
			}}
			onCheck={(checkedArr) => {
				setCheckedKeys(checkedArr);
				const returnArr = checkedArr?.map((i) => ({
					resource_id: i,
					data_scope: dataScope[i] || 'All',
				}));
				onChange?.(returnArr);
			}}
			titleRender={(item) => {
				const title = (
					<>
						{item[titleName]}
						{item?.resource_type ? (
							<span className="text-base-label pl-1">({item.resource_type})</span>
						) : (
							''
						)}
					</>
				);
				return item?.enable_data_scope ? (
					<div className="flex items-center justify-between px-2 border-b border-base-border">
						<span>{title}</span>
						<Radio.Group
							options={[
								{
									label: 'Personal data only',
									value: 'Own',
								},
								{
									label: 'All data',
									value: 'All',
								},
							]}
							value={dataScope[item[keyName]]}
							onChange={(e) => {
								setDataScope((v) => {
									v[item[keyName]] = e.target.value;
									return v;
								});
								const returnArr = checkedKeys?.map((i) => ({
									resource_id: i,
									data_scope: e.target.value,
								}));
								onChange?.(returnArr);
							}}
							size="small"
							className="text-base-gray"
						/>
					</div>
				) : (
					title
				);
			}}
			{...options}
		/>
	);
};
export default MenuPermission;
