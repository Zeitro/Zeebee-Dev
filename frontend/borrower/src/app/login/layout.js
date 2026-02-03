'use client';
import { useUserStore } from '@/store/user';
import Image from 'next/image';
import { usePathname, useRouter } from 'next/navigation';
import React, { useEffect } from 'react';
import { clearAllCookie } from '@repo/shared-methods/utils';
import { noTokenRouter } from '@config/routes';
import { flattenNavRouter, matchRoute } from '@config/common';
const LoginLayout = ({ children }) => {
	const router = useRouter();
	const pathname = usePathname();
	const { clearUserInfo } = useUserStore();
	useEffect(() => {
		clearUserInfo();
		clearAllCookie();
		const noNavList = flattenNavRouter(noTokenRouter);
		const isLogin = noNavList.some((item) => matchRoute(item.path, pathname));
		// localStorageUtil.clear();
		// bug：因为出现cookie到期，内容区域已经显示登录页面，但是侧边栏和面包屑都还存在，且路由还是之前的页面路由，并不是登录路由
		// 此bug还不知道怎么触发的，暂时这样处理
		if (!isLogin) {
			router.replace('/login');
		}
	}, [clearUserInfo, pathname, router]);
	return (
		<div className="h-full w-full flex flex-col">
			<div className="h-20 flex-shrink-0 flex items-center pl-6">
				<div className="h-7">
					<Image
						src="/images/logo.svg"
						width={95}
						height={30}
						alt="logo"
						className="!h-full"
					/>
				</div>
			</div>
			{children}
		</div>
	);
};
export default LoginLayout;
