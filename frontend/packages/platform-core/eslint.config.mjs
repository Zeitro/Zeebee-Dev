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

/**
 * 与 shared-methods 对齐；platform-core 无 Next 页面，仍复用同一套规则以免风格分叉。
 */
const eslintConfig = [
	...compat.config({
		extends: [
			'eslint:recommended',
			'plugin:@typescript-eslint/recommended',
			'plugin:react/recommended',
			'plugin:prettier/recommended',
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
			'@typescript-eslint/no-require-imports': 'off',
			'@typescript-eslint/no-var-requires': 'off',
			'no-sparse-arrays': 'off',
			'eol-last': ['error', 'always'],
			'no-unused-vars': 'warn',
			'prettier/prettier': 'error',
			'react/react-in-jsx-scope': 'off',
			'no-var': 'error',
			'react/no-unescaped-entities': 'off',
			'@typescript-eslint/no-explicit-any': 'off',
			'@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
			'@typescript-eslint/no-non-null-assertion': 'warn',
			'react/jsx-key': 'error',
			'react/jsx-no-duplicate-props': 'error',
			'react/jsx-no-undef': 'error',
			'react/no-children-prop': 'error',
			'react/no-danger-with-children': 'error',
			'react/no-deprecated': 'error',
			'react/no-direct-mutation-state': 'error',
			'react/no-find-dom-node': 'error',
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
			'node_modules/**/*',
			'dist/**/*',
			'build/**/*',
			'*.config.js',
			'*.config.mjs',
		],
	},
];

export default eslintConfig;
