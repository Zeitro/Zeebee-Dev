import { routers } from './routes';
import _ from 'lodash';
// 生成菜单栏
export const generateMenuBar = (list) => {
	return _.cloneDeep(list)
		.filter((i) => !i.hideMenu)
		.map((item) => {
			if (item.router) {
				item.router = generateMenuBar(item.router);
			}
			const newItem = {
				key: item.path,
				label: item.title,
			};
			if (item?.router?.length) {
				newItem.children = item.router;
			}
			return newItem;
		});
};

// 将多层级嵌套路由变为一层级路由，方便查找
export const flattenNavRouter = (arr, keyStr = 'router') => {
	return _.flatMap(arr, (item) =>
		item.router ? [_.omit(item, keyStr), ...flattenNavRouter(item.router)] : [item],
	);
};

/**
 * 可匹配动态路由
 * dynamicPath 动态路由/loans/123455
 * path 比较的固定路由/loans/:id
 * */
export const matchRoute = (path, dynamicPath) => {
	if (path === dynamicPath) return true;
	const pattern = '^' + path.replace(/:[^/]+/g, '[^/]+') + '$';
	return new RegExp(pattern).test(dynamicPath);
};

//根据路由生成面包屑
export const getCrumbs = (currentPath, routes) => {
	const segments = currentPath.split('/').filter(Boolean);
	let path = '';
	const breadcrumbs = [];

	for (let i = 0; i < segments.length; i++) {
		path += '/' + segments[i];
		// 依次尝试精确匹配和动态参数匹配
		let match = routes.find((r) => r.path === path);
		if (!match) {
			// 动态参数匹配
			match = routes.find((r) => matchRoute(r.path, path));
		}
		if (match) {
			// 动态参数面包屑可自定义显示
			let title = match.title;
			if (/:/.test(match.path)) {
				// 取实际 id 或自定义
				title = '详情'; // 或 segments[i]
			}
			let crumbItem = {
				title,
			};
			// 最后一个不需要path
			if (segments.length - 1 > i) {
				crumbItem.href = path;
			}
			breadcrumbs.push(crumbItem);
		}
	}
	return breadcrumbs;
};

// 根据pathname获取前端路由配置
export const getCurrentRouterItem = (currentPathname, allRouters = routers) => {
	const routers = flattenNavRouter(allRouters)
	const routerItem = routers.find((item) => matchRoute(item.path, currentPathname));
	return routerItem;
};
