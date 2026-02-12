'use client';
import { useEffect, useRef } from 'react';
import Script from 'next/script';
import { createRequest } from '@repo/shared-methods/utils';
const str = '/api/v1';
const baseURL = process.env.NEXT_PUBLIC_API_URL + str;
export const EmailAuthPage = () => {
	const { request } = createRequest();
	const client = useRef(null);

	// 初始化 GIS 客户端
	function initClient() {
		if (window.google && window.google.accounts && window.google.accounts.oauth2) {
			client.current = window.google.accounts.oauth2.initCodeClient({
				client_id: process.env.NEXT_PUBLIC_GOOGLE_OAUTH_CLIENT_ID,
				scope: 'https://www.googleapis.com/auth/gmail.send',
				redirect_uri: process.env.NEXT_PUBLIC_ENV !== 'local' ? `${baseURL}` : `${str}`,
				ux_mode: 'popup',
				callback: (response) => {
					// response.code 就是授权code
					// 发送到后端换取access_token
					request('/api/gmail-auth', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ code: response.code }),
					});
				},
			});
		}
	}

	useEffect(() => {
		// GIS 脚本加载后初始化
		if (window.google && window.google.accounts && window.google.accounts.oauth2) {
			initClient();
		}
	}, []);

	function handleAuth() {
		if (client.current) {
			client.current.requestCode();
		} else {
			alert('Google OAuth 客户端未初始化');
		}
	}

	return (
		<div>
			{/* 推荐用 Script 组件加载 GIS */}
			<Script
				src="https://accounts.google.com/gsi/client"
				strategy="afterInteractive"
				onLoad={initClient}
			/>
			<button onClick={handleAuth}>授权Gmail账号</button>
		</div>
	);
};
