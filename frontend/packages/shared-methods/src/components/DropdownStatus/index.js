'use client';
import { Dropdown, message, Tooltip } from 'antd';
import React, { useMemo, useState } from 'react';
import { DownOutlined } from '@ant-design/icons';
import styles from './index.module.css';
import { createRequest } from '../../lib/fetch';
/**
 * 状态修改
 * data 当前选中数据
 * successCallback 状态变更后的回调函数
 * items 状态列表
 * statusKey 状态字段（接口参数字段名）
 * disabled 是否禁用
 * api 接口
 * */
export const DropdownStatus = ({
	data,
	successCallback,
	items,
	statusKey,
	tooltipProps = {},
	disabled,
	api,
	...props
}) => {
	const [loading, setLoading] = useState(false);
	const [statusValue, setStatusValue] = useState(data[statusKey] || '');
	const { request, cancel } = createRequest();
	const rowData = useMemo(() => {
		return items.find((i) => i.key === statusValue) || {};
	}, [statusValue, items]);

	const handleSelect = async ({ key }) => {
		try {
			cancel(api);
			setLoading(true);
			setStatusValue(key);
			await request(api, {
				method: 'PATCH',
				body: { [statusKey]: key },
			});
			setLoading(false);
			message.open({
				type: 'success',
				content: 'Update successful！',
			});
			successCallback?.();
		} catch (error) {
			setLoading(false);
		}
	};

	return (
		<div
			onClick={(domEvent) => {
				domEvent?.stopPropagation();
				domEvent?.preventDefault();
			}}
		>
			<Tooltip {...tooltipProps} title={disabled ? tooltipProps?.title : null}>
				<Dropdown.Button
					icon={disabled ? null : <DownOutlined />}
					loading={loading}
					menu={{
						items,
						selectable: true,
						selectedKeys: [statusValue],
						onSelect: handleSelect,
						onClick: ({ domEvent }) => {
							domEvent?.stopPropagation();
							domEvent?.preventDefault();
						},
					}}
					size="small"
					disabled={disabled}
					{...props}
					className={[
						styles.loanStatus,
						disabled && styles.disabled,
						props?.className,
					].join(' ')}
					onClick={(domEvent) => {
						domEvent?.stopPropagation();
						domEvent?.preventDefault();
					}}
				>
					{statusValue ? rowData.label : 'Select Status'}
				</Dropdown.Button>
			</Tooltip>
		</div>
	);
};
