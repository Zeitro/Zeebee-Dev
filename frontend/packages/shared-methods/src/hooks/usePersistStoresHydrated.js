'use client';
import { useSyncExternalStore } from 'react';
import { usePermissionAccessStore, usePermissionPathsStore } from '../store/permission.js';
import { useUserStore } from '../store/user.js';

/** CustomizeAntdTheme 等首屏鉴权依赖的 persist store */
const defaultPersistStores = [usePermissionAccessStore, usePermissionPathsStore, useUserStore];

/**
 * 判断一组 zustand persist store 是否均已完成 storage 水合。
 * @param {Array<{ persist: { hasHydrated: () => boolean } }>} stores
 * @returns {boolean}
 */
function arePersistStoresHydrated(stores) {
	return stores.every((store) => store.persist.hasHydrated());
}

/**
 * 订阅多个 persist store 的水合完成事件。
 * @param {Array<{ persist: { onFinishHydration: (fn: () => void) => () => void } }>} stores
 * @param {() => void} onStoreChange
 * @returns {() => void}
 */
function subscribePersistStoresHydrated(stores, onStoreChange) {
	const unsubscribes = stores.map((store) => store.persist.onFinishHydration(onStoreChange));
	return () => {
		unsubscribes.forEach((unsubscribe) => unsubscribe?.());
	};
}

/**
 * 等待权限 / 用户相关 persist store 水合后再做客户端鉴权，避免刷新首帧误判未登录。
 * @param {Array<{ persist: { hasHydrated: () => boolean, onFinishHydration: (fn: () => void) => () => void } }>} [stores]
 * @returns {boolean}
 */
export function usePersistStoresHydrated(stores = defaultPersistStores) {
	return useSyncExternalStore(
		(onStoreChange) => subscribePersistStoresHydrated(stores, onStoreChange),
		() => arePersistStoresHydrated(stores),
		() => false,
	);
}
