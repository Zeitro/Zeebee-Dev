'use client';
import { useEffect, useRef } from 'react';

/**
 * 公共自定义hooks - useEffect相关
 * @param {*} fn
 * @param {*} inputs
 */
// 初始化渲染一次
export function useOnMountUnsafe(fn, inputs = []) {
	const initialized = useRef(false);
	useDisUpdateEffect(() => {
		if (!initialized.current) {
			initialized.current = true;
			fn();
		}
	}, inputs);
}

// 初始化不渲染
export function useDisUpdateEffect(fn, inputs = []) {
	const didMountRef = useRef(false);
	useEffect(() => {
		if (didMountRef.current) fn();
		else didMountRef.current = true;
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, inputs);
}
