import { create } from "zustand"
export const useGlobalStore = create((set) => ({
  siberWidth: 216,
  changeSiberWidth: (value) => set({ siberWidth: value })
}))