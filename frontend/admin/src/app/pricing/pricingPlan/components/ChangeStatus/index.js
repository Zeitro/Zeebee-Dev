'use client';

import { DropdownStatus } from '@repo/shared-methods';
import { planStatusOptions } from '../../constants';
import styles from './index.module.scss';

/**
 * 定价计划列表状态切换；样式见 index.module.scss（Active / Disabled）。
 */
const ChangeStatus = ({ data, successCallback, ...props }) => {
	const newItems = planStatusOptions.map((item) => ({
		label: item.label,
		key: item.value,
	}));

	return (
		<DropdownStatus
			data={data}
			successCallback={successCallback}
			items={newItems}
			statusKey="status"
			api={`/pricing_plan/${data.id}/status`}
			{...props}
			className={[styles?.[data?.status], props?.className].filter(Boolean).join(' ')}
		/>
	);
};

export default ChangeStatus;
