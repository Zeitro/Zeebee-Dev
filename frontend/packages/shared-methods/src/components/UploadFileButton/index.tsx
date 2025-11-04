'use client';

import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Upload, Button, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import type { UploadProps, UploadFile, RcFile } from 'antd/es/upload/interface';
import { uploadAwsS3 } from '@/api/common';
const MAX_SIZE_MB = 10;
export interface CommonUploadRef {
	// eslint-disable-next-line no-unused-vars
	triggerUpload: (options?: { onProgress?: (percent: number) => void }) => Promise<string[]>;
	clearFileList: () => void;
}

interface CommonUploadProps {
	immediateUpload?: boolean;
	// eslint-disable-next-line no-unused-vars
	onSuccess?: (fileUrl: string, file: UploadFile) => void;
	// eslint-disable-next-line no-unused-vars
	onProgress?: (percent: number, file: UploadFile) => void;
	listType?: UploadProps['listType'];
	maxCount?: number;
	accept?: string;
	children?: React.ReactNode;
	onGetUploadList?: () => UploadFile[];
}

export const UploadFileButton = forwardRef<CommonUploadRef, CommonUploadProps>(
	(
		{
			immediateUpload = true,
			onSuccess,
			onProgress,
			listType = 'text',
			maxCount = 1,
			accept,
			children,
			// eslint-disable-next-line no-unused-vars
			onGetUploadList = (fileList: UploadFile[]) => [],
			...moreProps
		},
		ref,
	) => {
		const [fileList, setFileList] = useState<UploadFile[]>([]);

		const checkFileSize = (file: RcFile) => {
			const isLt10M = file.size / 1024 / 1024 < MAX_SIZE_MB;
			if (!isLt10M) {
				message.destroy();
				message.open({
					type: 'error',
					content: `The file size cannot exceed ${MAX_SIZE_MB}MB!`,
				});
			}
			return isLt10M;
		};

		const handleUpload = async (
			file: UploadFile,
			progressCallback: (percent: number) => void,
		): Promise<string> => {
			const formData = new FormData();
			formData.append('file', file.originFileObj as RcFile);

			const url = await uploadAwsS3({
				url: '/storage/presign_url',
				file,
				onProgress: progressCallback,
			});

			if (url) {
				return url;
			} else {
				throw new Error('Upload failed');
			}
		};

		// 上传文件数组（可父级直接调用）
		const triggerUpload = async (options: any) => {
			// ... (这部分逻辑几乎不变，只是内部调用 handleUpload)
			if (immediateUpload) {
				console.warn('triggerUpload should not be called in immediateUpload mode.');
				return [];
			}

			const filesToUpload = fileList.filter((f) => f.status !== 'done');

			const uploadPromises = filesToUpload.map((file) => {
				setFileList((currentList) =>
					currentList.map((f) => {
						if (f.uid === file.uid) {
							f.status = 'uploading';
						}
						return f;
					}),
				);
				return handleUpload(file, (percent) => {
					setFileList((currentList) =>
						currentList.map((f) => {
							if (f.uid === file.uid) {
								f.percent = percent;
							}
							return f;
						}),
					);
					options?.onProgress?.(percent);
				})
					.then((url) => {
						setFileList((currentList) =>
							currentList.map((f) => {
								if (f.uid === file.uid) {
									f.status = 'done';
									f.url = url;
									f.percent = 100;
								}
								return f;
							}),
						);
						return url;
					})
					.catch(() => {
						setFileList((currentList) =>
							currentList.map((f) => {
								if (f.uid === file.uid) {
									f.status = 'error';
								}
								return f;
							}),
						);
						return null;
					});
			});

			const results = await Promise.all(uploadPromises);
			const successfulUrls = results.filter((url): url is string => url !== null);

			if (successfulUrls.length !== filesToUpload.length) {
				throw new Error('Some files failed to upload.');
			}
			return successfulUrls;
		};

		useImperativeHandle(ref, () => ({
			triggerUpload,
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
			// 模式一：customRequest 直接上传
			uploadProps.beforeUpload = (file: RcFile) => {
				return checkFileSize(file);
			};
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
						message.success(`${file.name} Upload successful`);
					})
					.catch((err) => {
						onError?.(err);
						message.error(`${file.name} Upload failed`);
					});
			};
		} else {
			// 模式二：阻止直接上传，拿到所有文件统一自定义上传
			uploadProps.beforeUpload = (file: RcFile) => {
				if (!checkFileSize(file)) return Upload.LIST_IGNORE;
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
UploadFileButton.displayName = 'UploadFileButton';
