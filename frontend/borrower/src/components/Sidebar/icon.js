'use client';
import LoanSvg from '@/assets/images/loan.svg';
import PropertySvg from '@/assets/images/property.svg';
import { overrideFill } from '@/lib/reactCommon';
export const sidebarIcon = {
	loan: overrideFill(<LoanSvg />, '#0958D9'),
	property: overrideFill(<PropertySvg />, '#474849'),
};
