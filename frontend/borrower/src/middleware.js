import { NextResponse } from 'next/server';
import { getCurrentRouterItem } from '@config/common';
import { routers } from '@config/routes';
export function middleware(request) {
	const { pathname } = request.nextUrl;

	// 2. 获取 cookie/token
	const token = request.cookies.get('token')?.value;
	// 路由白名单
	const whiteList = [
		'/login',
		'/login/resetPassword',
		'/login/register',
		'/login/sendRegisterTips',
	];
	const routerItem = getCurrentRouterItem(pathname, routers);

	if (!token && !whiteList.includes(pathname)) {
		console.log('noToken----------middleware---===', pathname);

		// 未登录且非白名单，重定向到登录页
		return NextResponse.redirect(new URL('/login', request.url));
	}

	const isAuthSuccessful = [].every((role) => routerItem?.roles?.includes(role));
	//  其他权限逻辑（如角色判断等）
	if (routerItem && routerItem?.roles?.length && !isAuthSuccessful) {
		return NextResponse.redirect(new URL('/not-found', request.url));
	}

	// 6. 允许通过
	return NextResponse.next();
}

export const config = {
	// 排除Next.js内部资源、静态文件、API路径
	matcher: [
		'/((?!_next/static|_next/image|api/|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
	],
};
