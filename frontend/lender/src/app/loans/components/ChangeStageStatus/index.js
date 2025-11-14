'use client';
import React from 'react';
import styles from './index.module.css';
import { transformArrayStructure } from '@repo/shared-methods/utils';
import { useOptionsStore } from '@/store/options';
import { DropdownStatus } from '@repo/shared-methods';

/**
 * 修改pending loans列表内的状态
 * */
export const LoanPendingStatus = ({ data, successCallback }) => {
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
