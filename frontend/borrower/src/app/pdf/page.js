'use client';
import React from 'react';
import dynamic from 'next/dynamic';

const ApplicationFormViewPDF = dynamic(
	() => import('@/components/ApplicationFormViewPDF'), // Adjust the path to your component
	{
		ssr: false, // This is the crucial part
		// Optional: Add a loading component to improve UX
		loading: () => <p>正在加载...</p>,
	},
);

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
	},
	{
		id: 'applicant2_signature',
		page: 3,
		top: '42.5%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	{
		id: 'applicant3_signature',
		page: 9,
		top: '87.5%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	{
		id: 'applicant4_signature',
		page: 9,
		top: '92%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	{
		id: 'applicant5_signature',
		page: 10,
		top: '87.5%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	{
		id: 'applicant6_signature',
		page: 10,
		top: '92%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	{
		id: 'applicant7_signature',
		page: 12,
		top: '69%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	{
		id: 'applicant8_signature',
		page: 12,
		top: '73.6%',
		left: '28%',
		width: '42%',
		height: '20px',
	},
	// ... 如果有更多签名区，继续在这里添加
];
const PDF = () => {
	return (
		<div className="h-full w-full">
			<ApplicationFormViewPDF url="/api/getApplication" zones={contractSignatureZones} />
		</div>
	);
};
export default PDF;
