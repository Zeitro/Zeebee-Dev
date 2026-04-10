'use client';
import {Dropdown, message, Tooltip} from 'antd';
import React, {useMemo, useState} from 'react';
import {DownOutlined} from '@ant-design/icons';
import styles from './index.module.css';
import {createRequest, transformArrayStructure} from '@repo/shared-methods/utils';
import {useOptionsStore} from '@repo/shared-methods/store';

let api = '/api/loan/status';

/**
 * loans列表内的状态
 * status 当前状态
 * successCallback 状态变更后的回调函数
 * is3rdParty true是第三方贷款,则可修改状态,否则不可修改
 * */
const LoanPendingStatus = ({ data, successCallback }) => {
	const { loanPendingStatus } = useOptionsStore();
	const [loading, setLoading] = useState(false);
	const [statusValue, setStatusValue] = useState(data.pending_status);
	const items = transformArrayStructure({
		sourceArray: loanPendingStatus,
		entryValue: 'value',
		outputValue: 'key',
	});

	const { request, cancel } = createRequest();
	const rowData = useMemo(() => {
		return items.find((i) => i.key === statusValue) || {};
	}, [items, statusValue]);
	const is3rdParty = data?.source === 'Borrower';

	const handleSelect = async ({ key }) => {
		try {
			cancel(api);
			setLoading(true);
			setStatusValue(key);
			await request(`/loan/${data.id}/pending_status`, {
				method: 'PUT',
				body: { pending_status: key },
			});
			setLoading(false);
			message.open({
				type: 'success',
				content: 'Update successful！',
			});
			successCallback?.();
		} catch {
			setLoading(false);
		}
	};

	return (
		<Tooltip
			title={
				is3rdParty
					? ''
					: 'Zeebee-connected loans are managed automatically and cannot be edited manually by users.'
			}
		>
			<Dropdown.Button
				icon={is3rdParty ? <DownOutlined /> : null}
				loading={loading}
				menu={{
					items,
					selectable: true,
					selectedKeys: [statusValue],
					onSelect: handleSelect,
				}}
				className={[
					styles.loanStatus,
					styles[statusValue],
					!is3rdParty && styles.disabled,
				].join(' ')}
				size="small"
				disabled={!is3rdParty}
			>
				{statusValue ? rowData.label : 'Select Status'}
			</Dropdown.Button>
		</Tooltip>
	);
};
export default LoanPendingStatus;
