import { create } from "zustand"
export const useAuthStore = create((set) => ({
  userInfo: {
    id:1
  },
  login: (userData) => set({ userInfo: userData })
}))