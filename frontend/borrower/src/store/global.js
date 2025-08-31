import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
export const useGlobalStore = create(
	persist(
		(set) => ({
			siberCollapsed: false, // 侧边栏是否折叠, true:折叠
			crumbs: [],
			changeSiberCollapsed: (value) => set({ siberCollapsed: value }),
			updateCrumbs: (v) => set({ crumbs: v }),
		}),
		{
			name: 'globalStore', // 存储唯一标识（必填）
			storage: createJSONStorage(() => sessionStorage),
			partialize: (state) => ({
				siberCollapsed: state.siberCollapsed,
			}),
		},
	),
);
