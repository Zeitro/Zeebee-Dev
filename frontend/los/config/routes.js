import {
	UserOutlined,
	BankOutlined,
} from '@ant-design/icons';
export const navRouter = [
	{
		path: '/',
		title: 'Home',
		hideCrumbs: true,
	},
	{
		path: '/loans',
		title: 'My loans',
		roles: ['loan'],
		icon:<BankOutlined />,
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
		icon:<UserOutlined />
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

export const routers = [
	...noNavRouter,
	...navRouter
];
