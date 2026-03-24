'use client';
import LoanSvg from '@/assets/images/loan.svg';
import BusinessSvg from '@/assets/images/business_center.svg';
import { FileTextOutlined } from '@ant-design/icons';
// import SettingSvg from '@repo/shared-methods/assets/images/setting.svg';
export const sidebarIcon = {
	loan: <LoanSvg fill="currentColor" className="common-fill" />,
	business: <BusinessSvg className="common-fill" />,
	template: <FileTextOutlined style={{ fontSize: 20, paddingLeft: 2 }} />,
};
