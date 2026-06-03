import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

export const useGlobalStore = create(
	persist(
		(set) => ({
			sidebarCollapsed: false, // 侧边栏是否折叠, true:折叠
			crumbs: [],
			changeSidebarCollapsed: (value) => set({ sidebarCollapsed: value }),
			updateCrumbs: (v) => set({ crumbs: v }),
		}),
		{
			name: 'globalStore', // 存储唯一标识（必填）
			storage: createJSONStorage(() => sessionStorage),
			partialize: (state) => ({
				sidebarCollapsed: state.sidebarCollapsed,
			}),
		},
	),
);
