import { dirname } from 'path';
import { fileURLToPath } from 'url';
import { FlatCompat } from '@eslint/eslintrc';
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
		plugins: ['@typescript-eslint', 'react', 'prettier'],
		rules: {
			// 允许所有require()形式的模块导入语句
			"@typescript-eslint/no-require-imports": "off",
			// 允许通过const/let/var变量赋值的require()调用允许
			'@typescript-eslint/no-var-requires': 'off',
			'no-sparse-arrays': 'off', //允许数组空插槽
			'eol-last': ['error', 'always'],
			'no-unused-vars': 'warn',// 未使用变量/参数只警告
			'prettier/prettier': 'error',
			'react/react-in-jsx-scope': 'off',
			// 'no-console': ['error', { allow: ['warn', 'error'] }],
			'no-var': 'error',
			// 'prefer-const': 'error',
			'react/no-unescaped-entities': 'off',
			'@typescript-eslint/no-explicit-any': 'off',
			'@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
			'@next/next/no-img-element': 'error',
			'@typescript-eslint/no-non-null-assertion': 'warn',
			'react/jsx-key': 'error',
			'react/jsx-no-duplicate-props': 'error',
			'react/jsx-no-undef': 'error',
			'react/no-children-prop': 'error',
			'react/no-danger-with-children': 'error',
			'react/no-deprecated': 'error',
			'react/no-direct-mutation-state': 'error',
			'react/no-find-dom-node': 'error',
			'react/no-is-mounted': 'error',
			'react/no-render-return-value': 'error',
			'react/no-string-refs': 'error',
			'react/no-unknown-property': 'error',
			'react/no-unsafe': 'warn',
			'react/prop-types': 'off',
			'react/require-render-return': 'error',
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
