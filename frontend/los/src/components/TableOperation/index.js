'use client';
import { Dropdown } from 'antd';
import React, { useEffect, useState } from 'react';
import { EditOutlined, EllipsisOutlined, DeleteOutlined } from '@ant-design/icons';
import styles from './index.module.scss';
const editItem = {
	label: 'Edit',
	key: 'edit',
	icon: <EditOutlined />,
};
const deleteItem = {
	label: 'Delete',
	key: 'delete',
	icon: <DeleteOutlined />,
};
/**
 * Table内的操作组件
 * list: 下拉列表，传值就不包含edit和delete，不传默认用edit/delete
 * getItems: 获取下拉列表的函数，可拿到edit/delete配置，重新组装顺序或在edit/delete前后添加其他项
 * onSelect：选择下拉项的回调函数
 */
const TableOperation = ({ list, className, onSelect, getItems }) => {
	const [selectedKey, setSelectedKey] = useState(null);
	const [items, setItems] = useState([]);

	const handleSelect = ({ key }) => {
		setSelectedKey(key);
		onSelect?.(key);
	};

	useEffect(() => {
		// 默认项edit和delete，getItems和list可任选一种传入，getItems可拿到edit和delete的配置重新组装
		const newItems = getItems?.(editItem, deleteItem) || [editItem, deleteItem];
		if (newItems.length && !list?.length) {
			setItems([...newItems]);
		} else {
			setItems([...(list || [])]);
		}
	}, [list]);

	return (
		<Dropdown
			menu={{
				items,
				selectable: true,
				selectedKeys: [selectedKey],
				onSelect: handleSelect,
			}}
			className={[styles.tableOperation, className].join(' ')}
		>
			<EllipsisOutlined />
		</Dropdown>
	);
};
export default TableOperation;
