export const statusInfo = {
	Pending: {
		border: '!border-[var(--orange-medium)]',
		background: '!bg-[var(--orange-light)]',
		color: '!text-[var(--orange-dark)]',
	},
	Cancelled: {
		border: '!border-[var(--error-medium)]',
		background: '!bg-[var(--error-light)]',
		color: '!text-base-error',
	},
	Closed: {
		border: '!border-base-border',
		background: '!bg-base-text2',
		color: '!text-foreground',
	},
};
export function getStatusTagClass(status) {
	if (status === 'Pending') {
		return `${statusInfo.Pending.border} ${statusInfo.Pending.background} ${statusInfo.Pending.color}`;
	}
	if (status === 'Cancelled') {
		return `${statusInfo.Cancelled.border} ${statusInfo.Cancelled.background} ${statusInfo.Cancelled.color}`;
	}
	return `${statusInfo.Closed.border} ${statusInfo.Closed.background} ${statusInfo.Closed.color}`;
}

/** TransactionStatusSelect.module.scss 导出的类名 */
const STATUS_SELECT_STYLE_KEY = {
	Pending: 'statusSelectPending',
	Cancelled: 'statusSelectCancelled',
	Closed: 'statusSelectClosed',
};

/**
 * @param {string} [status]
 * @returns {'statusSelectPending' | 'statusSelectCancelled' | 'statusSelectClosed'}
 */
export function getStatusSelectStyleKey(status) {
	return STATUS_SELECT_STYLE_KEY[status] ?? STATUS_SELECT_STYLE_KEY.Closed;
}
