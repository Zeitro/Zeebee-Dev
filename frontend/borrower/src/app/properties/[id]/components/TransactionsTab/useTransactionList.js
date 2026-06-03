'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';

/**
 * 房产详情 Transactions Tab：拉取 `/property/:id/borrower_property_transactions` 列表。
 * 供卡片头部进度、时间线展示；详情内上传/编辑后调用 `refreshTransactionList` 同步列表。
 *
 * @param {number|string|null|undefined} propertyId
 * @returns {{ list: object[], loading: boolean, refreshTransactionList: () => Promise<void> }}
 */
export function useTransactionList(propertyId) {
	const { request, cancelMatch } = useMemo(() => createRequest(), []);
	const transactionListApi = useMemo(
		() => `/property/${propertyId}/borrower_property_transactions`,
		[propertyId],
	);
	const [loading, setLoading] = useState(false);
	const [list, setList] = useState([]);

	const fetchTransactionList = useCallback(
		async ({ silent = false } = {}) => {
			if (!propertyId) {
				setList([]);
				return;
			}
			try {
				if (!silent) setLoading(true);
				cancelMatch(transactionListApi);
				const res = await request(transactionListApi);
				setList(Array.isArray(res?.list) ? res.list : []);
			} catch {
				setList([]);
			} finally {
				if (!silent) setLoading(false);
			}
		},
		[cancelMatch, propertyId, request, transactionListApi],
	);

	/** 详情内变更后静默刷新列表，避免整页 Spin 闪烁 */
	const refreshTransactionList = useCallback(
		() => fetchTransactionList({ silent: true }),
		[fetchTransactionList],
	);

	useEffect(() => {
		if (!propertyId) return;
		let cancelled = false;
		void (async () => {
			await Promise.resolve();
			if (cancelled) return;
			await fetchTransactionList();
		})();
		return () => {
			cancelled = true;
		};
	}, [fetchTransactionList, propertyId]);

	return { list, loading, refreshTransactionList };
}
