'use client';
import { Tooltip } from 'antd';
import { LogoutOutlined, RightOutlined, SwapOutlined, UserOutlined } from '@ant-design/icons';
import { DeleteBtn, SwitchOrganization } from '@repo/shared-methods';
import { useRouter } from 'next/navigation';
import { useUserStore } from '@/store/user';
import styles from './index.module.scss';
import { useState } from 'react';
import { eventBus } from '@repo/shared-methods/utils';
const useContent = () => {
	const router = useRouter();
	const { userInfo } = useUserStore();
	const [organizationRow, setOrganizationRow] = useState({});

	// 个人中心选中
	const handleClickPersonalItem = ({ key, domEvent }) => {
		switch (key) {
			case '/personalCenter/emailAuth':
				router.push('/personalCenter/emailAuth');
				break;
			case '/personalCenter/changePassword':
				router.push('/personalCenter/changePassword');
				break;
			case 'organization':
				console.log(domEvent);
				domEvent.stopPropagation();
				domEvent.preventDefault();
				break;
			default:
				break;
		}
	};
	// 退出
	const handleLogout = () => {
		router.push('/login');
	};

	// 切换组织
	const handleOrganizationClick = (e, item) => {
		e.stopPropagation();
		item.isDefault = true;
		setOrganizationRow({ ...item });
		eventBus.emit(
			'appendPortalComponent',
			<SwitchOrganization organization={item?.organization || 'ZeeBee'} />,
		);
	};

	const organizationContent = (
		<div className="w-[200px]">
			{(userInfo?.organizations || [{}, {}])?.map((item, index) => (
				<div
					key={index}
					className={[
						'px-6 py-2 rounded-lg text-foreground hover:text-base-primary common-ellipsis cursor-pointer',
						organizationRow?.isDefault && 'bg-[#E6F7FF]',
					].join(' ')}
					onClick={(e) => handleOrganizationClick(e, item)}
				>
					ZeeBee
				</div>
			))}
		</div>
	);

	const personalList = [
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
				<Tooltip
					title={organizationContent}
					placement="rightTop"
					align={{ offset: [17, -5] }}
					arrow={false}
					classNames={{
						root: styles.organizationsRoot,
					}}
					color="white"
				>
					<div className="flex justify-between items-center">
						<span>Change organization</span>
						<RightOutlined style={{ fontSize: 9 }} />
					</div>
				</Tooltip>
			),
			icon: <SwapOutlined />,
		},
		{
			type: 'divider',
		},
		{
			key: '/personalCenter/emailAuth',
			label: <div className="pl-5">Email auth</div>,
		},
		{
			type: 'divider',
		},
		{
			key: '/personalCenter/changePassword',
			label: <div className="pl-5">Change password</div>,
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
	return { personalList, handleClickPersonalItem };
};
export default useContent;
