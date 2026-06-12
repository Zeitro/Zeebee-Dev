import { NextResponse } from 'next/server';
import { flattenNavRouter, matchRoute } from '@repo/shared-methods/lib/commonRouter';
import { readPermissionFromCookieStore } from '@repo/shared-methods/lib/permissionStorage';
import { noTokenRouter } from '@config/routes';

export function middleware(request) {
	const { pathname } = request.nextUrl;

	const { access, accessPaths } = readPermissionFromCookieStore(request.cookies);
	// 路由白名单
	const noTokenList = flattenNavRouter(noTokenRouter);
	const whiteList = noTokenList.map((item) => item.path);

	if (!access && !whiteList.includes(pathname)) {
		console.log('noToken----------middleware---===', pathname, '++++', access);

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
	// 排除 Next 内部、API、常见静态后缀（与上方短路逻辑一致，减轻中间件编译/匹配差异带来的漏网）
	matcher: [
		'/((?!_next/|api/|favicon.ico|robots\\.txt|sitemap\\.xml|.*\\.(?:svg|png|jpe?g|gif|ico|webp|avif|mjs|js|css|map|json|txt|xml|webmanifest|woff2?|ttf|otf|wasm|pdf|xlsx?)$).*)',
	],
};
