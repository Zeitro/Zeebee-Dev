export const navRouter = [
	{
		path: '/',
		title: 'Home',
		hideCrumbs: true,
	},
	{
		path: '/loans',
		title: 'All Loans',
		icon: 'loan',
		router: [
			{
				path: '/loans/zeeBeeDetails',
				title: 'ZeeBee Loan Details',
				hideMenu: true,
			},
		],
	},
	{
		path: '/lenders',
		title: 'Lenders',
		icon: 'business',
	},
	{
		path: '/editorTemplate',
		title: 'Template',
		icon: 'template',
	},
	{
		path: '/role',
		title: 'App resource',
		icon: 'role',
		router: [
			{
				path: '/role/appApi',
				title: 'Api',
				// hideMenu: true,
			},
			{
				path: '/role/resource',
				title: 'Resource',
				// hideMenu: true,
			},
		],
	},
	{
		path: '/setting',
		title: 'Setting',
		icon: 'setting',
		router: [
			{
				path: '/setting/role',
				title: 'Role',
				// hideMenu: true,
			},
			{
				path: '/setting/user',
				title: 'User',
				// hideMenu: true,
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
			{
				path: '/login/joinOrganization',
				title: 'Join Organization',
				hideMenu: true,
				hideCrumbs: true,
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

// export const noNavRouter = [
// 	...noTokenRouter,
// 	{
// 		path: '/not-found',
// 		title: '404',
// 		hideCrumbs: true,
// 		hideMenu: true,
// 	},
// ];

// export const routers = [...noNavRouter, ...navRouter];
