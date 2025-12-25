'use client';
import { Button, Dropdown, Space } from 'antd';
import React, { useState } from 'react';
import { EditOutlined, EllipsisOutlined, DeleteOutlined, ProfileOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';

/**
 * Table内的操作组件
 * list: 下拉列表，传值就不包含edit和delete，不传默认用edit/delete
 * getItems: 获取下拉列表的函数，可拿到edit/delete配置，重新组装顺序或在edit/delete前后添加其他项
 * onSelect：选择下拉项的回调函数
 */
export const TableOperation = ({ list, className, onSelect, getItems, disabled }) => {
	const editItem = {
		label: 'Edit',
		key: 'edit',
		icon: <EditOutlined />,
	};
	const deleteItem = {
		label: 'Delete',
		key: 'delete',
		icon: <DeleteOutlined />,
		danger: true,
	};
	const detailsItem = {
		label: 'Details',
		key: 'details',
		icon: <ProfileOutlined />,
	};
	const [selectedKey, setSelectedKey] = useState(null);
	const router = useRouter();
	const newItems = getItems?.(editItem, deleteItem, detailsItem) || [editItem, deleteItem];

	const handleSelect = async ({ key, domEvent, item }) => {
		domEvent?.stopPropagation();
		domEvent?.preventDefault();
		// setSelectedKey(key);
		// const arr = list || newItems;
		// const optItem = arr.find((i) => i.key === key);
		await onSelect?.(key);
		if (item?.url) {
			router.push(item.url);
		}
	};

	return (
		// <Dropdown
		// 	menu={{
		// 		items: newItems.length && !list?.length ? newItems : [...(list || [])],
		// 		selectable: true,
		// 		selectedKeys: [selectedKey],
		// 		onClick: handleSelect,
		// 	}}
		// 	className={[styles.tableOperation, disabled && '!text-base-disabled', className].join(
		// 		' ',
		// 	)}
		// 	disabled={disabled}
		// 	onClick={(domEvent) => {
		// 		domEvent?.stopPropagation();
		// 		domEvent?.preventDefault();
		// 	}}
		// >
		// 	<EllipsisOutlined />
		// </Dropdown>
		<Space
			size={0}
			onClick={(domEvent) => {
				domEvent?.stopPropagation();
				domEvent?.preventDefault();
			}}
			className="[&_.ant-btn]:!gap-1"
		>
			{(list || newItems).map((item, index) => {
				const { key, icon, label, disabled, ...props } = item;
				return (
					item.children || (
						<Button
							size="small"
							type="link"
							key={index}
							icon={icon}
							disabled={disabled}
							onClick={(domEvent) => handleSelect({ key, domEvent, item })}
							{...props}
							className={['!h-[22px]', props.className].join(' ')}
						>
							{label}
						</Button>
					)
				);
			})}
		</Space>
	);
};
