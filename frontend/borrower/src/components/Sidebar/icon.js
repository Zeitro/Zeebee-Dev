'use client';
import styles from './index.module.scss';
import LoanSvg from '@/assets/images/loan.svg';
import PropertySvg from '@/assets/images/property.svg';
export const sidebarIcon = {
	loan: <LoanSvg fill="currentColor" className={styles.sidebar_icon} />,
	property: <PropertySvg fill="currentColor" className={styles.sidebar_icon} />,
};
