'use client';
import { Button, Dropdown, Space } from 'antd';
import React, { useMemo, useState } from 'react';
import { EditOutlined, EllipsisOutlined, DeleteOutlined, ProfileOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';

/**
 * Table内的操作组件
 * list: 下拉列表，传值就不包含edit和delete，不传默认用edit/delete
 * getItems: 获取下拉列表的函数，可拿到edit/delete配置，重新组装顺序或在edit/delete前后添加其他项
 * onSelect：选择下拉项的回调函数
 * sliceCount: null横向展示所有；传数字则横向展示指定个数，多出的用...hover显示下拉菜单
 */
export const TableOperation = ({ list, className, onSelect, getItems, disabled, sliceCount }) => {
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
	const info = useMemo(() => {
		const newItems = getItems?.(editItem, deleteItem, detailsItem) || [editItem, deleteItem];
		const arr = list || newItems;
		const count = sliceCount === null ? arr.length : sliceCount;
		const showList = arr.slice(0, count);
		return {
			showList,
			moreList: arr.slice(count, arr.length - showList.length),
		};
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [list, sliceCount]);

	const handleSelect = async ({ key }) => {
		setSelectedKey(key);
		const arr = info.showList.concat(info.moreList);
		const optItem = arr.find((i) => i.key === key);
		await onSelect?.(key);
		if (optItem?.url) {
			router.push(optItem.url);
		}
	};

	return (
		<Space
			size={0}
			onClick={(domEvent) => {
				domEvent?.stopPropagation();
				domEvent?.preventDefault();
			}}
			className="[&_.ant-btn]:!gap-1"
		>
			{info.showList.map((item, index) => {
				const { key, icon, label, disabled, ...props } = item;
				return (
					item.children ?? (
						<Button
							size="small"
							type="link"
							key={index}
							icon={icon}
							disabled={disabled}
							onClick={(domEvent) => handleSelect({ key, domEvent })}
							{...props}
							className={['!h-[22px]', props.className].join(' ')}
						>
							{label}
						</Button>
					)
				);
			})}
			{info.moreList?.length > 0 && (
				<Dropdown
					menu={{
						items: info.moreList,
						selectable: true,
						selectedKeys: [selectedKey],
						onClick: handleSelect,
					}}
					className={[disabled && '!text-base-disabled', className].join(' ')}
					disabled={disabled}
				>
					<EllipsisOutlined />
				</Dropdown>
			)}
		</Space>
	);
};
