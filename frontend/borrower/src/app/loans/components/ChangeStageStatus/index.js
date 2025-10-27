'use client';
import React from 'react';
import styles from './index.module.css';
import { transformArrayStructure } from '@/lib/common';
import { useOptionsStore } from '@/store/options';
import { DropdownStatus } from '@repo/shared-methods';

/**
 * 修改pending loans列表内的状态
 * */
export const LoanPendingStatus = ({ data, successCallback, disabled, tooltipProps }) => {
	const { loanPendingStatus } = useOptionsStore();
	const items = transformArrayStructure({
		sourceArray: loanPendingStatus,
		entryValue: 'value',
		outputValue: 'key',
	});
	return (
		<DropdownStatus
			data={data}
			successCallback={successCallback}
			items={items}
			statusKey="pending_status"
			api={`/loan/${data.id}/pending_status`}
			className={styles?.[data?.pending_status]}
			disabled={disabled}
			tooltipProps={tooltipProps}
		/>
	);
};

/**
 * 修改loans列表内的状态
 * */
export const LoanStage = ({ data, successCallback }) => {
	const { loanStatusOptions } = useOptionsStore();
	const items = transformArrayStructure({
		sourceArray: loanStatusOptions,
		entryValue: 'value',
		outputValue: 'key',
	});
	return (
		<DropdownStatus
			data={data}
			successCallback={successCallback}
			items={items}
			statusKey="stage"
			api={`/loan/${data.id}/stage`}
			className={styles?.[data?.stage]}
		/>
	);
};
