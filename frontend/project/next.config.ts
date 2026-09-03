import type { NextConfig } from 'next';
import path from 'path';
import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./src/i18n/request.js');

const isLocal = process.env.NEXT_PUBLIC_ENV === 'local';
const nextConfig: NextConfig = {
	/* config options here */
	// 仅在生产环境构建时启用 'standalone' 模式
	// output: process.env.NEXT_PUBLIC_ENV === 'local' ? undefined : 'standalone',
	reactStrictMode: !isLocal, // 开启 React 严格模式
	env: {
		NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL, // 注入环境变量
	},
	transpilePackages: ['@repo/platform-core'],
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

export default withNextIntl(nextConfig);
