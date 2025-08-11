'use client';
import styles from './index.module.scss';
import LoanSvg from '@/assets/images/loan.svg';
import PropertySvg from '@/assets/images/property.svg';
import { overrideFill } from '@/lib/reactCommon';
export const sidebarIcon = {
	loan: overrideFill(<LoanSvg fill="currentColor" className={styles.sidebar_icon} />, '#0958D9'),
	property: overrideFill(
		<PropertySvg fill="currentColor" className={styles.sidebar_icon} />,
		'#474849',
	),
};
