import type { NextConfig } from 'next';
const path = require('path');
const nextConfig: NextConfig = {
	/* config options here */
	reactStrictMode: true, // 开启 React 严格模式
	env: {
		NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL, // 注入环境变量
	},
	webpack: (config) => {
		config.resolve.alias['@public'] = path.resolve(__dirname, 'public');
		config.resolve.alias['@config'] = path.resolve(__dirname, 'config');
		return config;
	},
	async rewrites() {
		if (process.env.NEXT_PUBLIC_ENV === 'local') {
			return [
				{
					source: '/api/v1/:path*',
					destination: process.env.NEXT_PUBLIC_API_URL + '/api/v1/:path*',
				},
			];
		}
		return [];
	},
};

export default nextConfig;
