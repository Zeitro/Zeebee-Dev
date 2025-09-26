import { Spin } from 'antd';
import React from 'react';

const PageLoading = () => {
	return (
		<div className="h-full w-full flex items-center justify-center">
			<Spin />
		</div>
	);
};
export default PageLoading;
