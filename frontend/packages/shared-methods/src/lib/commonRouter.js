/**
 * 获取按钮权限标识
 * @param {Array} backendMenus - 后端返回的菜单列表
 * @returns {Array} 权限标识数组
 */
export const extractPermissionTags = (backendMenus) => {
	const tags = [];

	const extract = (menus) => {
		menus.forEach((menu) => {
			if (menu.tag) {
				tags.push(menu.tag);
			}
			if (menu.children) {
				extract(menu.children);
			}
		});
	};

	extract(backendMenus);
	return tags;
};

/**
 * 合并菜单配置
 * @param {Array} backendMenus - 后端返回的菜单列表
 * @param {*Object} staticMenusSettings 额外配置项（接口没有的字段，只有前端需要的字段）
 * @returns {Array} 合并后的菜单列表
 */
export const mergeMenu = (backendMenus, staticMenusSettings) => {
	return backendMenus.map((menu) => {
		// 递归查找静态配置
		const staticMenu = staticMenusSettings?.[menu.path] || {};
		let merged = { ...menu, ...staticMenu };
		if (menu?.children) {
			merged.children = mergeMenu(menu.children, staticMenusSettings);
		}
		return merged;
	});
};

/**
 * 一次性处理菜单：过滤 + 提取权限标签
 * @param {Array} menus - 菜单列表
 * @param {Object} menuSetting - 菜单额外配置项
 * @returns {Object} { filteredMenus, permissionTags }
 * permissionTags - 权限标识数组
 * filteredMenus - 过滤后的菜单列表
 * accessPaths - 访问权限路由
 */
export const processMenus = (menus, menuSetting) => {
	const permissionTags = [];
	const accessPaths = [];
	const filterAndExtract = (items) => {
		return items
			.filter((menu) => {
				// 提取权限标签（包括 Button 类型）
				if (menu.tag) {
					permissionTags.push(menu.tag);
				}
				if (menu.path) {
					accessPaths.push(menu.path);
				}
				// 只保留 Menu 类型用于导航
				return menu.resource_type === 'Menu' && menu.status === 'Active';
			})
			.map((menu) => {
				const filteredChildren = menu.children ? filterAndExtract(menu.children) : [];
				const setting = menuSetting?.[menu?.path] || {};
				return {
					...menu,
					...(filteredChildren.length > 0 ? { children: filteredChildren } : {}),
					...setting,
				};
			});
	};

	const filteredMenus = filterAndExtract(menus);

	return {
		filteredMenus,
		permissionTags,
		accessPaths,
	};
};

// 生成菜单栏，转成Menu组件所需字段（排除需隐藏的菜单）
export const generateMenuBar = (list, iconConfig, sidebarCollapsed) => {
	return JSON.parse(JSON.stringify(list))
		.filter((i) => i.is_show)
		.map((item, index) => {
			const newItem = {
				key: item.path,
				label: item?.name,
			};
			if (iconConfig) {
				newItem.icon = iconConfig[item.icon];
			}

			if (item.children?.length) {
				const children = generateMenuBar(item.children, iconConfig);
				if (children.length) {
					newItem.children = sidebarCollapsed
						? [
								{
									...newItem,
									key: `group-${newItem.key}-${index}`,
									type: 'group',
									children,
								},
							]
						: children;
				}
			}

			return newItem;
		});
};

// 将多层级嵌套路由变为一层级路由，方便查找
export function flattenNavRouter(routes = [], childName = 'router') {
	return routes?.reduce((acc, next) => {
		if (next[childName]) {
			return [...acc, next, ...flattenNavRouter(next[childName])];
		}
		return [...acc, next];
	}, []);
}

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
export const getCrumbs = (currentPath, routes, titleName = 'name') => {
	const homeRoute = routes.find((r) => r.path === '/');
	const segments = currentPath.split('/').filter(Boolean);
	let path = '';
	const breadcrumbs = [
		{
			title: homeRoute[titleName],
			path: homeRoute.path,
		},
	];
	// 两层侧边菜单时，第一层只是个组标题无实际路由，例如setting组标题，/setting/users
	const getUrl = (str) => {
		return str?.indexOf('/') === -1 ? `/${str}` : str;
	};

	for (let i = 0; i < segments.length; i++) {
		path += '/' + segments[i];

		// 依次尝试精确匹配和动态参数匹配
		let match = routes.find((r) => getUrl(r.path) === path);
		if (!match) {
			// 动态参数匹配
			match = routes.find((r) => matchRoute(getUrl(r.path), path));
		}
		if (match) {
			// 动态参数面包屑可自定义显示
			let title = match[titleName];
			// if (/:/.test(match.path)) {
			// 	// 取实际 id 或自定义
			// 	title = '详情'; // 或 segments[i]
			// }
			let crumbItem = {
				title,
			};
			// 最后一个不需要path
			if (segments.length - 1 !== i && match?.path?.indexOf('/') !== -1) {
				crumbItem.path = path;
			}
			breadcrumbs.push(crumbItem);
		}
	}
	return breadcrumbs;
};

// 根据pathname获取前端路由配置
export const getCurrentRouterItem = (currentPathname, allRouters, childName) => {
	const routers = flattenNavRouter(allRouters, childName);
	const routerItem = routers.find((item) => matchRoute(item.path, currentPathname));
	return routerItem;
};
