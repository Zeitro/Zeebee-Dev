'use client';
import { Tooltip } from 'antd';
import {
	LockOutlined,
	LogoutOutlined,
	RightOutlined,
	SwapOutlined,
	UserOutlined,
} from '@ant-design/icons';
import { DeleteBtn, SwitchOrganization } from '@repo/shared-methods';
import { useRouter } from 'next/navigation';
import { useUserStore } from '@/store/user';
// import styles from './index.module.scss';
import { useState } from 'react';
import { eventBus } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const useContent = () => {
	const router = useRouter();
	const { userInfo } = useUserStore();
	const [organizationRow, setOrganizationRow] = useState({});

	useDeferredSyncEffect({
		deferred: () => {
			setOrganizationRow(userInfo?.organizations?.find((item) => item.is_default));
		},
		deps: [userInfo],
		delay: 0,
	});

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
		if (item?.id === organizationRow?.id) return;
		item.is_default = true;
		setOrganizationRow({ ...item });
		eventBus.emit(
			'appendPortalComponent',
			<SwitchOrganization
				organization_name={item?.name}
				organization_id={item?.id}
				api="/user/switch_organization"
			/>,
		);
	};

	const organizationContent = (
		<div className="w-[200px]">
			{userInfo?.organizations?.map((item, index) => (
				<div
					key={index}
					className={[
						'px-6 py-2 rounded-lg hover:text-base-primary common-ellipsis cursor-pointer',
						organizationRow?.id === item.id
							? 'bg-[#E6F7FF] text-base-primary'
							: 'text-foreground',
					].join(' ')}
					onClick={(e) => handleOrganizationClick(e, item)}
				>
					{item?.name}
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
						root: '[&_.ant-tooltip-inner]:!py-1 !px-1 text-sm',
					}}
					color="white"
				>
					<div className="flex justify-between items-center">
						<span className="pr-10">Change organization</span>
						<RightOutlined style={{ fontSize: 9 }} />
					</div>
				</Tooltip>
			),
			icon: <SwapOutlined />,
		},
		{
			type: 'divider',
		},
		// {
		// 	key: '/personalCenter/changePassword',
		// 	label: 'Change password',
		// 	icon: <LockOutlined />,
		// },
		// {
		// 	type: 'divider',
		// },
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
