import { create } from "zustand"
export const useGlobalStore = create((set) => ({
  siberCollapsed: false,
  changeSiberCollapsed: (value) => set({ siberCollapsed: value })
}))