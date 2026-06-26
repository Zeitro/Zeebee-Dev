'use client';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { Checkbox, Tree } from 'antd';
import { useState } from 'react';

const MenuPermission = ({
	value,
	onChange,
	treeData = [],
	titleName = 'label',
	keyName = 'value',
	childrenName = 'children',
	...options
}) => {
	const [operationsValue, setOperationsValue] = useState([]);
	const [checkedKeys, setCheckedKeys] = useState(value);
	const [expandedKeys, setExpandedKeys] = useState([]);

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

	useDeferredSyncEffect({
		deferred: () => {
			const info = getAllValueAndExpandedKeys(treeData);
			if (operationsValue.includes('fold')) {
				setExpandedKeys([]);
			} else {
				setExpandedKeys(info.expandedArr);
			}
		},
		deps: [operationsValue, treeData],
		delay: 0,
	});

	return (
		<div>
			<Checkbox.Group
				options={[
					{
						label: 'Fold',
						value: 'fold',
					},
					{
						label: 'Select all',
						value: 'all',
					},
				]}
				value={operationsValue}
				onChange={(e) => {
					setOperationsValue(e);
					const info = getAllValueAndExpandedKeys(treeData);
					let arr = [];
					if (e.includes('all')) {
						arr = info?.valueArr || [];
					} else {
						arr = [];
					}
					setCheckedKeys(arr);
					onChange?.(arr);
				}}
			/>
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
				onCheck={(checkedKeys) => {
					setCheckedKeys(checkedKeys);
					onChange?.(checkedKeys);
				}}
				{...options}
			/>
		</div>
	);
};
export default MenuPermission;
