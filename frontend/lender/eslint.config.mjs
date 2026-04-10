import {dirname} from 'path';
import {fileURLToPath} from 'url';
import {FlatCompat} from '@eslint/eslintrc';
import js from '@eslint/js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
	baseDirectory: __dirname,
	recommendedConfig: js.configs.recommended,
});

const eslintConfig = [
	...compat.config({
		extends: [
			'eslint:recommended',
			'plugin:@typescript-eslint/recommended',
			'plugin:react/recommended',
			'plugin:react-hooks/recommended',
			'plugin:prettier/recommended',
			'next/core-web-vitals',
			'next/typescript',
			'prettier',
		],
		parser: '@typescript-eslint/parser',
		globals: {
			window: true,
		},
		parserOptions: {
			ecmaVersion: 2020,
			sourceType: 'module',
			ecmaFeatures: {
				arrowFunctions: true,
				classes: true,
				modules: true,
				defaultParams: true,
				experimentalObjectRestSpread: true,
			},
		},
		plugins: ['@typescript-eslint', 'react', 'react-hooks', 'prettier'],
		rules: {
			// 禁止使用 require() 的导入写法：off关闭
			'@typescript-eslint/no-require-imports': 'off',
			// 禁止在顶层用 var/let/const 赋值形式的 require（如 const x = require('x')）
			'@typescript-eslint/no-var-requires': 'off',
			'no-sparse-arrays': 'off', //允许数组空插槽
			// 强制文件末尾必须有一个换行
			'eol-last': ['error', 'always'],
			'no-unused-vars': 'warn', // 未使用变量/参数只警告
			'prettier/prettier': 'error', // 代码不符合 prettier 会报错
			'react/react-in-jsx-scope': 'off', // 不需要显式 import React
			// 'no-console': ['error', { allow: ['warn', 'error'] }],
			'no-var': 'error', // 禁止使用 var，要求用 let/const
			// 'prefer-const': 'error',
			'react/no-unescaped-entities': 'off', // 禁止 JSX 中未转义的字符：关闭
			'@typescript-eslint/no-explicit-any': 'off',
			'@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
			'@next/next/no-img-element': 'error', // next/image 替代 <img>
			'@typescript-eslint/no-non-null-assertion': 'warn', // 禁止使用 !（非空断言）
			'react/jsx-key': 'error', // 在数组渲染中要求每个元素有 key
			'react/jsx-no-duplicate-props': 'error', // 禁止重复的 JSX 属性
			'react/jsx-no-undef': 'error', // 禁止在 JSX 中使用未定义的变量/组件
			'react/no-children-prop': 'error', // 禁止直接在组件上使用 children 属性
			'react/no-danger-with-children': 'error', // 禁止同时使用 dangerouslySetInnerHTML 与 children
			'react/no-deprecated': 'error', //  警告使用 React 的已废弃 API
			'react/no-direct-mutation-state': 'error', // 禁止直接修改组件 state
			'react/no-find-dom-node': 'error', // 禁止使用 findDOMNode（不安全）
			'react/no-render-return-value': 'error', // 禁止依赖 ReactDOM.render 的返回值
			'react/no-string-refs': 'error', // 禁止字符串 refs（推荐使用回调或 createRef）
			'react/no-unknown-property': 'error', // 禁止使用未知的 DOM 属性
			'react/no-unsafe': 'warn', //  检查不安全的生命周期（如 UNSAFE_componentWillMount）
			'react/prop-types': 'off', // 要求使用 prop-types（在 TS 项目通常关闭）
			'react/require-render-return': 'error', // 要求 render 方法必须返回内容
			'react-hooks/rules-of-hooks': 'error', // 检查 Hook 的规则
			// 'react-hooks/exhaustive-deps': 'off', // useEffect 依赖检测
		},
		settings: {
			react: {
				version: 'detect',
			},
		},
	}),
	{
		ignores: [
			'.next/**/*',
			'node_modules/**/*',
			'dist/**/*',
			'build/**/*',
			'*.config.js',
			'*.config.mjs',
			'public/**/*',
		],
	},
];

export default eslintConfig;
