'use client';
import React, { useCallback, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { createRequest } from '@repo/shared-methods/utils';
import { Spin } from 'antd';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
// 授权 Google Drive，以组织身份同步交易文件<中转页,接收参数>
const DriveAuth = () => {
	const searchParams = useSearchParams();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [tips, setTips] = useState('');
	// 交换授权码并保存 Drive 连接
	const handleDriveAuth = useCallback(async () => {
		try {
			setLoading(true);
			setTips('');
			const code = searchParams.get('code');
			const state = searchParams.get('state');
			const integration = await request('/drive/authorize', {
				method: 'POST',
				body: { code, state },
			});
			// 同源下，新窗口通知旧窗口更新 Drive 连接状态
			const channel = new BroadcastChannel('drive-channel');
			channel.postMessage({ type: 'driveAuthorized', data: integration });
			setLoading(false);
			setTips('Authorization successful.');
			setTimeout(() => {
				window.close();
			}, 500);
		} catch (error) {
			setTips('Authorization failed: ' + error.message);
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [searchParams]);

	useDeferredSyncEffect({
		deferred: () => {
			handleDriveAuth();
		},
		deps: [handleDriveAuth],
		isUseMicrotask: true,
	});

	return (
		<div className="w-full h-full flex justify-center items-center text-base font-medium py-6 px-6">
			{loading ? <Spin /> : tips}
		</div>
	);
};
export default DriveAuth;
