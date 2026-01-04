'use client';

import React, { useState } from 'react';
import { Button, message } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import type { ButtonProps } from 'antd';
import { createRequest } from '@/lib/fetch';

// 定义 Props
interface CommonDownloadButtonProps extends Omit<ButtonProps, 'onProgress'> {
	fileUrl: string;
	fileName?: string;
	onProgress?: (percent: number) => void; // ✅ 新增：进度回调
}

const CommonDownloadButton: React.FC<CommonDownloadButtonProps> = ({
	fileUrl,
	fileName,
	onProgress,
	children,
	...rest
}) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);

	const handleDownload = async () => {
		setLoading(true);
		onProgress?.(0);
		message.loading({ content: '正在准备下载...', key: 'download' });

		try {
			const response = await request(fileUrl, { isExternalLink: true });

			// 1. 获取文件总大小
			const contentLength = response.headers.get('Content-Length');
			if (!contentLength) {
				console.warn('Content-Length header not found, cannot track progress.');
				// 如果没有总大小，就退回旧的下载方式
				const blob = await response.blob();
				onProgress?.(100);
				// ... (旧的下载逻辑)
				return;
			}
			const totalSize = parseInt(contentLength, 10);
			let loadedSize = 0;

			// 2. 创建一个可读流来处理数据块
			const reader = response.body.getReader();
			const chunks: ArrayBuffer[] = [];

			while (true) {
				const { done, value } = await reader.read();
				if (done) {
					break;
				}
				chunks.push(value.buffer); // Use value.buffer instead of just value
				loadedSize += value.byteLength;
				const percent = Math.round((loadedSize / totalSize) * 100);
				onProgress?.(percent); // ✅ 调用进度回调
			}

			// 3. 组合数据块并触发下载
			const blob = new Blob(chunks);
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = fileName || fileUrl.substring(fileUrl.lastIndexOf('/') + 1);
			document.body.appendChild(a);
			a.click();
			a.remove();
			window.URL.revokeObjectURL(url);

			message.success({ content: '下载已开始！', key: 'download', duration: 2 });
		} catch (error) {
			console.error('Download error:', error);
			message.error({ content: '下载失败，请稍后重试。', key: 'download', duration: 2 });
		} finally {
			setLoading(false);
			setTimeout(() => onProgress?.(0), 1000); // 结束后1秒重置进度
		}
	};

	return (
		<Button
			icon={<DownloadOutlined />}
			size="large"
			loading={loading}
			onClick={handleDownload}
			{...rest}
		>
			{children || 'Download'}
		</Button>
	);
};

export default CommonDownloadButton;
