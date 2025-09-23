'use client';

import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Upload, Button, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import type { UploadProps, UploadFile, RcFile } from 'antd/es/upload/interface';
import { createRequest } from '@/lib/fetch';

export interface CommonUploadRef {
	triggerUpload: (options?: { onProgress?: (percent: number) => void }) => Promise<string[]>;
	clearFileList: () => void;
}

interface CommonUploadProps {
	immediateUpload?: boolean;
	onSuccess?: (fileUrl: string, file: UploadFile) => void;
	onProgress?: (percent: number, file: UploadFile) => void;
	listType?: UploadProps['listType'];
	maxCount?: number;
	accept?: string;
	children?: React.ReactNode;
	onGetUploadList?: () => UploadFile[];
}

const CommonUpload = forwardRef<CommonUploadRef, CommonUploadProps>(
	(
		{
			immediateUpload = true,
			onSuccess,
			onProgress,
			listType = 'text',
			maxCount = 1,
			accept,
			children,
			onGetUploadList = (fileList: UploadFile[]) => [],
			...moreProps
		},
		ref,
	) => {
		const { uploadWithProgress } = createRequest();
		const [fileList, setFileList] = useState<UploadFile[]>([]);

		const handleUpload = async (
			file: UploadFile,
			progressCallback: (percent: number) => void,
		): Promise<string> => {
			const formData = new FormData();
			formData.append('file', file.originFileObj as RcFile);

			const response = await uploadWithProgress('/api/upload', {
				body: formData,
				onProgress: progressCallback,
			});

			if (response.success && response.url) {
				return response.url;
			} else {
				throw new Error(response.message || 'Upload failed');
			}
		};

		useImperativeHandle(ref, () => ({
			triggerUpload: async (options) => {
				// ... (这部分逻辑几乎不变，只是内部调用 handleUpload)
				if (immediateUpload) {
					console.warn('triggerUpload should not be called in immediateUpload mode.');
					return [];
				}

				const filesToUpload = fileList.filter((f) => f.status !== 'done');
				const uploadPromises = filesToUpload.map((file) => {
					setFileList((currentList) =>
						currentList.map((f) =>
							f.uid === file.uid ? { ...f, status: 'uploading' } : f,
						),
					);

					return handleUpload(file, (percent) => {
						setFileList((currentList) =>
							currentList.map((f) => (f.uid === file.uid ? { ...f, percent } : f)),
						);
						options?.onProgress?.(percent);
					})
						.then((url) => {
							setFileList((currentList) =>
								currentList.map((f) =>
									f.uid === file.uid
										? { ...f, status: 'done', url, percent: 100 }
										: f,
								),
							);
							return url;
						})
						.catch(() => {
							setFileList((currentList) =>
								currentList.map((f) =>
									f.uid === file.uid ? { ...f, status: 'error' } : f,
								),
							);
							return null;
						});
				});

				const results = await Promise.all(uploadPromises);
				const successfulUrls = results.filter((url): url is string => url !== null);

				if (successfulUrls.length !== filesToUpload.length) {
					throw new Error('部分文件上传失败');
				}
				return successfulUrls;
			},
			clearFileList: () => {
				setFileList([]);
			},
		}));

		const uploadProps: UploadProps = {
			accept,
			listType,
			fileList,
			maxCount,
			onRemove: (file) => {
				setFileList((currentList) => currentList.filter((f) => f.uid !== file.uid));
			},
			...moreProps,
		};

		if (immediateUpload) {
			// 模式一：customRequest 现在也只调用 handleUpload
			uploadProps.customRequest = ({
				file,
				onSuccess: onAntdSuccess,
				onError,
				onProgress: onAntdProgress,
			}: any) => {
				handleUpload(file, (percent) => {
					onAntdProgress({ percent });
					onProgress?.(percent, file);
				})
					.then((url) => {
						onAntdSuccess?.({ url });
						onSuccess?.(url, file);
						message.success(`${file.name} 上传成功`);
					})
					.catch((err) => {
						onError?.(err);
						message.error(`${file.name} 上传失败`);
					});
			};
		} else {
			// 模式二：逻辑不变
			uploadProps.beforeUpload = (file) => {
				setFileList((currentList) => {
					const newList = [...currentList, file];
					return maxCount ? newList.slice(-maxCount) : newList;
				});
				return false;
			};
		}

		useEffect(() => {
			onGetUploadList(fileList);
		}, [fileList]);

		return (
			<Upload {...uploadProps}>
				{children || <Button icon={<UploadOutlined />}>选择文件</Button>}
			</Upload>
		);
	},
);

CommonUpload.displayName = 'CommonUpload';
export default CommonUpload;
