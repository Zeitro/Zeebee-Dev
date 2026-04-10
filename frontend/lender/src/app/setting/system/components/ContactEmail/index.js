'use client';

import { MailOutlined, PlusCircleOutlined } from '@ant-design/icons';
import { Button, Space, Spin } from 'antd';
import { useCallback, useMemo, useState } from 'react';
import AddContactEmail from './AddContactEmail';
import { BtnPopUpModal } from '@repo/shared-methods';
import { authorizedColor } from '@repo/shared-methods/enum';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { useUserStore } from '@repo/shared-methods/store';
import { createRequest } from '@repo/shared-methods/utils';
const ContactEmail = () => {
	const { request } = createRequest();
	const { getOrganization, organization, setOrganization } = useUserStore();
	const [loading, setLoading] = useState(false);
	const [isInit, setIsInit] = useState(true);
	const [authLoading, setAuthLoading] = useState(false);
	const contactInfo = useMemo(() => {
		return organization?.organization_emails?.find((item) => item.purpose === 'Contact') || {};
	}, [organization]);
	const colorInfo = authorizedColor[contactInfo?.is_authorized ? 'authorized' : 'unauthorized'];

	const getData = useCallback(() => {
		setLoading(true);
		getOrganization()
			.then(() => {
				setLoading(false);
				setIsInit(false);
			})
			.catch(() => {
				setLoading(false);
				setIsInit(false);
			});
	}, [getOrganization]);

	const handleAuthorize = async () => {
		try {
			setAuthLoading(true);
			const params = new URLSearchParams({
				email_provider: 'Gmail',
				purpose: 'Contact',
			});
			const res = await request(`/organization_email/generate_auth_code_url?${params}`, {
				method: 'GET',
			});
			const popupWidth = 500;
			const popupHeight = 600;
			const left = (screen.availWidth - popupWidth) / 2;
			const top = (screen.availHeight - popupHeight) / 2;

			window.gmailAuthPopup = window.open(
				res.url,
				'_blank',
				`width=${popupWidth},height=${popupHeight},left=${left},top=${top},resizable=yes,scrollbars=yes`,
			);
			setAuthLoading(false);
		} catch {
			setAuthLoading(false);
		}
	};

	useDeferredSyncEffect({
		deferred: () => {
			// 同源下，等待新窗口通知更新数据
			const channel = new BroadcastChannel('org-channel');
			channel.onmessage = (event) => {
				// 更新组织信息
				if (event.data.type === 'orgInfoUpdated') {
					setOrganization(event.data.data);
				}
			};
			getData();
		},
		deps: [getData],
		delay: 0,
	});

	return (
		<Spin spinning={loading}>
			{isInit ? (
				<div className="h-[100px]"></div>
			) : contactInfo?.email ? (
				<div className="flex items-center justify-between py-6 px-6 border border-base-border rounded-[10px]">
					<div>
						<div className="text-xl font-medium text-black-dark pb-3">
							{contactInfo.email}
						</div>
						<div
							style={{
								color: colorInfo?.color,
								background: colorInfo?.bg,
								border: `1px solid ${colorInfo?.border}`,
								borderRadius: 2,
								padding: '0 8px',
								width: 'max-content',
							}}
						>
							{contactInfo?.is_authorized ? 'Authorized' : 'Unauthorized'}
						</div>
					</div>
					<Space>
						{!contactInfo?.is_authorized && (
							<Button type="primary" onClick={handleAuthorize} loading={authLoading}>
								Authorize
							</Button>
						)}
						<BtnPopUpModal
							type="default"
							size="medium"
							getContent={(setOpen) => (
								<AddContactEmail
									setOpen={setOpen}
									rowData={contactInfo}
									getData={getData}
								/>
							)}
						>
							Edit
						</BtnPopUpModal>
					</Space>
				</div>
			) : (
				<div className="border border-base-border bg-[rgba(0,0,0,0.04)] rounded-[10px] py-6 flex flex-col items-center justify-center">
					<div className="h-10 w-10 border border-base-border rounded-[50%] flex items-center justify-center">
						<MailOutlined style={{ fontSize: 24, color: 'rgba(0, 0, 0, 0.45)' }} />
					</div>
					<div className="font-medium text-black-dark text-sm pt-4 pb-1">
						No contact email added
					</div>
					<div className="text-[rgba(0,0,0,0.65)] text-sm pb-4">
						Add a Gmail address to get started.
					</div>
					<BtnPopUpModal
						getContent={(setOpen) => (
							<AddContactEmail setOpen={setOpen} getData={getData} />
						)}
						icon={<PlusCircleOutlined />}
						type="primary"
						size="medium"
					>
						Add email address
					</BtnPopUpModal>
				</div>
			)}
		</Spin>
	);
};
export default ContactEmail;
