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
				path: '/loans/createLoan',
				title: 'Add 3rd Party Loan',
				hideMenu: true,
			},
			{
				path: '/loans/:id',
				title: 'Edit 3rd Party Loan',
				hideMenu: true,
			},
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
	{
		path: '/estimateLoan',
		title: 'Estimate loan',
		hideMenu: true,
		hideCrumbs: true,
		router: [
			{
				path: '/estimateLoan/preApproved',
				title: 'Pre-Approved!',
				hideMenu: true,
				hideCrumbs: true,
			},
		],
	},
];

export const noNavRouter = [
	...noTokenRouter,
	{
		path: '/pdf',
		title: 'pdf',
		hideMenu: true,
		hideCrumbs: true,
	},
	{
		path: '/not-found',
		title: '404',
		hideCrumbs: true,
		hideMenu: true,
	},
];

export const routers = [...noNavRouter, ...navRouter];
