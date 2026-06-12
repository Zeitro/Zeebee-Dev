'use client';

import React, { forwardRef, useEffect, useImperativeHandle, useState } from 'react';
import type { UploadFile, UploadProps } from 'antd';
import { Button, message, Upload } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import type { RcFile } from 'antd/es/upload/interface';
import { uploadAwsS3 } from '../../api/common';

/** 未传 limitSize 时的默认单文件大小上限（MB） */
const DEFAULT_LIMIT_SIZE_MB = 20;

type UploadDetail = {
	// 成功上传后裁剪掉 query string 的最终可访问 url
	url: string;
	// 其余来自新版 getUploadAwsUrl 的额外字段（字段结构由后端决定）
	[key: string]: any;
};

type TriggerUploadOptions = {
	// eslint-disable-next-line no-unused-vars
	onProgress?: (_percent: number) => void;
	// true：返回包含 presign 额外字段的 detail；false：只返回 url（兼容老页面）
	returnDetail?: boolean;
};

export interface CommonUploadRef {
	triggerUpload: (
		// 兼容老用法：triggerUpload(fileList)；以及 triggerUpload({ returnDetail: true })
		// eslint-disable-next-line no-unused-vars
		fileListOrOptions?: UploadFile[] | TriggerUploadOptions,
		// eslint-disable-next-line no-unused-vars
		options?: TriggerUploadOptions,
	) => Promise<string[] | UploadDetail[]>;
	clearFileList: () => void;
}

interface CommonUploadProps {
	immediateUpload?: boolean;
	// eslint-disable-next-line no-unused-vars
	onSuccess?: (fileUrl: string, file: UploadFile) => void;
	// 组件级开关：true 时上传结果包含 presign 额外信息，false 时仅返回 url（兼容老页面）
	isReturnMoreInfo?: boolean;
	// eslint-disable-next-line no-unused-vars
	onProgress?: (percent: number, file: UploadFile) => void;
	listType?: UploadProps['listType'];
	maxCount?: number;
	accept?: string;
	/** 单文件大小上限（MB），默认 {@link DEFAULT_LIMIT_SIZE_MB} */
	limitSize?: number;
	children?: React.ReactNode;
	// eslint-disable-next-line no-unused-vars
	onGetUploadList?: (fileList: UploadFile[]) => void;
}

/** 透传 antd Upload 的其余 props（fileList / beforeUpload / customRequest 由组件内部管理） */
type UploadFileButtonProps = CommonUploadProps &
	Omit<UploadProps, keyof CommonUploadProps | 'fileList' | 'beforeUpload' | 'customRequest'>;

