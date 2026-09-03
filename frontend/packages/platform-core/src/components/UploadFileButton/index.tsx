'use client';

import React, { forwardRef, useEffect, useImperativeHandle, useState } from 'react';
import type { ReactNode } from 'react';
import type { UploadFile, UploadProps } from 'antd';
import { Button, message, Upload } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import type { RcFile } from 'antd/es/upload/interface';

/** 未传 limitSize 时的默认单文件大小上限（MB） */
const DEFAULT_LIMIT_SIZE_MB = 20;
const DEFAULT_PRESIGN_URL = '/storage/presign_url';

type UploadDetail = {
	url: string;
	[key: string]: unknown;
};

type UploadAwsS3Params = {
	url?: string;
	file?: RcFile | Blob | File;
	// eslint-disable-next-line no-unused-vars
	onProgress?: (percent: number) => void;
	type?: string | null;
	isReturnMoreInfo?: boolean;
};

// eslint-disable-next-line no-unused-vars
type UploadAwsS3Fn = (params: UploadAwsS3Params) => Promise<string | UploadDetail | null>;

type TriggerUploadOptions = {
	// eslint-disable-next-line no-unused-vars
	onProgress?: (percent: number) => void;
	returnDetail?: boolean;
};

export interface CommonUploadRef {
	triggerUpload: (
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
	onSuccess?: (fileUrl: string, file: UploadFile | RcFile) => void;
	isReturnMoreInfo?: boolean;
	// eslint-disable-next-line no-unused-vars
	onProgress?: (percent: number, file: UploadFile | RcFile) => void;
	listType?: UploadProps['listType'];
	maxCount?: number;
	accept?: string;
	/** 单文件大小上限（MB），默认 {@link DEFAULT_LIMIT_SIZE_MB} */
	limitSize?: number;
	/** 应用侧上传实现，签名对齐各项目自己的 uploadAwsS3 */
	uploadAwsS3: UploadAwsS3Fn;
	/** 预签名接口路径，传给 uploadAwsS3 的 url */
	presignUrl?: string;
	children?: ReactNode;
	// eslint-disable-next-line no-unused-vars
	onGetUploadList?: (fileList: UploadFile[]) => void;
}

/** 透传 antd Upload 的其余 props（fileList / beforeUpload / customRequest 由组件内部管理） */
type UploadFileButtonProps = CommonUploadProps &
	Omit<UploadProps, keyof CommonUploadProps | 'fileList' | 'beforeUpload' | 'customRequest'>;

function getRawFile(file: UploadFile | RcFile): RcFile | Blob | File {
	if ('originFileObj' in file && file.originFileObj) {
		return file.originFileObj;
	}
	return file as RcFile;
}

export const UploadFileButton = forwardRef<CommonUploadRef, UploadFileButtonProps>(
	(
		{
			immediateUpload = true,
			onSuccess,
			isReturnMoreInfo = false,
			onProgress,
			listType = 'text' as UploadProps['listType'],
			maxCount = 1,
			accept,
			limitSize = DEFAULT_LIMIT_SIZE_MB,
			uploadAwsS3,
			presignUrl = DEFAULT_PRESIGN_URL,
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
			file: UploadFile | RcFile,
			// eslint-disable-next-line no-unused-vars
			progressCallback: (percent: number) => void,
			isReturnDetail = false,
		): Promise<string | UploadDetail> => {
			const result = await uploadAwsS3({
				url: presignUrl,
				file: getRawFile(file),
				onProgress: progressCallback,
				isReturnMoreInfo: isReturnDetail,
			});

			if (!result) {
				throw new Error('Upload failed');
			}

			if (isReturnDetail) {
				if (typeof result === 'object' && result?.url) {
					return result;
				}
				throw new Error('Upload failed');
			}
			if (typeof result === 'string') {
				return result;
			}
			throw new Error('Upload failed');
		};

		const triggerUpload = async (
			fileListOrOptions?: UploadFile[] | TriggerUploadOptions,
			options?: TriggerUploadOptions,
		) => {
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
					currentList.map((f) =>
						f.uid === file.uid ? { ...f, status: 'uploading' } : f,
					),
				);

				return handleUpload(
					file,
					(percent) => {
						setFileList((currentList) =>
							currentList.map((f) => (f.uid === file.uid ? { ...f, percent } : f)),
						);
						mergedOptions?.onProgress?.(percent);
					},
					isReturnDetail,
				)
					.then((result) => {
						const fileUrl = typeof result === 'string' ? result : result.url;
						setFileList((currentList) =>
							currentList.map((f) =>
								f.uid === file.uid
									? { ...f, status: 'done', url: fileUrl, percent: 100 }
									: f,
							),
						);
						return result;
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
			uploadProps.beforeUpload = (file: RcFile) => {
				return checkFileSize(file);
			};
			uploadProps.customRequest = (options) => {
				const {
					file,
					onSuccess: onAntdSuccess,
					onError,
					onProgress: onAntdProgress,
				} = options;
				const uploadItem = file as RcFile;
				handleUpload(
					uploadItem,
					(percent) => {
						onAntdProgress?.({ percent });
						onProgress?.(percent, uploadItem);
					},
					isReturnMoreInfo,
				)
					.then((result) => {
						const fileUrl = typeof result === 'string' ? result : result.url;
						onAntdSuccess?.({
							url: fileUrl,
							detail: typeof result === 'string' ? null : result,
						});
						onSuccess?.(fileUrl, uploadItem);
						message.destroy();
						message.open({
							type: 'success',
							content: `${uploadItem.name} Upload successful`,
						});
					})
					.catch((err) => {
						onError?.(err);
						message.destroy();
						message.open({
							type: 'error',
							content: `${uploadItem.name} Upload failed`,
						});
					});
			};
		} else {
			uploadProps.beforeUpload = async (file: RcFile) => {
				if (!checkFileSize(file)) return Upload.LIST_IGNORE;
				const thumbUrl = await getBase64(file);
				const nextFile: UploadFile = {
					uid: file.uid,
					name: file.name,
					originFileObj: file,
					thumbUrl,
				};
				setFileList((currentList) => {
					const newList = [...currentList, nextFile];
					return maxCount ? newList.slice(-maxCount) : newList;
				});
				return false;
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
