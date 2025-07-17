export const navRouter = [
	{
		path: '/home',
		title: 'Home',
		hideCrumbs: true,
	},
	{
		path: '/loans',
		title: 'My loans',
		roles: [],
		router: [
			{
				path: '/details',
				title: 'details',
				hideMenu: true,
			},
		],
	},
	{
		path: '/properties',
		title: 'My properties',
		roles: [],
	},
];

export const noNavRouter = [
	{
		path: '/login',
		title: 'Login',
	},
];

export const routers = [
	...noNavRouter,
	{
		path: '/',
		router: navRouter,
	},
];

