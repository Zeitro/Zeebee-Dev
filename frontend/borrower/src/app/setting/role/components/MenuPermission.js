'use client';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { Checkbox, Tree } from 'antd';
import { useEffect, useState } from 'react';

const MenuPermission = ({
	value,
	onChange,
	treeData = [],
	titleName = 'label',
	keyName = 'value',
	childrenName = 'children',
	isFold = false,
	isSelectAll = false,
	...options
}) => {
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
			// const info = getAllValueAndExpandedKeys(treeData);
			// if (operationsValue.includes('fold')) {
			// 	setExpandedKeys([]);
			// } else {
			// 	setExpandedKeys(info.expandedArr);
			// }
			const info = getAllValueAndExpandedKeys(treeData);
			let arr = [];
			if (isSelectAll) {
				arr = info?.valueArr || [];
			} else {
				arr = [];
			}

			if (isFold) {
				setExpandedKeys([]);
			} else {
				setExpandedKeys(info.expandedArr);
			}
			setCheckedKeys(arr);
			onChange?.(arr);
		},
		deps: [isFold, isSelectAll, treeData],
		delay: 0,
	});

	// useEffect(() => {
	// 	// setOperationsValue(e);
	// 	const info = getAllValueAndExpandedKeys(treeData);
	// 	let arr = [];
	// 	if (isSelectAll) {
	// 		arr = info?.valueArr || [];
	// 	} else {
	// 		arr = [];
	// 	}

	// 	if (isFold) {
	// 		setExpandedKeys([]);
	// 	} else {
	// 		setExpandedKeys(info.expandedArr);
	// 	}
	// 	setCheckedKeys(arr);
	// 	onChange?.(arr);
	// }, [isFold, isSelectAll, treeData]);

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
			onCheck={(checkedKeys) => {
				setCheckedKeys(checkedKeys);
				onChange?.(checkedKeys);
			}}
			{...options}
		/>
	);
};
export default MenuPermission;
