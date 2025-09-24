'use client';
import { Button } from 'antd';
import React from 'react';
import ViewSvg from '@/assets/images/view.svg';
const ViewFile = ({ children, ...rest }) => {
	const handleView = () => {};
	return (
		<Button icon={<ViewSvg />} size="large" onClick={handleView} {...rest}>
			{children || 'View'}
		</Button>
	);
};

export default ViewFile;
