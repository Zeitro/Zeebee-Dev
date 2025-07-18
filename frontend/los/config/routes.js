export const navRouter = [
	{
		path: '/',
		title: 'Home',
		hideCrumbs: true,
	},
	{
		path: '/loans',
		title: 'My loans',
		roles: [],
		router: [
			{
				path: '/loans/details',
				title: 'details',
				hideMenu: true,
			},
			{
				path: '/loans/:id',
				title: 'details',
					hideMenu: true,
				router: [
					{
						path: '/loans/:id/type',
						title: 'types',
						hideMenu: true,
					},
				],
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
