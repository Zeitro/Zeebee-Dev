'use client';
import React, {useState} from 'react';
import styles from './index.module.scss';
import {InboxOutlined} from '@ant-design/icons';
import ArticleSvg from '@/assets/images/article.svg';
import {DownloadFileButton, UploadFileButton, ViewFile} from '@repo/shared-methods';
import {Space} from 'antd';

const Documents = () => {
	const [fileList, setFileList] = useState([]);
	return (
		<div>
			<div className="py-6 px-6 border border-base-border rounded-2xl mb-6">
				<div className="text-xl font-medium pb-1 pl-3 text-black-dark">Upload Document</div>
				<p
					className="text-sm/[22px] pl-3 pb-4 border-b border-base-border mb-4"
					style={{ color: 'rgba(0, 0, 0, 0.65)' }}
				>
					Drag and drop your document here, or click to select a file.
				</p>
				<UploadFileButton
					immediateUpload={false}
					accept=".pdf,.doc,.docx,.png,.jpg,.jpeg"
					onGetUploadList={setFileList}
					className={styles.documentUpload}
					limitSize={10}
					showUploadList={false}
					disabled={false}
				>
					<div className=" bg-gray-light py-4 w-full flex flex-col items-center border border-dashed border-base-border rounded-sm text-gray-light cursor-pointer hover:bg-[#E6F4FF]">
						<InboxOutlined style={{ fontSize: '48px', color: 'var(--base-primary)' }} />
						<div className="text-base pt-5 pb-1 text-foreground">
							Click or drag file to this area to upload
						</div>
						<p className="text-base-label">
							Supported formats: PDF, DOC, PNG, JPG, JPEG Maximum file size: 10MB
						</p>
					</div>
				</UploadFileButton>
			</div>
			{fileList?.length > 0 && (
				<div className="py-6 px-6 border border-base-border rounded-2xl">
					<div className="pb-4 pl-3  flex items-center justify-between border-b border-base-border mb-4">
						<span className="text-xl font-medium text-black-dark">
							Uploaded Document
						</span>
						<p className="text-xs/[20px] px-2 bg-gray-light rounded-sm border border-base-border">
							3 documents
						</p>
					</div>
					{fileList?.map((item, index) => {
						return (
							<div
								className="flex items-center py-4 px-4 border border-base-border rounded-sm gap-[10px] hover:bg-[#E6F4FF]"
								key={index}
							>
								<ArticleSvg style={{ fontSize: 48 }} />
								<div className="flex-1">
									<div className="text-base font-medium pb-2 text-black-dark">
										{item.name}
									</div>
									<p className="text-sm/[22px] font-medium opacity-65 text-black-dark">
										09/16/2025 | 9:45 PM
									</p>
								</div>
								<Space>
									<ViewFile disabled={false} file={item} />
									<DownloadFileButton disabled={false} fileUrl={item.url} />
								</Space>
							</div>
						);
					})}
				</div>
			)}
		</div>
	);
};

export default Documents;
