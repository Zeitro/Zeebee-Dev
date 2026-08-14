'use client';
import React, { createContext, useContext, useRef } from 'react';

/**
 * 应用侧 UI 注入（Table / Pagination / Modal 等带应用样式的封装）。
 * platform-core 不内置 NewTable / NewPagination / NewModal；由应用在根布局提供一次。
 * 新增槽位：AppUiProvider 多传一个 prop，消费侧 useResolvedUiComponent('SlotName')。
 * 嵌套 Provider 与上层合并，未传的槽位沿用上层。
 */

const AppUiContext = createContext({});

/**
 * @param {object} a
 * @param {object} b
 * @returns {boolean}
 */
function isSameSlots(a, b) {
	if (a === b) return true;
	const aKeys = Object.keys(a);
	const bKeys = Object.keys(b);
	if (aKeys.length !== bKeys.length) return false;
	return aKeys.every((key) => a[key] === b[key]);
}

/**
 * @param {object} props
 * @param {React.ReactNode} props.children
 * @param {Record<string, React.ComponentType>} [props] 其余均为 UI 槽位（Table / Pagination / Modal / ...）
 */
export function AppUiProvider({ children, ...slots }) {
	const parent = useContext(AppUiContext);
	const cacheRef = useRef(null);
	const prev = cacheRef.current;

	if (!prev || prev.parent !== parent || !isSameSlots(prev.slots, slots)) {
		const next = { ...parent };
		Object.entries(slots).forEach(([key, Comp]) => {
			if (Comp != null) next[key] = Comp;
		});
		cacheRef.current = { parent, slots, value: next };
	}

	return <AppUiContext.Provider value={cacheRef.current.value}>{children}</AppUiContext.Provider>;
}

/**
 * 读取已注入的 UI 槽位；缺省为 {}，不抛错。
 * @returns {Record<string, React.ComponentType|null|undefined>}
 */
export function useAppUi() {
	return useContext(AppUiContext);
}

/**
 * 解析某个 UI 槽位：props 覆盖优先，否则用 Context。新增组件只需传对应 slot 名。
 * @param {string} slot 如 'Table' | 'Pagination' | 'Modal'
 * @param {React.ComponentType|undefined|null} [override]
 * @returns {React.ComponentType}
 */
export function useResolvedUiComponent(slot, override) {
	const ui = useContext(AppUiContext);
	const Component = override ?? ui[slot];
	if (!Component) {
		throw new Error(
			`AppUi: pass ${slot} as prop or wrap with <AppUiProvider ${slot}={New${slot}} />`,
		);
	}
	return Component;
}

/** @deprecated 使用 AppUiProvider；保留以兼容已有 FilterTable 注入写法 */
export function FilterTableProvider({
	TableComponent,
	PaginationComponent,
	ModalComponent,
	children,
}) {
	return (
		<AppUiProvider Table={TableComponent} Pagination={PaginationComponent} Modal={ModalComponent}>
			{children}
		</AppUiProvider>
	);
}
