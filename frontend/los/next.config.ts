import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    /* config options here */
    reactStrictMode: true, // 开启 React 严格模式
    env: {
      NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL, // 注入环境变量
    },
};

export default nextConfig;
