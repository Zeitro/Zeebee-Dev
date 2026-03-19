'use client';
import React, { useCallback, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { createRequest } from '@repo/shared-methods/utils';
import { Spin } from 'antd';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { useUserStore } from '@repo/shared-methods/store';
// 授权gmail，以用户邮箱发送邮件<中转页,接收参数>
const GmailAuth = () => {
	const searchParams = useSearchParams();
	const { request } = createRequest();
	const { getOrganization } = useUserStore();
	const [loading, setLoading] = useState(false);
	const [tips, setTips] = useState('');
	// 授权使用用户邮箱发送邮件
	const handleGmailAuth = useCallback(async () => {
		try {
			setLoading(true);
			setTips('');
			const code = searchParams.get('code');
			const scope = searchParams.get('scope');
			const state = searchParams.get('state');
			await request('/organization_email/authorize', {
				method: 'POST',
				body: { code, email_provider: 'Gmail', scope, state },
			});
			const org = await getOrganization();
			// 同源下，新窗口通知旧窗口更新
			const channel = new BroadcastChannel('org-channel');
			channel.postMessage({ type: 'orgInfoUpdated', data: org });
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
			handleGmailAuth();
		},
		deps: [handleGmailAuth],
		isUseMicrotask: true,
	});

	return (
		<div className="w-full h-full flex justify-center items-center text-base font-medium py-6 px-6">
			{loading ? <Spin /> : tips}
		</div>
	);
};
export default GmailAuth;
