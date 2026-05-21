import { displayValue } from '@repo/shared-methods/utils';

export const TRANSACTION_DETAIL_TABS = ['documents', 'details', 'loans', 'log'];
export const TRANSACTION_DETAIL_DEFAULT_TAB = 'documents';

export function resolveTransactionDetailTab(tab) {
	return TRANSACTION_DETAIL_TABS.includes(tab) ? tab : TRANSACTION_DETAIL_DEFAULT_TAB;
}

export const TRANSACTION_STATUS_OPTIONS = [
	{ label: 'Pending', value: 'Pending' },
	{ label: 'Cancelled', value: 'Cancelled' },
	{ label: 'Closed', value: 'Closed' },
];

export function formatTransactionId(item) {
	if (item?.transaction_no) return item.transaction_no;
	if (item?.id != null && item?.id !== '') return `TC${item.id}`;
	return displayValue(null);
}

export function getTransactionTitle(item) {
	return item?.type ? `${item.type} transaction` : displayValue(null);
}

export function getStatusTagClass(status) {
	if (status === 'Pending') {
		return 'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]';
	}
	if (status === 'Cancelled') {
		return 'border-base-error bg-gray-light text-base-error';
	}
	return 'border-base-border bg-gray-light text-foreground';
}
