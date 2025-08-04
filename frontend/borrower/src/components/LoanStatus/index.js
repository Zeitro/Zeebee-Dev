'use client';
import { Dropdown, Tooltip } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';
import { DownOutlined } from '@ant-design/icons';
import styles from './index.module.css';
import { createRequest } from '@/lib/fetch';
import { transformArrayStructure } from '@/lib/common';
import {useLoanStore} from '@/store/loan';
// import { loanPendingStatus } from '@/enum/loan';
let api = '/api/loan/status';

/**
 * loans列表内的状态
 * status 当前状态
 * successCallback 状态变更后的回调函数
 * is3rdParty true是第三方贷款,则可修改状态,否则不可修改
 * */
const LoanStatus = ({ status, is3rdParty = true, successCallback }) => {
	const {loanPendingStatus} = useLoanStore();
	const [loading, setLoading] = useState(false);
	const [statusValue, setStatusValue] = useState(status);
	const items = transformArrayStructure({
		sourceArray: loanPendingStatus,
		entryValue: 'value',
		outputValue: 'key',
	});

	const { request, cancel } = createRequest();
	const rowData = useMemo(() => {
		return items.find((i) => i.key === statusValue) || {};
	}, [statusValue]);

	const handleSelect = ({ key }) => {
		try {
			cancel(api);
			setLoading(true);
			setStatusValue(key);
			// const res = request({
			// 	url: api,
			// 	method: 'POST',
			// 	data: { status: key },
			// });
			setLoading(false);
			successCallback?.();
			// console.log(res);
		} catch (error) {
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
				className={[styles.loanStatus, styles[statusValue]].join(' ')}
			>
				{statusValue ? rowData.label : 'Select Status'}
			</Dropdown.Button>
		</Tooltip>
	);
};
export default LoanStatus;
