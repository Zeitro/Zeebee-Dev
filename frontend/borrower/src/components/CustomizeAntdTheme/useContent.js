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
import { useUserStore } from '@repo/shared-methods/store';
// import styles from './index.module.scss';
import { useRef, useState } from 'react';
import { eventBus } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import useData from '@/app/login/components/useData';
const useContent = () => {
	const router = useRouter();
	const { getUserPermissions } = useData();
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
				domEvent.stopPropagation();
				domEvent.preventDefault();
				break;
			default:
				break;
		}
	};
	// 确认登出后等 Modal 关闭动画结束再跳转，避免弹框未关完就 replace 产生卡顿感
	const shouldLogoutRef = useRef(false);
	const handleLogout = () => {
		shouldLogoutRef.current = true;
	};
	const handleLogoutAfterClose = () => {
		if (!shouldLogoutRef.current) return;
		shouldLogoutRef.current = false;
		window.location.replace('/login');
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
				successCallback={getUserPermissions}
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

	const hasMultipleOrgs = (userInfo?.organizations?.length ?? 0) > 1;
	const currentOrgName = organizationRow?.name?.trim();

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
					title={hasMultipleOrgs ? organizationContent : null}
					placement="rightTop"
					align={{ offset: [17, -5] }}
					arrow={false}
					classNames={{
						root: '[&_.ant-tooltip-inner]:!py-1 !px-1 text-sm',
					}}
					color="white"
				>
					<div className="flex items-center justify-between gap-2">
						<div className="min-w-0 flex-1">
							<div className="leading-5">
								{hasMultipleOrgs ? 'Change organization' : 'Organization'}
							</div>
							{currentOrgName ? (
								<div className="truncate text-xs leading-4 text-base-label">
									{currentOrgName}
								</div>
							) : null}
						</div>
						{hasMultipleOrgs ? <RightOutlined style={{ fontSize: 9 }} /> : null}
					</div>
				</Tooltip>
			),
			icon: <SwapOutlined />,
		},
		{
			type: 'divider',
		},
		{
			key: '/personalCenter/changePassword',
			label: 'Change password',
			icon: <LockOutlined />,
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
					afterClose={handleLogoutAfterClose}
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
