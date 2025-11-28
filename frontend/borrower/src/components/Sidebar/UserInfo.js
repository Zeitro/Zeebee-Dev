'use client';
import React from 'react';
import { Dropdown, Popover } from 'antd';
import { CheckOutlined, LogoutOutlined, UserOutlined } from '@ant-design/icons';
import { DeleteBtn } from '@repo/shared-methods';
import { formatFirstLetterName } from '@repo/shared-methods/utils';
import { useRouter } from 'next/navigation';
const Userinfo = ({ userInfo, sidebarCollapsed }) => {
	const router = useRouter();
	const handleLogout = () => {
		router.push('/login');
	};

	const organizationContent = (
		<div className="flex items-center w-[150px]" onClick={(e) => e.stopPropagation()}>
			<div className="bg-base-primary text-white text-base font-medium px-2 rounded-sm">
				Z
			</div>
			<span className="flex-1 common-ellipsis pl-2">ZeeBee</span>
			<CheckOutlined />
		</div>
	);

	const items = [
		{
			key: 'personal',
			label: <span>{userInfo?.email}</span>,
			icon: <UserOutlined />,
		},
		{
			type: 'divider',
		},
		{
			key: 'organization',
			label: (
				<Popover
					content={organizationContent}
					placement="rightTop"
					align={{ offset: [17, -5] }}
					arrow={false}
				>
					<div className="pl-5">Change organization</div>
				</Popover>
			),
		},
		{
			type: 'divider',
		},
		{
			key: 'emailAuth',
			label: <div className="pl-5">Email auth</div>,
		},
		{
			type: 'divider',
		},
		{
			key: 'logout',
			label: (
				<DeleteBtn
					isButton={false}
					title="Are you sure you want to logout?"
					okText="Logout"
					onOk={handleLogout}
				>
					<LogoutOutlined style={{ fontSize: 14, paddingRight: 8 }} />
					Logout
				</DeleteBtn>
			),
		},
	];

	const handleMenuClick = ({ key }) => {
		switch (key) {
			case 'emailAuth':
				router.push('/personalCenter/emailAuth');
				break;

			default:
				break;
		}
	};
	return (
		<Dropdown
			menu={{ items, onClick: handleMenuClick }}
			placement="topLeft"
			trigger={['click']}
		>
			<div
				className={`h-16 border-t border-solid border-base-border text-base-primary flex items-center justify-center p-3.5 gap-2.5 overflow-hidden`}
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
					{formatFirstLetterName(userInfo)}
				</div>
				{!sidebarCollapsed && (
					<div className="flex-1">
						<p className={`text-base/4 text-black-dark font-medium pb-[3px]`}>
							{userInfo?.first_name} {userInfo?.last_name}
						</p>
						<span className="text-gray-medium text-ellipsis overflow-hidden whitespace-nowrap">
							{userInfo?.email}
						</span>
					</div>
				)}
			</div>
		</Dropdown>
	);
};

export default Userinfo;
