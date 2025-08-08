import { create } from 'zustand';
export const useGlobalStore = create((set) => ({
	siberCollapsed: false,
	crumbs: [],
	changeSiberCollapsed: (value) => set({ siberCollapsed: value }),
	updateCrumbs: (v) => set({ crumbs: v }),
}));
