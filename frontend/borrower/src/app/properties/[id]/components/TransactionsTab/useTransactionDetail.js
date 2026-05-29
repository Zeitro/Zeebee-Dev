'use client';

import { useCallback, useMemo, useState } from 'react';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { createRequest } from '@repo/shared-methods/utils';

/**
 * 展开卡片时拉取交易详情，与 myTransactions 详情页共用 `/borrower_property_transaction/:id`。
 * @param {number|string|null|undefined} transactionId
 * @param {boolean} enabled 为 false 时不请求
 * @param {object} [initialRowData] 列表项数据，请求完成前用于展示 checklist 等字段
 */
export function useTransactionDetail(transactionId, enabled, initialRowData = {}) {
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState(initialRowData);

	const getDetails = useCallback(async () => {
		if (transactionId == null || transactionId === '') return;
		const api = `/borrower_property_transaction/${transactionId}`;
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data ?? {});
		} catch {
			setRowData(initialRowData ?? {});
		} finally {
			setLoading(false);
		}
	}, [initialRowData, request, transactionId]);

	useDeferredSyncEffect({
		predicate: () => enabled && transactionId != null && transactionId !== '',
		deferred: getDetails,
		deps: [enabled, getDetails, transactionId],
		delay: 0,
	});

	return { rowData, loading, getDetails };
}
