'use client';
import React from 'react';
import { Dropdown } from 'antd';
import { LogoutOutlined, UserOutlined } from '@ant-design/icons';
import { DeleteBtn } from '../NewModal/BtnPopUpModal';
import { formatFirstLetterName } from '@/lib/format';
import { useRouter } from 'next/navigation';
const Userinfo = ({ userInfo, sidebarCollapsed }) => {
	const router = useRouter();
	const handleLogout = () => {
		router.push('/login');
	};

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
			case 'personal':
				break;

			default:
				break;
		}
	};
	return (
		<Dropdown menu={{ items, onClick: handleMenuClick }}>
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
