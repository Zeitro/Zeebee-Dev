'use client';

import React from 'react';
import { displayValue } from '@repo/shared-methods/utils';

const DetailField = ({ label, value, className = '' }) => {
	return (
		<div className={`flex min-w-0 flex-col gap-0.5 ${className}`}>
			<p className="m-0 text-xs leading-5 text-base-label">{label}</p>
			<p className="m-0 break-words text-sm leading-[22px] text-foreground">
				{displayValue(value)}
			</p>
		</div>
	);
};

export default DetailField;
