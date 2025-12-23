'use client';
import React, { useState } from 'react';
import { Dropdown } from 'antd';
import { formatFirstLetterName } from '@repo/shared-methods/utils';
import useContent from './useContent';
import { usePathname } from 'next/navigation';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const Userinfo = ({ userInfo, sidebarCollapsed }) => {
	const pathname = usePathname();
	const [selectedKey, setSelectedKey] = useState('');
	const { personalList, handleClickPersonalItem } = useContent();

	// 路由与个人中心的选项匹配，则设置选中
	useDeferredSyncEffect({
		deferred: () => {
			const selected = personalList.some((i) => i.key === pathname);
			if (selected) {
				setSelectedKey(pathname);
			}
		},
		deps: [pathname, personalList],
		delay: 0,
	});

	return (
		<Dropdown
			menu={{
				items: personalList,
				selectedKeys: [selectedKey],
				onClick: handleClickPersonalItem,
			}}
			placement="topLeft"
			trigger={['click']}
		>
			<div
				className={`h-16 border-t border-solid border-base-border text-base-primary p-3.5 gap-2.5 cursor-pointer flex items-center`}
				onClick={(e) => e.preventDefault()}
			>
				{/* <Image
					src="/globe.svg"
					width={32}
					height={32}
					alt="Picture of the author"
					className="shrink-0"
				/> */}
				<div
					className={[
						`w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-sm text-base-label bg-base-active`,
					].join(' ')}
				>
					{formatFirstLetterName(userInfo) || ''}
				</div>
				<div
					className="flex-1"
					style={{
						opacity: sidebarCollapsed ? 0 : 1,
						transition: 'opacity 0.3s cubic-bezier(0.2, 0, 0, 1) 0s',
					}}
				>
					<p className={`text-base/4 text-black-dark font-medium pb-[3px]`}>
						{userInfo?.first_name
							? `${userInfo?.first_name} ${userInfo?.last_name}`
							: '--'}
					</p>
					<span className="text-gray-medium text-ellipsis overflow-hidden whitespace-nowrap">
						{userInfo?.email || '--'}
					</span>
				</div>
			</div>
		</Dropdown>
	);
};

export default Userinfo;
