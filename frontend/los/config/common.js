// 生成菜单栏
export const generateMenuBar = (list) => {
	return list
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
