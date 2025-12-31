'use client';
import styles from './index.module.scss';
import LoanSvg from '@/assets/images/loan.svg';
import PropertySvg from '@/assets/images/property.svg';
import ManagementSvg from '@/assets/images/management.svg';
import SettingSvg from '@repo/shared-methods/assets/images/setting.svg';
import ShopOutlined from '@/assets/images/shopOutlined.svg';

export const sidebarIcon = {
	loan: <LoanSvg fill="currentColor" className={styles.sidebar_icon} />,
	property: <PropertySvg fill="currentColor" className={styles.sidebar_icon} />,
	management: <ManagementSvg fill="currentColor" className={styles.sidebar_icon} />,
	setting: <SettingSvg className="common-fill" />,
	marketTools: <ShopOutlined className="common-fill" />,
};
