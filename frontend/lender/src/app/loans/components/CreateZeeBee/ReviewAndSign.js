'use client';
import React, { useMemo } from 'react';
import dynamic from 'next/dynamic';
import { Form } from 'antd';

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
const h = '30px';
const w = '35%';
const left = '26%';
const height = '16px';
const width = '14%';
const right = '16%';

// 在这里定义你的合同的所有签名区
export const contractSignatureZones = [
	{
		id: 'applicant1_signature', // 唯一ID
		page: 3, // 所在的 PDF 页码
		// 下面是你在第一步中手动测量出的百分比值
		top: '32.2%',
		left,
		width: w,
		height: h,
		date: {
			top: '34%',
			right,
			width,
			height,
		},
		paramsKey: 'borrower_info',
	},
	{
		id: 'applicant2_signature',
		page: 3,
		top: '37%',
		left,
		width: w,
		height: h,
		date: {
			top: '38.8%',
			right,
			width,
			height,
		},
		paramsKey: 'borrower_info',
	},
	{
		id: 'applicant5_signature',
		page: 10,
		top: '81.8%',
		left,
		width: w,
		height: h,
		date: {
			top: '83.5%',
			right,
			width,
			height,
		},
		paramsKey: 'business_purpose',
	},
	{
		id: 'applicant6_signature',
		page: 10,
		top: '86.7%',
		left,
		width: w,
		height: h,
		date: {
			top: '88.4%',
			right,
			width,
			height,
		},
		paramsKey: 'business_purpose',
	},
	{
		id: 'applicant7_signature',
		page: 11,
		top: '71.3%',
		left,
		width: w,
		height: h,
		date: {
			top: '72.9%',
			right,
			width,
			height,
		},
		paramsKey: 'patriot_act',
	},
	{
		id: 'applicant8_signature',
		page: 11,
		top: '76.1%',
		left,
		width: w,
		height: h,
		date: {
			top: '77.8%',
			right,
			width,
			height,
		},
		paramsKey: 'patriot_act',
	},
	{
		id: 'applicant3_signature',
		page: 13,
		top: '58.5%',
		left,
		width: w,
		height: h,
		date: {
			top: '60.2%',
			right,
			width,
			height,
		},
		paramsKey: 'borrower_authorization',
	},
	{
		id: 'applicant4_signature',
		page: 13,
		top: '63.4%',
		left,
		width: w,
		height: h,
		date: {
			top: '65.1%',
			right,
			width,
			height,
		},
		paramsKey: 'borrower_authorization',
	},
	// ... 如果有更多签名区，继续在这里添加
];
// isNotRequest 用于已有合同的链接查看，不请求接口返回数据流
const ReviewAndSign = ({ form, rowData, isNotRequest = false }) => {
	const applicants = Form.useWatch('applicants', form);
	const params = useMemo(() => {
		// 因zeeBee详情会点击步骤跳到最后一步，导致getFieldsValue取不到全部数据，所以需将rowData直接传入
		const values = { ...(rowData || form.getFieldsValue() || {}) };
		if (values?.loan_borrower?.tin?.tin_type) {
			values.loan_borrower.tin_type = values?.loan_borrower?.tin?.tin_type;
		}
		if (values?.loan_borrower?.tin?.tin) {
			values.loan_borrower.tin = values?.loan_borrower?.tin?.tin;
		}
		if (values?.signatures) {
			delete values.signatures;
		}
		return values;
	}, [form, rowData]);

	// zeeBee详情时，拿到签订合同的外部链接
	const fullLink = useMemo(() => {
		const signedContract = rowData?.loan_documents?.find(
			(item) => item.file_type === 'Apply Form',
		);
		return signedContract?.files?.[0]?.file_url;
	}, [rowData]);

	// 根据申请人个数，返回需要渲染的签名区
	const getZones = () => {
		const count = applicants?.length ? applicants?.length : 1;
		let newZones = [];
		if (fullLink || isNotRequest) return newZones;
		contractSignatureZones.reduce((a, b) => {
			const arr = a?.filter((item) => item.page === b.page);
			if (arr.length < count) {
				newZones.push(b);
			}
			return newZones;
		}, []);
		return newZones;
	};

	const changeSignatures = (signatures) => {
		form.setFieldsValue({
			signatures,
		});
	};
	return (
		<div className="h-full w-full">
			<div className="text-3xl pb-14 font-medium text-center">
				Review and sign your application
			</div>
			<div className="h-[calc(100% - 120px)]">
				<ApplicationFormViewPDF
					url={fullLink || isNotRequest ? fullLink : '/pdf/apply_form'}
					zones={getZones()}
					params={params}
					changeSignatures={changeSignatures}
					isFullLink={!!fullLink || isNotRequest}
				/>
			</div>
			<Form.Item name="signatures" rules={[{ required: true }]} noStyle></Form.Item>
		</div>
	);
};
export default ReviewAndSign;
