'use client';

import { DropdownStatus } from '@repo/shared-methods';
import styles from './index.module.scss';
import { useOptionsStore } from '@repo/shared-methods/store';

/**
 * 组织订阅列表状态切换；PATCH `/organization_subscription/{id}/status`。
 */
const ChangeStatus = ({ data, successCallback, ...props }) => {
	const { subscriptionStatusOptions } = useOptionsStore();
	const newItems = subscriptionStatusOptions.map((item) => ({
		label: item.label,
		key: item.value,
	}));

	return (
		<DropdownStatus
			data={data}
			successCallback={successCallback}
			items={newItems}
			statusKey="status"
			api={`/organization_subscription/${data.id}/status`}
			{...props}
			className={[styles?.[data?.status], props?.className].filter(Boolean).join(' ')}
		/>
	);
};

export default ChangeStatus;
