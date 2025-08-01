import { dirname } from "path";
import { fileURLToPath } from "url";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
  ...compat.config({
    extends: [
      'eslint:recommended',
      'plugin:@typescript-eslint/recommended',
      'plugin:react/recommended',
      'prettier',
      'plugin:prettier/recommended',
      'next/core-web-vitals',
      'next/typescript',
    ],
    parser: '@typescript-eslint/parser',
    globals:{
      window: true
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
      "no-sparse-arrays": "off", //允许数组空插槽
      'eol-last': ['error', 'always'],
      'no-unused-vars': 'warn',
      'prettier/prettier': 'error',
      'react/react-in-jsx-scope': 'off',
      'no-console': ['error', { allow: ['warn', 'error'] }],
      'no-var': 'error',
      'prefer-const': 'error',
      'react/no-unescaped-entities': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_' },
      ],
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
