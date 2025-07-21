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
};

export default nextConfig;
