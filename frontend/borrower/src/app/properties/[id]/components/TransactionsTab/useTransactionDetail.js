'use client';

import { useCallback, useMemo, useState } from 'react';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { createRequest } from '@repo/shared-methods/utils';
import { buildMockTransactionDetail } from '@/app/myTransactions/[id]/components/mockTransactionDetail';

export function useTransactionDetail(transactionId, enabled) {
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [rowData, setRowData] = useState({});

	const getDetails = useCallback(async () => {
		if (transactionId == null || transactionId === '') return;
		const api = `/transaction/${transactionId}`;
		try {
			setLoading(true);
			const data = await request(api);
			setRowData(data);
		} catch {
			setRowData(buildMockTransactionDetail(transactionId));
		} finally {
			setLoading(false);
		}
	}, [request, transactionId]);

	useDeferredSyncEffect({
		predicate: () => enabled && transactionId != null && transactionId !== '',
		deferred: getDetails,
		deps: [enabled, getDetails, transactionId],
		delay: 0,
	});

	return { rowData, loading, getDetails };
}
