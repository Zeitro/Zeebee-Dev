import { create } from 'zustand';
export const useGlobalStore = create((set) => ({
	siberCollapsed: false, // 侧边栏是否折叠, true:折叠
	crumbs: [],
	changeSiberCollapsed: (value) => set({ siberCollapsed: value }),
	updateCrumbs: (v) => set({ crumbs: v }),
}));
