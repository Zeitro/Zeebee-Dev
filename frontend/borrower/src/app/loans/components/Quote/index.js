'use client';
import { FieldTips } from '@repo/shared-methods';
import { Button } from 'antd';
import { useEffect, useState } from 'react';

const Quote = () => {
	return <div>Quote</div>;
};
export default Quote;

export const QuoteBtn = ({ disabled }) => {
	const handleClick = () => {};
	return (
		<FieldTips title={disabled ? 'Only for pending stage & quoting status loan' : ''}>
			<Button type="primary" onClick={handleClick} ghost disabled={disabled}>
				Save&quote
			</Button>
		</FieldTips>
	);
};
