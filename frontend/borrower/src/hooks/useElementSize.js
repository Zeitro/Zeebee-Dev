import { useState, useEffect, useRef } from 'react';
/**
 * 监听元素的宽高变化
 * @param {*} rectDom 当已有dom元素时，可传入dom元素来监听
 * @returns 元素/元素的宽高和距离xy轴距离
 */
export default function useElementSize(rectDom) {
	const [size, setSize] = useState({ width: 0, height: 0 });
	const elementRef = useRef(null);

	useEffect(() => {
		if (elementRef?.current || rectDom) {
			const observer = new ResizeObserver((entries) => {
				setSize(entries[0].contentRect);
			});

			(elementRef?.current || rectDom) && observer.observe(elementRef?.current || rectDom);
		}
		// return () => observer.disconnect();
		return () => {
			(elementRef?.current || rectDom) && observer?.unobserve(elementRef?.current || rectDom);
		};
	}, [elementRef, rectDom]);

	return [size, elementRef || rectDom];
}
