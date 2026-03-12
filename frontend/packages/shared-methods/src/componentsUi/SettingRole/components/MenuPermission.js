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
	const [checkedKeys, setCheckedKeys] = useState([]);
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

	// 构建 parentMap 和 childrenMap，便于快速查父/子
	const { parentMap, childrenMap } = useMemo(() => {
		const pMap = {};
		const cMap = {};
		const walk = (nodes, parent) => {
			nodes.forEach((n) => {
				if (parent) {
					pMap[n[keyName]] = parent[keyName];
				}
				if (n[childrenName] && n[childrenName].length) {
					cMap[n[keyName]] = n[childrenName].map((c) => c[keyName]);
					walk(n[childrenName], n);
				} else {
					if (parent) {
						cMap[parent[keyName]] = cMap[parent[keyName]] || [];
					}
				}
			});
		};
		walk(treeData, null);
		return { parentMap: pMap, childrenMap: cMap };
	}, [treeData, keyName, childrenName]);

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

	// 勾选逻辑：使用 checkStrictly，让 antd 不自动联动父子；
	// 我们手动处理：
	// - 勾选父级：只勾选父级，不影响子级
	// - 勾选子级：同时确保该子级的父级被勾选
	const handleCheck = (checkedKeysParam, event) => {
		const nodeKey = event?.node?.[keyName];
		const isChecked = event?.checked;
		const currentCheckedArr = Array.isArray(checkedKeysParam)
			? checkedKeysParam
			: (checkedKeysParam && checkedKeysParam.checked) || [];

		// 基于 antd 返回的当前 checked 列表作为基础
		const newSet = new Set(currentCheckedArr || []);

		if (isChecked) {
			// 勾选节点：添加节点本身
			if (nodeKey !== undefined && nodeKey !== null) newSet.add(nodeKey);
			// 如果是子节点，确保父节点也被勾选
			const parent = parentMap[nodeKey];
			if (parent) newSet.add(parent);
		} else {
			// 取消勾选节点：移除该节点
			if (nodeKey !== undefined && nodeKey !== null) newSet.delete(nodeKey);

			// 如果取消的是父节点，则需要取消其所有后代子节点
			const children = childrenMap[nodeKey] || [];
			if (children && children.length) {
				const stack = [...children];
				while (stack.length) {
					const c = stack.pop();
					newSet.delete(c);
					const cc = childrenMap[c];
					if (cc && cc.length) stack.push(...cc);
				}
			}

			// 如果取消的是子节点，不再影响父节点（满足“取消最后一个子级不影响父级”的要求）
		}

		const finalArr = Array.from(newSet);
		setCheckedKeys(finalArr);
		const returnArr = finalArr?.map((i) => ({
			resource_id: i,
			data_scope: dataScope[i] || 'All',
		}));
		onChange?.(returnArr);
	};

	return (
		<Tree
			checkable
			checkStrictly={true}
			multiple
			expandedKeys={expandedKeys}
			checkedKeys={checkedKeys}
			fieldNames={{ title: titleName, key: keyName, children: childrenName }}
			treeData={treeData}
			onExpand={(expandedKeys) => {
				setExpandedKeys(expandedKeys);
			}}
			onCheck={handleCheck}
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