export const UploadFileButton = forwardRef<CommonUploadRef, UploadFileButtonProps>(
	(
		{
			immediateUpload = true,
			onSuccess,
			isReturnMoreInfo = false,
			onProgress,
			listType = 'text',
			maxCount = 1,
			accept,
			limitSize = DEFAULT_LIMIT_SIZE_MB,
			children,
			// eslint-disable-next-line no-unused-vars
			onGetUploadList = (_fileList: UploadFile[]) => [],
			...moreProps
		},
		ref,
	) => {
		const [fileList, setFileList] = useState<UploadFile[]>([]);

		const checkFileSize = (file: RcFile) => {
			const isWithinLimit = file.size / 1024 / 1024 < limitSize;
			if (!isWithinLimit) {
				message.destroy();
				message.open({
					type: 'error',
					content: `The file size cannot exceed ${limitSize}MB!`,
				});
			}
			return isWithinLimit;
		};

		const getBase64 = (file: RcFile): Promise<string> =>
			new Promise((resolve, reject) => {
				const reader = new FileReader();
				reader.readAsDataURL(file);
				reader.onload = () => resolve(reader.result as string);
				reader.onerror = (error) => reject(error);
			});

		const handleUpload = async (
			file: UploadFile,
			// eslint-disable-next-line no-unused-vars
			progressCallback: (_percent: number) => void,
			isReturnDetail = false,
		): Promise<string | UploadDetail> => {
			const result = await uploadAwsS3({
				url: '/storage/presign_url',
				file: file.originFileObj,
				onProgress: progressCallback,
				isReturnMoreInfo: isReturnDetail,
			});

			if (!result) {
				throw new Error('Upload failed');
			}

			if (isReturnDetail) {
				if (typeof result === 'object' && result?.url) {
					return result as UploadDetail;
				}
				throw new Error('Upload failed');
			} else {
				if (typeof result === 'string') {
					return result;
				}
				throw new Error('Upload failed');
			}
		};

		// 上传文件数组（可父级直接调用）
		const triggerUpload = async (fileListOrOptions: any, options?: TriggerUploadOptions) => {
			if (immediateUpload) {
				console.warn('triggerUpload should not be called in immediateUpload mode.');
				return [];
			}

			const filesArg = Array.isArray(fileListOrOptions) ? fileListOrOptions : undefined;
			const mergedOptions: TriggerUploadOptions | undefined = Array.isArray(fileListOrOptions)
				? options
				: fileListOrOptions;
			const isReturnDetail = mergedOptions?.returnDetail ?? isReturnMoreInfo;

			const filesToUpload = (filesArg ?? fileList).filter((f) => f.status !== 'done');

			const uploadPromises = filesToUpload.map((file) => {
				setFileList((currentList) =>
					currentList.map((f) => {
						if (f.uid === file.uid) {
							f.status = 'uploading';
						}
						return f;
					}),
				);

				return handleUpload(
					file,
					(percent) => {
						setFileList((currentList) =>
							currentList.map((f) => {
								if (f.uid === file.uid) {
									f.percent = percent;
								}
								return f;
							}),
						);
						mergedOptions?.onProgress?.(percent);
					},
					isReturnDetail,
				)
					.then((result: any) => {
						setFileList((currentList) =>
							currentList.map((f) => {
								if (f.uid === file.uid) {
									f.status = 'done';
									// 老逻辑：url 是 string；新逻辑：detail.url 是最终访问 url
									f.url = isReturnDetail ? result.url : result;
									f.percent = 100;
								}
								return f;
							}),
						);
						return result;
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

			if (isReturnDetail) {
				const successfulDetails = results.filter((d): d is UploadDetail => d !== null);
				if (successfulDetails.length !== filesToUpload.length) {
					throw new Error('Some files failed to upload.');
				}
				return successfulDetails;
			}

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

		const uploadProps: Partial<UploadProps> = {
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
				handleUpload(
					file,
					(percent) => {
						onAntdProgress({ percent });
						onProgress?.(percent, file);
					},
					isReturnMoreInfo,
				)
					.then((result) => {
						const fileUrl = typeof result === 'string' ? result : result.url;
						onAntdSuccess?.({
							url: fileUrl,
							detail: typeof result === 'string' ? null : result,
						});
						onSuccess?.(fileUrl, file);
						message.destroy();
						message.open({
							type: 'success',
							content: `${file.name} Upload successful`,
						});
					})
					.catch((err) => {
						onError?.(err);
						message.destroy();
						message.open({
							type: 'error',
							content: `${file.name} Upload failed`,
						});
					});
			};
		} else {
			// 模式二：阻止直接上传，拿到所有文件统一自定义上传
			uploadProps.beforeUpload = async (file: RcFile) => {
				if (!checkFileSize(file)) return Upload.LIST_IGNORE;
				const thumbUrl = await getBase64(file);
				const uploadFile = {
					uid: file.uid,
					name: file.name,
					status: 'uploading' as const,
					originFileObj: file,
					thumbUrl,
				} as UploadFile;
				setFileList((currentList) => {
					const newList = [...currentList, uploadFile];
					return maxCount ? newList.slice(-maxCount) : newList;
				});
				return false; // 阻止自动上传
			};
		}

		useEffect(() => {
			onGetUploadList(fileList);
		}, [fileList, onGetUploadList]);

		return (
			<Upload {...uploadProps}>
				{children || <Button icon={<UploadOutlined />}>选择文件</Button>}
			</Upload>
		);
	},
);
UploadFileButton.displayName = 'UploadFileButton';
