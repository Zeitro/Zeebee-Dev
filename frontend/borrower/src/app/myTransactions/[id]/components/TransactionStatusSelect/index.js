'use client';

import React, { useCallback, useMemo, useState } from 'react';
import { Select, message } from 'antd';
import { createRequest } from '@repo/shared-methods/utils';
import { getStatusSelectStyleKey } from '../../../components/utils';
import styles from './index.module.scss';

/**
 * PATCH 交易 status；成功后 toast 并执行 onSuccess（如刷新详情）。
 *
 * @param {string|number} transactionId
 * @param {{ onSuccess?: () => void | Promise<void> }} [options]
 * @returns {{ updateStatus: (status: string) => Promise<void>, loading: boolean }}
 */
export function useTransactionStatusUpdate(transactionId, { onSuccess } = {}) {
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);

	const updateStatus = useCallback(
		async (status) => {
			try {
				setLoading(true);
				await request(`/borrower_property_transaction/${transactionId}/status`, {
					method: 'PATCH',
					body: { status },
				});
				message.open({
					type: 'success',
					content: 'Status updated',
				});
				await onSuccess?.();
			} catch {
				// 错误提示由 request 层统一处理
			} finally {
				setLoading(false);
			}
		},
		[onSuccess, request, transactionId],
	);

	return { updateStatus, loading };
}

/**
 * 交易状态下拉；选项来自 borrower_property_transaction_status，样式随当前 status 变色。
 *
 * @param {object} props
 * @param {string} [props.value] 当前 status 文案（如 Pending）
 * @param {boolean} [props.loading] 外部 loading（与 useTransactionStatusUpdate 配合）
 * @param {(status: string) => void} [props.onChange]
 * @param {string} [props.className]
 * @param {string} [props.placeholder]
 * @param {boolean} [props.disabled]
 */
export default function TransactionStatusSelect({
	value,
	onSuccess,
	className = '',
	placeholder = 'Status',
	disabled = false,
	id,
	options = [],
	...props
}) {
	const { updateStatus, loading } = useTransactionStatusUpdate(id, {
		onSuccess,
	});
	const statusStyleKey = getStatusSelectStyleKey(value);

	return (
		<Select
			className={[styles[statusStyleKey], className].filter(Boolean).join(' ')}
			value={value}
			placeholder={placeholder}
			options={options}
			loading={loading}
			onChange={updateStatus}
			allowClear={false}
			disabled={disabled}
			{...props}
		/>
	);
}
