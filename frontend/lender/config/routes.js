export const navRouter = [
	{
		path: '/',
		title: 'Home',
		hideCrumbs: true,
	},
	{
		path: '/loans',
		title: 'All loans',
		icon: 'loan',
		roles: [],
		router: [
			{
				path: '/loans/createZeeBeeLoan',
				title: 'Apply ZeeBee Loan',
				hideMenu: true,
			},
			{
				path: '/loans/zeeBeeDetails',
				title: 'ZeeBee Loan Details',
				hideMenu: true,
			},
		],
	},
];

export const noTokenRouter = [
	{
		path: '/login',
		title: 'Login',
		hideMenu: true,
		hideCrumbs: true,
		router: [
			{
				path: '/login/resetPassword',
				title: 'Reset Password',
				hideMenu: true,
				hideCrumbs: true,
			},
			{
				path: '/login/register',
				title: 'Register',
				hideMenu: true,
				hideCrumbs: true,
			},
			{
				path: '/login/sendRegisterTips',
				title: 'Send Register Tips',
				hideMenu: true,
				hideCrumbs: true,
			},
			{
				path: '/login/verifyEmail',
				title: 'Verify Email',
				hideMenu: true,
				hideCrumbs: true,
			},
		],
	},
];

export const noNavRouter = [
	...noTokenRouter,
	{
		path: '/not-found',
		title: '404',
		hideCrumbs: true,
		hideMenu: true,
	},
];

export const routers = [...noNavRouter, ...navRouter];
