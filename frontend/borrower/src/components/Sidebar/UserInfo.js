'use client';
import React from 'react';
import { Dropdown } from 'antd';
import { formatFirstLetterName } from '@repo/shared-methods/utils';
import useContent from './useContent';
const Userinfo = ({ userInfo, sidebarCollapsed }) => {
	const { personalList, handleClickPersonalItem } = useContent();

	return (
		<Dropdown
			menu={{ items: personalList, onClick: handleClickPersonalItem }}
			placement="topLeft"
			trigger={['click']}
		>
			<div
				className={`h-16 border-t border-solid border-base-border text-base-primary flex items-center justify-center p-3.5 gap-2.5 overflow-hidden cursor-pointer`}
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
					{formatFirstLetterName(userInfo) || '-'}
				</div>
				{!sidebarCollapsed && (
					<div className="flex-1">
						<p className={`text-base/4 text-black-dark font-medium pb-[3px]`}>
							{userInfo?.first_name
								? `${userInfo?.first_name} ${userInfo?.last_name}`
								: '--'}
						</p>
						<span className="text-gray-medium text-ellipsis overflow-hidden whitespace-nowrap">
							{userInfo?.email || '--'}
						</span>
					</div>
				)}
			</div>
		</Dropdown>
	);
};

export default Userinfo;
