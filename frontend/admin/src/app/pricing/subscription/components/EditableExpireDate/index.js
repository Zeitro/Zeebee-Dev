'use client';

import React, { useMemo, useState } from 'react';
import { DatePicker, message } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import {
	createRequest,
	dayjs,
	displayValue,
	formatDate,
	formatStr,
} from '@repo/shared-methods/utils';

/**
 * 订阅过期日行内编辑；PATCH `/organization_subscription/{id}/expire_date`，body.expire_date 为 MM/DD/YYYY。
 *
 * @param {object} props
 * @param {object} props.record 当前行数据，需含 id、expire_date
 * @param {() => void} [props.successCallback] 保存成功后刷新列表
 */
const EditableExpireDate = ({ record, successCallback }) => {
	const [isEditing, setIsEditing] = useState(false);
	const [open, setOpen] = useState(false);
	const [loading, setLoading] = useState(false);
	const [messageApi, contextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const expireDate = record?.expire_date;

	const pickerValue = useMemo(() => {
		if (!expireDate) return null;
		const parsed = dayjs(expireDate, formatStr, true);
		if (parsed.isValid()) return parsed;
		const fallback = dayjs(expireDate);
		return fallback.isValid() ? fallback : null;
	}, [expireDate]);

	const stopEdit = () => {
		setOpen(false);
		setIsEditing(false);
	};

	const startEdit = () => {
		setIsEditing(true);
		setOpen(true);
	};

	const handleOpenChange = (nextOpen) => {
		setOpen(nextOpen);
		if (!nextOpen) {
			setIsEditing(false);
		}
	};

	const handleChange = async (_, dateString) => {
		if (!dateString) return;
		const currentFormatted = expireDate ? formatDate(expireDate) : '';
		if (dateString === currentFormatted) {
			stopEdit();
			return;
		}
		try {
			setLoading(true);
			await request(`/organization_subscription/${record.id}/expire_date`, {
				method: 'PATCH',
				body: { expire_date: dateString },
			});
			messageApi.success('Expire date updated successfully.');
			stopEdit();
			successCallback?.();
		} catch {
			// 请求失败保留编辑态，便于重选
		} finally {
			setLoading(false);
		}
	};

	return (
		<>
			{contextHolder}
			<div
				onClick={(e) => {
					e.stopPropagation();
				}}
			>
				{isEditing ? (
					<DatePicker
						open={open}
						onOpenChange={handleOpenChange}
						format={formatStr}
						value={pickerValue}
						onChange={handleChange}
						allowClear={false}
						disabled={loading}
						autoFocus
						size="small"
						style={{ width: 120 }}
						getPopupContainer={() => document.body}
					/>
				) : (
					<div
						className="inline-flex items-center gap-1 cursor-pointer"
						onClick={startEdit}
					>
						<span>{displayValue(expireDate, formatDate)}</span>
						<EditOutlined
							className="text-sm shrink-0"
							style={{ color: 'var(--base-primary)' }}
						/>
					</div>
				)}
			</div>
		</>
	);
};

export default EditableExpireDate;
