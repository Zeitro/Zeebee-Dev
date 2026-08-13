module.exports = {
	content: [
		'./src/**/*.{js,ts,jsx,tsx}', // 根据你的项目结构调整
		'./src/**/*.{css,scss}',
		'../packages/shared-methods/src/**/*.{js,jsx,ts,tsx}'
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
					label: 'var(--base-label)',
					gray: 'var(--base-gray)',
					active: 'var(--base-active)',
					text: 'var(--base-text)',
					text2: 'var(--base-text2)',
				},
			},
			fontFamily: {
				roboto: ['var(--font-roboto)', 'sans-serif'],
			},
		},
	},
	plugins: [],
};
