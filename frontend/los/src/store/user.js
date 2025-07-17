import { create } from "zustand"
export const useUserStore = create((set) => ({
  userInfo: {
    id:1
  },
  login: (userData) => set({ userInfo: userData })
}))