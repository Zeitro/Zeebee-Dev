'use client';
import { Loading3QuartersOutlined } from '@ant-design/icons';
import { eventBus, createRequest } from '@repo/shared-methods/utils';
import { useUserStore } from '@repo/shared-methods/store';
import { message } from 'antd';
import { useCallback, useEffect } from 'react';

// 切换组织中的页面
export const SwitchOrganization = ({
	organization_name,
	organization_id,
	api,
	successCallback,
}) => {
	const { request } = createRequest();
	const { getUserInfo, updateAccess } = useUserStore();

	// 切换组织
	const handleSwitch = useCallback(async () => {
		if (!(organization_name && organization_id && api)) return;
		const access = await request(api, {
			method: 'POST',
			body: {
				organization_id,
			},
		});
		successCallback?.(access, () => {
			message.open({
				type: 'success',
				content: `Successfully switched organization to ${organization_name}`,
			});
		});
		setTimeout(() => {
			eventBus.emit('removePortalComponent');
		}, 300);
		// setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		// await updateAccess(access);
		// await getUserInfo();
		// router.replace('/');
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [api, getUserInfo, organization_id, organization_name, updateAccess]);

	useEffect(() => {
		handleSwitch();
	}, [handleSwitch]);
	return (
		<div
			className="flex items-center justify-center"
			style={{
				background: '#7a7a7a',
				zIndex: 9999,
				position: 'absolute',
				top: 0,
				left: 0,
				right: 0,
				bottom: 0,
			}}
		>
			<div
				style={{ minWidth: 280, maxWidth: 500 }}
				className="flex items-center justify-center flex-col bg-white rounded-lg px-6 py-6 shadow-[0_6px_16px_0_rgba(0,0,0,0.08), 0_3px_6px_-4px_rgba(0,0,0,0.12), 0_9px_28px_8px_rgba(0,0,0,0.05)]"
			>
				<Loading3QuartersOutlined spin style={{ color: '#8C8C8C', fontSize: '24px' }} />
				<div className="text-base font-medium text-black-dark pt-6">
					Switch to {organization_name}
				</div>
			</div>
		</div>
	);
};
