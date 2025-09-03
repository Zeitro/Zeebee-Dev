'use client';
import React from 'react';
import dynamic from 'next/dynamic';

const ApplicationFormViewPDF = dynamic(
	() => import('@/components/ApplicationFormViewPDF'), // Adjust the path to your component
	{
		ssr: false, // This is the crucial part
		// Optional: Add a loading component to improve UX
		loading: () => (
			<p className="text-center w-full h-[200px] text-base">Asynchronous loading...</p>
		),
	},
);
const height = '20px';
const width = '12.5%';
const right = '6%';

// 在这里定义你的合同的所有签名区
const contractSignatureZones = [
	{
		id: 'applicant1_signature', // 唯一ID
		page: 3, // 所在的 PDF 页码
		// 下面是你在第一步中手动测量出的百分比值
		top: '38%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '38%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant2_signature',
		page: 3,
		top: '42.5%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '42.5%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant3_signature',
		page: 9,
		top: '87.5%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '87.5%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant4_signature',
		page: 9,
		top: '92%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '92%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant5_signature',
		page: 10,
		top: '87.5%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '87.5%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant6_signature',
		page: 10,
		top: '92%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '92%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant7_signature',
		page: 12,
		top: '69%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '69%',
			right,
			width,
			height,
		},
	},
	{
		id: 'applicant8_signature',
		page: 12,
		top: '73.6%',
		left: '28%',
		width: '42%',
		height: '20px',
		date: {
			top: '73.6%',
			right,
			width,
			height,
		},
	},
	// ... 如果有更多签名区，继续在这里添加
];
const ReviewAndSign = () => {
	return (
		<div className="h-full w-full">
			<div className="text-3xl pb-14 font-medium text-center">
				Review and sign your application
			</div>
			<div className="h-[calc(100% - 120px)]">
				<ApplicationFormViewPDF url="/api/getApplication" zones={contractSignatureZones} />
			</div>
		</div>
	);
};
export default ReviewAndSign;
