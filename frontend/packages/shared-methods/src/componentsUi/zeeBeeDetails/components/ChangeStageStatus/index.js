'use client';
import React from 'react';
import styles from './index.module.css';
import { transformArrayStructure } from '@repo/shared-methods/utils';
import { useOptionsStore } from '@repo/shared-methods/store';
import { DropdownStatus } from '@repo/shared-methods';

/**
 * 修改pending loans列表内的状态
 * */
export const LoanPendingStatus = ({ data, successCallback, ...props }) => {
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
			{...props}
			className={[styles?.[data?.pending_status], props?.className].join(' ')}
		/>
	);
};

/**
 * 修改loans列表内的状态
 * */
export const LoanStage = ({ data, successCallback, ...props }) => {
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
			{...props}
			className={[styles.stage, styles?.[data?.stage], props?.className].join(' ')}
		/>
	);
};
