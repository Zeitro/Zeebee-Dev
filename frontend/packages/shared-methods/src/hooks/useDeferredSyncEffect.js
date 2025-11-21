import { useEffect, useRef } from 'react';

/**
 * 单个判断条件
 * useDeferredSyncEffect({ predicate, immediate, deferred, deps, delay = 0, isUseMicrotask = false })
 *	解决：useEffect 的异步执行，导致报错
 * - predicate: () => boolean  可选，返回 true 则执行
 * - notPassedCallback: ()=>{} predicate返回false时执行
 * - immediate: () => void     立即执行（同步）回调，适合 form.setFieldsValue 等
 * - deferred: () => void      异步执行的回调，适合 setState（会触发 React 渲染）
 * - deps: Array               effect 依赖
 * - delay: number             延迟 ms（默认 0，使用 setTimeout）
 * - isUseMicrotask: boolean     若为 true，使用 Promise.resolve().then(...)（微任务）
 */
export function useDeferredSyncEffect({
	predicate = () => true,
	notPassedCallback = () => {},
	immediate,
	deferred,
	deps = [],
	delay = 0,
	isUseMicrotask = false,
}) {
	const cancelledRef = useRef(false);

	useEffect(() => {
		cancelledRef.current = false;

		if (!predicate?.()) {
			notPassedCallback?.();
			return;
		}

		try {
			immediate?.();
		} catch {
			// ignore immediate errors or rethrow as needed
			// console.error(e);
		}

		if (!deferred) return;

		let timerId = null;
		const runDeferred = () => {
			if (!cancelledRef.current) {
				try {
					deferred?.();
				} catch {
					// console.error(e);
				}
			}
		};

		if (isUseMicrotask) {
			Promise.resolve().then(runDeferred);
		} else {
			timerId = setTimeout(runDeferred, delay);
		}

		return () => {
			cancelledRef.current = true;
			if (timerId) clearTimeout(timerId);
		};
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, deps);
}

/**
 * 多个判断条件
 * useDeferredSyncEffects({ rules: [{ predicate, notPassedCallback, immediate, deferred, delay, isUseMicrotask }], deps })
 */
export function useDeferredSyncEffects({ rules = [], deps = [] }) {
	const cancelledRef = useRef(false);

	useEffect(() => {
		cancelledRef.current = false;
		const timers = [];

		for (const rule of rules) {
			try {
				const shouldRun =
					typeof rule.predicate === 'function'
						? rule?.predicate?.()
						: Boolean(rule.predicate);
				if (!shouldRun) {
					rule?.notPassedCallback?.();
					continue;
				}

				// immediate 执行（同步），适合 form.setFieldsValue 等外部系统更新
				try {
					rule?.immediate?.();
				} catch {
					// ignore or log
				}

				// deferred 延迟执行（适合 setState）
				const delay = Number(rule.delay ?? 0);
				const isUseMicrotask = Boolean(rule.isUseMicrotask);
				// 微任务
				if (isUseMicrotask) {
					Promise.resolve().then(() => {
						if (!cancelledRef.current) {
							try {
								rule?.deferred?.();
							} catch {
								// ignore or log
							}
						}
					});
				} else {
					// 宏任务
					const id = setTimeout(() => {
						if (!cancelledRef.current) {
							try {
								rule?.deferred?.();
							} catch {
								// ignore or log
							}
						}
					}, delay);
					timers.push(id);
				}
			} catch {
				// ignore predicate errors
			}
		}

		return () => {
			cancelledRef.current = true;
			for (const t of timers) clearTimeout(t);
		};
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, deps);
}
