import { NextResponse } from 'next/server';
import { flattenNavRouter, matchRoute } from '@repo/shared-methods/lib/commonRouter';
import { noTokenRouter } from '@config/routes';

export function middleware(request) {
	const { pathname } = request.nextUrl;

	// cookie获取token/权限路由
	const permissionStore = JSON.parse(
		request.cookies.get(`permissionStore_${process.env.NEXT_PUBLIC_SYSTEM_TYPE}`)?.value ||
			'{}',
	);
	const accessPaths = permissionStore?.state?.accessPaths || [];
	const access = permissionStore?.state?.access;
	// 路由白名单
	const noTokenList = flattenNavRouter(noTokenRouter);
	const whiteList = noTokenList.map((item) => item.path);

	if (!access && !whiteList.includes(pathname)) {
		console.log(
			'noToken----------middleware---===',
			pathname,
			'++++',
			permissionStore?.state?.access,
		);

		// token到期且非白名单路由，重定向到登录页
		return NextResponse.redirect(new URL('/login', request.url));
	}

	// 权限判断-排除特殊404页面
	const isAuthSuccessful = accessPaths.some((path) => matchRoute(path, pathname));
	if (!isAuthSuccessful && !whiteList.includes(pathname)) {
		console.log('no permission---', pathname, isAuthSuccessful);

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
