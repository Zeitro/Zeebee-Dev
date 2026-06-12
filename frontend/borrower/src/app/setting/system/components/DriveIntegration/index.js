'use client';

import { GoogleOutlined } from '@ant-design/icons';
import { Button, Input, Select, Space, Spin, message } from 'antd';
import { useCallback, useMemo, useRef, useState } from 'react';
import { authorizedColor } from '@repo/shared-methods/enum';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { createRequest } from '@repo/shared-methods/utils';

const DriveIntegration = () => {
	// createRequest() 每次渲染都会返回新的 request 函数；若直接在渲染期调用，
	// 会让依赖它的 getStatus/loadDrives 每次渲染都变更身份，进而让下方的
	// useDeferredSyncEffect 反复触发、无限请求 /drive/auth_status。用 useMemo
	// 固定单一实例，既消除请求风暴，也让每个组件实例拥有稳定的取消控制器。
	const { request } = useMemo(() => createRequest(), []);
	const [loading, setLoading] = useState(false);
	const [isInit, setIsInit] = useState(true);
	const [authLoading, setAuthLoading] = useState(false);
	const [saving, setSaving] = useState(false);
	const [integration, setIntegration] = useState(null);
	const [drives, setDrives] = useState([]);
	const [folderName, setFolderName] = useState('');

	const colorInfo = authorizedColor[integration?.is_authorized ? 'authorized' : 'unauthorized'];

	const getStatus = useCallback(async () => {
		try {
			setLoading(true);
			const status = await request('/drive/auth_status', { method: 'GET' });
			setIntegration(status);
			setFolderName(status?.folder_name || '');
		} catch {
			// ignore
		} finally {
			setLoading(false);
			setIsInit(false);
		}
	}, [request]);

	const loadDrives = useCallback(async () => {
		try {
			const res = await request('/drive/drives', { method: 'GET' });
			setDrives(res?.lists || []);
		} catch {
			// ignore
		}
	}, [request]);

	const handleAuthorize = async () => {
		try {
			setAuthLoading(true);
			const res = await request('/drive/auth_url', { method: 'GET' });
			const popupWidth = 500;
			const popupHeight = 600;
			const left = (screen.availWidth - popupWidth) / 2;
			const top = (screen.availHeight - popupHeight) / 2;
			window.driveAuthPopup = window.open(
				res.url,
				'_blank',
				`width=${popupWidth},height=${popupHeight},left=${left},top=${top},resizable=yes,scrollbars=yes`,
			);
		} catch {
			message.error('Failed to start Google Drive authorization.');
		} finally {
			setAuthLoading(false);
		}
	};

	const handleSelectDrive = async (value, option) => {
		try {
			setSaving(true);
			await request('/drive/selection', {
				method: 'PUT',
				body: { drive_id: value || '', drive_name: option?.label || '' },
			});
			setIntegration((prev) => ({
				...prev,
				drive_id: value || '',
				drive_name: option?.label || '',
			}));
			message.success('Target drive updated.');
		} catch {
			message.error('Failed to update target drive.');
		} finally {
			setSaving(false);
		}
	};

	const handleSaveFolder = async () => {
		const name = folderName.trim();
		if (!name) {
			message.warning('Folder name cannot be empty.');
			return;
		}
		try {
			setSaving(true);
			await request('/drive/folder', { method: 'PUT', body: { folder_name: name } });
			setIntegration((prev) => ({ ...prev, folder_name: name }));
			message.success('Root folder name updated.');
		} catch {
			message.error('Failed to update root folder name.');
		} finally {
			setSaving(false);
		}
	};

	const handleDisconnect = async () => {
		try {
			setSaving(true);
			await request('/drive/disconnect', { method: 'POST' });
			await getStatus();
			message.success('Google Drive disconnected.');
		} catch {
			message.error('Failed to disconnect Google Drive.');
		} finally {
			setSaving(false);
		}
	};

	const channelRef = useRef(null);
	useDeferredSyncEffect({
		deferred: () => {
			// 同源下，等待授权弹窗回调后刷新连接状态
			const channel = new BroadcastChannel('drive-channel');
			channelRef.current = channel;
			channel.onmessage = (event) => {
				if (event.data?.type === 'driveAuthorized') {
					getStatus();
					loadDrives();
				}
			};
			getStatus();
		},
		destroyCallback: () => {
			channelRef.current?.close();
			channelRef.current = null;
		},
		deps: [getStatus, loadDrives],
		delay: 0,
	});

	return (
		<Spin spinning={loading}>
			{isInit ? (
				<div className="h-[100px]"></div>
			) : integration?.is_authorized ? (
				<div className="border border-base-border rounded-[10px] py-6 px-6">
					<div className="flex items-center justify-between">
						<div>
							<div className="text-xl font-medium text-black-dark pb-3">
								{integration.auth_email}
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
								Authorized
							</div>
						</div>
						<Button danger onClick={handleDisconnect} loading={saving}>
							Disconnect
						</Button>
					</div>

					<div className="pt-6">
						<div className="text-sm font-medium text-black-dark pb-2">Target drive</div>
						<Select
							className="w-full"
							value={integration.drive_id || ''}
							onChange={handleSelectDrive}
							onDropdownVisibleChange={(open) => open && loadDrives()}
							options={drives.map((d) => ({ value: d.id, label: d.name }))}
							placeholder="Select a drive"
						/>
					</div>

					<div className="pt-6">
						<div className="text-sm font-medium text-black-dark pb-2">
							Root folder name
						</div>
						<Space.Compact className="w-full">
							<Input
								value={folderName}
								onChange={(e) => setFolderName(e.target.value)}
								placeholder="Zeebee"
							/>
							<Button type="primary" onClick={handleSaveFolder} loading={saving}>
								Save
							</Button>
						</Space.Compact>
					</div>
				</div>
			) : (
				<div className="border border-base-border bg-[rgba(0,0,0,0.04)] rounded-[10px] py-6 flex flex-col items-center justify-center">
					<div className="h-10 w-10 border border-base-border rounded-[50%] flex items-center justify-center">
						<GoogleOutlined style={{ fontSize: 24, color: 'rgba(0, 0, 0, 0.45)' }} />
					</div>
					<div className="font-medium text-black-dark text-sm pt-4 pb-1">
						Google Drive not connected
					</div>
					<div className="text-[rgba(0,0,0,0.65)] text-sm pb-4">
						Connect Google Drive to sync transaction documents.
					</div>
					<Button type="primary" onClick={handleAuthorize} loading={authLoading}>
						Connect Google Drive
					</Button>
				</div>
			)}
		</Spin>
	);
};

export default DriveIntegration;
