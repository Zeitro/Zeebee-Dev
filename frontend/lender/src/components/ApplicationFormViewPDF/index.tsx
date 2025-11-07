'use client';

import { useState, useEffect } from 'react';
import { Document, Page, pdfjs } from 'react-pdf';
import SignatureModal from '../SignatureModal';
import Image from 'next/image';
import { Modal } from 'antd';
import { formatDate } from '@/lib/date';
import { createRequest } from '@/lib/fetch';
import EditSvg from '@/assets/images/edit.svg';
pdfjs.GlobalWorkerOptions.workerSrc = '/pdf.worker.min.mjs';

// 定义签名区 (Hotspot/Zone) 的类型
export interface SignatureZone {
	id: string; // 每个区域的唯一标识
	page: number;
	// 使用百分比字符串
	left: string;
	top: string;
	width: string;
	height: string;
	paramsKey: string;
	date: {
		top: string;
		left?: string;
		right?: string;
		width: string;
		height: string;
	};
}

export interface SignatureInfo {
	image: string;
	date: string;
	paramsKey: string;
}

// 定义签名数据存储的类型，使用 zoneId 作为 key
interface SignatureData {
	[zoneId: string]: SignatureInfo; // key 是 zoneId, value 是签名图片的 base64 数据
}

interface ContractViewerProps {
	url: string;
	zones?: SignatureZone[]; // 接收一个签名区配置数组
	params?: {
		[key: string]: unknown;
	};
	changeSignatures?: (signatures: SignatureData) => void;
	isFullLink?: boolean;
}
const config = {
	title: 'Signature prompt!',
	content: 'You have already signed, do you want to sign again?',
};
export default function ContractViewer({
	url,
	zones,
	params,
	isFullLink = false,
	changeSignatures,
}: ContractViewerProps) {
	const [pdfFile, setPdfFile] = useState<File | null>(null);
	const [numPages, setNumPages] = useState<number>(0);
	const [isModalOpen, setIsModalOpen] = useState(false);

	// 新的状态：存储签名数据和当前正在操作的签名区
	const [signatures, setSignatures] = useState<SignatureData>({});
	const [activeZoneId, setActiveZoneId] = useState<string | null>(null);

	const [isLoading, setIsLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [modal, contextHolder] = Modal.useModal();
	const { request } = createRequest();

	useEffect(() => {
		const fetchPdf = async () => {
			try {
				if (!url) return;
				setIsLoading(true);
				setError(null);
				const config: any = {
					isExternalLink: isFullLink,
				};
				if (!isFullLink) {
					config.body = params;
					config.method = 'POST';
				}
				const response = await request(url, config);
				const blob = await response.blob();
				if (blob.size === 0) {
					throw new Error('The obtained PDF file is empty!');
				}
				const file = new File([blob], 'contract.pdf', { type: 'application/pdf' });
				setPdfFile(file);
			} catch (e: any) {
				setError(e.message);
				console.error('获取或处理 PDF 时出错:', e);
			} finally {
				setIsLoading(false);
			}
		};
		fetchPdf();
	}, [url]);

	useEffect(() => {
		const info: any = {};
		const signaturesArr = Object.values(signatures);
		signaturesArr?.map((item: any) => {
			info[item.paramsKey] = [
				...(info?.[item.paramsKey] || []),
				{
					signature_image: item.image,
					signature_date: item.date,
				},
			];
		});
		// 剩下的签名数量
		info.remainingSignatureCount = (zones?.length || 0) - signaturesArr?.length;
		changeSignatures?.(info);
	}, [signatures]);

	const onDocumentLoadSuccess = ({ numPages }: { numPages: number }) => {
		setNumPages(numPages);
	};

	const onDocumentLoadError = (error: Error) => {
		setError(`react-pdf Parsing the file failed: ${error.message}`);
	};

	// 当点击一个签名区时触发
	const handleZoneClick = async (zoneId: string) => {
		// 如果这个区域已经签过名了，可以选择重新签名或不做任何事
		if (signatures[zoneId]) {
			const confirmed = await modal.warning(config);
			if (confirmed) {
				setActiveZoneId(zoneId);
				setIsModalOpen(true);
			}
		} else {
			setActiveZoneId(zoneId);
			setIsModalOpen(true);
		}
	};

	// 当在模态框中确认签名后触发
	const handleSignatureConfirm = (signatureImage: string) => {
		if (activeZoneId) {
			const row = zones?.find((zone) => zone.id === activeZoneId);
			setSignatures({
				...signatures,
				[activeZoneId]: {
					image: signatureImage,
					date: formatDate(new Date()),
					paramsKey: row?.paramsKey || '',
				},
			});
		}
		setIsModalOpen(false);
		setActiveZoneId(null);
	};

	if (isLoading || error) {
		return (
			<div className="w-full h-[200px] text-center text-base">
				{isLoading && 'Data loading, please wait...'}
				{error && `Loading failed: ${error}`}
			</div>
		);
	}

	return (
		<>
			{contextHolder}
			<div className="w-full h-full overflow-y-auto">
				<Document
					file={pdfFile}
					onLoadSuccess={onDocumentLoadSuccess}
					onLoadError={onDocumentLoadError}
					loading={
						<div className="w-full h-[200px] text-center text-base">Loading PDF...</div>
					}
				>
					{Array.from(new Array(numPages), (el, index) => {
						const pageNumber = index + 1;
						// 筛选出当前页面需要渲染的签名区
						const pageZones = zones?.filter((z) => z.page === pageNumber);

						return (
							<div
								key={`page_wrapper_${pageNumber}`}
								className="flex items-center justify-center mb-2"
							>
								<div className="w-max relative">
									<Page pageNumber={pageNumber} />

									{/* 4. 现在，签名区相对于这个“贴身”的父级来定位，位置将永远准确 */}
									{pageZones?.map((zone) => (
										<div
											key={zone.id}
											onClick={() => handleZoneClick(zone.id)}
											className="absolute cursor-pointer z-[999] flex items-end justify-center"
											style={{
												left: zone.left,
												top: zone.top,
												width: zone.width,
												height: zone.height,
												// backgroundColor: 'pink',
											}}
										>
											{signatures[zone.id] ? (
												<Image
													unoptimized
													src={signatures[zone.id]?.image}
													alt={`${zone.id} signature`}
													layout="fill"
													objectFit="contain"
												/>
											) : (
												<EditSvg className="w-6 h-6" />
											)}
										</div>
									))}
									{pageZones?.map((zone, index) => (
										<div
											key={index}
											className="absolute z-[999] text-xs text-center"
											style={{
												right: zone?.date?.right,
												left: zone?.date?.left,
												top: zone?.date?.top,
												width: zone?.date?.width,
												height: zone?.date?.height,
											}}
										>
											{signatures[zone.id]?.date}
										</div>
									))}
								</div>
							</div>
						);
					})}
				</Document>
			</div>

			<SignatureModal
				isOpen={isModalOpen}
				onClose={() => setIsModalOpen(false)}
				onConfirm={handleSignatureConfirm}
			/>
		</>
	);
}
