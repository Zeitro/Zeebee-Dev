'use client';
import { Menu, Button } from 'antd';
import React from 'react';
import { useGlobalStore } from '@/store/global';
import styles from './index.module.scss';
import { navRouter } from '@config/routes';
import { generateMenuBar } from '@config/common';
import Image from 'next/image';
const SidebarDom = () => {
	const { siberWidth, changeSiberWidth } = useGlobalStore();

	return (
		<div
			style={{ width: siberWidth }}
			className={[
				'shrink-0 h-full flex flex-col border-r border-base-border',
				styles.sidebar,
			].join(' ')}
		>
			<div
				className={['h-14 mx-3.5 shrink-0 border-b border-base-border', styles.logo].join(
					' ',
				)}
			>
				Zeebee
			</div>
			<Menu
				mode="vertical"
				items={generateMenuBar(navRouter)}
				className={['flex-1 overflow-hidden border-r-0', styles.menuList].join(' ')}
			/>
			<div className={`h-16 border-t border-solid border-base-border text-base-primary flex items-center justify-center p-3.5 gap-2.5`}>
				<Image src="/globe.svg" width={32} height={32} alt="Picture of the author" className='shrink-0'/>
                <div className='flex-1'>
                    <p className={`text-base/4 text-black-dark font-medium pb-[3px]`}>username</p>
                    <span className='text-gray-medium'>kevinmo@email.com</span>
                </div>
            </div>
		</div>
	);
};
export default SidebarDom;
