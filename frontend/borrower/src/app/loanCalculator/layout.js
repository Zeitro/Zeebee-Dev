'use client';
import Image from 'next/image';
import React from 'react';
import styles from './index.module.scss';

const EstimateLoanLayout = ({ children }) => {
	return (
		<div className="h-full w-full flex flex-col">
			<div className={['h-20 flex-shrink-0 flex items-center pl-6', styles.header].join(' ')}>
				<div className="h-7">
					<Image src="/images/logo.svg" width={95} height={30} alt="logo" />
				</div>
			</div>
			<div className="flex-1 bg-[#E4F3FF] overflow-hidden flex justify-center">
				{children}
			</div>
		</div>
	);
};
export default EstimateLoanLayout;
