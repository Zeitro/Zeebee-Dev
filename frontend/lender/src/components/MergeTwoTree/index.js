import React, { useEffect, useRef, useState } from 'react';
import { Col, Row, Tree } from 'antd';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

/**
 * 多层嵌套多选的tree
 * （可横向展示两个tree，值合并为一个数组，只有一个options时，默认展示一个tree）
 * */
const MergeTwoTree = ({ options, options2, onChange, value }) => {
	// 根据传入的options列表过滤值
	const filterArr = (arr = [], options = []) => {
		let newValue = [];
		const getOptions = (list) => {
			list.map((i) => {
				if (i.children) {
					getOptions(i.children);
				}
				newValue.push(i.key);
				return i;
			});
			return newValue;
		};
		const keys = getOptions(options);
		return arr?.filter((i) => keys.includes(i)) || [];
	};

	// 组装tree数据
	const getTreeList = (list) => {
		return list?.map((i) => {
			i.key = i.key || i?.label?.replace(/ /g, '');
			if (i.children) {
				i.children = getTreeList(i.children);
			}
			return i;
		});
	};

	// 封装tree
	const getTree = (treeData = [], options = {}) => {
		return (
			<Tree
				checkable
				defaultExpandAll
				multiple
				autoExpandParent={true}
				fieldNames={{ title: 'label' }}
				switcherIcon=""
				titleRender={({ label, other, key }) => {
					const isChecked = options.checkedKeys.includes(key);
					return (
						<div className="flex items-center justify-center gap-8">
							<span>{label}</span>
							{other && other(isChecked)}
						</div>
					);
				}}
				treeData={treeData}
				{...options}
			/>
		);
	};

	const newOptions = getTreeList(options);
	const newOptions2 = getTreeList(options2);

	const [leftKeys, setLeftKeys] = useState([]);
	const [rightKeys, setRightKeys] = useState([]);
	const isInit = useRef(true);
	// 只有初始化拿到传入值操作
	useDeferredSyncEffect({
		predicate: () => isInit.current,
		deferred: () => {
			isInit.current = false;
			setLeftKeys(filterArr(value, newOptions) || []);
			setRightKeys(filterArr(value, newOptions2) || []);
		},
		deps: [newOptions, newOptions2, value],
		isUseMicrotask: true,
	});

	useEffect(() => {
		onChange([...new Set([...leftKeys, ...rightKeys])]);
	}, [leftKeys, onChange, rightKeys]);

	const onSelect = (key, checkedKeys, setFun) => {
		let newCheckedKeys = [...checkedKeys];
		if (checkedKeys.includes(key)) {
			// 已勾选则取消
			newCheckedKeys = newCheckedKeys.filter((k) => k !== key);
		} else {
			// 未勾选则添加
			newCheckedKeys.push(key);
		}
		setFun(newCheckedKeys);
	};

	const treeDom = getTree(newOptions, {
		checkedKeys: leftKeys,
		selectedKeys: leftKeys,
		onCheck: (keys) => {
			setLeftKeys(keys);
		},
		onSelect: (selectedKeys, info) => onSelect(info.node.key, leftKeys, setLeftKeys),
	});
	return options?.length && !options2?.length ? (
		treeDom
	) : (
		<Row>
			<Col md={12} xs={24}>
				{treeDom}
			</Col>
			<Col md={12} xs={24}>
				{getTree(newOptions2, {
					checkedKeys: rightKeys,
					selectedKeys: rightKeys,
					onCheck: (keys) => {
						setRightKeys(keys);
					},
					onSelect: (selectedKeys, info) =>
						onSelect(info.node.key, rightKeys, setRightKeys),
				})}
			</Col>
		</Row>
	);
};
export default MergeTwoTree;
