'use client';
import { useState, useEffect, useRef } from 'react';
/**
 * 监听元素的宽高变化
 * @param {*} rectDom 当已有dom元素时，可传入dom元素来监听
 * @returns 元素/元素的宽高和距离xy轴距离
 */
export function useElementSize(rectDom) {
	const [size, setSize] = useState({ width: 0, height: 0 });
	const elementRef = useRef(null);

	useEffect(() => {
		let observer = null;
		const targetElement = elementRef?.current || rectDom;
		if (targetElement) {
			observer = new ResizeObserver((entries) => {
				setSize(entries[0].contentRect);
			});

			observer.observe(targetElement);
		}
		return () => {
			if (observer && targetElement) {
				observer?.unobserve(targetElement);
			}
		};
	}, [rectDom]);

	return [size, rectDom || elementRef];
}
