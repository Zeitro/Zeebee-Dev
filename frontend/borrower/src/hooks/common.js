import { useState, useEffect, useRef } from 'react';

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
		// eslint-disable-next-line react-hooks/exhaustive-deps
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

// 获取上一次更新的值，可用在effect内比较新旧值
export function usePrevious(value) {
	const ref = useRef();
	useEffect(() => {
		ref.current = value;
	});
	return ref.current;
}

export const useUpdate = (fn, inputs = []) => {
	const [count, setCount] = useState(0);
	useEffect(() => {
		setCount((num) => num + 1);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [...inputs]);

	useEffect(() => {
		if (count > 1) fn();
	}, [count, fn]);
};
