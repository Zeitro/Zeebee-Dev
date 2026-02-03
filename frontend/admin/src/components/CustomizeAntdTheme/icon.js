'use client';
import BusinessSvg from '@/assets/images/business_center.svg';
import { FileTextOutlined } from '@ant-design/icons';
// import SettingSvg from '@repo/shared-methods/assets/images/setting.svg';
export const sidebarIcon = {
	business: <BusinessSvg className="common-fill" />,
	template: <FileTextOutlined style={{ fontSize: 20, paddingLeft: 2 }} />,
};
