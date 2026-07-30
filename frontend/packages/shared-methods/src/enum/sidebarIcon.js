// 汇总admin/lender/borrower三端的icon
import LoanSvg from '@repo/shared-methods/assets/images/menu/loan.svg';
import BusinessSvg from '@repo/shared-methods/assets/images/menu/business_center.svg';
import PropertySvg from '@repo/shared-methods/assets/images/menu/property.svg';
import ManagementSvg from '@repo/shared-methods/assets/images/menu/management.svg';
import SettingSvg from '@repo/shared-methods/assets/images/menu/setting.svg';
import ShopOutlined from '@repo/shared-methods/assets/images/menu/shopOutlined.svg';
import Dashboard from '@repo/shared-methods/assets/images/menu/dashboard.svg';
import Transaction from '@repo/shared-methods/assets/images/menu/sync.svg';
import { ApartmentOutlined, FileTextOutlined } from '@ant-design/icons';
export const iconData = {
	loan: <LoanSvg fill="currentColor" className="common-fill" />,
	business: <BusinessSvg className="common-fill" />,
	template: <FileTextOutlined style={{ fontSize: 20, paddingLeft: 2 }} />,
	resource: <ApartmentOutlined style={{ fontSize: 20, paddingLeft: 2 }} />,
	property: <PropertySvg fill="currentColor" className="common-fill" />,
	management: <ManagementSvg fill="currentColor" className="common-fill" />,
	setting: <SettingSvg className="common-fill" />,
	marketTools: <ShopOutlined className="common-fill" />,
	dashboard: <Dashboard fill="currentColor" className="common-fill" />,
	transaction: <Transaction className="common-fill" />,
};
