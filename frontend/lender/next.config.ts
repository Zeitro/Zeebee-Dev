import type { NextConfig } from 'next';
import path from 'path';

const isLocal = process.env.NEXT_PUBLIC_ENV === 'local';
const nextConfig: NextConfig = {
	/* config options here */
	output: isLocal ? undefined : 'standalone',
	reactStrictMode: !isLocal, // 开启 React 严格模式
	serverExternalPackages: ['react-pdf', 'pdfjs-dist'],
	env: {
		NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL, // 注入环境变量
	},
	webpack: (config) => {
		config.module.rules.push({
			test: /\.svg$/,
			use: ['@svgr/webpack'],
		});
		config.resolve.alias['@config'] = path.resolve(__dirname, 'config');
		return config;
	},
	async rewrites() {
		if (isLocal) {
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
