export const navRouter = [
	{
		path: '/',
		title: 'Home',
		hideCrumbs: true,
	},
	{
		path: '/loans',
		title: 'My loans',
		icon: 'loan',
		roles: [],
		router: [
			{
				path: '/createLoan',
				title: 'Create A New Loan',
				hideMenu: true,
			},
			{
				path: '/loans/:id',
				title: 'Edit Loan',
				hideMenu: true,
			},
		],
	},
	{
		path: '/properties',
		title: 'My properties',
		icon: 'property',
		roles: [],
		router: [
			{
				path: '/properties/:id',
				title: 'Details',
				hideMenu: true,
			},
		],
	},
	{
		path: '/not-found',
		title: '404',
		hideCrumbs: true,
		hideMenu: true,
	},
];

export const noNavRouter = [
	{
		path: '/login',
		title: 'Login',
	},
];

export const routers = [...noNavRouter, ...navRouter];
