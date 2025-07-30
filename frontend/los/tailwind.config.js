module.exports = {
	content: [
		'./src/**/*.{js,ts,jsx,tsx}', // 根据你的项目结构调整
		'./src/**/*.{css,scss}',
	],
	theme: {
		extend: {
			colors: {
				background: 'var(--background)',
				foreground: 'var(--foreground)',
				base: {
					primary: 'var(--base-primary)',
					success: 'var(--base-success)',
					warning: 'var(--base-warning)',
					error: 'var(--base-error)',
					disabled: 'var(--base-disabled)',
					border: 'var(--base-border)',
					label:'var(--base-label)'
				},
				black: {
					light: 'var(--black-light)', // 浅黑色
					medium: 'var(--black-medium)', // 中黑色
					dark: 'var(--black-dark)', // 深黑色
				},
				gray: {
					// 灰色调
					light: 'var(--gray-light)', // 浅灰色
					medium: 'var(--gray-medium)', // 中灰色
					dark: 'var(--gray-dark)', // 深灰色
				},
				blue: {
					light: 'var(--blue-light)',
					medium: 'var(--blue-medium)',
					dark: 'var(--blue-dark)',
				},
				green:{
					light: 'var(--green-light)',
					medium:'var(--green-medium)'
				},
				purple:{
					light:'var(--purple-light)',
					medium:'var(--purple-medium)',
				}
				
			},
			fontFamily: {
				roboto: ['var(--font-roboto)', 'sans-serif'],
			},
		},
	},
	plugins: [],
};
