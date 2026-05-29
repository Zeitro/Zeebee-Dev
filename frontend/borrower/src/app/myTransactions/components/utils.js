export function getStatusTagClass(status) {
	if (status === 'Pending') {
		return 'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]';
	}
	if (status === 'Cancelled') {
		return 'border-base-error bg-gray-light text-base-error';
	}
	return 'border-base-border bg-gray-light text-foreground';
}
